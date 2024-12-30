import { defineStore } from 'pinia';

export const usePreferencesStore = defineStore('preferences', {
  state: () => {
    return {
      isDark: false,
    };
  },
  actions: {
    toDark() {
      this.isDark = true;
    },
    toLight() {
      this.isDark = false;
    },
  },
});
