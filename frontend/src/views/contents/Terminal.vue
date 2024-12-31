<template>
  <div class="xterm-container">
    <div ref="xterm" class="xterm-wrapper" />
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

const xterm = ref<HTMLElement>();
const socket = ref<WebSocket>();

const fitAddon = ref<FitAddon>();
const webLinksAddon = ref<WebLinksAddon>();
const webglAddon = ref<WebglAddon>();
const canvasAddon = ref<CanvasAddon>();

const terminal = new Terminal({
  convertEol: true,
  disableStdin: false,
  fontSize: 16,
  cursorBlink: true,
  cursorStyle: 'bar',
  theme: baseTheme,
});

const fitXterm = throttle(() => {
  const dims = fitAddon.value?.proposeDimensions();
  if (!dims || !terminal || !dims.cols || !dims.rows) return;
  if (terminal.rows !== dims.rows || terminal.cols !== dims.cols) {
    terminal.resize(dims.cols, dims.rows);
  }
}, 50);

const loadAddon = async (addonFactory: () => any, addonName: string): Promise<any> => {
  try {
    const addon = addonFactory();
    terminal.loadAddon(addon);
    console.log(`${addonName} addon loaded successfully`);
    return addon;
  } catch (e) {
    console.warn(`Failed to load ${addonName} addon`, e);
    return null;
  }
};

const initializeXterm = async () => {
  if (!xterm.value) {
    return;
  }

  fitAddon.value = await loadAddon(() => new FitAddon(), 'Fit');
  webLinksAddon.value = await loadAddon(() => new WebLinksAddon(), 'WebLinks');
  webglAddon.value = await loadAddon(() => new WebglAddon(), 'WebGL');
  if (!webglAddon.value) {
    canvasAddon.value = await loadAddon(() => new CanvasAddon(), 'Canvas');
  }

  terminal.attachCustomKeyEventHandler(arg => {
    if (arg.code === 'PageUp' && arg.type === 'keydown') {
      terminal.scrollPages(-1);
      return false;
    } else if (arg.code === 'PageDown' && arg.type === 'keydown') {
      terminal.scrollPages(1);
      return false;
    }
    return true;
  });

  terminal.open(xterm.value);
  terminal.write('\x1B[32m[INFO] Connecting···\x1B[0m\r\n\r\n');
  terminal.onData(data => socket.value?.send(JSON.stringify({ type: 'cmd', cmd: data })));
  terminal.onResize(({ cols, rows }) => {
    if (socket.value?.readyState === WebSocket.OPEN) {
      socket.value?.send(JSON.stringify({ type: 'resize', cols, rows }));
    }
  });
};

const initializeWebsocket = async () => {
  const port = await WebSocketPort();
  socket.value = new WebSocket(`ws://localhost:${port}/ws/terminal`);
  socket.value.binaryType = 'arraybuffer';
  if (!socket.value) return;
  socket.value.onopen = () => {
    socket.value?.send('\n');
    terminal.focus();
    fitXterm();
  };
  socket.value.onmessage = (event: MessageEvent) => {
    if (event.data instanceof ArrayBuffer) {
      const content = new Uint8Array(event.data);
      terminal.write(content);
    }
  };
  socket.value.onerror = () => {
    terminal.write('\x1B[31m[ERR] App internal error.\x1B[0m');
  };
  socket.value.onclose = () => {
    terminal.write('\r\n\r\n\x1B[33m[WARN] Connection is closed.\x1B[0m');
  };
};

onMounted(async () => {
  await initializeXterm();
  await initializeWebsocket();
  window.addEventListener('resize', fitXterm);
});

onUnmounted(() => {
  terminal?.dispose();
  webLinksAddon.value?.dispose();
  fitAddon.value?.dispose();
  webglAddon.value?.dispose();
  canvasAddon.value?.dispose();
  socket.value?.close();
  window.removeEventListener('resize', fitXterm);
});
</script>

<style lang="less" scoped>
.xterm-container {
  height: calc(100vh - 38px);
  width: 100%;
  display: flex;

  .xterm-wrapper {
    flex-grow: 1;
    display: flex;
    width: 100%;
  }

  :deep(.xterm) {
    .xterm-viewport {
      overflow: hidden;
    }
    flex-grow: 1;
    display: flex;
    padding: 8px 0px 8px 8px;
  }
}
</style>
