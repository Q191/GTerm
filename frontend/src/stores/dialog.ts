import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useDialogStore = defineStore('dialog', () => {
  const aboutDialogVisible = ref(false);
  const preferencesDialogVisible = ref(false);

  const openAboutDialog = () => {
    aboutDialogVisible.value = true;
  };

  const closeAboutDialog = () => {
    aboutDialogVisible.value = false;
  };

  const openPreferencesDialog = () => {
    preferencesDialogVisible.value = true;
  };

  const closePreferencesDialog = () => {
    preferencesDialogVisible.value = false;
  };

  return {
    aboutDialogVisible,
    preferencesDialogVisible,
    openAboutDialog,
    closeAboutDialog,
    openPreferencesDialog,
    closePreferencesDialog,
  };
});
