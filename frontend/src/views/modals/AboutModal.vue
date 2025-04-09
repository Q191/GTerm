<template>
  <n-modal
    v-model:show="dialogStore.aboutDialogVisible"
    preset="dialog"
    :close-on-esc="true"
    :title="$t('frontend.about.name')"
    transform-origin="center"
    :show-icon="false"
    style="width: 460px"
    :on-close="dialogStore.closeAboutDialog"
    :auto-focus="false"
  >
    <n-space :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
      <n-avatar :size="80" :src="iconUrl" color="#0000" />
      <div class="about-app-title">GTerm</div>
      <n-button text @click="onOpenVersionURL">{{ version }}</n-button>
      <n-space :size="1" :wrap="false" :wrap-item="false" align="center">
        <n-button text @click="onOpenSource">{{ $t('frontend.about.github') }}</n-button>
        <n-divider vertical />
        <n-button text @click="onOpenWebsite">{{ $t('frontend.about.website') }}</n-button>
      </n-space>
      <div :style="{ color: themeVars.textColor3 }" class="about-copyright">
        <span>{{ copyright }}</span>
      </div>
    </n-space>
  </n-modal>
</template>

<script setup lang="ts">
import iconUrl from '@/assets/images/icon.png';
import { useDialogStore } from '@/stores/dialog';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { Copyright, Version, VersionURL } from '@wailsApp/go/services/PreferencesSrv';
import { NAvatar, NDivider, NModal, NSpace, useThemeVars } from 'naive-ui';

const version = ref('');
const copyright = ref('');
const versionUrl = ref('');

onMounted(async () => {
  version.value = await Version();
  copyright.value = await Copyright();
  versionUrl.value = await VersionURL();
});

const dialogStore = useDialogStore();
const themeVars = useThemeVars();

const onOpenVersionURL = () => {
  BrowserOpenURL(versionUrl.value);
};

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
}

.about-copyright {
  font-size: 12px;
  text-align: center;
  white-space: pre-wrap;
}
</style>
