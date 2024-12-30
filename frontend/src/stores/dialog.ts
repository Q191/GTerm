import { defineStore } from 'pinia';

export const useDialogStore = defineStore('dialog', {
  state: () => {
    return {
      addServerDialogVisible: false,
      aboutDialogVisible: false,
      settingsDialogVisible: false,
    };
  },
  actions: {
    openAddServerDialog() {
      this.addServerDialogVisible = true;
    },
    closeAddServerDialog() {
      this.addServerDialogVisible = false;
    },
    openAboutDialog() {
      this.aboutDialogVisible = true;
    },
    closeAboutDialog() {
      this.aboutDialogVisible = false;
    },
    openSettingsDialog() {
      this.settingsDialogVisible = true;
    },
    closeSettingsDialog() {
      this.settingsDialogVisible = false;
    },
  },
});
