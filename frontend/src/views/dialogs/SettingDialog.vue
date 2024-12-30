<template>
  <NModal
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
    <NTabs key="Oasis" animated placement="left" type="line">
      <NTabPane name="oasis" tab="Oasis">
        <NSpace>
          <NButton type="primary" @click="prefStore.toLight"> 明亮模式 </NButton>
          <NButton type="primary" @click="prefStore.toDark"> 暗黑模式 </NButton>
        </NSpace>
      </NTabPane>
      <NTabPane name="the beatles" tab="the Beatles"> Hey Jude </NTabPane>
    </NTabs>
  </NModal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';
import { NButton, NModal, NSpace, NTabPane, NTabs } from 'naive-ui';

const dialogStore = useDialogStore();
const prefStore = usePreferencesStore();
</script>

<style lang="less" scoped></style>
