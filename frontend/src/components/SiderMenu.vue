<template>
  <div class="flex flex-col h-full menu-background">
    <!-- 主导航菜单 -->
    <div class="flex-grow">
      <NMenu
        :options="menuOptions"
        :collapsed="true"
        :collapsed-width="60"
        :value="selectedKey"
        class="border-0"
        @update:value="handleSelect"
      />
    </div>

    <!-- 底部功能按钮区 -->
    <div class="pb-2 flex flex-col justify-center items-center">
      <!-- 设置下拉菜单 -->
      <NDropdown :options="settingsOptions" trigger="click" :width="150" @select="handleSettingsSelect">
        <div class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer">
          <NIcon size="large">
            <Icon icon="ph:gear-six-duotone" />
          </NIcon>
        </div>
      </NDropdown>

      <!-- 主题切换按钮 -->
      <div
        class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer"
        @click="toggleTheme"
      >
        <NIcon size="large">
          <Icon :icon="prefStore.isDark ? 'ph:sun-duotone' : 'ph:moon-duotone'" />
        </NIcon>
      </div>

      <!-- Github链接 -->
      <div class="w-11 h-9 mt-2 flex justify-center items-center cursor-pointer" @click="openGithub">
        <NIcon size="large">
          <Icon icon="ph:github-logo-duotone" />
        </NIcon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NMenu, NIcon, NDropdown } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { useRouter } from 'vue-router';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { useDialogStore } from '@/stores/dialog';

const router = useRouter();
const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const selectedKey = ref('host');

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});

// 主导航菜单选项
const menuOptions = [
  {
    label: '资产清单',
    key: 'host',
    icon: renderIcon('ph:hard-drives-duotone'),
  },
  {
    label: '登录凭证',
    key: 'credentials',
    icon: renderIcon('ph:key-duotone'),
  },
  {
    label: '端口转发',
    key: 'port-forwarding',
    icon: renderIcon('ph:arrows-left-right-duotone'),
  },
  {
    label: '代码片段',
    key: 'snippets',
    icon: renderIcon('ph:code-duotone'),
  },
  {
    label: '历史记录',
    key: 'history',
    icon: renderIcon('ph:clock-counter-clockwise-duotone'),
  },
];

// 设置下拉菜单选项
const settingsOptions = [
  {
    label: '偏好设置',
    key: 'settings',
    icon: renderIcon('ph:gear-six-duotone'),
  },
  {
    label: '检查更新',
    key: 'check-update',
    icon: renderIcon('ph:arrow-circle-up-duotone'),
  },
  {
    type: 'divider',
    key: 'd1',
  },
  {
    label: '关于',
    key: 'about',
    icon: renderIcon('ph:info-duotone'),
  },
];

// 渲染图标
function renderIcon(name: string) {
  return () => h(NIcon, null, { default: () => h(Icon, { icon: name }) });
}

// 处理主菜单选择
const handleSelect = (key: string) => {
  selectedKey.value = key;
  router.push({ name: key });
};

// 处理设置菜单选择
const handleSettingsSelect = (key: string) => {
  switch (key) {
    case 'settings':
      dialogStore.openSettingsDialog();
      break;
    case 'about':
      dialogStore.openAboutDialog();
      break;
    case 'check-update':
      // TODO: 实现检查更新功能
      break;
  }
};

// 切换主题
const toggleTheme = () => {
  prefStore.toDark();
};

// 打开 Github
const openGithub = () => {
  BrowserOpenURL('https://github.com/DLinkProjects/DLink');
};
</script>

<style lang="less" scoped>
.menu-background {
  background-color: v-bind('gtermThemeVars.sidebarColor');
  height: 100vh;
}

.hover\:bg-custom-hover:hover {
  background-color: v-bind('gtermThemeVars.splitColor');
}
</style>
