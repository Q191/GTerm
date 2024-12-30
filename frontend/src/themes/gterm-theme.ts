export const gtermLightTheme = {
  titleColor: '#FFFFFF',
  ribbonColor: '#F8FAFC',
  ribbonActiveColor: '#F1F5F9',
  sidebarColor: '#FFFFFF',
  splitColor: '#E2E8F0',
  cardColor: '#FFFFFF',
  cardHoverColor: '#F8FAFC',
  borderColor: '#E2E8F0',
  primaryColor: '#4098FC',
};

export const gtermDarkTheme = {
  titleColor: '#0F172A',
  ribbonColor: '#1E293B',
  ribbonActiveColor: '#334155',
  sidebarColor: '#0F172A',
  splitColor: '#334155',
  cardColor: '#1E293B',
  borderColor: '#1E293B',
  cardHoverColor: '#334155',
  primaryColor: '#60A5FA',
};

export const gtermTheme = (isDark: boolean) => {
  return {
    ...(isDark ? gtermDarkTheme : gtermLightTheme),
    menuHoverColor: isDark ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.04)',
  };
};
