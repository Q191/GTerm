<template>
  <n-modal
    v-model:show="dialogStore.hostDialogVisible"
    :close-on-esc="true"
    :negative-text="$t('hostDialog.cancel')"
    :on-close="resetForm"
    :positive-text="$t('hostDialog.confirm')"
    :show-icon="false"
    :title="dialogStore.isEditMode ? $t('hostDialog.editTitle') : $t('hostDialog.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-tabs animated type="line" placement="left" v-model:value="activeTab">
      <n-tab-pane name="basic" :tab="$t('hostDialog.basicConfig')">
        <n-form ref="formRef" :model="formValue" :rules="rules" :show-label="false">
          <n-form-item path="name">
            <n-input v-model:value="formValue.name" clearable :placeholder="$t('hostDialog.placeholder.name')" />
          </n-form-item>

          <n-form-item path="groupID">
            <n-select
              v-model:value="formValue.groupID"
              :options="groupOptions"
              clearable
              tag
              :placeholder="$t('hostDialog.placeholder.group')"
            />
          </n-form-item>

          <div class="flex items-center w-full gap-2">
            <n-form-item path="host" class="flex-1">
              <n-input v-model:value="formValue.host" clearable :placeholder="$t('hostDialog.placeholder.host')" />
            </n-form-item>
            <n-form-item path="port" class="port-input">
              <n-input-number
                v-model:value="formValue.port"
                :min="1"
                :max="65535"
                :show-button="false"
                :placeholder="$t('hostDialog.placeholder.port')"
              />
            </n-form-item>
          </div>

          <n-form-item :label="$t('hostDialog.authType')" path="credential.authMethod">
            <div class="flex items-center justify-between w-full">
              <n-button-group>
                <n-button
                  :type="
                    !formValue.isCommonCredential && formValue.credential?.authMethod === AuthMethod.Password
                      ? 'primary'
                      : 'default'
                  "
                  @click="handleCredentialTypeChange(CredentialType.Password)"
                >
                  <template #icon>
                    <Icon icon="ph:password" />
                  </template>
                  {{ $t('hostDialog.password') }}
                </n-button>
                <n-button
                  :type="
                    !formValue.isCommonCredential && formValue.credential?.authMethod === AuthMethod.PrivateKey
                      ? 'primary'
                      : 'default'
                  "
                  @click="handleCredentialTypeChange(CredentialType.PrivateKey)"
                >
                  <template #icon>
                    <Icon icon="ph:key" />
                  </template>
                  {{ $t('hostDialog.privateKey') }}
                </n-button>
                <n-button
                  :type="formValue.isCommonCredential ? 'primary' : 'default'"
                  @click="handleCredentialTypeChange(CredentialType.Common)"
                >
                  <template #icon>
                    <Icon icon="ph:vault" />
                  </template>
                  {{ $t('hostDialog.commonCredentialLib') }}
                </n-button>
              </n-button-group>
              <n-tooltip v-if="!formValue.isCommonCredential && formValue.credential" trigger="hover" placement="right">
                <template #trigger>
                  <n-switch v-model:value="formValue.credential.isCommonCredential">
                    <template #checked>{{ $t('hostDialog.commonCredential') }}</template>
                    <template #unchecked>{{ $t('hostDialog.privateCredential') }}</template>
                  </n-switch>
                </template>
                <span class="tooltip-text">{{ $t('hostDialog.credentialTooltip') }}</span>
              </n-tooltip>
            </div>
          </n-form-item>

          <template v-if="formValue.isCommonCredential">
            <n-form-item path="credentialID">
              <n-select
                v-model:value="formValue.credentialID"
                :options="credentialOptions"
                clearable
                :placeholder="$t('hostDialog.placeholder.selectCredential')"
                @update:value="handleSelectCredential"
              />
            </n-form-item>
          </template>

          <template v-else>
            <n-form-item path="credential.username">
              <n-input
                v-model:value="formValue.credential!.username"
                clearable
                :placeholder="$t('hostDialog.placeholder.username')"
              />
            </n-form-item>

            <template v-if="formValue.credential!.authMethod === AuthMethod.Password">
              <n-form-item path="credential.password">
                <n-input
                  v-model:value="formValue.credential!.password"
                  type="password"
                  show-password-on="click"
                  clearable
                  :placeholder="$t('hostDialog.placeholder.password')"
                />
              </n-form-item>
            </template>

            <template v-if="formValue.credential!.authMethod === AuthMethod.PrivateKey">
              <n-form-item path="credential.privateKey">
                <n-input
                  v-model:value="formValue.credential!.privateKey"
                  type="textarea"
                  :rows="3"
                  clearable
                  :placeholder="$t('hostDialog.placeholder.privateKey')"
                />
              </n-form-item>
              <n-form-item path="credential.passphrase">
                <n-input
                  v-model:value="formValue.credential!.passphrase"
                  type="password"
                  show-password-on="click"
                  clearable
                  :placeholder="$t('hostDialog.placeholder.passphrase')"
                />
              </n-form-item>
            </template>
          </template>
        </n-form>
      </n-tab-pane>

      <n-tab-pane name="advanced" :tab="$t('hostDialog.advancedConfig')">
        <n-empty size="small" :description="$t('hostDialog.developing')">
          <template #icon>
            <Icon icon="ph:code" />
          </template>
        </n-empty>
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { model } from '@wailsApp/go/models';
import { ListGroup } from '@wailsApp/go/services/GroupSrv';
import { CreateHost, UpdateHost } from '@wailsApp/go/services/HostSrv';
import { Icon } from '@iconify/vue';

import {
  FormInst,
  FormRules,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NButton,
  NButtonGroup,
  NSelect,
  NTabPane,
  NTabs,
  useMessage,
  NSwitch,
  NTooltip,
  NEmpty,
} from 'naive-ui';
import { SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const dialogStore = useDialogStore();
const formRef = ref<FormInst | null>(null);
const activeTab = ref('basic');
const message = useMessage();

const AuthMethod = {
  Password: 0,
  PrivateKey: 1,
} as const;

const CredentialType = {
  Password: 0,
  PrivateKey: 1,
  Common: 2,
} as const;

const defaultCredential = {
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  isCommonCredential: false,
  authMethod: AuthMethod.Password,
};

const defaultHost = {
  name: '',
  host: '',
  port: 22,
  connProtocol: 0,
  credentialID: null,
  isCommonCredential: false,
  groupID: null,
};

const createCredential = (authMethod: number) =>
  model.Credential.createFrom({
    ...defaultCredential,
    authMethod,
  });

const createHost = (isCommon = false) => {
  const host = model.Host.createFrom({
    ...defaultHost,
    isCommonCredential: isCommon,
  });
  if (!isCommon) {
    host.credential = createCredential(AuthMethod.Password);
  }
  return host;
};

const formValue = ref(createHost());

const tempCachedCredentials = ref({
  password: createCredential(AuthMethod.Password),
  privateKey: createCredential(AuthMethod.PrivateKey),
});

const useCommonCredential = ref(false);
const credentialOptions = ref<SelectMixedOption[]>([]);

watch(
  () => dialogStore.editHost,
  newHost => {
    if (newHost) {
      formValue.value = model.Host.createFrom(newHost);
      if (newHost.credential && !newHost.isCommonCredential) {
        if (newHost.credential.authMethod === AuthMethod.Password) {
          tempCachedCredentials.value.password = model.Credential.createFrom(newHost.credential);
        } else if (newHost.credential.authMethod === AuthMethod.PrivateKey) {
          tempCachedCredentials.value.privateKey = model.Credential.createFrom(newHost.credential);
        }
      }
      useCommonCredential.value = newHost.isCommonCredential;
      return;
    }
    formValue.value = createHost();
    tempCachedCredentials.value = {
      password: createCredential(AuthMethod.Password),
      privateKey: createCredential(AuthMethod.PrivateKey),
    };
  },
  { immediate: true },
);

const handleCredentialTypeChange = (credentialType: number) => {
  if (credentialType === CredentialType.Common) {
    // 切换到凭据库模式
    formValue.value.isCommonCredential = true;
    formValue.value.credentialID = undefined;
    return;
  }
  // 切换到密码或私钥认证
  formValue.value.isCommonCredential = false;
  formValue.value.credentialID = undefined;
  formValue.value.credential =
    credentialType === CredentialType.Password
      ? tempCachedCredentials.value.password
      : tempCachedCredentials.value.privateKey;
  formValue.value.credential.authMethod =
    credentialType === CredentialType.Password ? AuthMethod.Password : AuthMethod.PrivateKey;
};

const handleSelectCredential = async (id: number) => {
  if (!id) {
    formValue.value.credentialID = undefined;
    return;
  }
  formValue.value.credentialID = id;
  formValue.value.isCommonCredential = true;
};

const rules = computed<FormRules>(() => ({
  name: {
    required: true,
    message: t('hostDialog.validation.nameRequired'),
    trigger: 'blur',
  },
  host: {
    required: true,
    message: t('hostDialog.validation.hostRequired'),
    trigger: 'blur',
  },
  port: {
    required: true,
    type: 'number',
    message: t('hostDialog.validation.portRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (typeof value !== 'number' || value < 1 || value > 65535) {
        return new Error(t('hostDialog.validation.portInvalid'));
      }
    },
  },
  'credential.username': {
    required: !formValue.value.isCommonCredential,
    message: t('hostDialog.validation.usernameRequired'),
    trigger: 'blur',
  },
  'credential.password': {
    required: !formValue.value.isCommonCredential && formValue.value.credential?.authMethod === AuthMethod.Password,
    message: t('hostDialog.validation.passwordRequired'),
    trigger: 'blur',
  },
  'credential.privateKey': {
    required: !formValue.value.isCommonCredential && formValue.value.credential?.authMethod === AuthMethod.PrivateKey,
    message: t('hostDialog.validation.privateKeyRequired'),
    trigger: 'blur',
  },
  credentialID: {
    required: formValue.value.isCommonCredential,
    message: t('hostDialog.validation.credentialRequired'),
    trigger: ['blur', 'change'],
  },
}));

const groupOptions = ref<SelectMixedOption[]>([]);

const fetchGroups = async () => {
  const resp = await ListGroup();
  if (!resp.ok) {
    message.error(resp.msg);
    return [];
  }
  return resp.data;
};

const fetchCredentials = async () => {
  // const resp = await ListCredentials();
  // if (!resp.ok) {
  // message.error(resp.msg);
  return [];
  // }
  // return resp.data || [];
};

const initOptions = async () => {
  const [groups, credentials] = await Promise.all([fetchGroups(), fetchCredentials()]);
  groupOptions.value = groups.map((group: model.Group) => ({
    label: group.name,
    value: group.id,
  }));
  credentialOptions.value = credentials.map((credential: model.Credential) => ({
    label: credential.name || credential.username,
    value: credential.id,
  }));
};

onMounted(initOptions);

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();
    const resp = dialogStore.isEditMode ? await UpdateHost(formValue.value) : await CreateHost(formValue.value);

    if (!resp.ok) {
      message.error(resp.msg);
      return false;
    }

    message.success(dialogStore.isEditMode ? '更新成功' : '创建成功');
    dialogStore.closeAddHostDialog();
  } catch (errors) {
    return false;
  }
};

const resetForm = () => {
  formValue.value = createHost();
  tempCachedCredentials.value = {
    password: createCredential(AuthMethod.Password),
    privateKey: createCredential(AuthMethod.PrivateKey),
  };
  activeTab.value = 'basic';
  dialogStore.closeAddHostDialog();
};
</script>

<style lang="less" scoped>
:deep(.n-form-item .n-form-item-label) {
  font-size: 13px;
}

:deep(.n-input) {
  font-size: 13px;
}

.tooltip-text {
  white-space: pre-line;
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
