<template>
  <div class="connection-tabs">
    <div
      v-for="tab in tabs"
      :key="tab.id"
      class="tab-item"
      :class="{ active: tab.id === activeTab }"
      @click="switchTab(tab.id)"
    >
      <div class="status-dot" :class="tabStatus[tab.id]" />
      <NIcon size="16" class="terminal-icon">
        <Icon icon="ph:terminal" />
      </NIcon>
      <NTooltip trigger="hover">
        <template #trigger>
          <span class="tab-label">{{ tab.label }}</span>
        </template>
        {{ tab.label }}
      </NTooltip>
      <NButton circle text size="tiny" class="close-btn" @click.stop="closeTab(tab.id)">
        <template #icon>
          <Icon icon="ph:x-bold" />
        </template>
      </NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NIcon, NTooltip, useThemeVars } from 'naive-ui';
import { useRouter } from 'vue-router';
import { useConnectionStore } from '@/stores/connection';

const router = useRouter();

const connectionStore = useConnectionStore();
const activeTab = computed(() => connectionStore.activeConnectionId);
const themeVars = useThemeVars();

const tabs = computed(() => connectionStore.connections);

const terminalRefs = ref<Map<number, any>>(new Map());

const tabStatus = ref<Record<number, string>>({});

const updateTabStatus = (id: number, status: 'connected' | 'error' | 'connecting' | 'warning') => {
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
  if (router.currentRoute.value.name !== 'Terminal') {
    router.push({ name: 'Terminal' });
  }
};

const closeTab = (id: number) => {
  const terminal = terminalRefs.value.get(id);
  if (terminal?.closeTerminal) {
    terminal.closeTerminal();
  }
  terminalRefs.value.delete(id);
  connectionStore.removeConnection(id);

  if (connectionStore.connections.length === 0) {
    router.push({ name: 'Connection' });
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
  display: flex;
  align-items: center;
  height: 40px;

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
  width: auto;
  //min-width: 160px;
  height: 40px;
  padding: 0 12px;
  cursor: pointer;

  .tab-label {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    color: v-bind('themeVars.textColorBase');
    font-size: 14px;
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
    margin-left: 4px;
    flex-shrink: 0;
  }

  &:hover .close-btn {
    opacity: 1;
  }

  .terminal-icon {
    margin-right: 4px;
    flex-shrink: 0;
  }
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 4px;

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
