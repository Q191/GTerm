<template>
  <n-modal
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
    <n-tabs animated type="line" placement="left" v-model:value="activeTab">
      <n-tab-pane name="basic" :tab="$t('host_dialog.basic_config')">
        <n-form ref="formRef" :model="formValue" :rules="rules" :show-label="false">
          <n-form-item spath="name">
            <n-input v-model:value="formValue.name" clearable :placeholder="$t('host_dialog.placeholder.name')" />
          </n-form-item>

          <n-form-item path="group_id">
            <n-select
              v-model:value="formValue.group_id"
              :options="groupOptions"
              clearable
              tag
              :placeholder="$t('host_dialog.placeholder.group')"
            />
          </n-form-item>

          <div class="flex items-center w-full gap-2">
            <n-form-item path="host" class="flex-1">
              <n-input v-model:value="formValue.host" clearable :placeholder="$t('host_dialog.placeholder.host')" />
            </n-form-item>
            <n-form-item path="port" class="port-input">
              <n-input-number
                v-model:value="formValue.port"
                :min="1"
                :max="65535"
                :show-button="false"
                :placeholder="$t('host_dialog.placeholder.port')"
              />
            </n-form-item>
          </div>

          <n-form-item :label="$t('host_dialog.auth_type')" path="credential.auth_type">
            <div class="flex items-center justify-between w-full">
              <n-button-group>
                <n-button
                  :type="formValue.credential!.auth_type === CredentialType.Password ? 'primary' : 'default'"
                  @click="handleAuthTypeChange(CredentialType.Password)"
                >
                  <template #icon>
                    <Icon icon="ph:password" />
                  </template>
                  {{ $t('host_dialog.password') }}
                </n-button>
                <n-button
                  :type="formValue.credential!.auth_type === CredentialType.PrivateKey ? 'primary' : 'default'"
                  @click="handleAuthTypeChange(CredentialType.PrivateKey)"
                >
                  <template #icon>
                    <Icon icon="ph:key" />
                  </template>
                  {{ $t('host_dialog.private_key') }}
                </n-button>
                <n-button
                  :type="useCommonCredential ? 'primary' : 'default'"
                  @click="handleAuthTypeChange(CredentialType.Common)"
                >
                  <template #icon>
                    <Icon icon="ph:vault" />
                  </template>
                  {{ $t('host_dialog.common_credential_lib') }}
                </n-button>
              </n-button-group>
              <n-tooltip v-if="!useCommonCredential" trigger="hover" placement="right">
                <template #trigger>
                  <n-switch v-model:value="formValue.credential!.is_common_credential">
                    <template #checked>{{ $t('host_dialog.common_credential') }}</template>
                    <template #unchecked>{{ $t('host_dialog.private_credential') }}</template>
                  </n-switch>
                </template>
                <span class="tooltip-text">{{ $t('host_dialog.credential_tooltip') }}</span>
              </n-tooltip>
            </div>
          </n-form-item>

          <template v-if="useCommonCredential">
            <n-form-item path="credential_id">
              <n-select
                v-model:value="formValue.credential_id"
                :options="credentialOptions"
                clearable
                :placeholder="$t('host_dialog.placeholder.select_credential')"
                @update:value="handleCredentialChange"
              />
            </n-form-item>
          </template>

          <template v-else>
            <n-form-item path="credential.username">
              <n-input
                v-model:value="formValue.credential!.username"
                clearable
                :placeholder="$t('host_dialog.placeholder.username')"
              />
            </n-form-item>

            <template v-if="formValue.credential!.auth_type === CredentialType.Password">
              <n-form-item path="credential.password">
                <n-input
                  v-model:value="formValue.credential!.password"
                  type="password"
                  show-password-on="click"
                  clearable
                  :placeholder="$t('host_dialog.placeholder.password')"
                />
              </n-form-item>
            </template>

            <template v-if="formValue.credential!.auth_type === CredentialType.PrivateKey">
              <n-form-item path="credential.private_key">
                <n-input
                  v-model:value="formValue.credential!.private_key"
                  type="textarea"
                  :rows="3"
                  clearable
                  :placeholder="$t('host_dialog.placeholder.private_key')"
                />
              </n-form-item>
              <n-form-item path="credential.key_password">
                <n-input
                  v-model:value="formValue.credential!.key_password"
                  type="password"
                  show-password-on="click"
                  clearable
                  :placeholder="$t('host_dialog.placeholder.key_password')"
                />
              </n-form-item>
            </template>
          </template>
        </n-form>
      </n-tab-pane>

      <n-tab-pane name="advanced" :tab="$t('host_dialog.advanced_config')">
        <n-empty size="small" :description="$t('host_dialog.developing')">
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

// 认证类型常量
const CredentialType = {
  Password: 0,
  PrivateKey: 1,
  Common: 2,
} as const;

// 默认凭据信息
const defaultCredential = {
  username: '',
  password: '',
  private_key: '',
  key_password: '',
  is_common_credential: false,
};

// 默认主机信息
const defaultHost = {
  name: '',
  host: '',
  port: 22,
  description: '',
  credential_id: 0,
  credential_type: 0,
  conn_protocol: 0,
};

// 创建新的凭据对象
const createCredential = (type: number) =>
  model.Credential.createFrom({
    ...defaultCredential,
    auth_type: type,
  });

// 创建新的主机对象
const createHost = (credentialType = 0) =>
  model.Host.createFrom({
    ...defaultHost,
    credential_type: credentialType,
    credential: createCredential(credentialType),
  });

const formValue = ref(createHost());

// 保存不同认证类型的凭据信息
const savedCredentials = ref({
  password: createCredential(0),
  privateKey: createCredential(1),
});

const useCommonCredential = ref(false);
const credentialOptions = ref<SelectMixedOption[]>([]);

watch(
  () => dialogStore.editHost,
  newHost => {
    if (newHost) {
      formValue.value = model.Host.createFrom(newHost);
      // 保存认证信息到对应类型
      if (newHost.credential) {
        if (newHost.credential_type === CredentialType.Password) {
          savedCredentials.value.password = model.Credential.createFrom(newHost.credential);
        } else if (newHost.credential_type === CredentialType.PrivateKey) {
          savedCredentials.value.privateKey = model.Credential.createFrom(newHost.credential);
        }
      }
      useCommonCredential.value = newHost.credential_type === 2;
    } else {
      formValue.value = createHost();
      savedCredentials.value = {
        password: createCredential(0),
        privateKey: createCredential(1),
      };
    }
  },
  { immediate: true },
);

const handleAuthTypeChange = (type: number) => {
  if (type === CredentialType.Common) {
    // 切换到凭据库之前，保存当前认证信息
    const currentType = formValue.value.credential_type;
    if (currentType === CredentialType.Password || currentType === CredentialType.PrivateKey) {
      const key = currentType === CredentialType.Password ? 'password' : 'privateKey';
      savedCredentials.value[key] = model.Credential.createFrom(formValue.value.credential!);
    }

    useCommonCredential.value = true;
    formValue.value.credential_type = CredentialType.Common; // 凭据库
    formValue.value.credential_id = undefined as any; // 清空凭据选择
    formValue.value.credential = createCredential(CredentialType.Common);
  } else {
    useCommonCredential.value = false;
    formValue.value.credential_type = type; // 密码或私钥
    formValue.value.credential_id = 0; // 清空凭据ID
    // 恢复之前保存的认证信息
    formValue.value.credential = model.Credential.createFrom(
      type === CredentialType.Password ? savedCredentials.value.password : savedCredentials.value.privateKey,
    );
  }
};

const handleCredentialChange = async (id: string) => {
  if (!id) return;
  formValue.value.credential_id = parseInt(id);
  formValue.value.credential_type = 2; // 设置为凭据库类型
  // const resp = await GetCredential(id);
  // if (!resp.ok) {
  // message.error(resp.msg);
  // return;
  // }
  // formValue.value.credential = resp.data;
};

const rules = computed<FormRules>(() => ({
  name: {
    required: true,
    message: t('host_dialog.validation.name_required'),
    trigger: 'blur',
  },
  host: {
    required: true,
    message: t('host_dialog.validation.host_required'),
    trigger: 'blur',
  },
  port: {
    required: true,
    type: 'number',
    message: t('host_dialog.validation.port_required'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (typeof value !== 'number' || value < 1 || value > 65535) {
        return new Error(t('host_dialog.validation.port_invalid'));
      }
    },
  },
  'credential.username': {
    required: !useCommonCredential.value,
    message: t('host_dialog.validation.username_required'),
    trigger: 'blur',
  },
  'credential.password': {
    required: !useCommonCredential.value && formValue.value.credential?.auth_type === 0,
    message: t('host_dialog.validation.password_required'),
    trigger: 'blur',
  },
  'credential.private_key': {
    required: !useCommonCredential.value && formValue.value.credential?.auth_type === 1,
    message: t('host_dialog.validation.private_key_required'),
    trigger: 'blur',
  },
  credential_id: {
    required: useCommonCredential.value,
    message: t('host_dialog.validation.credential_required'),
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
