<template>
  <NModal
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
    <NForm
      ref="formRef"
      :model="formValue"
      :rules="rules"
      label-placement="left"
      label-width="80"
    >
      <NFormItem :label="$t('group_dialog.name')" path="name">
        <NInput v-model:value="formValue.name" clearable />
      </NFormItem>

      <NFormItem :label="$t('group_dialog.description')" path="description">
        <NInput
          v-model:value="formValue.description"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 5 }"
          clearable
        />
      </NFormItem>
    </NForm>
  </NModal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { FormInst, FormRules, NForm, NFormItem, NInput, NModal } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const dialogStore = useDialogStore();
const formRef = ref<FormInst | null>(null);

interface FormState {
  name: string;
  description: string;
}

const formValue = ref<FormState>({
  name: '',
  description: '',
});

const rules: FormRules = {
  name: {
    required: true,
    message: t('group_dialog.validation.name_required'),
    trigger: 'blur',
  },
};

const handleConfirm = () => {
  formRef.value?.validate(errors => {
    if (!errors) {
      console.log('验证通过');
      // TODO: 处理表单提交
      dialogStore.closeAddGroupDialog();
    }
  });
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