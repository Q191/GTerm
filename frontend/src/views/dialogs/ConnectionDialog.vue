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
              :placeholder="$t('connDialog.placeholder.group')"
            />
          </n-form-item>

          <n-form-item path="connProtocol">
            <n-select
              v-model:value="formValue.connProtocol"
              :options="connProtocolOptions"
              :placeholder="$t('connDialog.placeholder.connProtocol')"
            />
          </n-form-item>

          <template v-if="formValue.connProtocol === enums.ConnProtocol.SERIAL">
            <n-form-item path="serialPort">
              <n-select
                v-model:value="formValue.serialPort"
                :options="serialPortsOptions"
                :placeholder="$t('connDialog.placeholder.serialPort')"
              />
            </n-form-item>

            <n-form-item path="baudRate">
              <n-select
                v-model:value="formValue.baudRate"
                :options="baudRateOptions"
                :placeholder="$t('connDialog.placeholder.baudRate')"
              />
            </n-form-item>

            <div class="form-row">
              <n-form-item path="dataBits" class="form-item">
                <n-select
                  v-model:value="formValue.dataBits"
                  :options="dataBitsOptions"
                  :placeholder="$t('connDialog.placeholder.dataBits')"
                />
              </n-form-item>

              <n-form-item path="stopBits" class="form-item">
                <n-select
                  v-model:value="formValue.stopBits"
                  :options="stopBitsOptions"
                  :placeholder="$t('connDialog.placeholder.stopBits')"
                />
              </n-form-item>

              <n-form-item path="parity" class="form-item">
                <n-select
                  v-model:value="formValue.parity"
                  :options="parityOptions"
                  :placeholder="$t('connDialog.placeholder.parity')"
                />
              </n-form-item>
            </div>
          </template>

          <template v-if="formValue.connProtocol === enums.ConnProtocol.SSH">
            <div class="form-row">
              <n-form-item path="host" class="form-item">
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
              <div class="auth-type-container">
                <n-button-group>
                  <n-button
                    :type="
                      !formValue.useCommonCredential && formValue.credential?.authMethod === enums.AuthMethod.PASSWORD
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
                      !formValue.useCommonCredential && formValue.credential?.authMethod === enums.AuthMethod.PRIVATEKEY
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
                    :type="formValue.useCommonCredential ? 'primary' : 'default'"
                    @click="handleCredentialTypeChange(CredentialType.Common)"
                  >
                    <template #icon>
                      <Icon icon="ph:vault" />
                    </template>
                    {{ $t('connDialog.commonCredentialLib') }}
                  </n-button>
                </n-button-group>
              </div>
            </n-form-item>

            <template v-if="formValue.useCommonCredential">
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

              <template v-if="formValue.credential!.authMethod === enums.AuthMethod.PASSWORD">
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

              <template v-if="formValue.credential!.authMethod === enums.AuthMethod.PRIVATEKEY">
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
import { enums, model } from '@wailsApp/go/models';
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

const CredentialType = {
  Password: 0,
  PrivateKey: 1,
  Common: 2,
} as const;

const connProtocolOptions = [
  { label: enums.ConnProtocol.SSH, value: enums.ConnProtocol.SSH },
  { label: enums.ConnProtocol.SERIAL, value: enums.ConnProtocol.SERIAL },
];

const defaultCredential = {
  label: '',
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  isCommonCredential: false,
  authMethod: enums.AuthMethod.PASSWORD,
};

const defaultConnection = {
  label: '',
  host: '',
  port: 22,
  serialPort: null,
  connProtocol: null,
  credentialID: null,
  useCommonCredential: false,
  groupID: null,
  baudRate: 9600,
  dataBits: 8,
  stopBits: 0, // OneStopBit
  parity: 0, // NoParity
};

const createCredential = (authMethod: enums.AuthMethod) =>
  model.Credential.createFrom({
    ...defaultCredential,
    authMethod,
  });

const createConnection = (isCommon = false) => {
  const conn = model.Connection.createFrom({
    ...defaultConnection,
    useCommonCredential: isCommon,
  });
  if (!isCommon) {
    conn.credential = createCredential(enums.AuthMethod.PASSWORD);
  }
  return conn;
};

const formValue = ref(createConnection());

const tempCachedCredentials = ref({
  password: createCredential(enums.AuthMethod.PASSWORD),
  privateKey: createCredential(enums.AuthMethod.PRIVATEKEY),
});

const useCommonCredential = ref(false);
const credentialOptions = ref<SelectMixedOption[]>([]);

watch(
  () => dialogStore.editConnection,
  newConnection => {
    if (newConnection) {
      formValue.value = model.Connection.createFrom(newConnection);
      if (newConnection.credential && !newConnection.useCommonCredential) {
        if (newConnection.credential.authMethod === enums.AuthMethod.PASSWORD) {
          tempCachedCredentials.value.password = model.Credential.createFrom(newConnection.credential);
        } else if (newConnection.credential.authMethod === enums.AuthMethod.PRIVATEKEY) {
          tempCachedCredentials.value.privateKey = model.Credential.createFrom(newConnection.credential);
        }
      }
      useCommonCredential.value = newConnection.useCommonCredential;
      return;
    }
    formValue.value = createConnection();
    tempCachedCredentials.value = {
      password: createCredential(enums.AuthMethod.PASSWORD),
      privateKey: createCredential(enums.AuthMethod.PRIVATEKEY),
    };
  },
  { immediate: true },
);

const handleCredentialTypeChange = (credentialType: number) => {
  if (credentialType === CredentialType.Common) {
    // 切换到凭据库模式
    formValue.value.useCommonCredential = true;
    formValue.value.credentialID = undefined;
    return;
  }
  // 切换到密码或私钥认证
  formValue.value.useCommonCredential = false;
  formValue.value.credentialID = undefined;
  formValue.value.credential =
    credentialType === CredentialType.Password
      ? tempCachedCredentials.value.password
      : tempCachedCredentials.value.privateKey;
  formValue.value.credential.authMethod =
    credentialType === CredentialType.Password ? enums.AuthMethod.PASSWORD : enums.AuthMethod.PRIVATEKEY;
};

const handleSelectCredential = async (id: number) => {
  if (!id) {
    formValue.value.credentialID = undefined;
    return;
  }
  formValue.value.credentialID = id;
  formValue.value.useCommonCredential = true;
};

const baudRateOptions = [
  { label: '9600', value: 9600 },
  { label: '19200', value: 19200 },
  { label: '38400', value: 38400 },
  { label: '57600', value: 57600 },
  { label: '115200', value: 115200 },
];

const parityOptions = [
  { label: 'None', value: 0 }, // NoParity
  { label: 'Odd', value: 1 }, // OddParity
  { label: 'Even', value: 2 }, // EvenParity
  { label: 'Mark', value: 3 }, // MarkParity
  { label: 'Space', value: 4 }, // SpaceParity
];

const dataBitsOptions = [
  { label: '5', value: 5 },
  { label: '6', value: 6 },
  { label: '7', value: 7 },
  { label: '8', value: 8 },
];

const stopBitsOptions = [
  { label: '1', value: 0 }, // OneStopBit
  { label: '1.5', value: 1 }, // OnePointFiveStopBits
  { label: '2', value: 2 }, // TwoStopBits
];

const rules = computed<FormRules>(() => ({
  label: {
    required: true,
    message: t('connDialog.validation.labelRequired'),
    trigger: 'blur',
  },
  host: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SSH,
    message: t('connDialog.validation.hostRequired'),
    trigger: 'blur',
  },
  port: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SSH,
    type: 'number',
    message: t('connDialog.validation.portRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (typeof value !== 'number' || value < 1 || value > 65535) {
        return new Error(t('connDialog.validation.portInvalid'));
      }
    },
  },
  serialPort: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SERIAL,
    message: t('connDialog.validation.serialPortRequired'),
    trigger: ['blur', 'change'],
  },
  baudRate: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SERIAL,
    type: 'number',
    message: t('connDialog.validation.baudRateRequired'),
    trigger: ['blur', 'change'],
  },
  dataBits: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SERIAL,
    type: 'number',
    message: t('connDialog.validation.dataBitsRequired'),
    trigger: ['blur', 'change'],
  },
  stopBits: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SERIAL,
    type: 'number',
    message: t('connDialog.validation.stopBitsRequired'),
    trigger: ['blur', 'change'],
  },
  parity: {
    required: formValue.value.connProtocol === enums.ConnProtocol.SERIAL,
    type: 'number',
    message: t('connDialog.validation.parityRequired'),
    trigger: ['blur', 'change'],
  },
  'credential.username': {
    required: !formValue.value.useCommonCredential,
    message: t('connDialog.validation.usernameRequired'),
    trigger: 'blur',
  },
  'credential.password': {
    required:
      !formValue.value.useCommonCredential && formValue.value.credential?.authMethod === enums.AuthMethod.PASSWORD,
    message: t('connDialog.validation.passwordRequired'),
    trigger: 'blur',
  },
  'credential.privateKey': {
    required:
      !formValue.value.useCommonCredential && formValue.value.credential?.authMethod === enums.AuthMethod.PRIVATEKEY,
    message: t('connDialog.validation.privateKeyRequired'),
    trigger: 'blur',
  },
  credentialID: {
    required: formValue.value.useCommonCredential,
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
    label: credential.label,
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
    password: createCredential(enums.AuthMethod.PASSWORD),
    privateKey: createCredential(enums.AuthMethod.PRIVATEKEY),
  };
  activeTab.value = 'basic';
  dialogStore.closeConnDialog();
};
</script>

<style lang="less" scoped>
.form-row {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 8px;

  .form-item {
    flex: 1;
  }
}

.auth-type-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.port-input {
  width: 120px;
}

.tooltip-text {
  font-size: 12px;
  line-height: 1.5;
}
</style>
