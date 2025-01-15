<template>
  <n-modal
    v-model:show="dialogStore.preferencesDialogVisible"
    :close-on-esc="true"
    :on-close="dialogStore.closePreferencesDialog"
    :show-icon="false"
    negative-text="取消"
    positive-text="确定"
    preset="dialog"
    style="width: 650px"
    title="偏好设置"
    transform-origin="center"
    @positive-click="handleSubmit"
  >
    <n-tabs key="settings" animated placement="left" type="line" style="min-height: 400px">
      <n-tab-pane name="appearance" tab="常规设置" style="padding: 0 20px">
        <n-form ref="formRef" label-placement="top" size="small">
          <n-form-item label="主题">
            <n-button-group size="medium">
              <n-button
                :type="prefStore.themeMode === 'auto' ? 'primary' : 'default'"
                @click="prefStore.updateThemeMode('auto')"
              >
                <template #icon>
                  <icon icon="ph:desktop" />
                </template>
                跟随系统
              </n-button>
              <n-button
                :type="prefStore.themeMode === 'light' ? 'primary' : 'default'"
                @click="prefStore.updateThemeMode('light')"
              >
                <template #icon>
                  <icon icon="ph:sun" />
                </template>
                明亮模式
              </n-button>
              <n-button
                :type="prefStore.themeMode === 'dark' ? 'primary' : 'default'"
                @click="prefStore.updateThemeMode('dark')"
              >
                <template #icon>
                  <icon icon="ph:moon" />
                </template>
                暗黑模式
              </n-button>
            </n-button-group>
          </n-form-item>

          <n-form-item label="语言">
            <n-select
              v-model:value="prefStore.language"
              :options="languageOptions"
              size="medium"
              style="width: 100%"
              @update:value="prefStore.updateLanguage"
            />
          </n-form-item>

          <n-form-item label="分组栏宽度">
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
                <template #suffix>px</template>
              </n-input-number>
              <n-button size="medium" @click="prefStore.resetSidebarWidth">重置</n-button>
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
import { NButton, NButtonGroup, NModal, NTabPane, NTabs, NInputNumber, NSelect, NForm, NFormItem } from 'naive-ui';
import { onMounted, onUnmounted } from 'vue';
import { Icon } from '@iconify/vue';
import type { FormInst } from 'naive-ui';

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
