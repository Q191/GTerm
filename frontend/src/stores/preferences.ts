import { defineStore } from 'pinia';
import type { SelectOption } from 'naive-ui';

const DEFAULT_SIDEBAR_WIDTH = 260;
const DEFAULT_LANGUAGE = 'auto';

export const languageOptions: SelectOption[] = [
  {
    label: '跟随系统',
    value: 'auto',
  },
  {
    label: '简体中文',
    value: 'zh-CN',
  },
  {
    label: '繁體中文',
    value: 'zh-TW',
  },
  {
    label: 'English',
    value: 'en-US',
  },
  {
    label: '日本語',
    value: 'ja-JP',
  },
  {
    label: '한국어',
    value: 'ko-KR',
  },
];

export const usePreferencesStore = defineStore('preferences', {
  state: () => {
    return {
      isDark: false,
      themeMode: localStorage.getItem('themeMode') || 'auto',
      sidebarWidth: Number(localStorage.getItem('sidebarWidth')) || DEFAULT_SIDEBAR_WIDTH,
      language: localStorage.getItem('language') || 'auto',
      activeLanguage: DEFAULT_LANGUAGE,
    };
  },
  actions: {
    toDark() {
      this.isDark = true;
    },
    toLight() {
      this.isDark = false;
    },
    updateThemeMode(mode: 'auto' | 'light' | 'dark') {
      this.themeMode = mode;
      localStorage.setItem('themeMode', mode);

      switch (mode) {
        case 'light':
          this.toLight();
          break;
        case 'dark':
          this.toDark();
          break;
        case 'auto':
          this.updateThemeBySystem();
          break;
      }
    },
    updateThemeBySystem() {
      if (this.themeMode === 'auto') {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        if (prefersDark) {
          this.toDark();
        } else {
          this.toLight();
        }
      }
    },
    updateSidebarWidth(width: number) {
      this.sidebarWidth = width;
      localStorage.setItem('sidebarWidth', width.toString());
      window.dispatchEvent(new CustomEvent('sidebar-width-change', { detail: width }));
    },
    resetSidebarWidth() {
      localStorage.removeItem('sidebarWidth');
      this.updateSidebarWidth(DEFAULT_SIDEBAR_WIDTH);
    },
    updateLanguage(lang: string) {
      this.language = lang;
      localStorage.setItem('language', lang);

      if (lang === 'auto') {
        this.updateLanguageBySystem();
      } else {
        this.activeLanguage = lang;
        this.applyLanguage();
      }
    },
    updateLanguageBySystem() {
      if (this.language === 'auto') {
        const systemLanguages = navigator.languages;
        let detectedLang = DEFAULT_LANGUAGE;

        // 遍历系统语言列表，找到第一个支持的语言
        for (const lang of systemLanguages) {
          const normalizedLang = this.normalizeLanguage(lang);
          if (languageOptions.some(opt => opt.value === normalizedLang && opt.value !== 'auto')) {
            detectedLang = normalizedLang;
            break;
          }
        }

        this.activeLanguage = detectedLang;
        this.applyLanguage();
      }
    },
    normalizeLanguage(lang: string): string {
      // 标准化语言代码
      const langMap: { [key: string]: string } = {
        zh: 'zh-CN',
        'zh-CN': 'zh-CN',
        'zh-TW': 'zh-TW',
        'zh-HK': 'zh-TW',
        en: 'en-US',
        ja: 'ja-JP',
        ko: 'ko-KR',
      };

      const code = lang.split('-')[0].toLowerCase();
      return langMap[lang] || langMap[code] || DEFAULT_LANGUAGE;
    },
    applyLanguage() {
      // TODO: 实现语言切换功能，使用 this.activeLanguage
      console.log('切换语言到:', this.activeLanguage);
    },
  },
});
