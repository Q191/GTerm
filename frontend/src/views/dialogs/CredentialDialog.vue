<template>
  <n-modal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('credentialDialog.cancel')"
    :on-close="resetForm"
    :positive-text="$t('credentialDialog.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('credentialDialog.editTitle') : $t('credentialDialog.title')"
    preset="dialog"
    transform-origin="center"
    style="width: 600px"
    @positive-click="handleConfirm"
  >
    <n-scrollbar style="max-height: 70vh; padding-right: 12px">
      <n-form ref="formRef" :model="formValue" :rules="rules">
        <n-form-item path="label" :label="$t('credentialDialog.label')">
          <n-input v-model:value="formValue.label" clearable :placeholder="$t('credentialDialog.placeholder.label')" />
        </n-form-item>

        <n-form-item :label="$t('credentialDialog.authType')" path="authMethod">
          <div class="auth-type-container">
            <n-button-group>
              <n-button
                :type="formValue.authMethod === enums.AuthMethod.PASSWORD ? 'primary' : 'default'"
                @click="handleAuthTypeChange(enums.AuthMethod.PASSWORD)"
              >
                <template #icon>
                  <Icon icon="ph:password" />
                </template>
                {{ $t('credentialDialog.password') }}
              </n-button>
              <n-button
                :type="formValue.authMethod === enums.AuthMethod.PRIVATEKEY ? 'primary' : 'default'"
                @click="handleAuthTypeChange(enums.AuthMethod.PRIVATEKEY)"
              >
                <template #icon>
                  <Icon icon="ph:key" />
                </template>
                {{ $t('credentialDialog.privateKey') }}
              </n-button>
            </n-button-group>
          </div>
        </n-form-item>

        <n-form-item path="username" :label="$t('credentialDialog.username')">
          <n-input
            v-model:value="formValue.username"
            clearable
            :placeholder="$t('credentialDialog.placeholder.username')"
          />
        </n-form-item>

        <template v-if="formValue.authMethod === enums.AuthMethod.PASSWORD">
          <n-form-item path="password" :label="$t('credentialDialog.password')">
            <n-input
              v-model:value="formValue.password"
              type="password"
              show-password-on="click"
              clearable
              :placeholder="$t('credentialDialog.placeholder.password')"
            />
          </n-form-item>
        </template>

        <template v-if="formValue.authMethod === enums.AuthMethod.PRIVATEKEY">
          <n-form-item path="privateKey" :label="$t('credentialDialog.privateKey')">
            <n-input
              v-model:value="formValue.privateKey"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 3 }"
              clearable
              :placeholder="$t('credentialDialog.placeholder.privateKey')"
            />
          </n-form-item>
          <n-form-item path="passphrase" :label="$t('credentialDialog.passphrase')">
            <n-input
              v-model:value="formValue.passphrase"
              type="password"
              show-password-on="click"
              clearable
              :placeholder="$t('credentialDialog.placeholder.passphrase')"
            />
          </n-form-item>
        </template>
      </n-form>
    </n-scrollbar>
  </n-modal>
</template>

<script lang="ts" setup>
import { enums, model } from '@wailsApp/go/models';
import { CreateCredential, UpdateCredential } from '@wailsApp/go/services/CredentialSrv';
import { Icon } from '@iconify/vue';

import { FormInst, FormRules, NForm, NFormItem, NInput, NModal, NButton, NButtonGroup, useMessage } from 'naive-ui';
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  credential?: model.Credential;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const message = useMessage();

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value),
});

const defaultCredential = {
  label: '',
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  authMethod: enums.AuthMethod.PASSWORD,
};

const formValue = ref(model.Credential.createFrom(defaultCredential));

watch(
  () => props.credential,
  newCredential => {
    if (newCredential) {
      formValue.value = model.Credential.createFrom(newCredential);
      return;
    }
    formValue.value = model.Credential.createFrom(defaultCredential);
  },
  { immediate: true },
);

const handleAuthTypeChange = (authMethod: enums.AuthMethod) => {
  formValue.value.authMethod = authMethod;
  if (authMethod === enums.AuthMethod.PASSWORD) {
    formValue.value.privateKey = '';
    formValue.value.passphrase = '';
  } else {
    formValue.value.password = '';
  }
};

const rules = computed<FormRules>(() => ({
  label: {
    required: true,
    message: t('credentialDialog.validation.labelRequired'),
    trigger: 'blur',
  },
  username: {
    required: true,
    message: t('credentialDialog.validation.usernameRequired'),
    trigger: 'blur',
  },
  password: {
    required: formValue.value.authMethod === enums.AuthMethod.PASSWORD,
    message: t('credentialDialog.validation.passwordRequired'),
    trigger: 'blur',
  },
  privateKey: {
    required: formValue.value.authMethod === enums.AuthMethod.PRIVATEKEY,
    message: t('credentialDialog.validation.privateKeyRequired'),
    trigger: 'blur',
  },
}));

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();
    console.log('表单验证通过，准备提交数据');
    const resp = props.isEdit ? await UpdateCredential(formValue.value) : await CreateCredential(formValue.value);

    if (!resp.ok) {
      message.error(resp.msg);
      return false;
    }

    console.log('操作成功，准备关闭对话框并触发刷新');
    message.success(props.isEdit ? t('message.updateSuccess') : t('message.createSuccess'));
    emit('update:show', false);
    emit('success');
  } catch (errors) {
    console.log('操作失败:', errors);
    return false;
  }
};

const resetForm = () => {
  console.log('重置表单');
  formValue.value = model.Credential.createFrom(defaultCredential);
  emit('update:show', false);
};
</script>

<style lang="less" scoped>
.auth-type-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}
</style>
