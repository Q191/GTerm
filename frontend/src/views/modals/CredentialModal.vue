<template>
  <n-modal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('credentialModal.cancel')"
    :on-close="resetForm"
    :positive-text="$t('credentialModal.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('credentialModal.editTitle') : $t('credentialModal.title')"
    preset="dialog"
    transform-origin="center"
    style="width: 600px"
    @positive-click="handleConfirm"
  >
    <n-scrollbar style="max-height: 70vh; padding-right: 12px">
      <n-form ref="formRef" :model="formValue" :rules="rules">
        <n-form-item path="label" :label="$t('credentialModal.label')">
          <n-input
            v-model:value="formValue.label"
            clearable
            :placeholder="$t('credentialModal.placeholder.label')"
            :allow-input="value => !/\s/.test(value)"
          />
        </n-form-item>

        <n-form-item :label="$t('credentialModal.authType')" path="authMethod">
          <div class="auth-type-container">
            <n-button-group>
              <n-button
                :type="formValue.authMethod === AuthMethod.PASSWORD ? 'primary' : 'default'"
                @click="handleAuthTypeChange(AuthMethod.PASSWORD)"
              >
                <template #icon>
                  <Icon icon="ph:password" />
                </template>
                {{ $t('credentialModal.password') }}
              </n-button>
              <n-button
                :type="formValue.authMethod === AuthMethod.PRIVATEKEY ? 'primary' : 'default'"
                @click="handleAuthTypeChange(AuthMethod.PRIVATEKEY)"
              >
                <template #icon>
                  <Icon icon="ph:key" />
                </template>
                {{ $t('credentialModal.privateKey') }}
              </n-button>
            </n-button-group>
          </div>
        </n-form-item>

        <n-form-item path="username" :label="$t('credentialModal.username')">
          <n-input
            v-model:value="formValue.username"
            clearable
            :placeholder="$t('credentialModal.placeholder.username')"
            :allow-input="value => !/\s/.test(value)"
          />
        </n-form-item>

        <template v-if="formValue.authMethod === AuthMethod.PASSWORD">
          <n-form-item path="password" :label="$t('credentialModal.password')">
            <n-input
              v-model:value="formValue.password"
              type="password"
              show-password-on="click"
              clearable
              :placeholder="$t('credentialModal.placeholder.password')"
              :allow-input="value => !/\s/.test(value)"
            />
          </n-form-item>
        </template>

        <template v-if="formValue.authMethod === AuthMethod.PRIVATEKEY">
          <n-form-item path="privateKey" :label="$t('credentialModal.privateKey')">
            <n-input
              v-model:value="formValue.privateKey"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 3 }"
              clearable
              :placeholder="$t('credentialModal.placeholder.privateKey')"
            />
          </n-form-item>
          <n-form-item path="passphrase" :label="$t('credentialModal.passphrase')">
            <n-input
              v-model:value="formValue.passphrase"
              type="password"
              show-password-on="click"
              clearable
              :placeholder="$t('credentialModal.placeholder.passphrase')"
              :allow-input="value => !/\s/.test(value)"
            />
          </n-form-item>
        </template>
      </n-form>
    </n-scrollbar>
  </n-modal>
</template>

<script lang="ts" setup>
import { enums, model } from '@wailsApp/go/models';
import { CreateCredential, FindCredentialByID, UpdateCredential } from '@wailsApp/go/services/CredentialSrv';
import { Icon } from '@iconify/vue';

import {
  FormInst,
  FormRules,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NButton,
  NButtonGroup,
  useMessage,
  NScrollbar,
} from 'naive-ui';
import { ref, computed, onMounted, onUpdated } from 'vue';
import { useI18n } from 'vue-i18n';

const { AuthMethod } = enums;

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  credentialId?: number;
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

const defaultCredential: Partial<model.Credential> = {
  label: '',
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  authMethod: AuthMethod.PASSWORD,
};

function createCredentialObject(): model.Credential {
  return { ...defaultCredential } as model.Credential;
}

const formValue = ref<model.Credential>(createCredentialObject());

const handleAuthTypeChange = (authMethod: enums.AuthMethod) => {
  formValue.value.authMethod = authMethod;
  if (authMethod === AuthMethod.PASSWORD) {
    formValue.value.privateKey = '';
    formValue.value.passphrase = '';
  } else {
    formValue.value.password = '';
  }
};

const rules = computed<FormRules>(() => ({
  label: {
    required: true,
    message: t('credentialModal.validation.labelRequired'),
    trigger: 'blur',
  },
  username: {
    required: true,
    message: t('credentialModal.validation.usernameRequired'),
    trigger: 'blur',
  },
  password: {
    required: formValue.value.authMethod === AuthMethod.PASSWORD,
    message: t('credentialModal.validation.passwordRequired'),
    trigger: 'blur',
  },
  privateKey: {
    required: formValue.value.authMethod === AuthMethod.PRIVATEKEY,
    message: t('credentialModal.validation.privateKeyRequired'),
    trigger: 'blur',
  },
}));

const initModalData = async () => {
  if (props.credentialId && props.credentialId > 0 && props.isEdit) {
    const resp = await FindCredentialByID(props.credentialId);
    if (resp.ok) {
      formValue.value = { ...resp.data } as model.Credential;
      return;
    }
    message.error(resp.msg);
    emit('update:show', false);
  } else {
    formValue.value = createCredentialObject();
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
    const resp = props.isEdit ? await UpdateCredential(formValue.value) : await CreateCredential(formValue.value);

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
  formValue.value = createCredentialObject();
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
