<template>
  <NModal
    v-model:show="dialogStore.hostDialogVisible"
    :close-on-esc="true"
    :negative-text="$t('host_dialog.cancel')"
    :on-close="dialogStore.closeAddHostDialog"
    :positive-text="$t('host_dialog.confirm')"
    :show-icon="false"
    :title="dialogStore.isEditMode ? $t('host_dialog.edit_title') : $t('host_dialog.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <NTabs animated type="line" v-model:value="activeTab">
      <NTabPane name="basic" :tab="$t('host_dialog.basic_config')">
        <NForm ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="80">
          <NFormItem :label="$t('host_dialog.label')" path="label">
            <NInput v-model:value="formValue.label" clearable />
          </NFormItem>

          <NFormItem :label="$t('host_dialog.group')" path="group">
            <NSelect v-model:value="formValue.group" :options="groupOptions" clearable tag />
          </NFormItem>

          <div class="flex items-center w-full gap-2">
            <NFormItem :label="$t('host_dialog.host')" path="host" class="flex-1">
              <NInput v-model:value="formValue.host" clearable />
            </NFormItem>
            <NFormItem path="port" class="port-input">
              <NInputNumber v-model:value="formValue.port" :min="1" :max="65535" />
            </NFormItem>
          </div>

          <NFormItem :label="$t('host_dialog.auth_type')" path="authType">
            <NRadioGroup v-model:value="formValue.authType">
              <NSpace>
                <NRadio value="password">{{ $t('host_dialog.password') }}</NRadio>
                <NRadio value="privateKey">{{ $t('host_dialog.private_key') }}</NRadio>
              </NSpace>
            </NRadioGroup>
          </NFormItem>

          <NFormItem :label="$t('host_dialog.username')" path="username">
            <NInput v-model:value="formValue.username" clearable />
          </NFormItem>

          <template v-if="formValue.authType === 'password'">
            <NFormItem :label="$t('host_dialog.password')" path="password">
              <NInput v-model:value="formValue.password" type="password" show-password-on="click" clearable />
            </NFormItem>
          </template>

          <template v-else>
            <NFormItem :label="$t('host_dialog.private_key')" path="privateKey">
              <NInput v-model:value="formValue.privateKey" type="textarea" :rows="3" clearable />
            </NFormItem>
            <NFormItem :label="$t('host_dialog.passphrase')" path="passphrase">
              <NInput v-model:value="formValue.passphrase" type="password" show-password-on="click" clearable />
            </NFormItem>
          </template>
        </NForm>
      </NTabPane>

      <NTabPane name="advanced" :tab="$t('host_dialog.advanced_config')">
        <NForm :model="formValue" label-placement="left" label-width="100">
          <NFormItem :label="$t('host_dialog.charset')" path="charset">
            <NSelect v-model:value="formValue.charset" :options="charsetOptions" />
          </NFormItem>
          <NFormItem :label="$t('host_dialog.term_type')" path="termType">
            <NSelect v-model:value="formValue.termType" :options="termTypeOptions" />
          </NFormItem>
        </NForm>
      </NTabPane>
    </NTabs>
  </NModal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import {
  FormInst,
  FormRules,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NRadio,
  NRadioGroup,
  NSelect,
  NSpace,
  NTabPane,
  NTabs,
} from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const dialogStore = useDialogStore();
const formRef = ref<FormInst | null>(null);
const activeTab = ref('basic');

interface FormState {
  label: string;
  group: string;
  host: string;
  port: number;
  authType: 'password' | 'privateKey';
  username: string;
  password: string;
  privateKey: string;
  passphrase: string;
  charset: string;
  termType: string;
}

const formValue = ref<FormState>({
  label: '',
  group: '',
  host: '',
  port: 22,
  authType: 'password',
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  charset: 'UTF-8',
  termType: 'xterm-256color',
});

const rules: FormRules = {
  label: {
    required: true,
    message: t('host_dialog.validation.label_required'),
    trigger: 'blur',
  },
  host: {
    required: true,
    message: t('host_dialog.validation.host_required'),
    trigger: 'blur',
  },
  port: {
    required: true,
    message: t('host_dialog.validation.port_required'),
    trigger: 'blur',
  },
  username: {
    required: true,
    message: t('host_dialog.validation.username_required'),
    trigger: 'blur',
  },
};

const groupOptions = [
  { label: '开发环境', value: 'dev' },
  { label: '测试环境', value: 'test' },
  { label: '生产环境', value: 'prod' },
];

const charsetOptions = [
  { label: 'UTF-8', value: 'UTF-8' },
  { label: 'GBK', value: 'GBK' },
];

const termTypeOptions = [
  { label: 'xterm', value: 'xterm' },
  { label: 'xterm-256color', value: 'xterm-256color' },
];

const handleConfirm = () => {
  formRef.value?.validate(errors => {
    if (!errors) {
      console.log('验证通过');
      // TODO: 处理表单提交
      dialogStore.closeAddHostDialog();
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

.port-input {
  width: 140px;
  flex-shrink: 0;

  :deep(.n-input-number) {
    width: 100%;
  }

  :deep(.n-form-item-label) {
    width: 40px !important;
  }
}
</style>
