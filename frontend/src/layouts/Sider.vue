<template>
  <div class="sider-container menu-background">
    <div class="top-menu">
      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="menu-item top-menu-item"
            :class="{ active: selectedKey === 'Connection' }"
            @click="handleSelect('Connection')"
          >
            <n-icon size="x-large">
              <icon icon="ph:hard-drives" />
            </n-icon>
          </div>
        </template>
        {{ $t('frontend.sider.assets') }}
      </n-tooltip>

      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div
            class="menu-item top-menu-item"
            :class="{ active: selectedKey === 'Credential' }"
            @click="handleSelect('Credential')"
          >
            <n-icon size="x-large">
              <icon icon="ph:vault" />
            </n-icon>
          </div>
        </template>
        {{ $t('frontend.sider.credentials') }}
      </n-tooltip>
    </div>

    <div class="bottom-menu">
      <n-dropdown
        :options="[
          {
            label: $t('frontend.sider.menu.preferences'),
            key: 'preferences',
            icon: () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: 'ph:sliders-horizontal' }) }),
          },
          {
            label: $t('frontend.sider.menu.check_update'),
            key: 'check-update',
            icon: () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: 'ph:arrow-circle-up' }) }),
          },
          {
            label: $t('frontend.sider.menu.about'),
            key: 'about',
            icon: () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: 'ph:info' }) }),
          },
        ]"
        trigger="click"
        :width="180"
        placement="right"
        @select="handleSettingsSelect"
      >
        <n-tooltip placement="right" trigger="hover">
          <template #trigger>
            <div class="menu-item bottom-menu-item">
              <n-icon size="x-large">
                <icon icon="ph:gear-six" />
              </n-icon>
            </div>
          </template>
          {{ $t('frontend.sider.settings') }}
        </n-tooltip>
      </n-dropdown>

      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div class="menu-item bottom-menu-item" @click="toggleTheme">
            <n-icon size="x-large">
              <icon :icon="prefStore.isDark ? 'ph:sun' : 'ph:moon'" />
            </n-icon>
          </div>
        </template>
        {{ prefStore.isDark ? $t('frontend.sider.theme.toggle_light') : $t('frontend.sider.theme.toggle_dark') }}
      </n-tooltip>

      <n-tooltip placement="right" trigger="hover">
        <template #trigger>
          <div class="menu-item bottom-menu-item" @click="openGithub">
            <n-icon size="x-large">
              <icon icon="ph:github-logo" />
            </n-icon>
          </div>
        </template>
        {{ $t('frontend.sider.github') }}
      </n-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NIcon, NDropdown, NTooltip, useThemeVars } from 'naive-ui';
import { usePreferencesStore } from '@/stores/preferences';
import { useRouter, useRoute } from 'vue-router';
import { BrowserOpenURL } from '@wailsApp/runtime';
import { useDialogStore } from '@/stores/dialog';
import { onMounted, watch } from 'vue';

const router = useRouter();
const route = useRoute();
const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const selectedKey = ref('Connection');
const themeVars = useThemeVars();

const updateSelectedKey = () => {
  const routeName = route.name as string;
  if (routeName === 'Connection' || routeName === 'Credential') {
    selectedKey.value = routeName;
  }
};

watch(() => route.name, updateSelectedKey);

onMounted(() => {
  updateSelectedKey();
});

const handleSelect = (key: string) => {
  selectedKey.value = key;
  router.push({ name: key });
};

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

const toggleTheme = () => {
  prefStore.isDark ? prefStore.toLight() : prefStore.toDark();
};

const openGithub = () => {
  BrowserOpenURL('https://github.com/DLinkProjects/DLink');
};
</script>

<style lang="less" scoped>
.sider-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.top-menu {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 8px;
}

.bottom-menu {
  padding-bottom: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.menu-item {
  width: 44px;
  height: 38px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 4px;
  cursor: pointer;

  &:hover {
    background-color: v-bind('themeVars.hoverColor');
  }
}

.top-menu-item {
  margin-bottom: 4px;

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
  height: 36px;
  margin-top: 8px;

  &:hover {
    color: v-bind('themeVars.primaryColor');
  }
}

.menu-background {
  height: calc(100vh - 38px);
}
</style>
