<template>
  <div class="connection-tabs flex items-center">
    <div
      v-for="tab in tabs"
      :key="tab.id"
      class="tab-item flex items-center px-3 h-8 cursor-pointer"
      :class="{ active: tab.id === activeTab }"
      @click="switchTab(tab.id)"
    >
      <NIcon size="16" class="mr-1">
        <Icon icon="ph:terminal-duotone" />
      </NIcon>
      <span class="text-sm">{{ tab.name }}</span>
      <NButton circle text size="tiny" class="ml-1 close-btn" @click.stop="closeTab(tab.id)">
        <template #icon>
          <Icon icon="ph:x-bold" />
        </template>
      </NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NIcon, NButton } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { useConnectionStore } from '@/stores/connection';

const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const activeTab = computed(() => connectionStore.activeConnectionId);

const tabs = computed(() => connectionStore.connections);

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});

const switchTab = (id: number) => {
  connectionStore.setActiveConnection(id);
};

const closeTab = (id: number) => {
  connectionStore.removeConnection(id);
};
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
</style>
