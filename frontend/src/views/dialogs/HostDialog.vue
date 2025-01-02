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
          <NFormItem :label="$t('host_dialog.name')" path="name">
            <NInput v-model:value="formValue.name" clearable />
          </NFormItem>

          <NFormItem :label="$t('host_dialog.group')" path="group_id">
            <NSelect v-model:value="formValue.group_id" :options="groupOptions" clearable tag />
          </NFormItem>

          <div class="flex items-center w-full gap-2">
            <NFormItem :label="$t('host_dialog.host')" path="host" class="flex-1">
              <NInput v-model:value="formValue.host" clearable />
            </NFormItem>
            <NFormItem path="port" class="port-input">
              <NInputNumber v-model:value="formValue.port" :min="1" :max="65535" />
            </NFormItem>
          </div>

          <NFormItem :label="$t('host_dialog.auth_type')" path="credential.auth_type">
            <NRadioGroup v-model:value="formValue.credential!.auth_type">
              <NSpace>
                <NRadio :value="0">{{ $t('host_dialog.password') }}</NRadio>
                <NRadio :value="1">{{ $t('host_dialog.private_key') }}</NRadio>
              </NSpace>
            </NRadioGroup>
          </NFormItem>

          <NFormItem :label="$t('host_dialog.username')" path="credential.username">
            <NInput v-model:value="formValue.credential!.username" clearable />
          </NFormItem>

          <template v-if="formValue.credential!.auth_type === 0">
            <NFormItem :label="$t('host_dialog.password')" path="credential.password">
              <NInput
                v-model:value="formValue.credential!.password"
                type="password"
                show-password-on="click"
                clearable
              />
            </NFormItem>
          </template>

          <template v-else>
            <NFormItem :label="$t('host_dialog.private_key')" path="credential.private_key">
              <NInput v-model:value="formValue.credential!.private_key" type="textarea" :rows="3" clearable />
            </NFormItem>
            <NFormItem :label="$t('host_dialog.passphrase')" path="credential.key_password">
              <NInput
                v-model:value="formValue.credential!.key_password"
                type="password"
                show-password-on="click"
                clearable
              />
            </NFormItem>
          </template>

          <NFormItem :label="$t('host_dialog.description')" path="description">
            <NInput
              v-model:value="formValue.description"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              clearable
            />
          </NFormItem>
        </NForm>
      </NTabPane>

      <NTabPane name="advanced" :tab="$t('host_dialog.advanced_config')"> </NTabPane>
    </NTabs>
  </NModal>
</template>

<script lang="ts" setup>
import { useDialogStore } from '@/stores/dialog';
import { model } from '@wailsApp/go/models';
import { ListGroup } from '@wailsApp/go/services/GroupSrv';
import { CreateHost } from '@wailsApp/go/services/HostSrv';

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
  useMessage,
} from 'naive-ui';
import { SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const dialogStore = useDialogStore();
const formRef = ref<FormInst | null>(null);
const activeTab = ref('basic');
const message = useMessage();

const formValue = ref(
  model.Host.createFrom({
    name: '',
    host: '',
    port: 22,
    description: '',
    credential: model.Credential.createFrom({
      username: '',
      password: '',
      auth_type: 0,
      private_key: '',
      key_password: '',
      is_common_credential: false,
    }),
  }),
);

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
    required: true,
    message: t('host_dialog.validation.username_required'),
    trigger: 'blur',
  },
  'credential.password': {
    required: formValue.value.credential?.auth_type === 0,
    message: t('host_dialog.validation.password_required'),
    trigger: 'blur',
  },
  'credential.private_key': {
    required: formValue.value.credential?.auth_type === 1,
    message: t('host_dialog.validation.private_key_required'),
    trigger: 'blur',
  },
}));

const groupOptions = ref<SelectMixedOption[]>([]);

onMounted(async () => {
  const groups = await fetchGroups();
  groupOptions.value = groups.map((group: model.Group) => ({
    label: group.name,
    value: group.id,
  }));
});

const fetchGroups = async () => {
  const resp = await ListGroup();
  if (!resp.ok) {
    message.error(resp.msg);
  }
  return resp.data;
};

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();
    const resp = await CreateHost(formValue.value);
    if (!resp.ok) {
      message.error(resp.msg);
    } else {
      message.success('创建成功');
    }
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
