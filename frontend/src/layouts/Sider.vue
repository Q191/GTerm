<template>
  <div class="flex flex-col h-full menu-background">
    <!-- 主导航按钮组 -->
    <div class="flex-grow flex flex-col items-center pt-2">
      <!-- 资产清单 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-9 mb-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer"
            :class="{ 'bg-active': selectedKey === 'Host' }"
            @click="handleSelect('Host')"
          >
            <n-icon size="large" class="icon-hover">
              <icon icon="ph:hard-drives-duotone" />
            </n-icon>
          </div>
        </template>
        资产清单
      </n-tooltip>

      <!-- 登录凭证 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-9 mb-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer"
            :class="{ 'bg-active': selectedKey === 'Credentials' }"
            @click="handleSelect('Credentials')"
          >
            <n-icon size="large" class="icon-hover">
              <icon icon="ph:key-duotone" />
            </n-icon>
          </div>
        </template>
        登录凭证
      </n-tooltip>
    </div>

    <!-- 底部功能按钮区 -->
    <div class="pb-2 flex flex-col justify-center items-center">
      <!-- 设置下拉菜单 -->
      <n-dropdown :options="settingsOptions" trigger="click" :width="150" @select="handleSettingsSelect">
        <n-tooltip placement="right" trigger="hover">
          <template #trigger>
            <div class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer">
              <n-icon size="large" class="icon-hover">
                <icon icon="ph:gear-six-duotone" />
              </n-icon>
            </div>
          </template>
          设置
        </n-tooltip>
      </n-dropdown>

      <!-- 主题切换按钮 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer"
            @click="toggleTheme"
          >
            <n-icon size="large" class="icon-hover">
              <icon :icon="prefStore.isDark ? 'ph:sun-duotone' : 'ph:moon-duotone'" />
            </n-icon>
          </div>
        </template>
        {{ prefStore.isDark ? '切换亮色主题' : '切换暗色主题' }}
      </n-tooltip>

      <!-- Github链接 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer"
            @click="openGithub"
          >
            <n-icon size="large" class="icon-hover">
              <icon icon="ph:github-logo-duotone" />
            </n-icon>
          </div>
        </template>
        访问 Github
      </n-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NIcon, NDropdown, NTooltip } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { useRouter } from 'vue-router';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { useDialogStore } from '@/stores/dialog';

const router = useRouter();
const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const selectedKey = ref('Host');

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});

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
  prefStore.isDark ? prefStore.toLight() : prefStore.toDark();
};

// 打开 Github
const openGithub = () => {
  BrowserOpenURL('https://github.com/DLinkProjects/DLink');
};
</script>

<style lang="less" scoped>
.menu-background {
  background-color: v-bind('gtermThemeVars.sidebarColor');
  height: calc(100vh - 38px);
}

.hover\:bg-custom-hover:hover {
  background-color: v-bind('gtermThemeVars.splitColor');
}

.bg-active {
  background-color: v-bind('gtermThemeVars.splitColor');
}

.icon-hover {
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(1.1);
  }
}
</style>
