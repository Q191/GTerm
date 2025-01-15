import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useDialogStore = defineStore('dialog', () => {
  const hostDialogVisible = ref(false);
  const aboutDialogVisible = ref(false);
  const preferencesDialogVisible = ref(false);
  const groupDialogVisible = ref(false);
  const isEditMode = ref(false);

  const openAddHostDialog = (edit = false) => {
    isEditMode.value = edit;
    hostDialogVisible.value = true;
  };

  const closeAddHostDialog = () => {
    hostDialogVisible.value = false;
  };

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

  const openAddGroupDialog = () => {
    groupDialogVisible.value = true;
  };

  const closeAddGroupDialog = () => {
    groupDialogVisible.value = false;
  };

  return {
    hostDialogVisible,
    aboutDialogVisible,
    preferencesDialogVisible,
    groupDialogVisible,
    openAddHostDialog,
    closeAddHostDialog,
    openAboutDialog,
    closeAboutDialog,
    openPreferencesDialog,
    closePreferencesDialog,
    openAddGroupDialog,
    closeAddGroupDialog,
    isEditMode,
  };
});
