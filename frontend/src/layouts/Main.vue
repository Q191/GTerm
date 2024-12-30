<template>
  <NLayout class="content" has-sider>
    <NLayoutSider :width="60">
      <Sider class="settings" />
    </NLayoutSider>
    <NLayoutContent>
      <NScrollbar>
        <router-view />
      </NScrollbar>
    </NLayoutContent>
  </NLayout>
</template>

<script lang="ts" setup>
import Sider from '@/layouts/Sider.vue';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { NLayout, NLayoutContent, NLayoutSider, NScrollbar } from 'naive-ui';

const prefStore = usePreferencesStore();

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});
</script>

<style lang="less" scoped>
.content {
  height: calc(100vh - 38px);
  border-top: 1px solid v-bind('gtermThemeVars.splitColor');
}

.settings {
  height: 100%;
  border-right: 1px solid v-bind('gtermThemeVars.splitColor');
}
</style>
