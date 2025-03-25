<template>
  <n-message-provider placement="bottom-right">
    <n-config-provider :theme="currentTheme" :theme-overrides="currentThemeOverrides" :hljs="hljs">
      <router-view />
      <about-dialog />
      <preferences-dialog />
    </n-config-provider>
  </n-message-provider>
</template>

<script lang="ts" setup>
import { usePreferencesStore } from '@/stores/preferences';
import { darkThemeOverrides, themeOverrides } from '@/themes/naive-theme';
import AboutDialog from '@/views/dialogs/AboutDialog.vue';
import PreferencesDialog from '@/views/dialogs/PreferencesDialog.vue';
import { darkTheme, NConfigProvider, NMessageProvider } from 'naive-ui';
import hljs from 'highlight.js/lib/core';
import bash from 'highlight.js/lib/languages/bash';

hljs.registerLanguage('bash', bash);

const prefStore = usePreferencesStore();

const currentTheme = computed(() => (prefStore.isDark ? darkTheme : null));

const currentThemeOverrides = computed(() => (prefStore.isDark ? darkThemeOverrides : themeOverrides));
</script>
