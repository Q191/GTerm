import { defineStore } from 'pinia';
import { languageMap, languageOptions } from '@/locales';
import { i18n } from '@/utils/i18n';

export { languageOptions };

type Language = 'zh' | 'en';
type ThemeMode = 'light' | 'dark' | 'auto';

const getInitialLanguage = (): Language => {
  const storedLang = localStorage.getItem('language');
  if (storedLang === 'zh' || storedLang === 'en') {
    return storedLang;
  }
  const browserLang = navigator.language.toLowerCase();
  return browserLang.startsWith('zh') ? 'zh' : 'en';
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
      const newLang: Language = browserLang.startsWith('zh') ? 'zh' : 'en';
      if (!localStorage.getItem('language')) {
        this.updateLanguage(newLang);
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
