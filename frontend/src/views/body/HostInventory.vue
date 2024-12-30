<template>
  <div class="p-2">
    <NInput placeholder="Find a host or user@hostname...">
      <template #suffix>
        <NButton size="tiny" type="primary" disabled>connect</NButton>
      </template>
    </NInput>
  </div>
  <div class="mx-2 flex space-x-2">
    <NButton size="small" strong secondary @click="dialogStore.openAddServerDialog">
      <Icon icon="ph:hard-drives-duotone" class="mr-1" />
      添加服务器
    </NButton>
  </div>

  <NH3 class="mx-5" prefix="bar" align-text>分组</NH3>

  <div class="mx-5 mt-2 mp-2">
    <NGrid x-gap="8" y-gap="8" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
      <NGi v-for="index in 3" :key="index">
        <div class="p-2 flex justify-between items-center card-style" @click="toTerminal">
          <div class="flex items-center">
            <Icon icon="ph:squares-four-duotone" class="bg-blue-500 rounded-lg p-1" style="font-size: 34px" />
            <div class="flex flex-col ml-5">
              <p>测试分组</p>
              <div class="type-size">7 Hosts</div>
            </div>
          </div>
          <div class="hover:bg-custom-hover rounded-lg p-1"><Edit /></div>
        </div>
      </NGi>
    </NGrid>
  </div>

  <NH3 class="mx-5" prefix="bar" align-text>主机列表</NH3>

  <div class="mx-5 mt-2 mb-5">
    <NGrid x-gap="8" y-gap="8" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
      <NGi v-for="index in 100" :key="index">
        <div class="p-2 flex justify-between items-center card-style" @click="toTerminal">
          <div class="flex items-center">
            <Centos size="34" color="#FFF" class="bg-blue-500 rounded-lg p-1" />
            <div class="flex flex-col ml-5">
              <p>测试服务器</p>
              <div class="type-size">ssh,root</div>
            </div>
          </div>
          <div class="hover:bg-custom-hover rounded-lg p-1"><Edit /></div>
        </div>
      </NGi>
    </NGrid>
  </div>
</template>

<script setup lang="ts">
import Centos from '@/assets/icons/Centos.vue';
import Edit from '@/assets/icons/Edit.vue';
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { Icon } from '@iconify/vue';
import { NButton, NGi, NGrid, NH3, NInput } from 'naive-ui';

const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const router = useRouter();

const toTerminal = () => {
  router.push({ name: 'Terminal' });
};

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});
</script>

<style lang="less" scoped>
.card-style {
  height: 60px;
  border-radius: 1rem;
  transition: all 0.3s ease;
  cursor: pointer;
  background-color: v-bind('gtermThemeVars.cardColor');
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid v-bind('gtermThemeVars.borderColor');
}

.card-style:hover {
  background-color: v-bind('gtermThemeVars.cardHoverColor');
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.08);
  border-color: v-bind('gtermThemeVars.splitColor');
}

.type-size {
  font-size: 12px;
  color: #666;
}

.hover\:bg-custom-hover:hover {
  background-color: v-bind('gtermThemeVars.splitColor');
}
</style>
