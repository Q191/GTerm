<template>
  <NConfigProvider :theme="currentTheme" :theme-overrides="currentThemeOverrides">
    <router-view />
    <AboutDialog />
    <HostDialog />
    <GroupDialog />
    <SettingDialog />
  </NConfigProvider>
</template>

<script lang="ts" setup>
import { usePreferencesStore } from '@/stores/preferences';
import { darkThemeOverrides, themeOverrides } from '@/themes/naive-theme';
import AboutDialog from '@/views/dialogs/AboutDialog.vue';
import HostDialog from '@/views/dialogs/HostDialog.vue';
import SettingDialog from '@/views/dialogs/SettingDialog.vue';
import { darkTheme, NConfigProvider } from 'naive-ui';
import GroupDialog from '@/views/dialogs/GroupDialog.vue';

const prefStore = usePreferencesStore();

// 使用计算属性来动态设置主题
const currentTheme = computed(() => (prefStore.isDark ? darkTheme : null));

// 使用计算属性来动态设置主题覆盖
const currentThemeOverrides = computed(() => (prefStore.isDark ? darkThemeOverrides : themeOverrides));
</script>
