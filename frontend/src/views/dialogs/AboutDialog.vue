<template>
  <n-modal
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
    <n-space :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
      <n-avatar :size="80" :src="iconUrl" color="#0000" />
      <div class="about-app-title">GTerm</div>
      <n-text>{{ gTermVer }}</n-text>
      <n-space :size="5" :wrap="false" :wrap-item="false" align="center">
        <n-text class="about-link" @click="onOpenSource">{{ $t('about.github') }}</n-text>
        <n-divider vertical />
        <n-text class="about-link" @click="onOpenWebsite">{{ $t('about.website') }}</n-text>
      </n-space>
      <div :style="{ color: themeVars.textColor3 }" class="about-copyright">
        Copyright © 2024 - {{ currentYear }} MisakaTAT All rights reserved
      </div>
    </n-space>
  </n-modal>
</template>

<script setup lang="ts">
import iconUrl from '@/assets/images/icon.png';
import { useDialogStore } from '@/stores/dialog';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { GTermVer } from '@wailsApp/go/services/PreferencesSrv';
import { NAvatar, NDivider, NModal, NSpace, NText, useThemeVars } from 'naive-ui';

const gTermVer = ref('Unknown');
const currentYear = computed(() => new Date().getFullYear());

onMounted(async () => {
  gTermVer.value = await GTermVer();
});

const dialogStore = useDialogStore();
const themeVars = useThemeVars();

const onOpenSource = () => {
  BrowserOpenURL('https://github.com/MisakaTAT/GTerm');
};

const onOpenWebsite = () => {
  BrowserOpenURL('https://github.com/MisakaTAT');
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
