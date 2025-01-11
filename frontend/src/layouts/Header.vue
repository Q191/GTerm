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

    <div v-if="!isDarwin" class="flex items-center ml-auto">
      <div class="window-control-btn" @click="WindowMinimise">
        <n-icon size="16"><icon icon="ph:minus-bold" /></n-icon>
      </div>
      <div class="window-control-btn" @click="reduction">
        <n-icon size="16">
          <span v-if="windowIsMaximised"><icon icon="ph:corners-in-bold" /></span>
          <span v-else><icon icon="ph:corners-out-bold" /></span>
        </n-icon>
      </div>
      <div class="window-control-btn close-btn rounded-tr-lg" @click="Quit">
        <n-icon size="16"><icon icon="ph:x-bold" /></n-icon>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { WindowMinimise, Quit, WindowIsMaximised, WindowMaximise, WindowUnmaximise } from '@wailsApp/runtime';
import { NIcon } from 'naive-ui';
import { Icon } from '@iconify/vue';
import { IsDarwin } from '@wailsApp/go/services/PreferencesSrv';
import ConnectionTabs from '@/layouts/ConnectionTabs.vue';
import { useConnectionStore } from '@/stores/connection';

const connectionTabsRef = ref();

defineExpose({
  connectionTabsRef,
});

const connectionStore = useConnectionStore();
const hasConnections = computed(() => connectionStore.hasConnections);
const router = useRouter();

const isDarwin = ref(false);

const toHost = () => {
  router.push({ name: 'Host' });
};

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
  width: 100%;
  height: 38px;
}

.window-control-btn {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: rgba(0, 0, 0, 0.06);
  }
}

.close-btn {
  &:hover {
    background-color: #e54d42;
    color: #fff;
  }
}
</style>
