<template>
  <div class="credential-container">
    <div class="header">
      <div class="title">
        <span>凭据</span>
      </div>
      <div class="actions">
        <n-input-group>
          <n-input placeholder="搜索凭据...">
            <template #prefix>
              <icon icon="ph:magnifying-glass" />
            </template>
          </n-input>
          <n-button type="primary" ghost @click="handleAdd">
            <template #icon>
              <icon icon="ph:plus-bold" />
            </template>
            添加凭据
          </n-button>
        </n-input-group>
      </div>
    </div>

    <n-scrollbar class="content">
      <div class="content-wrapper">
        <n-list v-if="creds?.length" hoverable clickable class="credential-list">
          <n-list-item v-for="v in creds" :key="v.id">
            <div class="credential-item">
              <n-thing>
                <template #avatar>
                  <div
                    class="credential-type"
                    :class="v.authMethod === enums.AuthMethod.PASSWORD ? 'success' : 'warning'"
                  >
                    <icon :icon="v.authMethod === enums.AuthMethod.PASSWORD ? 'ph:password' : 'ph:key'" />
                  </div>
                </template>
                <template #header>
                  <div class="credential-header">
                    <span class="name">{{ v.label }}</span>
                  </div>
                </template>
                <template #description>
                  <div class="credential-info">
                    <span class="info-item">
                      <icon icon="ph:user" />
                      {{ v.username }}
                    </span>
                    <span class="info-item">
                      <icon icon="ph:clock" />
                      {{ formatTime(v.createdAt) }}
                    </span>
                  </div>
                </template>
              </n-thing>
              <div class="credential-actions">
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button circle text @click="handleCopy(v)">
                      <template #icon>
                        <icon :icon="v.authMethod === enums.AuthMethod.PASSWORD ? 'ph:copy' : 'ph:file-text'" />
                      </template>
                    </n-button>
                  </template>
                  {{ v.authMethod === enums.AuthMethod.PASSWORD ? '复制密码' : '查看密钥' }}
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button circle text @click="handleEdit(v)">
                      <template #icon>
                        <icon icon="ph:pencil-simple" />
                      </template>
                    </n-button>
                  </template>
                  编辑
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button circle text type="error" @click="handleDelete(v)">
                      <template #icon>
                        <icon icon="ph:trash" />
                      </template>
                    </n-button>
                  </template>
                  删除
                </n-tooltip>
              </div>
            </div>
          </n-list-item>
        </n-list>
        <n-empty v-else class="credential-empty" description="暂无可用凭据" />
      </div>
    </n-scrollbar>

    <credential-dialog
      v-model:show="showDialog"
      :is-edit="isEdit"
      :credential="editCredential"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import {
  NButton,
  NInput,
  NInputGroup,
  NList,
  NListItem,
  NThing,
  NScrollbar,
  useThemeVars,
  useMessage,
  NEmpty,
} from 'naive-ui';
import { ref, onMounted } from 'vue';
import { ListCredential, DeleteCredential } from '@wailsApp/go/services/CredentialSrv';
import { enums, model } from '@wailsApp/go/models';
import dayjs from 'dayjs';
import CredentialDialog from '@/views/dialogs/CredentialDialog.vue';

const message = useMessage();
const showDialog = ref(false);
const isEdit = ref(false);
const editCredential = ref<model.Credential | undefined>(undefined);

const handleCopy = async (credential: model.Credential) => {
  if (credential.authMethod === enums.AuthMethod.PASSWORD) {
    try {
      await navigator.clipboard.writeText(credential.password);
      message.success('密码已复制到剪贴板');
    } catch (err) {
      message.error('复制失败');
    }
  }
};

const handleEdit = (credential: model.Credential) => {
  console.log('准备编辑凭据:', credential);
  isEdit.value = true;
  editCredential.value = credential;
  showDialog.value = true;
};

const handleDelete = async (credential: model.Credential) => {
  console.log('准备删除凭据:', credential);
  const resp = await DeleteCredential(credential.id);
  if (!resp.ok) {
    message.error(resp.msg);
    return;
  }
  message.success('删除成功');
  await fetchCredentials();
};

const handleAdd = () => {
  isEdit.value = false;
  editCredential.value = undefined;
  showDialog.value = true;
};

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm');
};

const creds = ref<model.Credential[]>();

const fetchCredentials = async () => {
  console.log('开始获取凭据列表');
  const resp = await ListCredential();
  if (!resp.ok) {
    message.error(resp.msg);
  }
  console.log('获取凭据列表完成:', resp.data);
  creds.value = resp.data;
  return resp.data;
};

const handleSuccess = () => {
  console.log('收到凭据操作成功事件');
  fetchCredentials();
};

onMounted(async () => {
  const [credsData] = await Promise.all([fetchCredentials()]);
  creds.value = credsData;
});

const themeVars = useThemeVars();
</script>

<style lang="less" scoped>
.credential-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
}

.header {
  position: sticky;
  top: 0;
  padding: 16px 24px;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;

  .title {
    width: 100%;
    max-width: 800px;
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    font-size: 18px;
    font-weight: 500;
    color: v-bind('themeVars.textColorBase');
  }

  .actions {
    width: 100%;
    max-width: 800px;

    :deep(.n-input-group) {
      width: 100%;
      display: flex;
      gap: 8px;

      .n-input {
        flex: 1;
      }
    }
  }
}

.content {
  flex: 1;
  min-height: 0;
}

.content-wrapper {
  padding: 16px;

  .credential-list {
    margin: 0 auto;
    width: 100%;
    max-width: 800px;
  }
}

.credential-type {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  transition: all 0.2s ease;

  &.success {
    background-color: rgba(63, 185, 80, 0.1);
    color: #3fb950;
  }

  &.warning {
    background-color: rgba(210, 153, 34, 0.1);
    color: #d29922;
  }
}

.credential-header {
  display: flex;
  align-items: center;
  gap: 8px;

  .name {
    font-weight: 500;
    color: v-bind('themeVars.textColorBase');
  }
}

.credential-info {
  display: flex;
  gap: 16px;
  margin-top: 4px;
  font-size: 13px;
  color: v-bind('themeVars.textColor3');

  .info-item {
    display: flex;
    align-items: center;
    gap: 4px;

    .iconify {
      font-size: 16px;
    }
  }
}

.credential-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 16px;

  :deep(.n-thing) {
    flex: 1;
    min-width: 0;
  }
}

.credential-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
  margin-right: -8px;

  :deep(.n-button) {
    width: 32px;
    height: 32px;

    .n-button__icon {
      font-size: 16px;
    }

    &:hover {
      background-color: v-bind('themeVars.hoverColor');
    }
  }
}

.credential-empty {
  margin: 120px auto;
  width: 100%;
  max-width: 800px;
}
</style>
