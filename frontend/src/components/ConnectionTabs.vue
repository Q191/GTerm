<template>
  <div class="connection-tabs flex items-center">
    <div
      v-for="tab in tabs"
      :key="tab.id"
      class="tab-item flex items-center px-3 h-8 cursor-pointer"
      :class="{ active: tab.id === activeTab }"
      @click="switchTab(tab.id)"
    >
      <div class="status-dot mr-1" :class="tabStatus[tab.id]" />
      <n-icon size="16" class="mr-1">
        <icon icon="ph:terminal-duotone" />
      </n-icon>
      <span class="text-sm">{{ tab.name }}</span>
      <n-button circle text size="tiny" class="ml-1 close-btn" @click.stop="closeTab(tab.id)">
        <template #icon>
          <icon icon="ph:x-bold" />
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NIcon, NButton } from 'naive-ui';
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
.tab-item {
  border-right: 1px solid v-bind('gtermThemeVars.splitColor');

  &:hover {
    background-color: v-bind('gtermThemeVars.splitColor');
  }

  &.active {
    background-color: v-bind('gtermThemeVars.splitColor');
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
  width: 8px;
  height: 8px;
  border-radius: 50%;

  &.connected {
    background-color: #18a058;
  }

  &.error {
    background-color: #d03050;
  }

  &.connecting {
    background-color: #2080f0;
    animation: pulse 1.5s infinite;
  }
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
  100% {
    opacity: 1;
  }
}
</style>
