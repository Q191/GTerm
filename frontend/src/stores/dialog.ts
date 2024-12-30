import { defineStore } from 'pinia';

export const useDialogStore = defineStore('dialog', {
  state: () => {
    return {
      addHostDialogVisible: false,
      aboutDialogVisible: false,
      settingsDialogVisible: false,
    };
  },
  actions: {
    openAddHostDialog() {
      this.addHostDialogVisible = true;
    },
    closeAddHostDialog() {
      this.addHostDialogVisible = false;
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
