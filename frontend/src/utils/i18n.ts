import { createI18n } from 'vue-i18n';
import { lang } from '@/locales';

const getBrowserLanguage = (): 'zh' | 'en' => {
  const storedLang = localStorage.getItem('language');
  if (storedLang === 'zh' || storedLang === 'en') {
    return storedLang;
  }

  const browserLang = navigator.language.toLowerCase();
  return browserLang.startsWith('zh') ? 'zh' : 'en';
};

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: getBrowserLanguage(),
  fallbackLocale: 'en',
  messages: lang,
  missingWarn: false,
  fallbackWarn: false,
});
