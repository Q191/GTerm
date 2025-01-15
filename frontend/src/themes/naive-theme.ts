import { merge } from 'lodash';
import { GlobalThemeOverrides } from 'naive-ui';

export const themeOverrides: GlobalThemeOverrides = {
  common: {
    borderColor: 'rgb(239, 239, 245)',
  },
  Divider: {
    color: '#A8A8A9',
  },
  Form: {
    labelFontSizeTopSmall: '12px',
    labelFontSizeTopMedium: '13px',
    labelFontSizeTopLarge: '13px',
    labelHeightSmall: '18px',
    labelHeightMedium: '18px',
    labelHeightLarge: '18px',
    labelPaddingVertical: '0 0 5px 2px',
    feedbackHeightSmall: '18px',
    feedbackHeightMedium: '18px',
    feedbackHeightLarge: '20px',
    feedbackFontSizeSmall: '11px',
    feedbackFontSizeMedium: '12px',
    feedbackFontSizeLarge: '12px',
    labelFontWeight: '450',
  },
};

const _darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    borderColor: 'rgba(255, 255, 255, 0.09)',
  },
};

export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides);
