import { merge } from 'lodash';
import { GlobalThemeOverrides } from 'naive-ui';

export const themeOverrides: GlobalThemeOverrides = {
  common: {
    borderColor: 'rgb(239, 239, 245)',
  },
};

const _darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    borderColor: 'rgba(255, 255, 255, 0.09)',
  },
};

export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides);
