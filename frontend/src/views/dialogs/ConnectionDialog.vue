<template>
  <n-modal
    v-model:show="dialogStore.connDialogVisible"
    :close-on-esc="true"
    :negative-text="$t('connDialog.cancel')"
    :on-close="resetForm"
    :positive-text="$t('connDialog.confirm')"
    :show-icon="false"
    :title="dialogStore.isEditMode ? $t('connDialog.editTitle') : $t('connDialog.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-tabs animated type="line" placement="left" v-model:value="activeTab">
      <n-tab-pane name="basic" :tab="$t('connDialog.basicConfig')">
        <n-form ref="formRef" :model="formValue" :rules="rules" :show-label="false">
          <n-form-item path="label">
            <n-input v-model:value="formValue.label" clearable :placeholder="$t('connDialog.placeholder.label')" />
          </n-form-item>

          <n-form-item path="groupID">
            <n-select
              v-model:value="formValue.groupID"
              :options="groupOptions"
              clearable
              tag
              :placeholder="$t('connDialog.placeholder.group')"
            />
          </n-form-item>

          <n-form-item path="serialPort">
            <n-select
              v-model:value="formValue.serialPort"
              :options="serialPortsOptions"
              clearable
              tag
              :placeholder="$t('connDialog.placeholder.serialPort')"
            />
          </n-form-item>

          <div class="flex items-center w-full gap-2">
            <n-form-item path="host" class="flex-1">
              <n-input v-model:value="formValue.host" clearable :placeholder="$t('connDialog.placeholder.host')" />
            </n-form-item>
            <n-form-item path="port" class="port-input">
              <n-input-number
                v-model:value="formValue.port"
                :min="1"
                :max="65535"
                :show-button="false"
                :placeholder="$t('connDialog.placeholder.port')"
              />
            </n-form-item>
          </div>

          <n-form-item :label="$t('connDialog.authType')" path="credential.authMethod">
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
                  {{ $t('connDialog.password') }}
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
                  {{ $t('connDialog.privateKey') }}
                </n-button>
                <n-button
                  :type="formValue.isCommonCredential ? 'primary' : 'default'"
                  @click="handleCredentialTypeChange(CredentialType.Common)"
                >
                  <template #icon>
                    <Icon icon="ph:vault" />
                  </template>
                  {{ $t('connDialog.commonCredentialLib') }}
                </n-button>
              </n-button-group>
              <n-tooltip v-if="!formValue.isCommonCredential && formValue.credential" trigger="hover" placement="right">
                <template #trigger>
                  <n-switch v-model:value="formValue.credential.isCommonCredential">
                    <template #checked>{{ $t('connDialog.commonCredential') }}</template>
                    <template #unchecked>{{ $t('connDialog.privateCredential') }}</template>
                  </n-switch>
                </template>
                <span class="tooltip-text">{{ $t('connDialog.credentialTooltip') }}</span>
              </n-tooltip>
            </div>
          </n-form-item>

          <template v-if="formValue.isCommonCredential">
            <n-form-item path="credentialID">
              <n-select
                v-model:value="formValue.credentialID"
                :options="credentialOptions"
                clearable
                :placeholder="$t('connDialog.placeholder.selectCredential')"
                @update:value="handleSelectCredential"
              />
            </n-form-item>
          </template>

          <template v-else>
            <n-form-item path="credential.username">
              <n-input
                v-model:value="formValue.credential!.username"
                clearable
                :placeholder="$t('connDialog.placeholder.username')"
              />
            </n-form-item>

            <template v-if="formValue.credential!.authMethod === AuthMethod.Password">
              <n-form-item path="credential.password">
                <n-input
                  v-model:value="formValue.credential!.password"
                  type="password"
                  show-password-on="click"
                  clearable
                  :placeholder="$t('connDialog.placeholder.password')"
                />
              </n-form-item>
            </template>

            <template v-if="formValue.credential!.authMethod === AuthMethod.PrivateKey">
              <n-form-item path="credential.privateKey">
                <n-input
                  v-model:value="formValue.credential!.privateKey"
                  type="textarea"
                  :row="3"
                  clearable
                  :placeholder="$t('connDialog.placeholder.privateKey')"
                />
              </n-form-item>
              <n-form-item path="credential.passphrase">
                <n-input
                  v-model:value="formValue.credential!.passphrase"
                  type="password"
                  show-password-on="click"
                  clearable
                  :placeholder="$t('connDialog.placeholder.passphrase')"
                />
              </n-form-item>
            </template>
          </template>
        </n-form>
      </n-tab-pane>

      <n-tab-pane name="advanced" :tab="$t('connDialog.advancedConfig')">
        <n-empty size="small" :description="$t('connDialog.developing')">
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
import { CreateConnection, UpdateConnection } from '@wailsApp/go/services/ConnectionSrv';
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
import { ref, computed, watch, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { SerialPorts } from '@wailsApp/go/services/TerminalSrv';

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

const defaultConnection = {
  label: '',
  host: '',
  port: 22,
  serialPort: null,
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

const createConnection = (isCommon = false) => {
  const conn = model.Connection.createFrom({
    ...defaultConnection,
    isCommonCredential: isCommon,
  });
  if (!isCommon) {
    conn.credential = createCredential(AuthMethod.Password);
  }
  return conn;
};

const formValue = ref(createConnection());

const tempCachedCredentials = ref({
  password: createCredential(AuthMethod.Password),
  privateKey: createCredential(AuthMethod.PrivateKey),
});

const useCommonCredential = ref(false);
const credentialOptions = ref<SelectMixedOption[]>([]);

watch(
  () => dialogStore.editConnection,
  newConnection => {
    if (newConnection) {
      formValue.value = model.Connection.createFrom(newConnection);
      if (newConnection.credential && !newConnection.isCommonCredential) {
        if (newConnection.credential.authMethod === AuthMethod.Password) {
          tempCachedCredentials.value.password = model.Credential.createFrom(newConnection.credential);
        } else if (newConnection.credential.authMethod === AuthMethod.PrivateKey) {
          tempCachedCredentials.value.privateKey = model.Credential.createFrom(newConnection.credential);
        }
      }
      useCommonCredential.value = newConnection.isCommonCredential;
      return;
    }
    formValue.value = createConnection();
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
  label: {
    required: true,
    message: t('connDialog.validation.labelRequired'),
    trigger: 'blur',
  },
  host: {
    required: true,
    message: t('connDialog.validation.hostRequired'),
    trigger: 'blur',
  },
  port: {
    required: true,
    type: 'number',
    message: t('connDialog.validation.portRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (typeof value !== 'number' || value < 1 || value > 65535) {
        return new Error(t('connDialog.validation.portInvalid'));
      }
    },
  },
  'credential.username': {
    required: !formValue.value.isCommonCredential,
    message: t('connDialog.validation.usernameRequired'),
    trigger: 'blur',
  },
  'credential.password': {
    required: !formValue.value.isCommonCredential && formValue.value.credential?.authMethod === AuthMethod.Password,
    message: t('connDialog.validation.passwordRequired'),
    trigger: 'blur',
  },
  'credential.privateKey': {
    required: !formValue.value.isCommonCredential && formValue.value.credential?.authMethod === AuthMethod.PrivateKey,
    message: t('connDialog.validation.privateKeyRequired'),
    trigger: 'blur',
  },
  credentialID: {
    required: formValue.value.isCommonCredential,
    message: t('connDialog.validation.credentialRequired'),
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

const serialPortsOptions = ref<SelectMixedOption[]>([]);

const fetchSerialPorts = async () => {
  const resp = await SerialPorts();
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
  const [groups, credentials, serialPorts] = await Promise.all([fetchGroups(), fetchCredentials(), fetchSerialPorts()]);
  groupOptions.value = groups.map((group: model.Group) => ({
    label: group.name,
    value: group.id,
  }));
  credentialOptions.value = credentials.map((credential: model.Credential) => ({
    label: credential.name,
    value: credential.id,
  }));
  serialPortsOptions.value = serialPorts.map((serialPort: string[]) => ({
    label: serialPort,
    value: serialPort,
  }));
};

onMounted(initOptions);

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();
    const resp = dialogStore.isEditMode
      ? await UpdateConnection(formValue.value)
      : await CreateConnection(formValue.value);

    if (!resp.ok) {
      message.error(resp.msg);
      return false;
    }

    message.success(dialogStore.isEditMode ? t('message.updateSuccess') : t('message.createSuccess'));
    dialogStore.closeConnDialog();
  } catch (errors) {
    return false;
  }
};

const resetForm = () => {
  formValue.value = createConnection();
  tempCachedCredentials.value = {
    password: createCredential(AuthMethod.Password),
    privateKey: createCredential(AuthMethod.PrivateKey),
  };
  activeTab.value = 'basic';
  dialogStore.closeConnDialog();
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
