<template>
  <div class="flex flex-col h-full menu-background">
    <!-- 主导航按钮组 -->
    <div class="flex-grow flex flex-col items-center pt-2">
      <!-- 资产清单 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-[38px] mb-1 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer top-menu-item"
            :class="{ active: selectedKey === 'Host' }"
            @click="handleSelect('Host')"
          >
            <n-icon size="x-large">
              <icon icon="ph:hard-drives" />
            </n-icon>
          </div>
        </template>
        资产清单
      </n-tooltip>

      <!-- 凭据 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-[38px] mb-1 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer top-menu-item"
            :class="{ active: selectedKey === 'Credential' }"
            @click="handleSelect('Credential')"
          >
            <n-icon size="x-large">
              <icon icon="ph:vault" />
            </n-icon>
          </div>
        </template>
        凭据
      </n-tooltip>
    </div>

    <!-- 底部功能按钮区 -->
    <div class="pb-2 flex flex-col justify-center items-center">
      <!-- 设置下拉菜单 -->
      <n-dropdown :options="settingsOptions" trigger="click" :width="150" @select="handleSettingsSelect">
        <n-tooltip placement="right" trigger="hover">
          <template #trigger>
            <div
              class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer bottom-menu-item"
            >
              <n-icon size="x-large">
                <icon icon="ph:gear-six" />
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
            class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer bottom-menu-item"
            @click="toggleTheme"
          >
            <n-icon size="x-large">
              <icon :icon="prefStore.isDark ? 'ph:sun' : 'ph:moon'" />
            </n-icon>
          </div>
        </template>
        {{ prefStore.isDark ? '切换亮色主题' : '切换暗色主题' }}
      </n-tooltip>

      <!-- Github链接 -->
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="w-11 h-9 mt-2 flex justify-center items-center hover:bg-custom-hover rounded cursor-pointer bottom-menu-item"
            @click="openGithub"
          >
            <n-icon size="x-large">
              <icon icon="ph:github-logo" />
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
import { NIcon, NDropdown, NTooltip, useThemeVars } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { useRouter } from 'vue-router';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { useDialogStore } from '@/stores/dialog';

const router = useRouter();
const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const selectedKey = ref('Host');
const themeVars = useThemeVars();

// 设置下拉菜单选项
const settingsOptions = [
  {
    label: '偏好设置',
    key: 'preferences',
    icon: renderIcon('ph:sliders-horizontal'),
  },
  {
    label: '检查更新',
    key: 'check-update',
    icon: renderIcon('ph:arrow-circle-up'),
  },
  {
    label: '关于',
    key: 'about',
    icon: renderIcon('ph:info'),
  },
];

// 渲染图标
function renderIcon(name: string) {
  return () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: name }) });
}

// 处理主菜单选择
const handleSelect = (key: string) => {
  selectedKey.value = key;
  router.push({ name: key });
};

// 处理设置菜单选择
const handleSettingsSelect = (key: string) => {
  switch (key) {
    case 'preferences':
      dialogStore.openPreferencesDialog();
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
  height: calc(100vh - 38px);
}

.hover\:bg-custom-hover:hover {
  background-color: v-bind('themeVars.hoverColor');
}

.top-menu-item {
  &:hover {
    color: v-bind('themeVars.primaryColor');
  }

  &.active {
    color: v-bind('themeVars.primaryColor');
    position: relative;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 3px;
      height: 20px;
      border-radius: 3px;
      background-color: v-bind('themeVars.primaryColor');
    }
  }
}

.bottom-menu-item {
  &:hover {
    color: v-bind('themeVars.primaryColor');
  }
}
</style>
