import zh from './zh-cn.json';
import en from './en-us.json';
import ja from './ja-jp.json';
import ko from './ko-kr.json';

export const lang = {
  zh,
  en,
  ja,
  ko,
};

export type Language = 'zh' | 'en' | 'ja' | 'ko';

export const languageOptions = [
  {
    label: '简体中文',
    value: 'zh',
  },
  {
    label: 'English',
    value: 'en',
  },
  {
    label: '日本語',
    value: 'ja',
  },
  {
    label: '한국어',
    value: 'ko',
  },
];

export const languageMap = {
  'zh-CN': 'zh',
  'en-US': 'en',
  'ja-JP': 'ja',
  'ko-KR': 'ko',
};
