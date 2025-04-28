<template>
  <NConfigProvider :theme="currentTheme" :theme-overrides="currentThemeOverrides" :hljs="hljs">
    <NMessageProvider placement="bottom-right">
      <NDialogProvider>
        <router-view />
      </NDialogProvider>
      <AboutModal />
    </NMessageProvider>
  </NConfigProvider>
</template>

<script lang="ts" setup>
import hljs from 'highlight.js/lib/core';
import bash from 'highlight.js/lib/languages/bash';
import { darkTheme, NConfigProvider, NDialogProvider, NMessageProvider } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { darkThemeOverrides, themeOverrides } from '@/themes/naive-theme';
import AboutModal from '@/views/modals/AboutModal.vue';

hljs.registerLanguage('bash', bash);

const prefStore = usePreferencesStore();

const currentTheme = computed(() => (prefStore.isDark ? darkTheme : null));

const currentThemeOverrides = computed(() => (prefStore.isDark ? darkThemeOverrides : themeOverrides));
</script>
