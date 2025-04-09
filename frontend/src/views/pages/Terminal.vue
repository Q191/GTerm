<template>
  <div class="xterm-container">
    <template v-for="conn in connectionStore.connections" :key="conn.id">
      <n-result
        v-if="conn.errorCausedClosed"
        :status="conn.errorCausedClosed ? 'error' : 'warning'"
        :title="conn.message"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #footer>
          <n-button @click="() => reconnect(conn.id)" type="primary">{{ $t('frontend.terminal.reconnect') }}</n-button>
        </template>
        <template #default>
          <div v-if="conn.details">
            <n-collapse>
              <n-collapse-item :title="$t('frontend.terminal.error.details')" name="details">
                <n-code :code="conn.details" language="bash" :word-wrap="true" />
              </n-collapse-item>
            </n-collapse>
          </div>
        </template>
      </n-result>
      <n-result
        v-else-if="conn.isFingerprintConfirm"
        status="warning"
        :title="$t('frontend.terminal.fingerprint.confirm')"
        :description="$t('frontend.terminal.fingerprint.description')"
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
              {{ $t('frontend.terminal.fingerprint.reject') }}
            </n-button>
            <n-button @click="() => acceptFingerprint(conn.id)" type="primary">
              {{ $t('frontend.terminal.fingerprint.accept') }}
            </n-button>
          </n-space>
        </template>
        <template #default>
          <div class="fingerprint-details">
            <p>
              <strong>{{ $t('frontend.terminal.fingerprint.host') }}:</strong> {{ conn.hostAddress }}
            </p>
            <p>
              <strong>{{ $t('frontend.terminal.fingerprint.fingerprint') }}:</strong>
            </p>
            <n-code :code="conn.hostFingerprint || ''" language="bash" :word-wrap="true" />
          </div>
        </template>
      </n-result>
      <n-result
        v-else-if="conn.isConnecting || !connectedTerminals[conn.id]"
        status="info"
        :title="$t('frontend.terminal.connecting')"
        :description="$t('frontend.terminal.connecting_desc')"
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
import { LogWarning, LogError, LogInfo } from '@wailsApp/runtime/runtime';
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
import { getTranslated } from '@/utils/call';

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
const resizeObservers = ref<Record<number, ResizeObserver | undefined>>({});

const updateStatus = (
  id: number,
  status: Partial<{
    isConnecting: boolean;
    errorCausedClosed: boolean;
    isFingerprintConfirm: boolean;
    message: string;
    details: string;
    hostAddress: string;
    hostFingerprint: string;
  }>,
) => {
  connectionStore.updateConnectionStatus(id, status);
  if (connectionTabs?.value) {
    let tabStatus = 'connecting';
    if (status.errorCausedClosed) {
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
    isConnecting: true,
  });
};

const fitXterm = (id: number) => {
  if (!fitAddons.value[id] || !terminals.value[id]) return;
  fitAddons.value[id].fit();
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
    LogWarning(`Failed to activate WebGL addon, falling back to Canvas: ${e}`);
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

  const resizeObserver = new ResizeObserver(() => {
    if (fitAddons.value[id]) {
      fitAddons.value[id].fit();
    }
  });
  resizeObserver.observe(terminalEl);
  resizeObservers.value[id] = resizeObserver;

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
      errorCausedClosed: false,
      isFingerprintConfirm: false,
      message: '',
      details: '',
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
            errorCausedClosed: true,
            message: getTranslated(data.code, data.message),
            details: data.details,
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
        errorCausedClosed: true,
        message: t('frontend.terminal.error.connection'),
      });
      connectedTerminals.value[id] = false;
      socket?.close();
    };

    socket.onclose = event => {
      if (!connectionStore.connections.find(c => c.id === id)?.errorCausedClosed) {
        updateStatus(id, {
          isConnecting: false,
          errorCausedClosed: true,
          // TODO: reason 需要返回 message 原文。用于在 code 不存在的时候 fallback
          message: getTranslated(event.reason, t('frontend.terminal.error.disconnected')),
        });
      }
      connectedTerminals.value[id] = false;
      sockets.value[id] = undefined;
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : t('frontend.terminal.error.failed');
    updateStatus(id, {
      isConnecting: false,
      errorCausedClosed: true,
      message: errorMessage,
    });
    sockets.value[id]?.close();
    sockets.value[id] = undefined;
  }
};

const reconnect = async (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn) return;

  LogInfo(`Reconnecting to terminal ID: ${id}, connId: ${conn.connId}`);
  updateStatus(id, {
    errorCausedClosed: false,
    isConnecting: true,
    message: '',
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

    if (resizeObservers.value[id]) {
      resizeObservers.value[id].disconnect();
      resizeObservers.value[id] = undefined;
    }

    if (terminals.value[id]) {
      terminals.value[id].onData(() => {});
      terminals.value[id].onResize(() => {});
      try {
        terminals.value[id].dispose();
      } catch (e) {
        LogWarning(`Error disposing terminal: ${e}`);
      }
      terminals.value[id] = undefined;
    }

    fitAddons.value[id] = undefined;
    webLinksAddons.value[id] = undefined;
    webglAddons.value[id] = undefined;
    canvasAddons.value[id] = undefined;

    updateStatus(id, {
      isConnecting: false,
      errorCausedClosed: false,
      message: '',
      details: '',
    });
  } catch (e) {
    LogError(`Error in close terminal: ${e}`);
  }
};

const registerToTabs = async () => {
  await nextTick();
  if (connectionTabs?.value && activeConn.value?.id) {
    const id = activeConn.value.id;
    const isConnected =
      !connectionStore.connections.find(c => c.id === id)?.errorCausedClosed &&
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
  await registerToTabs();
});

onUnmounted(() => {
  Object.keys(terminals.value).forEach(id => closeTerminal(Number(id)));
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
