<template>
  <n-modal
    v-model:show="dialogStore.preferencesDialogVisible"
    :close-on-esc="true"
    :on-close="dialogStore.closePreferencesDialog"
    :show-icon="false"
    :negative-text="$t('preferencesModal.cancel')"
    :positive-text="$t('preferencesModal.confirm')"
    preset="dialog"
    style="width: 600px"
    :title="$t('preferencesModal.title')"
    transform-origin="center"
    @positive-click="handleSubmit"
  >
    <n-tabs key="settings" animated placement="left" type="line">
      <n-tab-pane name="appearance" :tab="$t('preferencesModal.tabs.general')">
        <n-form ref="formRef" label-placement="top" size="small">
          <n-form-item :label="$t('preferencesModal.theme.title')">
            <div style="display: flex; gap: 12px">
              <div
                :class="['theme-option', prefStore.themeMode === 'light' && 'theme-option--active']"
                @click="prefStore.updateThemeMode('light')"
              >
                <div class="theme-preview light" v-html="PreviewSvg" />
                <div class="theme-label">{{ $t('preferencesModal.theme.light') }}</div>
              </div>
              <div
                :class="['theme-option', prefStore.themeMode === 'dark' && 'theme-option--active']"
                @click="prefStore.updateThemeMode('dark')"
              >
                <div class="theme-preview dark" v-html="PreviewSvg" />
                <div class="theme-label">{{ $t('preferencesModal.theme.dark') }}</div>
              </div>
              <div
                :class="['theme-option', prefStore.themeMode === 'auto' && 'theme-option--active']"
                @click="prefStore.updateThemeMode('auto')"
              >
                <div class="theme-preview auto">
                  <div class="split-preview">
                    <div class="split-half light" v-html="PreviewSvg" />
                    <div class="split-half dark" v-html="PreviewSvg" />
                  </div>
                </div>
                <div class="theme-label">{{ $t('preferencesModal.theme.auto') }}</div>
              </div>
            </div>
          </n-form-item>

          <n-form-item :label="$t('preferencesModal.language.title')">
            <n-select
              v-model:value="prefStore.language"
              :options="languageOptions"
              size="medium"
              style="width: 100%"
              @update:value="prefStore.updateLanguage"
            />
          </n-form-item>

          <n-form-item :label="$t('preferencesModal.sidebar.title')">
            <div style="display: flex; gap: 12px; width: 100%">
              <n-input-number
                v-model:value="prefStore.sidebarWidth"
                :max="380"
                :min="260"
                :step="10"
                size="medium"
                style="flex: 1; width: 0"
                @update:value="value => value && prefStore.updateSidebarWidth(value)"
              >
                <template #suffix>{{ $t('preferencesModal.sidebar.width') }}</template>
              </n-input-number>
              <n-button size="medium" @click="prefStore.resetSidebarWidth">{{
                $t('preferencesModal.sidebar.reset')
              }}</n-button>
            </div>
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore, languageOptions } from '@/stores/preferences';
import { NButton, NModal, NTabPane, NTabs, NInputNumber, NSelect, NForm, NFormItem } from 'naive-ui';
import { onMounted, onUnmounted } from 'vue';
import type { FormInst } from 'naive-ui';
import PreviewSvg from '@/assets/images/preview.svg?raw';

const dialogStore = useDialogStore();
const prefStore = usePreferencesStore();
const formRef = ref<FormInst | null>(null);

const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

onMounted(() => {
  mediaQuery.addEventListener('change', prefStore.updateThemeBySystem);
  prefStore.updateThemeBySystem();
  prefStore.updateLanguageBySystem();
});

onUnmounted(() => {
  mediaQuery.removeEventListener('change', prefStore.updateThemeBySystem);
});

const handleSubmit = () => {
  dialogStore.closePreferencesDialog();
};
</script>

<style scoped>
.theme-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 106px;
  position: relative;
}

.theme-preview :deep(svg) {
  width: 120px;
  height: 90px;
  border-radius: 8px;
  overflow: hidden;
  transition: border-color 0.2s;
}

.theme-option {
  flex: 1;
  cursor: pointer;
}

.theme-option--active .theme-preview::after {
  content: '✓';
  position: absolute;
  top: 0;
  right: -8px;
  width: 20px;
  height: 20px;
  background: #18a058;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  border: 2px solid #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 1;
}

.theme-label {
  text-align: center;
  font-size: 13px;
  font-weight: 500;
  color: var(--n-text-color);
  margin-top: 4px;
}

.split-preview {
  display: flex;
  width: 120px;
  height: 90px;
  border-radius: 8px;
  overflow: hidden;
}

.split-half {
  width: 60px;
  height: 90px;
  position: relative;
  overflow: hidden;
}

.split-half:first-child {
  border-right: 1px solid var(--n-border-color);
}

.split-half.light :deep(svg) {
  position: absolute;
  left: 0;
  top: 0;
  width: 120px;
  height: 90px;
  --bg-color: #ffffff;
  --header-color: #f5f5f5;
  --sidebar-color: #f0f0f0;
  --icon-color: #d9d9d9;
  --text-color: #000000;
}

.split-half.dark :deep(svg) {
  position: absolute;
  right: 0;
  top: 0;
  width: 120px;
  height: 90px;
  --bg-color: #121212;
  --header-color: #1f1f1f;
  --sidebar-color: #262626;
  --icon-color: #404040;
  --text-color: #ffffff;
}

.theme-preview.light :deep(svg) {
  --bg-color: #ffffff;
  --header-color: #f5f5f5;
  --sidebar-color: #f0f0f0;
  --icon-color: #d9d9d9;
  --text-color: #000000;
}

.theme-preview.dark :deep(svg) {
  --bg-color: #121212;
  --header-color: #1f1f1f;
  --sidebar-color: #262626;
  --icon-color: #404040;
  --text-color: #ffffff;
}

.theme-option--active .theme-preview :deep(svg) {
  border-color: #18a058;
}

:deep(.n-form-item) {
  margin-bottom: 18px;
}
</style>
