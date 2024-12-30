<template>
  <NModal
    v-model:show="dialogStore.aboutDialogVisible"
    preset="dialog"
    :close-on-esc="true"
    :title="$t('about.name')"
    transform-origin="center"
    :show-icon="false"
    style="width: 460px"
    :on-close="dialogStore.closeAboutDialog"
    :auto-focus="false"
  >
    <NSpace :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
      <NAvatar :size="80" :src="iconUrl" color="#0000" />
      <div class="about-app-title">GTerm</div>
      <NText>{{ gTermVer }}</NText>
      <NSpace :size="5" :wrap="false" :wrap-item="false" align="center">
        <NText class="about-link" @click="onOpenSource">{{ $t('about.github') }}</NText>
        <NDivider vertical />
        <NText class="about-link" @click="onOpenWebsite">{{ $t('about.website') }}</NText>
      </NSpace>
      <div :style="{ color: themeVars.textColor3 }" class="about-copyright">
        Copyright © 2024 OpenToolkitLab All rights reserved
      </div>
    </NSpace>
  </NModal>
</template>

<script setup lang="ts">
import iconUrl from '@/assets/images/icon.png';
import { useDialogStore } from '@/stores/dialog';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { GTermVer } from '@wailsApp/go/services/PreferencesSrv';
import { NAvatar, NDivider, NModal, NSpace, NText, useThemeVars } from 'naive-ui';

const gTermVer = ref('Unknown');

onMounted(async () => {
  gTermVer.value = await GTermVer();
});

const dialogStore = useDialogStore();
const themeVars = useThemeVars();

const onOpenSource = () => {
  BrowserOpenURL('https://github.com/OpenToolkitLab/GTerm');
};

const onOpenWebsite = () => {
  BrowserOpenURL('https://github.com/OpenToolkitLab');
};
</script>

<style lang="less" scoped>
.about-app-title {
  font-weight: bold;
  font-size: 18px;
  margin: 5px;
  padding-top: 20px;
}

.about-link {
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
}

.about-copyright {
  font-size: 12px;
}
</style>
