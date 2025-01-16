import { defineStore } from 'pinia';
import { ref } from 'vue';
import { model } from '@wailsApp/go/models';

export const useDialogStore = defineStore('dialog', () => {
  const hostDialogVisible = ref(false);
  const aboutDialogVisible = ref(false);
  const preferencesDialogVisible = ref(false);
  const groupDialogVisible = ref(false);
  const isEditMode = ref(false);
  const editHost = ref<model.Host | null>(null);

  const openAddHostDialog = (edit = false, host?: model.Host) => {
    isEditMode.value = edit;
    editHost.value = host || null;
    hostDialogVisible.value = true;
  };

  const closeAddHostDialog = () => {
    hostDialogVisible.value = false;
    editHost.value = null;
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
    editHost,
  };
});
