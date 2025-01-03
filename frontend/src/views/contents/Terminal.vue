<template>
  <div class="xterm-container">
    <template v-for="conn in connectionStore.connections" :key="conn.id">
      <n-result
        v-if="conn.connectionError"
        status="error"
        :title="conn.errorMessage"
        v-show="conn.id === activeConn?.id"
      >
        <template #footer>
          <n-button @click="() => reconnect(conn.id)" type="primary">重新连接</n-button>
        </template>
        <template #default>
          <div class="error-details" v-if="conn.errorDetails">
            <n-collapse class="error-collapse">
              <n-collapse-item title="详细信息" name="details">
                <n-code :code="conn.errorDetails" language="bash" :word-wrap="true" />
              </n-collapse-item>
            </n-collapse>
          </div>
        </template>
      </n-result>
      <n-result
        v-else-if="conn.isConnecting || !connectedTerminals[conn.id]"
        status="info"
        title="正在连接"
        description="少女祈祷中..."
        v-show="conn.id === activeConn?.id"
      >
        <template #icon>
          <n-spin size="large" />
        </template>
      </n-result>
      <div
        v-show="
          conn.id === activeConn?.id && !conn.connectionError && !conn.isConnecting && connectedTerminals[conn.id]
        "
        :ref="el => el && (terminalRefs[conn.id] = el as HTMLElement)"
        class="xterm-wrapper"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { baseTheme } from '@/themes/xtermjs-theme';
import { WebSocketPort } from '@wailsApp/go/services/TerminalSrv';
import { CanvasAddon } from '@xterm/addon-canvas';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { WebglAddon } from '@xterm/addon-webgl';
import { Terminal } from '@xterm/xterm';
import { throttle } from 'lodash';
import '@xterm/xterm/css/xterm.css';
import { useConnectionStore } from '@/stores/connection';
import { NCode, NCollapse, NCollapseItem, NResult, NSpin, NButton } from 'naive-ui';
import { onActivated } from 'vue';

const connectionStore = useConnectionStore();
const connectionTabs = inject<any>('connectionTabs');
const activeConn = computed(() => connectionStore.activeConnection);

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
    errorMessage: string;
    errorDetails: string;
  }>,
) => {
  connectionStore.updateConnectionStatus(id, status);
  if (connectionTabs?.value) {
    let tabStatus = 'connecting';
    if (status.connectionError) {
      tabStatus = 'error';
    } else if (!status.isConnecting && connectedTerminals.value[id]) {
      tabStatus = 'connected';
    }
    connectionTabs.value.updateTabStatus(id, tabStatus);
  }
};

const fitXterm = throttle((id: number) => {
  const dims = fitAddons.value[id]?.proposeDimensions();
  if (!dims || !terminals.value[id] || !dims.cols || !dims.rows) return;
  if (terminals.value[id].rows !== dims.rows || terminals.value[id].cols !== dims.cols) {
    terminals.value[id].resize(dims.cols, dims.rows);
  }
}, 50);

const initializeTerminal = (id: number) => {
  if (terminals.value[id]) return;

  terminals.value[id] = new Terminal({
    convertEol: true,
    disableStdin: false,
    fontSize: 16,
    cursorBlink: true,
    cursorStyle: 'bar',
    theme: baseTheme,
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
  terminal.onData(data => sockets.value[id]?.send(JSON.stringify({ type: 'cmd', cmd: data })));
  terminal.onResize(({ cols, rows }) => {
    if (sockets.value[id]?.readyState === WebSocket.OPEN) {
      sockets.value[id]?.send(JSON.stringify({ type: 'resize', cols, rows }));
    }
  });

  nextTick(() => fitXterm(id));
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
      errorMessage: '',
      errorDetails: '',
    });

    const port = await WebSocketPort();
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
        case 'error':
          updateStatus(id, {
            isConnecting: false,
            connectionError: true,
            errorMessage: data.error,
            errorDetails: data.details,
          });
          connectedTerminals.value[id] = false;
          socket?.close();
          break;
        case 'connected':
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
        case 'data':
          terminals.value[id]?.write(data.content);
          break;
      }
    };

    socket.onerror = () => {
      updateStatus(id, {
        isConnecting: false,
        connectionError: true,
        errorMessage: '连接发生错误，请检查应用是否正常运行',
      });
      connectedTerminals.value[id] = false;
      socket?.close();
    };

    socket.onclose = event => {
      if (!connectionStore.connections.find(c => c.id === id)?.connectionError) {
        updateStatus(id, {
          isConnecting: false,
          connectionError: true,
          errorMessage: event.reason || '连接已断开',
        });
      }
      connectedTerminals.value[id] = false;
      sockets.value[id] = undefined;
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : '无法建立连接';
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
  await initializeWebsocket(id, conn.hostId);
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
    console.error('Error in closeTerminal:', e);
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
      status: isConnected ? 'connected' : 'error',
    });
  }
};

watchEffect(async () => {
  const connections = connectionStore.connections;
  for (const conn of connections) {
    if (!terminals.value[conn.id]) {
      initializeTerminal(conn.id);
      await nextTick();
      await initializeXterm(conn.id);
      if (!sockets.value[conn.id]) {
        await initializeWebsocket(conn.id, conn.hostId);
      }
    }
  }
});

onMounted(async () => {
  window.addEventListener('resize', () => activeConn.value && fitXterm(activeConn.value.id));
  await registerToTabs();
});

onUnmounted(() => {
  Object.keys(terminals.value).forEach(id => closeTerminal(Number(id)));
  window.removeEventListener('resize', () => activeConn.value && fitXterm(activeConn.value.id));
});

onActivated(() => {
  Object.keys(terminals.value).forEach(id => {
    const numId = Number(id);
    if (terminals.value[numId]) {
      nextTick(() => {
        fitXterm(numId);
        if (activeConn.value?.id === numId) {
          terminals.value[numId]?.focus();
        }
      });
    }
  });
});

defineExpose({ closeTerminal });
defineOptions({ name: 'Terminal' });
</script>

<style lang="less" scoped>
.xterm-container {
  height: calc(100vh - 38px);

  :deep(.xterm) {
    height: calc(100vh - 38px);
    padding: 8px 0px 8px 8px;

    .xterm-viewport {
      overflow: hidden;
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
}
</style>
