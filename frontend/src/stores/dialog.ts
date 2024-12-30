import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useDialogStore = defineStore('dialog', () => {
  const hostDialogVisible = ref(false);
  const aboutDialogVisible = ref(false);
  const settingsDialogVisible = ref(false);
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

  const openSettingsDialog = () => {
    settingsDialogVisible.value = true;
  };

  const closeSettingsDialog = () => {
    settingsDialogVisible.value = false;
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
    settingsDialogVisible,
    groupDialogVisible,
    openAddHostDialog,
    closeAddHostDialog,
    openAboutDialog,
    closeAboutDialog,
    openSettingsDialog,
    closeSettingsDialog,
    openAddGroupDialog,
    closeAddGroupDialog,
    isEditMode,
  };
});
