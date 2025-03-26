import { defineStore } from 'pinia';
import { languageMap, languageOptions } from '@/locales';
import { i18n } from '@/utils/i18n';

export { languageOptions };

type Language = 'zh';
type ThemeMode = 'light' | 'dark' | 'auto';

const getInitialLanguage = (): Language => {
  const storedLang = localStorage.getItem('language');
  if (storedLang === 'zh') {
    return storedLang;
  }
  const browserLang = navigator.language.toLowerCase();
  // 尝试匹配浏览器语言
  if (browserLang.startsWith('zh')) {
    return 'zh';
  }
  // 默认使用中文
  return 'zh';
};

export const usePreferencesStore = defineStore('preferences', {
  state: () => ({
    language: getInitialLanguage(),
    themeMode: (localStorage.getItem('theme') || 'auto') as ThemeMode,
    isDark: false,
    sidebarWidth: Number(localStorage.getItem('sidebarWidth')) || 260,
  }),

  actions: {
    updateLanguage(lang: Language) {
      this.language = lang;
      i18n.global.locale.value = lang;
      localStorage.setItem('language', lang);
    },

    updateLanguageBySystem() {
      const browserLang = navigator.language.toLowerCase();
      if (browserLang.startsWith('zh') && !localStorage.getItem('language')) {
        this.updateLanguage('zh');
      }
    },

    updateThemeMode(mode: ThemeMode) {
      this.themeMode = mode;
      localStorage.setItem('theme', mode);
      this.updateThemeBySystem();
    },

    updateThemeBySystem() {
      if (this.themeMode === 'auto') {
        this.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      } else {
        this.isDark = this.themeMode === 'dark';
      }
    },

    toDark() {
      this.updateThemeMode('dark');
    },

    toLight() {
      this.updateThemeMode('light');
    },

    updateSidebarWidth(width: number) {
      this.sidebarWidth = width;
      localStorage.setItem('sidebarWidth', width.toString());
      window.dispatchEvent(new CustomEvent('sidebar-width-change', { detail: width }));
    },

    resetSidebarWidth() {
      this.updateSidebarWidth(260);
    },
  },
});
