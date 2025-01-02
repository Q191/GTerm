<template>
  <n-modal
    v-model:show="dialogStore.settingsDialogVisible"
    :close-on-esc="true"
    :on-close="dialogStore.closeSettingsDialog"
    :show-icon="false"
    negative-text="取消"
    positive-text="确定"
    preset="dialog"
    style="width: 550px"
    title="全局设置"
    transform-origin="center"
  >
    <n-tabs key="Oasis" animated placement="left" type="line">
      <n-tab-pane name="oasis" tab="Oasis">
        <n-space>
          <n-button type="primary" @click="prefStore.toLight"> 明亮模式 </n-button>
          <n-button type="primary" @click="prefStore.toDark"> 暗黑模式 </n-button>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="the beatles" tab="the Beatles"> Hey Jude </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';
import { NButton, NModal, NSpace, NTabPane, NTabs } from 'naive-ui';

const dialogStore = useDialogStore();
const prefStore = usePreferencesStore();
</script>

<style lang="less" scoped></style>
