<template>
  <n-modal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('groupDialog.cancel')"
    :on-close="resetForm"
    :positive-text="$t('groupDialog.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('groupDialog.editTitle') : $t('groupDialog.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-form ref="formRef" :model="formValue" :rules="rules">
      <n-form-item path="name" :label="$t('groupDialog.name')">
        <n-input v-model:value="formValue.name" clearable :placeholder="$t('groupDialog.placeholder.name')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { enums, model } from '@wailsApp/go/models';
import { CreateGroup, UpdateGroup } from '@wailsApp/go/services/GroupSrv';
import { FormInst, FormRules, NForm, NFormItem, NInput, NModal, useMessage } from 'naive-ui';
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const message = useMessage();

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  group?: model.Group;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value),
});

const defaultGroup = {
  name: '',
};

const formValue = ref(model.Group.createFrom(defaultGroup));

watch(
  () => props.group,
  newGroup => {
    if (newGroup) {
      formValue.value = model.Group.createFrom(newGroup);
      return;
    }
    formValue.value = model.Group.createFrom(defaultGroup);
  },
  { immediate: true },
);

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
    const resp = props.isEdit ? await UpdateGroup(formValue.value) : await CreateGroup(formValue.value);

    if (!resp.ok) {
      message.error(resp.msg);
      return false;
    }

    message.success(props.isEdit ? t('message.updateSuccess') : t('message.createSuccess'));
    emit('update:show', false);
    emit('success');
  } catch (errors) {
    return false;
  }
};

const resetForm = () => {
  formValue.value = model.Group.createFrom(defaultGroup);
  emit('update:show', false);
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
