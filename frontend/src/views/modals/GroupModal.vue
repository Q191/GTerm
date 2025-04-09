<template>
  <n-modal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('frontend.groupModal.cancel')"
    :on-close="resetForm"
    :positive-text="$t('frontend.groupModal.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('frontend.groupModal.editTitle') : $t('frontend.groupModal.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-form ref="formRef" :model="formValue" :rules="rules">
      <n-form-item path="name" :label="$t('frontend.groupModal.name')">
        <n-input
          v-model:value="formValue.name"
          clearable
          :placeholder="$t('frontend.groupModal.placeholder.name')"
          :allow-input="value => !/\s/.test(value)"
        />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { model } from '@wailsApp/go/models';
import { CreateGroup, UpdateGroup } from '@wailsApp/go/services/GroupSrv';
import { FormInst, FormRules, NForm, NFormItem, NInput, NModal } from 'naive-ui';
import { ref, computed, onMounted, onUpdated } from 'vue';
import { useI18n } from 'vue-i18n';
import { useCall } from '@/utils/call';

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const { call } = useCall();

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

const defaultGroup: Partial<model.Group> = {
  name: '',
};

function createGroupObject(): model.Group {
  return { ...defaultGroup } as model.Group;
}

const formValue = ref<model.Group>(createGroupObject());

const rules: FormRules = {
  name: {
    required: true,
    message: t('frontend.groupModal.validation.nameRequired'),
    trigger: 'blur',
  },
};

const initModalData = () => {
  if (props.group) {
    formValue.value = { ...props.group } as model.Group;
  } else {
    formValue.value = createGroupObject();
  }
};

onUpdated(() => {
  if (props.show) {
    initModalData();
  }
});

onMounted(() => {
  if (props.show) {
    initModalData();
  }
});

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();

    const backendFunc = props.isEdit ? UpdateGroup : CreateGroup;
    const result = await call(backendFunc, {
      args: [formValue.value],
    });

    if (result.ok) {
      emit('update:show', false);
      emit('success');
    }

    return result.ok;
  } catch (errors) {
    return false;
  }
};

const resetForm = () => {
  formValue.value = createGroupObject();
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
