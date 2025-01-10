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
import { NIcon, NButton, NTooltip } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { useConnectionStore } from '@/stores/connection';
import { useRouter } from 'vue-router';

const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const activeTab = computed(() => connectionStore.activeConnectionId);

const tabs = computed(() => connectionStore.connections);

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});

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
  background-color: v-bind('gtermThemeVars.cardColor');
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
  border-right: 1px solid rgba(127, 127, 127, 0.1);
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
  }

  &:hover {
    background-color: rgba(127, 127, 127, 0.08);
  }

  &.active {
    background-color: rgba(127, 127, 127, 0.15);
    border-right-color: transparent;
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
    background-color: #18a058;
    box-shadow: 0 0 4px rgba(24, 160, 88, 0.4);
  }

  &.error {
    background-color: #d03050;
    box-shadow: 0 0 4px rgba(208, 48, 80, 0.4);
  }

  &.connecting {
    background-color: #2080f0;
    box-shadow: 0 0 4px rgba(32, 128, 240, 0.4);
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
