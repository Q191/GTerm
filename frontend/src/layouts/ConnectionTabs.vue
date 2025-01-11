<template>
  <div class="connection-tabs flex items-center h-10">
    <div
      v-for="tab in tabs"
      :key="tab.id"
      class="tab-item flex items-center px-3 h-10 cursor-pointer relative"
      :class="{ active: tab.id === activeTab }"
      @click="switchTab(tab.id)"
    >
      <div class="status-dot mr-1" :class="tabStatus[tab.id]" />
      <n-icon size="16" class="mr-1 flex-shrink-0">
        <icon icon="ph:terminal-duotone" />
      </n-icon>
      <n-tooltip trigger="hover">
        <template #trigger>
          <span class="text-sm tab-name">{{ tab.name }}</span>
        </template>
        {{ tab.name }}
      </n-tooltip>
      <n-button circle text size="tiny" class="ml-1 close-btn flex-shrink-0" @click.stop="closeTab(tab.id)">
        <template #icon>
          <icon icon="ph:x-bold" />
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NIcon, NButton, NTooltip, useThemeVars } from 'naive-ui';
import { useConnectionStore } from '@/stores/connection';
import { useRouter } from 'vue-router';

const connectionStore = useConnectionStore();
const activeTab = computed(() => connectionStore.activeConnectionId);
const themeVars = useThemeVars();

const tabs = computed(() => connectionStore.connections);

const terminalRefs = ref<Map<number, any>>(new Map());

const tabStatus = ref<Record<number, string>>({});

const updateTabStatus = (id: number, status: 'connected' | 'error' | 'connecting') => {
  tabStatus.value[id] = status;
};

const registerTerminal = (id: number, terminal: any) => {
  terminalRefs.value.set(id, terminal);
  if (terminal.status) {
    tabStatus.value[id] = terminal.status;
  }
};

const switchTab = (id: number) => {
  connectionStore.setActiveConnection(id);
  // 确保切换到终端页面
  if (router.currentRoute.value.name !== 'Terminal') {
    router.push({ name: 'Terminal' });
  }
};

const router = useRouter();

const closeTab = (id: number) => {
  const terminal = terminalRefs.value.get(id);
  if (terminal?.closeTerminal) {
    terminal.closeTerminal();
  }
  terminalRefs.value.delete(id);
  connectionStore.removeConnection(id);

  // 如果关闭后没有剩余标签，导航回Host页面
  if (connectionStore.connections.length === 0) {
    router.push({ name: 'Host' });
  }
};

defineExpose({
  registerTerminal,
  updateTabStatus,
});
</script>

<style lang="less" scoped>
.connection-tabs {
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: none;
  -ms-overflow-style: none;
  &::-webkit-scrollbar {
    display: none;
  }
}

.tab-item {
  position: relative;
  border-right: 1px solid v-bind('themeVars.borderColor');
  transition: background-color 0.2s ease;
  display: flex;
  align-items: center;
  width: 160px;
  min-width: 160px;

  .tab-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    color: v-bind('themeVars.textColorBase');
  }

  &:hover {
    background-color: v-bind('themeVars.hoverColor');
  }

  &.active {
    background-color: v-bind('themeVars.hoverColor');
  }

  .close-btn {
    opacity: 0;
    transition: opacity 0.2s;
  }

  &:hover .close-btn {
    opacity: 1;
  }
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;

  &.connected {
    background-color: v-bind('themeVars.successColor');
    box-shadow: 0 0 4px v-bind('`${themeVars.successColor}40`');
  }

  &.error {
    background-color: v-bind('themeVars.errorColor');
    box-shadow: 0 0 4px v-bind('`${themeVars.errorColor}40`');
  }

  &.connecting {
    background-color: v-bind('themeVars.infoColor');
    box-shadow: 0 0 4px v-bind('`${themeVars.infoColor}40`');
    animation: pulse 1.5s infinite;
  }
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.6;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>
