<template>
  <NMessageProvider placement="bottom-right">
    <NConfigProvider :theme="currentTheme" :theme-overrides="currentThemeOverrides" :hljs="hljs">
      <router-view />
      <AboutDialog />
      <HostDialog />
      <GroupDialog />
      <SettingDialog />
    </NConfigProvider>
  </NMessageProvider>
</template>

<script lang="ts" setup>
import { usePreferencesStore } from '@/stores/preferences';
import { darkThemeOverrides, themeOverrides } from '@/themes/naive-theme';
import AboutDialog from '@/views/dialogs/AboutDialog.vue';
import HostDialog from '@/views/dialogs/HostDialog.vue';
import SettingDialog from '@/views/dialogs/SettingDialog.vue';
import { darkTheme, NConfigProvider, NMessageProvider } from 'naive-ui';
import GroupDialog from '@/views/dialogs/GroupDialog.vue';
import hljs from 'highlight.js/lib/core';
import bash from 'highlight.js/lib/languages/bash';

hljs.registerLanguage('bash', bash);

const prefStore = usePreferencesStore();

const currentTheme = computed(() => (prefStore.isDark ? darkTheme : null));

const currentThemeOverrides = computed(() => (prefStore.isDark ? darkThemeOverrides : themeOverrides));
</script>
