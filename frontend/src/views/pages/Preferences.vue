<template>
  <div class="page-container">
    <div ref="sidebarRef" class="settings-sidebar">
      <div class="menu-items">
        <div
          v-for="(item, index) in menuItems"
          :key="index"
          class="menu-item"
          :class="{ active: selectedMenu === item.key }"
          @click="selectedMenu = item.key"
        >
          <Icon :icon="item.icon" class="menu-icon" />
          <span>{{ item.label }}</span>
        </div>
      </div>
    </div>

    <div class="content-container">
      <NScrollbar>
        <div v-if="selectedMenu === 'general'" class="settings-panel">
          <div class="setting-section">
            <div>
              <h2 class="section-title">{{ $t('frontend.preferencesModal.theme.title') }}</h2>
              <p class="section-description">选择适合您使用场景的主题</p>
            </div>

            <div class="theme-cards">
              <div @click="prefStore.updateThemeMode('light')">
                <div class="theme-preview" :class="{ active: prefStore.themeMode === 'light' }">
                  <div class="theme-preview light" v-html="ThemeSvg" />
                </div>
                <div class="theme-card-footer">
                  <span>{{ $t('frontend.preferencesModal.theme.light') }}</span>
                </div>
              </div>

              <div @click="prefStore.updateThemeMode('dark')">
                <div class="theme-preview" :class="{ active: prefStore.themeMode === 'dark' }">
                  <div class="theme-preview dark" v-html="ThemeSvg" />
                </div>
                <div class="theme-card-footer">
                  <span>{{ $t('frontend.preferencesModal.theme.dark') }}</span>
                </div>
              </div>

              <div @click="prefStore.updateThemeMode('auto')">
                <div class="theme-preview" :class="{ active: prefStore.themeMode === 'auto' }">
                  <div class="auto-preview" v-html="AutoThemeSvg" />
                </div>
                <div class="theme-card-footer">
                  <span>{{ $t('frontend.preferencesModal.theme.auto') }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="setting-section">
            <div>
              <h2 class="section-title">{{ $t('frontend.preferencesModal.language.title') }}</h2>
              <p class="section-description">设置应用界面的显示语言</p>
            </div>

            <NSelect
              v-model:value="prefStore.language"
              :options="languageOptions"
              size="medium"
              class="language-selector"
              @update:value="prefStore.updateLanguage"
            />
          </div>

          <div class="setting-section">
            <div>
              <h2 class="section-title">{{ $t('frontend.preferencesModal.sidebar.title') }}</h2>
              <p class="section-description">自定义侧边栏宽度</p>
            </div>

            <div class="sidebar-width-control">
              <NSlider
                v-model:value="prefStore.sidebarWidth"
                :min="260"
                :max="380"
                :step="1"
                class="slider"
                @update:value="value => prefStore.updateSidebarWidth(value)"
              />
              <div class="sidebar-width-actions">
                <div class="width-display">{{ prefStore.sidebarWidth }}px</div>
                <NButton text type="primary" size="small" @click="prefStore.resetSidebarWidth">
                  {{ $t('frontend.preferencesModal.sidebar.reset') }}
                </NButton>
              </div>
            </div>
          </div>
        </div>
      </NScrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NScrollbar, NSelect, NSlider, useThemeVars } from 'naive-ui';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import AutoThemeSvg from '@/assets/images/auto_theme.svg?raw';
import ThemeSvg from '@/assets/images/theme.svg?raw';
import { languageOptions, usePreferencesStore } from '@/stores/preferences';

const prefStore = usePreferencesStore();
const themeVars = useThemeVars();

const sidebarRef = ref<HTMLElement | null>(null);
const selectedMenu = ref('general');

const menuItems = computed(() => [
  {
    key: 'general',
    label: '通用设置',
    icon: 'ph:gear',
  },
  {
    key: 'terminal',
    label: '终端',
    icon: 'ph:terminal',
  },
  {
    key: 'shortcuts',
    label: '快捷键',
    icon: 'ph:keyboard',
  },
]);

const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

onMounted(() => {
  mediaQuery.addEventListener('change', prefStore.updateThemeBySystem);
  prefStore.updateThemeBySystem();
  prefStore.updateLanguageBySystem();
});

onUnmounted(() => {
  mediaQuery.removeEventListener('change', prefStore.updateThemeBySystem);
});
</script>

<style lang="less" scoped>
.page-container {
  height: 100%;
  display: flex;
  position: relative;
  color: v-bind('themeVars.textColorBase');
}

.settings-sidebar {
  width: 220px;
  border-right: 1px solid v-bind('themeVars.borderColor');
  padding: 8px;
  overflow-y: auto;
}

.menu-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-item {
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: v-bind('themeVars.textColor3');
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;

  &:hover {
    color: v-bind('themeVars.textColorBase');
    background-color: v-bind('themeVars.hoverColor');
  }

  &.active {
    background-color: v-bind('`${themeVars.primaryColor}15`');
    color: v-bind('themeVars.primaryColor');
  }

  .menu-icon {
    font-size: 18px;
  }
}

.content-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.settings-panel {
  display: flex;
  flex-direction: column;
  padding: 24px;
  gap: 32px;
}

.setting-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: v-bind('themeVars.textColorBase');
  margin: 0 0 4px 0;
}

.section-description {
  font-size: 14px;
  color: v-bind('themeVars.textColor3');
  margin: 0 0 16px 0;
  opacity: 0.85;
  line-height: 1.5;
}

.theme-cards {
  display: flex;
  gap: 16px;
}

.theme-preview {
  padding: 2px;
  border-radius: 8px;
  position: relative;
  box-sizing: border-box;
  &.active {
    outline: 2px solid v-bind('themeVars.primaryColor');
  }
}

.theme-preview :deep(svg) {
  border-radius: 8px;
  overflow: hidden;
  width: 100%;
  height: 100%;
  display: block;
}

.theme-preview.light :deep(svg) {
  --bg-color: #ffffff;
  --header-color: #f5f5f5;
  --sidebar-color: #f0f0f0;
  --icon-color: #d9d9d9;
  --text-color: #000000;
}

.theme-preview.dark :deep(svg) {
  --bg-color: #121212;
  --header-color: #1f1f1f;
  --sidebar-color: #262626;
  --icon-color: #404040;
  --text-color: #ffffff;
}

.theme-card-footer {
  display: flex;
  justify-content: center;
}

.active .theme-card-footer {
  color: v-bind('themeVars.primaryColor');
}

.language-selector {
  width: 100%;
}

.sidebar-width-control {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  margin-top: 8px;
}

.sidebar-width-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.width-display {
  background-color: v-bind('`${themeVars.primaryColor}10`');
  color: v-bind('themeVars.primaryColor');
  font-weight: 500;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 14px;
}

.slider {
  width: 100%;
}

.auto-preview {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.auto-preview :deep(svg) {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}
</style>
