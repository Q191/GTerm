<template>
  <n-config-provider :theme="currentTheme" :theme-overrides="currentThemeOverrides" :hljs="hljs">
    <n-message-provider placement="bottom-right">
      <n-dialog-provider>
        <router-view />
      </n-dialog-provider>
      <about-modal />
    </n-message-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import { usePreferencesStore } from '@/stores/preferences';
import { darkThemeOverrides, themeOverrides } from '@/themes/naive-theme';
import AboutModal from '@/views/modals/AboutModal.vue';
import { darkTheme, NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui';
import hljs from 'highlight.js/lib/core';
import bash from 'highlight.js/lib/languages/bash';

hljs.registerLanguage('bash', bash);

const prefStore = usePreferencesStore();

const currentTheme = computed(() => (prefStore.isDark ? darkTheme : null));

const currentThemeOverrides = computed(() => (prefStore.isDark ? darkThemeOverrides : themeOverrides));
</script>
