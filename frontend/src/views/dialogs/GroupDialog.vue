<template>
  <n-modal
    v-model:show="dialogStore.groupDialogVisible"
    :close-on-esc="true"
    :negative-text="$t('groupDialog.cancel')"
    :on-close="resetForm"
    :positive-text="$t('groupDialog.confirm')"
    :show-icon="false"
    :title="$t('groupDialog.title')"
    preset="dialog"
    style="width: 500px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="80">
      <n-form-item :label="$t('groupDialog.name')" path="name">
        <n-input v-model:value="formValue.name" clearable />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { FormInst, FormRules, NForm, NFormItem, NInput, NModal, useMessage } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { CreateGroup } from '@wailsApp/go/services/GroupSrv';
import { model } from '@wailsApp/go/models';

const { t } = useI18n();
const dialogStore = useDialogStore();
const formRef = ref<FormInst | null>(null);
const message = useMessage();

const formValue = ref(
  model.Group.createFrom({
    name: '',
  }),
);

const resetForm = () => {
  formValue.value = model.Group.createFrom({
    name: '',
  });
  dialogStore.closeGroupDialog();
};

const rules: FormRules = {
  name: {
    required: true,
    message: t('groupDialog.validation.nameRequired'),
    trigger: 'blur',
  },
};

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();
    const resp = await CreateGroup(formValue.value);
    if (!resp.ok) {
      message.error(resp.msg);
      return false;
    }
    message.success(t('message.createSuccess'));
    dialogStore.closeGroupDialog();
  } catch (errors) {
    return false;
  }
};
</script>

<style lang="less" scoped>
:deep(.n-form-item .n-form-item-label) {
  font-size: 13px;
}

:deep(.n-input) {
  font-size: 13px;
}
</style>
