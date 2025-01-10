import { merge } from 'lodash';

export const themeOverrides: any = {
  common: {},
};

const _darkThemeOverrides: any = {
  common: {},
};
export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides);
