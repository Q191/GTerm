<template>
  <n-modal
    v-model:show="dialogStore.groupDialogVisible"
    :close-on-esc="true"
    :negative-text="$t('group_dialog.cancel')"
    :on-close="dialogStore.closeAddGroupDialog"
    :positive-text="$t('group_dialog.confirm')"
    :show-icon="false"
    :title="$t('group_dialog.title')"
    preset="dialog"
    style="width: 500px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="80">
      <n-form-item :label="$t('group_dialog.name')" path="name">
        <n-input v-model:value="formValue.name" clearable />
      </n-form-item>

      <n-form-item :label="$t('group_dialog.description')" path="description">
        <n-input
          v-model:value="formValue.description"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 5 }"
          clearable
        />
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
    description: '',
  }),
);

const rules: FormRules = {
  name: {
    required: true,
    message: t('group_dialog.validation.name_required'),
    trigger: 'blur',
  },
};

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();
    const resp = await CreateGroup(formValue.value);
    if (!resp.ok) {
      message.error(resp.msg);
    } else {
      message.success('创建成功');
    }
    dialogStore.closeAddHostDialog();
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
