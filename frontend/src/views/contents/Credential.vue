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
              <icon icon="ph:magnifying-glass-duotone" />
            </template>
          </n-input>
          <n-button type="primary" ghost>
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
        <n-list hoverable clickable class="credential-list">
          <n-list-item v-for="item in mockCredential" :key="item.id">
            <div class="credential-item">
              <n-thing>
                <template #avatar>
                  <div class="credential-type" :class="item.authType === 'password' ? 'success' : 'warning'">
                    <icon :icon="item.authType === 'password' ? 'ph:password-duotone' : 'ph:key-duotone'" />
                  </div>
                </template>
                <template #header>
                  <div class="credential-header">
                    <span class="name">{{ item.name }}</span>
                  </div>
                </template>
                <template #description>
                  <div class="credential-info">
                    <span class="info-item">
                      <icon icon="ph:user-duotone" />
                      {{ item.username }}
                    </span>
                    <span class="info-item">
                      <icon icon="ph:clock-duotone" />
                      {{ item.createdAt }} 创建
                    </span>
                  </div>
                </template>
              </n-thing>
              <div class="credential-actions">
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button circle text>
                      <template #icon>
                        <icon :icon="item.authType === 'password' ? 'ph:copy-duotone' : 'ph:file-text-duotone'" />
                      </template>
                    </n-button>
                  </template>
                  {{ item.authType === 'password' ? '复制密码' : '查看密钥' }}
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button circle text>
                      <template #icon>
                        <icon icon="ph:pencil-simple-duotone" />
                      </template>
                    </n-button>
                  </template>
                  编辑
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button circle text type="error">
                      <template #icon>
                        <icon icon="ph:trash-duotone" />
                      </template>
                    </n-button>
                  </template>
                  删除
                </n-tooltip>
              </div>
            </div>
          </n-list-item>
        </n-list>
      </div>
    </n-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NInput, NInputGroup, NList, NListItem, NTag, NThing, NScrollbar, useThemeVars } from 'naive-ui';
import { ref } from 'vue';

interface CredentialItem {
  id: number;
  name: string;
  username: string;
  authType: 'password' | 'private_key';
  createdAt: string;
}

const mockCredential = ref<CredentialItem[]>([
  {
    id: 1,
    name: 'root@192.168.1.100',
    username: 'root',
    authType: 'password',
    createdAt: '2024-01-20',
  },
  {
    id: 2,
    name: 'admin@example.com',
    username: 'admin',
    authType: 'private_key',
    createdAt: '2024-01-19',
  },
  {
    id: 3,
    name: 'dev@staging-server',
    username: 'developer',
    authType: 'private_key',
    createdAt: '2024-01-18',
  },
  {
    id: 4,
    name: 'ubuntu@aws-ec2',
    username: 'ubuntu',
    authType: 'private_key',
    createdAt: '2024-01-17',
  },
  {
    id: 5,
    name: 'jenkins@ci-server',
    username: 'jenkins',
    authType: 'password',
    createdAt: '2024-01-16',
  },
  {
    id: 6,
    name: 'gitlab@runner',
    username: 'gitlab-runner',
    authType: 'private_key',
    createdAt: '2024-01-15',
  },
  {
    id: 7,
    name: 'mysql@db-server',
    username: 'mysql',
    authType: 'password',
    createdAt: '2024-01-14',
  },
  {
    id: 8,
    name: 'nginx@web-server',
    username: 'www-data',
    authType: 'password',
    createdAt: '2024-01-13',
  },
  {
    id: 9,
    name: 'redis@cache-server',
    username: 'redis',
    authType: 'password',
    createdAt: '2024-01-12',
  },
  {
    id: 10,
    name: 'elastic@es-node',
    username: 'elastic',
    authType: 'private_key',
    createdAt: '2024-01-11',
  },
]);

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
</style>
