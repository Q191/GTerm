<template>
  <div ref="widthRef" class="flex items-center justify-between h-12 header rounded-t-lg">
    <div
      class="flex items-center absolute transition-all duration-300 ease-in-out"
      :class="{
        'pl-2': !isDarwin,
        'pl-20': isDarwin,
        'left-1/2 -translate-x-1/2': !hasConnections,
        'left-0 translate-x-0': hasConnections,
      }"
    >
      <div class="flex items-center" @click="toHost">
        <img src="@/assets/images/icon.png" alt="Logo" class="w-6 h-6" />
        <span class="pl-2 text-base font-semibold" :class="{ hidden: hasConnections }">GTerm</span>
      </div>
      <ConnectionTabs v-if="hasConnections" ref="connectionTabsRef" class="ml-4" />
    </div>

    <div v-if="!isDarwin" class="flex items-center mt-0">
      <div class="w-12 h-12 flex items-center justify-center hover:bg-custom-hover" @click="WindowMinimise">
        <n-icon size="20"><Remove /></n-icon>
      </div>
      <div class="w-12 h-12 flex items-center justify-center hover:bg-custom-hover" @click="reduction">
        <n-icon size="20">
          <span v-if="windowIsMaximised"><Contract /></span>
          <span v-else><Expand /></span>
        </n-icon>
      </div>
      <div class="w-12 h-12 flex items-center justify-center hover:bg-custom-hover hover:rounded-tr-lg" @click="Quit">
        <n-icon size="20"><Close /></n-icon>
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

const connectionTabsRef = ref();

defineExpose({
  connectionTabsRef,
});

const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const hasConnections = computed(() => connectionStore.hasConnections);
const router = useRouter();

const isDarwin = ref(false);

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

onMounted(async () => {
  isDarwin.value = await IsDarwin();
});
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
</style>
