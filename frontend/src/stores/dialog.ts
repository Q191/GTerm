import { defineStore } from 'pinia';
import { ref } from 'vue';
import { model } from '@wailsApp/go/models';

export const useDialogStore = defineStore('dialog', () => {
  const connDialogVisible = ref(false);
  const aboutDialogVisible = ref(false);
  const preferencesDialogVisible = ref(false);
  const groupDialogVisible = ref(false);
  const isEditMode = ref(false);
  const editConnection = ref<model.Connection | null>(null);

  const openConnDialog = (edit = false, connection?: model.Connection) => {
    isEditMode.value = edit;
    editConnection.value = connection || null;
    connDialogVisible.value = true;
  };

  const closeConnDialog = () => {
    connDialogVisible.value = false;
    editConnection.value = null;
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

  const openGroupDialog = () => {
    groupDialogVisible.value = true;
  };

  const closeGroupDialog = () => {
    groupDialogVisible.value = false;
  };

  return {
    connDialogVisible,
    aboutDialogVisible,
    preferencesDialogVisible,
    groupDialogVisible,
    openConnDialog,
    closeConnDialog,
    openAboutDialog,
    closeAboutDialog,
    openPreferencesDialog,
    closePreferencesDialog,
    openGroupDialog,
    closeGroupDialog,
    isEditMode,
    editConnection,
  };
});
