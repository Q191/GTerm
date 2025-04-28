<template>
  <NModal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('frontend.credentialModal.cancel')"
    :on-close="resetForm"
    :positive-text="$t('frontend.credentialModal.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('frontend.credentialModal.editTitle') : $t('frontend.credentialModal.title')"
    preset="dialog"
    transform-origin="center"
    style="width: 600px"
    @positive-click="handleConfirm"
  >
    <NScrollbar style="max-height: 70vh; padding-right: 12px">
      <NForm ref="formRef" :model="formValue" :rules="rules">
        <NFormItem path="label" :label="$t('frontend.credentialModal.label')">
          <NInput
            v-model:value="formValue.label"
            clearable
            :placeholder="$t('frontend.credentialModal.placeholder.label')"
            :allow-input="value => !/\s/.test(value)"
          />
        </NFormItem>

        <NFormItem :label="$t('frontend.credentialModal.authType')" path="authMethod">
          <div class="auth-type-container">
            <NButtonGroup>
              <NButton
                :type="formValue.authMethod === AuthMethod.PASSWORD ? 'primary' : 'default'"
                @click="handleAuthTypeChange(AuthMethod.PASSWORD)"
              >
                <template #icon>
                  <Icon icon="ph:password" />
                </template>
                {{ $t('frontend.credentialModal.password') }}
              </NButton>
              <NButton
                :type="formValue.authMethod === AuthMethod.PRIVATEKEY ? 'primary' : 'default'"
                @click="handleAuthTypeChange(AuthMethod.PRIVATEKEY)"
              >
                <template #icon>
                  <Icon icon="ph:key" />
                </template>
                {{ $t('frontend.credentialModal.privateKey') }}
              </NButton>
            </NButtonGroup>
          </div>
        </NFormItem>

        <NFormItem path="username" :label="$t('frontend.credentialModal.username')">
          <NInput
            v-model:value="formValue.username"
            clearable
            :placeholder="$t('frontend.credentialModal.placeholder.username')"
            :allow-input="value => !/\s/.test(value)"
          />
        </NFormItem>

        <template v-if="formValue.authMethod === AuthMethod.PASSWORD">
          <NFormItem path="password" :label="$t('frontend.credentialModal.password')">
            <NInput
              v-model:value="formValue.password"
              type="password"
              show-password-on="click"
              clearable
              :placeholder="$t('frontend.credentialModal.placeholder.password')"
              :allow-input="value => !/\s/.test(value)"
            />
          </NFormItem>
        </template>

        <template v-if="formValue.authMethod === AuthMethod.PRIVATEKEY">
          <NFormItem path="privateKey" :label="$t('frontend.credentialModal.privateKey')">
            <NInput
              v-model:value="formValue.privateKey"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 3 }"
              clearable
              :placeholder="$t('frontend.credentialModal.placeholder.privateKey')"
            />
          </NFormItem>
          <NFormItem path="passphrase" :label="$t('frontend.credentialModal.passphrase')">
            <NInput
              v-model:value="formValue.passphrase"
              type="password"
              show-password-on="click"
              clearable
              :placeholder="$t('frontend.credentialModal.placeholder.passphrase')"
              :allow-input="value => !/\s/.test(value)"
            />
          </NFormItem>
        </template>
      </NForm>
    </NScrollbar>
  </NModal>
</template>

<script lang="ts" setup>
import { Icon } from '@iconify/vue';
import type { model } from '@wailsApp/go/models';
import { enums } from '@wailsApp/go/models';
import { CreateCredential, FindCredentialByID, UpdateCredential } from '@wailsApp/go/services/CredentialSrv';
import type { FormInst, FormRules } from 'naive-ui';
import { NButton, NButtonGroup, NForm, NFormItem, NInput, NModal, NScrollbar } from 'naive-ui';
import { computed, onMounted, onUpdated, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useCall } from '@/utils/call';

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  credentialId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();

const { AuthMethod } = enums;

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const { call } = useCall();

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
    message: t('frontend.credentialModal.validation.labelRequired'),
    trigger: 'blur',
  },
  username: {
    required: true,
    message: t('frontend.credentialModal.validation.usernameRequired'),
    trigger: 'blur',
  },
  password: {
    required: formValue.value.authMethod === AuthMethod.PASSWORD,
    message: t('frontend.credentialModal.validation.passwordRequired'),
    trigger: 'blur',
  },
  privateKey: {
    required: formValue.value.authMethod === AuthMethod.PRIVATEKEY,
    message: t('frontend.credentialModal.validation.privateKeyRequired'),
    trigger: 'blur',
  },
}));

const initModalData = async () => {
  if (props.credentialId && props.credentialId > 0 && props.isEdit) {
    const result = await call(FindCredentialByID, {
      args: [props.credentialId],
    });

    if (result.ok) {
      formValue.value = { ...result.data } as model.Credential;
    } else {
      emit('update:show', false);
    }
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

    const backendFunc = props.isEdit ? UpdateCredential : CreateCredential;
    const result = await call(backendFunc, {
      args: [formValue.value],
    });

    if (result.ok) {
      emit('update:show', false);
      emit('success');
    }

    return result.ok;
  } catch {
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
