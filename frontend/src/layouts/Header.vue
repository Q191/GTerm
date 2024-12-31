<template>
  <div ref="widthRef" class="flex items-center justify-between h-12 header rounded-t-lg">
    <div
      class="flex items-center"
      :class="{
        'justify-center flex-1': !hasConnections,
        'pl-2': !IsDarwin(),
        'pl-20': IsDarwin(),
      }"
    >
      <div class="flex items-center" @click="toHost">
        <img src="@/assets/images/icon.png" alt="Logo" class="w-6 h-6" />
        <span class="pl-2 text-base font-semibold" :class="{ hidden: hasConnections }">GTerm</span>
      </div>
      <ConnectionTabs v-if="hasConnections" class="ml-4" />
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
import ConnectionTabs from '@/components/ConnectionTabs.vue';
import { useConnectionStore } from '@/stores/connection';

const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const hasConnections = computed(() => connectionStore.hasConnections);
const router = useRouter();

const toHost = () => {
  router.push({ name: 'Host' });
};

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
  height: 38px;
}
.hover\:bg-custom-hover:hover {
  background-color: v-bind('gtermThemeVars.splitColor');
}

:deep(.connection-tabs) {
  --wails-draggable: none;
}
</style>
