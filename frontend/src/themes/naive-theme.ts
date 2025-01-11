import { merge } from 'lodash';
import { GlobalThemeOverrides } from 'naive-ui';

export const themeOverrides: GlobalThemeOverrides = {
  common: {},
};

const _darkThemeOverrides: GlobalThemeOverrides = {
  common: {},
};
export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides);
