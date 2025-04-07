<template>
  <div class="xterm-container">
    <template v-for="conn in connectionStore.connections" :key="conn.id">
      <n-result
        v-if="conn.connectionError"
        status="error"
        :title="conn.errorMessage"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #footer>
          <n-button @click="() => reconnect(conn.id)" type="primary">{{ $t('terminal.reconnect') }}</n-button>
        </template>
        <template #default>
          <div class="error-details" v-if="conn.errorDetails">
            <n-collapse class="error-collapse">
              <n-collapse-item :title="$t('terminal.error.details')" name="details">
                <n-code :code="conn.errorDetails" language="bash" :word-wrap="true" />
              </n-collapse-item>
            </n-collapse>
          </div>
        </template>
      </n-result>
      <n-result
        v-else-if="conn.isFingerprintConfirm"
        status="warning"
        :title="$t('terminal.fingerprint.confirm')"
        :description="$t('terminal.fingerprint.description')"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #icon>
          <n-icon size="48">
            <icon icon="ph:fingerprint" />
          </n-icon>
        </template>
        <template #footer>
          <n-space>
            <n-button @click="() => rejectFingerprint(conn.id)" type="error">
              {{ $t('terminal.fingerprint.reject') }}
            </n-button>
            <n-button @click="() => acceptFingerprint(conn.id)" type="primary">
              {{ $t('terminal.fingerprint.accept') }}
            </n-button>
          </n-space>
        </template>
        <template #default>
          <div class="fingerprint-details">
            <p>
              <strong>{{ $t('terminal.fingerprint.host') }}:</strong> {{ conn.hostAddress }}
            </p>
            <p>
              <strong>{{ $t('terminal.fingerprint.fingerprint') }}:</strong>
            </p>
            <n-code :code="conn.hostFingerprint || ''" language="bash" :word-wrap="true" />
          </div>
        </template>
      </n-result>
      <n-result
        v-else-if="conn.isConnecting || !connectedTerminals[conn.id]"
        status="info"
        :title="$t('terminal.connecting')"
        :description="$t('terminal.connecting_desc')"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #icon>
          <n-spin size="large" />
        </template>
      </n-result>
      <div
        v-show="isTerminalVisible(conn)"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
        :ref="el => el && (terminalRefs[conn.id] = el as HTMLElement)"
        class="xterm-wrapper"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { loadTheme } from '@/themes/xtermjs';
import { WebsocketPort } from '@wailsApp/go/services/TerminalSrv';
import { CanvasAddon } from '@xterm/addon-canvas';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { WebglAddon } from '@xterm/addon-webgl';
import { Terminal } from '@xterm/xterm';
import '@xterm/xterm/css/xterm.css';
import { useConnectionStore } from '@/stores/connection';
import { NCode, NCollapse, NCollapseItem, NResult, NSpin, NButton, NSpace, NIcon } from 'naive-ui';
import { onActivated } from 'vue';
import { useI18n } from 'vue-i18n';
import { enums } from '@wailsApp/go/models';

const { t } = useI18n();

const connectionStore = useConnectionStore();
const connectionTabs = inject<any>('connectionTabs');
const activeConn = computed(() => connectionStore.activeConnection);

const isTerminalHidden = (connId: number) => {
  return connId !== activeConn.value?.id;
};

const isTerminalVisible = (conn: any) => {
  return !conn.connectionError && !conn.isConnecting && !conn.isFingerprintConfirm && connectedTerminals.value[conn.id];
};

const terminalRefs = ref<Record<number, HTMLElement | null>>({});
const terminals = ref<Record<number, Terminal | undefined>>({});
const sockets = ref<Record<number, WebSocket | undefined>>({});
const fitAddons = ref<Record<number, FitAddon | undefined>>({});
const webLinksAddons = ref<Record<number, WebLinksAddon | undefined>>({});
const webglAddons = ref<Record<number, WebglAddon | undefined>>({});
const canvasAddons = ref<Record<number, CanvasAddon | undefined>>({});
const connectedTerminals = ref<Record<number, boolean>>({});

const updateStatus = (
  id: number,
  status: Partial<{
    isConnecting: boolean;
    connectionError: boolean;
    isFingerprintConfirm: boolean;
    errorMessage: string;
    errorDetails: string;
    hostAddress: string;
    hostFingerprint: string;
  }>,
) => {
  connectionStore.updateConnectionStatus(id, status);
  if (connectionTabs?.value) {
    let tabStatus = 'connecting';
    if (status.connectionError) {
      tabStatus = 'error';
    } else if (status.isFingerprintConfirm) {
      tabStatus = 'warning';
    } else if (!status.isConnecting && connectedTerminals.value[id]) {
      tabStatus = 'connected';
    }
    connectionTabs.value.updateTabStatus(id, tabStatus);
  }
};

const acceptFingerprint = (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn || !sockets.value[id]) return;

  sockets.value[id]?.send(
    JSON.stringify({
      type: enums.TerminalType.FINGERPRINTCONFIRM,
      accept: true,
    }),
  );

  updateStatus(id, {
    isFingerprintConfirm: false,
    isConnecting: true,
  });
};

const rejectFingerprint = (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn || !sockets.value[id]) return;

  sockets.value[id]?.send(
    JSON.stringify({
      type: enums.TerminalType.FINGERPRINTCONFIRM,
      accept: false,
    }),
  );

  updateStatus(id, {
    isFingerprintConfirm: false,
    connectionError: true,
    errorMessage: t('terminal.fingerprint.rejected'),
  });

  sockets.value[id]?.close();
  sockets.value[id] = undefined;
};

const fitXterm = (id: number) => {
  const dims = fitAddons.value[id]?.proposeDimensions();
  if (!dims || !terminals.value[id] || !dims.cols || !dims.rows) return;
  if (terminals.value[id].rows !== dims.rows || terminals.value[id].cols !== dims.cols) {
    terminals.value[id].resize(dims.cols, dims.rows);
  }
};

const fitAllTerminals = () => {
  requestAnimationFrame(() => {
    Object.keys(terminals.value).forEach(id => {
      const numId = Number(id);
      if (terminals.value[numId]) {
        fitXterm(numId);
      }
    });
  });
};

const initializeTerminal = async (id: number) => {
  if (terminals.value[id]) return;

  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn) return;

  const theme = await loadTheme(conn.theme || 'Default');
  terminals.value[id] = new Terminal({
    convertEol: true,
    disableStdin: false,
    fontSize: 16,
    cursorBlink: true,
    cursorStyle: 'bar',
    theme: theme,
  });
};

const initializeXterm = async (id: number) => {
  const terminalEl = terminalRefs.value[id] as HTMLElement;
  const terminal = terminals.value[id];
  if (!terminalEl || !terminal) return;

  fitAddons.value[id] = new FitAddon();
  webLinksAddons.value[id] = new WebLinksAddon();
  webglAddons.value[id] = new WebglAddon();

  try {
    fitAddons.value[id].activate(terminal);
    webLinksAddons.value[id].activate(terminal);
    webglAddons.value[id].activate(terminal);
  } catch (e) {
    console.warn('Failed to activate WebGL addon, falling back to Canvas:', e);
    webglAddons.value[id] = undefined;
    canvasAddons.value[id] = new CanvasAddon();
    canvasAddons.value[id].activate(terminal);
  }

  terminal.attachCustomKeyEventHandler(arg => {
    if (arg.code === 'PageUp' && arg.type === 'keydown') {
      terminal?.scrollPages(-1);
      return false;
    } else if (arg.code === 'PageDown' && arg.type === 'keydown') {
      terminal?.scrollPages(1);
      return false;
    }
    return true;
  });

  terminal.open(terminalEl);
  terminal.onData(data => sockets.value[id]?.send(JSON.stringify({ type: enums.TerminalType.CMD, cmd: data })));
  terminal.onResize(({ cols, rows }) => {
    if (sockets.value[id]?.readyState === WebSocket.OPEN) {
      sockets.value[id]?.send(JSON.stringify({ type: enums.TerminalType.RESIZE, cols, rows }));
    }
  });

  await nextTick();
  fitXterm(id);
};

const initializeWebsocket = async (id: number, hostId: number) => {
  try {
    if (sockets.value[id]?.readyState === WebSocket.OPEN) {
      updateStatus(id, { isConnecting: false });
      return;
    }

    if (sockets.value[id]) {
      sockets.value[id].close();
      sockets.value[id] = undefined;
    }

    updateStatus(id, {
      isConnecting: true,
      connectionError: false,
      isFingerprintConfirm: false,
      errorMessage: '',
      errorDetails: '',
      hostAddress: '',
      hostFingerprint: '',
    });

    const port = await WebsocketPort();
    sockets.value[id] = new WebSocket(`ws://localhost:${port}/ws/terminal?hostId=${hostId}`);
    const socket = sockets.value[id];
    if (!socket) return;

    socket.onopen = () => {
      updateStatus(id, { isConnecting: false });
      nextTick(() => {
        terminals.value[id]?.focus();
        fitXterm(id);
      });
    };

    socket.onmessage = async (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      switch (data.type) {
        case enums.TerminalType.ERROR:
          updateStatus(id, {
            isConnecting: false,
            connectionError: true,
            errorMessage: data.error,
            errorDetails: data.details,
          });
          connectedTerminals.value[id] = false;
          socket?.close();
          break;
        case enums.TerminalType.FINGERPRINTCONFIRM:
          updateStatus(id, {
            isConnecting: false,
            isFingerprintConfirm: true,
            hostAddress: data.host,
            hostFingerprint: data.fingerprint,
          });
          break;
        case enums.TerminalType.CONNECTED:
          updateStatus(id, { isConnecting: false });
          connectedTerminals.value[id] = true;
          if (connectionTabs?.value) {
            connectionTabs.value.updateTabStatus(id, 'connected');
          }
          nextTick(() => {
            terminals.value[id]?.focus();
            fitXterm(id);
          });
          break;
        case enums.TerminalType.DATA:
          terminals.value[id]?.write(data.content);
          break;
      }
    };

    socket.onerror = () => {
      updateStatus(id, {
        isConnecting: false,
        connectionError: true,
        errorMessage: t('terminal.error.connection'),
      });
      connectedTerminals.value[id] = false;
      socket?.close();
    };

    socket.onclose = event => {
      if (!connectionStore.connections.find(c => c.id === id)?.connectionError) {
        updateStatus(id, {
          isConnecting: false,
          connectionError: true,
          errorMessage: event.reason || t('terminal.error.disconnected'),
        });
      }
      connectedTerminals.value[id] = false;
      sockets.value[id] = undefined;
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : t('terminal.error.failed');
    updateStatus(id, {
      isConnecting: false,
      connectionError: true,
      errorMessage,
    });
    sockets.value[id]?.close();
    sockets.value[id] = undefined;
  }
};

const reconnect = async (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn) return;

  updateStatus(id, {
    connectionError: false,
    isConnecting: true,
    errorMessage: '',
  });
  await initializeWebsocket(id, conn.connId);
  await nextTick();
  fitXterm(id);
};

const closeTerminal = (id: number) => {
  try {
    sockets.value[id]?.close();
    sockets.value[id] = undefined;
    connectedTerminals.value[id] = false;

    if (terminals.value[id]) {
      terminals.value[id].onData(() => {});
      terminals.value[id].onResize(() => {});
      try {
        terminals.value[id].dispose();
      } catch (e) {
        console.warn('Error disposing terminal:', e);
      }
      terminals.value[id] = undefined;
    }

    fitAddons.value[id] = undefined;
    webLinksAddons.value[id] = undefined;
    webglAddons.value[id] = undefined;
    canvasAddons.value[id] = undefined;

    updateStatus(id, {
      isConnecting: false,
      connectionError: false,
      errorMessage: '',
      errorDetails: '',
    });
  } catch (e) {
    console.error('Error in close terminal:', e);
  }
};

const registerToTabs = async () => {
  await nextTick();
  if (connectionTabs?.value && activeConn.value?.id) {
    const id = activeConn.value.id;
    const isConnected =
      !connectionStore.connections.find(c => c.id === id)?.connectionError &&
      !connectionStore.connections.find(c => c.id === id)?.isConnecting &&
      connectedTerminals.value[id];
    connectionTabs.value.registerTerminal(id, {
      closeTerminal: () => closeTerminal(id),
      status: isConnected ? enums.TerminalType.CONNECTED : enums.TerminalType.ERROR,
    });
  }
};

watchEffect(async () => {
  const connections = connectionStore.connections;
  for (const conn of connections) {
    if (conn.id === activeConn.value?.id && !terminals.value[conn.id]) {
      await initializeTerminal(conn.id);
      await nextTick();
      await initializeXterm(conn.id);
      if (!sockets.value[conn.id]) {
        await initializeWebsocket(conn.id, conn.connId);
      }
    }
  }
});

onMounted(async () => {
  window.addEventListener('resize', fitAllTerminals);
  await registerToTabs();
});

onUnmounted(() => {
  Object.keys(terminals.value).forEach(id => closeTerminal(Number(id)));
  window.removeEventListener('resize', fitAllTerminals);
});

onActivated(() => {
  const activeId = activeConn.value?.id;
  if (activeId && terminals.value[activeId]) {
    nextTick(() => {
      terminals.value[activeId]?.focus();
    });
  }
});

defineExpose({ closeTerminal });
defineOptions({ name: 'Terminal' });
</script>

<style lang="less" scoped>
.xterm-container {
  height: calc(100vh - 38px);
  position: relative;

  .terminal-hidden {
    position: absolute;
    opacity: 0;
    pointer-events: none;
    z-index: -1;
  }

  .xterm-wrapper {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }

  :deep(.xterm) {
    height: calc(100vh - 54px);
    padding: 8px 16px 8px 8px;

    .xterm-screen {
      width: 100% !important;
      height: 100% !important;
    }
    canvas {
      width: 100% !important;
      height: 100% !important;
    }

    .xterm-viewport {
      &::-webkit-scrollbar {
        width: 8px;
      }
      &::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.2);
        border-radius: 4px;

        &:hover {
          background: rgba(255, 255, 255, 0.3);
        }
      }
    }
  }

  .n-result {
    margin: auto;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .n-button {
    min-width: 120px;
    margin-top: 16px;
  }

  .fingerprint-details {
    text-align: left;
  }
}
</style>
