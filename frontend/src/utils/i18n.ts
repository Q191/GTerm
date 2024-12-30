import { lang } from '@/locales';
import { createI18n } from 'vue-i18n';

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  messages: { ...lang },
  locale: 'zh',
});
