import { merge } from 'lodash';

export const themeOverrides: any = {
  common: {
    primaryColor: '#4098FC',
    primaryColorHover: '#60A5FA',
  },
};

const _darkThemeOverrides: any = {
  common: {
    modalColor: '#1E293B',
  },
};
export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides);

