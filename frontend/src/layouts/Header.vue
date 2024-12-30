<template>
  <div ref="widthRef" class="flex items-center justify-between h-12 header rounded-t-lg">
    <div class="pl-2 flex-1 flex justify-center items-center">
      <div class="flex items-center">
        <img src="@/assets/images/icon.png" alt="Logo" class="w-6 h-6" />
        <span class="pl-2 text-base font-semibold">GTerm</span>
      </div>
    </div>

    <div v-if="!IsDarwin()" class="flex items-center mt-0">
      <div class="w-12 h-12 flex items-center justify-center hover:bg-custom-hover" @click="WindowMinimise">
        <NIcon size="20"><Remove /></NIcon>
      </div>
      <div class="w-12 h-12 flex items-center justify-center hover:bg-custom-hover" @click="reduction">
        <NIcon size="20">
          <span v-if="windowIsMaximised"><Contract /></span>
          <span v-else><Expand /></span>
        </NIcon>
      </div>
      <div class="w-12 h-12 flex items-center justify-center hover:bg-custom-hover hover:rounded-tr-lg" @click="Quit">
        <NIcon size="20"><Close /></NIcon>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Close, Contract, Expand, Remove } from '@vicons/ionicons5';
import { WindowMinimise, Quit, WindowIsMaximised, WindowMaximise, WindowUnmaximise } from '@wailsApp/runtime';
import { NIcon } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { IsDarwin } from '@wailsApp/go/services/PreferencesSrv';

const prefStore = usePreferencesStore();

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});

const windowIsMaximised = ref(false);

const reduction = () => {
  WindowIsMaximised().then(isMaximised => {
    if (isMaximised) {
      windowIsMaximised.value = false;
      WindowUnmaximise();
    } else {
      windowIsMaximised.value = true;
      WindowMaximise();
    }
  });
};
</script>

<style lang="less" scoped>
.header {
  --wails-draggable: drag;
  background: v-bind('gtermThemeVars.titleColor');
  width: 100%;
}
.hover\:bg-custom-hover:hover {
  background-color: v-bind('gtermThemeVars.splitColor');
}
</style>
