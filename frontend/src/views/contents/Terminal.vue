<template>
  <div class="xterm-container">
    <n-result v-if="activeConnection?.connectionError" status="error" :title="activeConnection.errorMessage">
      <template #footer>
        <n-button @click="reconnect" type="primary">重新连接</n-button>
      </template>
      <template #default>
        <div class="error-details">
          <n-collapse v-if="activeConnection.errorDetails" class="error-collapse">
            <n-collapse-item title="详细信息" name="details">
              <n-code :code="activeConnection.errorDetails" language="bash" :word-wrap="true" />
            </n-collapse-item>
          </n-collapse>
        </div>
      </template>
    </n-result>
    <n-result v-else-if="activeConnection?.isConnecting" status="info" title="正在连接" description="少女祈祷中...">
      <template #icon>
        <n-spin size="large" />
      </template>
    </n-result>
    <div v-else ref="xterm" class="xterm-wrapper" />
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
import { NCode, NCollapse, NCollapseItem } from 'naive-ui';

const xterm = ref<HTMLElement>();
const socket = ref<WebSocket>();
const fitAddon = ref<FitAddon>();
const webLinksAddon = ref<WebLinksAddon>();
const webglAddon = ref<WebglAddon>();
const canvasAddon = ref<CanvasAddon>();
const terminal = ref<Terminal>();

const connectionStore = useConnectionStore();
const connectionTabs = inject<any>('connectionTabs');
const activeConnection = computed(() => connectionStore.activeConnection);

const updateStatus = (
  status: Partial<{
    isConnecting: boolean;
    connectionError: boolean;
    errorMessage: string;
    errorDetails: string;
  }>,
) => {
  if (activeConnection.value?.id) {
    connectionStore.updateConnectionStatus(activeConnection.value.id, status);
  }
};

const initializeTerminal = () => {
  terminal.value = new Terminal({
    convertEol: true,
    disableStdin: false,
    fontSize: 16,
    cursorBlink: true,
    cursorStyle: 'bar',
    theme: baseTheme,
  });
};

const fitXterm = throttle(() => {
  const dims = fitAddon.value?.proposeDimensions();
  if (!dims || !terminal.value || !dims.cols || !dims.rows) return;
  if (terminal.value.rows !== dims.rows || terminal.value.cols !== dims.cols) {
    terminal.value.resize(dims.cols, dims.rows);
  }
}, 50);

const loadAddon = async (addonFactory: () => any, addonName: string): Promise<any> => {
  try {
    const addon = addonFactory();
    terminal.value?.loadAddon(addon);
    return addon;
  } catch (e) {
    console.warn(`Failed to load ${addonName} addon:`, e);
    return null;
  }
};

const initializeXterm = async () => {
  if (!xterm.value || !terminal.value) return;

  fitAddon.value = await loadAddon(() => new FitAddon(), 'Fit');
  webLinksAddon.value = await loadAddon(() => new WebLinksAddon(), 'WebLinks');
  webglAddon.value = await loadAddon(() => new WebglAddon(), 'WebGL');
  if (!webglAddon.value) {
    canvasAddon.value = await loadAddon(() => new CanvasAddon(), 'Canvas');
  }

  terminal.value.attachCustomKeyEventHandler(arg => {
    if (arg.code === 'PageUp' && arg.type === 'keydown') {
      terminal.value?.scrollPages(-1);
      return false;
    } else if (arg.code === 'PageDown' && arg.type === 'keydown') {
      terminal.value?.scrollPages(1);
      return false;
    }
    return true;
  });

  terminal.value.open(xterm.value);
  terminal.value.onData(data => socket.value?.send(JSON.stringify({ type: 'cmd', cmd: data })));
  terminal.value.onResize(({ cols, rows }) => {
    if (socket.value?.readyState === WebSocket.OPEN) {
      socket.value?.send(JSON.stringify({ type: 'resize', cols, rows }));
    }
  });
};

const initializeWebsocket = async () => {
  try {
    const port = await WebSocketPort();
    socket.value = new WebSocket(`ws://localhost:${port}/ws/terminal?hostId=${activeConnection.value?.hostId}`);
    if (!socket.value) return;

    socket.value.onopen = () => {
      terminal.value?.focus();
      fitXterm();
    };

    socket.value.onmessage = (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      switch (data.type) {
        case 'error':
          updateStatus({
            isConnecting: false,
            connectionError: true,
            errorMessage: data.error,
            errorDetails: data.details,
          });
          socket.value?.close();
          break;
        case 'connected':
          updateStatus({ isConnecting: false });
          terminal.value?.focus();
          fitXterm();
          break;
        case 'data':
          terminal.value?.write(data.content);
          break;
        default:
          console.warn('Unknown message type:', data.type);
      }
    };

    socket.value.onerror = () => {
      updateStatus({
        isConnecting: false,
        connectionError: true,
        errorMessage: '连接发生错误，请检查应用是否正常运行',
      });
      socket.value?.close();
    };

    socket.value.onclose = event => {
      if (!activeConnection.value?.connectionError) {
        updateStatus({
          isConnecting: false,
          connectionError: true,
          errorMessage: event.reason || '连接已断开',
        });
      }
      socket.value = undefined;
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : '无法建立连接';
    updateStatus({
      isConnecting: false,
      connectionError: true,
      errorMessage,
    });
    socket.value?.close();
    socket.value = undefined;
  }
};

const reconnect = async () => {
  updateStatus({
    connectionError: false,
    isConnecting: true,
    errorMessage: '',
  });
  await initializeWebsocket();
};

const closeTerminal = () => {
  socket.value?.close();
  socket.value = undefined;
  if (terminal.value) {
    fitAddon.value?.dispose();
    webLinksAddon.value?.dispose();
    webglAddon.value?.dispose();
    canvasAddon.value?.dispose();
    terminal.value.dispose();
    terminal.value = undefined;
  }
  if (activeConnection.value?.id) {
    updateStatus({
      isConnecting: false,
      connectionError: false,
      errorMessage: '',
      errorDetails: '',
    });
  }
};

const registerToTabs = async () => {
  await nextTick();
  if (connectionTabs?.value && connectionStore.activeConnectionId) {
    connectionTabs.value.registerTerminal(connectionStore.activeConnectionId, {
      closeTerminal,
    });
  }
};

onMounted(async () => {
  initializeTerminal();
  await initializeXterm();
  await initializeWebsocket();
  window.addEventListener('resize', fitXterm);
  await registerToTabs();
});

onUnmounted(() => {
  closeTerminal();
  window.removeEventListener('resize', fitXterm);
});

defineExpose({
  closeTerminal,
});
</script>

<style lang="less" scoped>
.xterm-container {
  height: calc(100vh - 38px);
  width: 100%;
  display: flex;

  :deep(.n-spin-container) {
    width: 100%;
    height: 100%;
    display: flex;
  }

  :deep(.n-spin-content) {
    width: 100%;
    height: 100%;
    display: flex;
  }

  .xterm-wrapper {
    flex-grow: 1;
    display: flex;
    width: 100%;
    height: 100%;
  }

  :deep(.xterm) {
    .xterm-viewport {
      overflow: hidden;
    }
    flex-grow: 1;
    display: flex;
    padding: 8px 0px 8px 8px;
  }

  :deep(.n-result) {
    margin: auto;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    .n-result-header {
      margin-bottom: 12px;
    }

    .n-result-icon {
      font-size: 48px;
    }

    .n-result-title {
      font-size: 16px;
      margin-top: 12px;
    }

    .error-details {
      width: 100%;
      max-width: 600px;
      margin-top: 12px;
      font-size: 13px;
    }

    .error-collapse {
      width: 100%;
      background-color: transparent;

      :deep(.n-collapse-item .n-collapse-item__header) {
        font-size: 13px;
        padding: 8px 0;
      }

      :deep(.n-code) {
        font-size: 12px;
      }
    }

    .n-button {
      min-width: 120px;
      margin-top: 16px;
    }
  }
}
</style>
