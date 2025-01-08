<template>
  <div class="inventory-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <n-input placeholder="搜索主机或用户@主机名...">
        <template #prefix>
          <icon icon="ph:magnifying-glass-duotone" />
        </template>
      </n-input>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar">
      <n-button size="small" type="primary" ghost @click="() => dialogStore.openAddHostDialog()">
        <template #icon>
          <icon icon="ph:plus-circle-duotone" />
        </template>
        添加主机
      </n-button>
      <n-button size="small" ghost @click="dialogStore.openAddGroupDialog">
        <template #icon>
          <icon icon="ph:folder-simple-plus-duotone" />
        </template>
        新建分组
      </n-button>
    </div>

    <!-- 分组列表 -->
    <div class="section">
      <div class="section-header">
        <n-h3 prefix="bar">
          <icon icon="ph:folders-duotone" class="section-icon" />
          分组
        </n-h3>
        <span class="text-count">共 {{ groups?.length }} 个分组</span>
      </div>

      <n-grid x-gap="12" y-gap="12" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
        <n-gi v-for="v in groups" :key="v.id">
          <div class="card group">
            <div class="card-content">
              <div class="card-icon">
                <icon icon="ph:folders-duotone" />
              </div>
              <div class="card-info">
                <div class="card-title">{{ v.name }}</div>
                <div class="card-subtitle">1 台主机</div>
              </div>
            </div>
            <div class="card-actions">
              <n-button text circle class="action-btn">
                <template #icon>
                  <icon icon="ph:pencil-simple-duotone" />
                </template>
              </n-button>
            </div>
          </div>
        </n-gi>
      </n-grid>
    </div>

    <!-- 主机列表 -->
    <div class="section">
      <div class="section-header">
        <n-h3 prefix="bar">
          <icon icon="ph:computer-tower-duotone" class="section-icon" />
          主机
        </n-h3>
        <span class="text-count">共 {{ hosts?.length }} 台主机</span>
      </div>

      <n-grid x-gap="12" y-gap="12" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
        <n-gi v-for="v in hosts" :key="v.id">
          <div class="card host" @click="toTerminal(v)">
            <div class="card-content">
              <div class="card-left">
                <div class="card-icon">
                  <icon icon="ph:computer-tower-duotone" />
                </div>
                <div class="card-info">
                  <div class="card-header">
                    <div class="card-title">{{ v.name }}</div>
                  </div>
                  <div class="card-subtitle">
                    <span>ssh://{{ v.credential?.username }}@{{ v.host }}</span>
                  </div>
                </div>
              </div>
              <div class="card-stats">
                <div class="stat-item">
                  <div class="stat-label">CPU</div>
                  <div class="stat-value">
                    <icon icon="ph:cpu-duotone" class="stat-icon" />
                    <span>{{ v.metadata?.processors }} 核心</span>
                  </div>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                  <div class="stat-label">内存</div>
                  <div class="stat-value">
                    <icon icon="ph:memory-duotone" class="stat-icon" />
                    <span>{{ v.metadata?.mem_total }} GB</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="card-actions">
              <n-button text circle class="action-btn" @click.stop="handleEditHost">
                <template #icon>
                  <icon icon="ph:dots-three-outline-duotone" />
                </template>
              </n-button>
            </div>
          </div>
        </n-gi>
      </n-grid>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NGi, NGrid, NH3, NInput, NTag, useMessage } from 'naive-ui';
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { ListGroup } from '@wailsApp/go/services/GroupSrv';
import { ListHost } from '@wailsApp/go/services/HostSrv';
import { model } from '@wailsApp/go/models';
import { useConnectionStore } from '@/stores/connection';
import { useRouter } from 'vue-router';

const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const router = useRouter();
const message = useMessage();
const connectionStore = useConnectionStore();

const groups = ref<model.Group[]>();
const hosts = ref<model.Host[]>();

const toTerminal = (host: model.Host) => {
  const connection = {
    id: Date.now(),
    hostId: host.id,
    name: `${host.name} (${connectionStore.connections.filter(c => c.host === host.host).length + 1})`,
    host: host.host,
    username: host.credential?.username || '',
  };
  connectionStore.addConnection(connection);
  router.push({ name: 'Terminal' });
};

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
});

const handleEditHost = (event: MouseEvent) => {
  event.preventDefault();
  dialogStore.openAddHostDialog(true);
};

const fetchGroups = async () => {
  const resp = await ListGroup();
  if (!resp.ok) {
    message.error(resp.msg);
  }
  return resp.data;
};

const fetchHosts = async () => {
  const resp = await ListHost();
  if (!resp.ok) {
    message.error(resp.msg);
  }
  return resp.data;
};

onMounted(async () => {
  const [groupsData, hostsData] = await Promise.all([fetchGroups(), fetchHosts()]);
  groups.value = groupsData;
  hosts.value = hostsData;
});
</script>

<style lang="less" scoped>
.inventory-container {
  padding: 16px 24px;
  height: 100%;
}

.search-bar {
  margin-bottom: 16px;
}

.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.section {
  margin-bottom: 24px;

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;

    :deep(.n-h3) {
      margin: 0;
      display: flex;
      align-items: center;
      gap: 8px;

      .section-icon {
        font-size: 20px;
        opacity: 0.8;
      }
    }
  }
}

// 分组卡片样式
.card.group {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-radius: 12px;
  background-color: v-bind('gtermThemeVars.cardColor');
  border: 1px solid v-bind('gtermThemeVars.borderColor');
  transition: all 0.2s ease;

  &:hover {
    background-color: v-bind('gtermThemeVars.cardHoverColor');
    border-color: v-bind('gtermThemeVars.primaryColor');
  }

  .card-content {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .card-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 52px;
    height: 52px;
    border-radius: 10px;
    background-color: v-bind('gtermThemeVars.primaryColor');
    color: white;
    font-size: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);

    :deep(svg) {
      width: 28px;
      height: 28px;
      filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
    }
  }

  .card-info {
    .card-title {
      font-weight: 500;
      margin-bottom: 2px;
    }

    .card-subtitle {
      font-size: 12px;
      color: v-bind('gtermThemeVars.secondaryText');
      opacity: 0.7;
    }
  }

  .card-actions {
    opacity: 0;
    transition: opacity 0.2s ease;

    .action-btn {
      width: 32px;
      height: 32px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: v-bind('gtermThemeVars.textColor');
      opacity: 0.7;
      transition: all 0.2s ease;

      &:hover {
        opacity: 1;
        background-color: v-bind('gtermThemeVars.splitColor');
        transform: scale(1.05);
      }

      &:active {
        transform: scale(0.95);
      }

      :deep(.n-icon) {
        font-size: 18px;
      }
    }
  }

  &:hover .card-actions {
    opacity: 1;
  }
}

// 主机卡片样式
.card.host {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 16px;
  border-radius: 12px;
  background-color: v-bind('gtermThemeVars.cardColor');
  border: 1px solid v-bind('gtermThemeVars.borderColor');
  transition: all 0.2s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;

  &:hover {
    background-color: v-bind('gtermThemeVars.cardHoverColor');
    border-color: v-bind('gtermThemeVars.primaryColor');
  }

  .card-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
    flex: 1;
    min-width: 0;
  }

  .card-left {
    display: flex;
    align-items: flex-start;
    gap: 16px;
  }

  .card-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 52px;
    height: 52px;
    border-radius: 10px;
    background: v-bind('gtermThemeVars.primaryColor');
    color: white;
    font-size: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);

    :deep(svg) {
      width: 28px;
      height: 28px;
      filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
    }
  }

  .card-info {
    flex: 1;
    min-width: 0;

    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 6px;
    }

    .card-title {
      font-weight: 600;
      font-size: 16px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      color: v-bind('gtermThemeVars.textColor');
    }

    .card-subtitle {
      font-size: 13px;
      color: v-bind('gtermThemeVars.secondaryText');
      margin-bottom: 0;
      display: flex;
      align-items: center;
      gap: 4px;

      .info-icon {
        font-size: 14px;
        opacity: 0.7;
      }

      span {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
  }

  .card-stats {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px;
    background-color: v-bind('gtermThemeVars.cardHoverColor');
    border-radius: 8px;
    border: 1px solid v-bind('gtermThemeVars.borderColor');

    .stat-item {
      flex: 1;

      .stat-label {
        font-size: 12px;
        color: v-bind('gtermThemeVars.secondaryText');
        margin-bottom: 4px;
      }

      .stat-value {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 14px;
        color: v-bind('gtermThemeVars.textColor');
        font-weight: 500;

        .stat-icon {
          font-size: 16px;
          color: v-bind('gtermThemeVars.primaryColor');
        }
      }
    }

    .stat-divider {
      width: 1px;
      height: 32px;
      background-color: v-bind('gtermThemeVars.borderColor');
      opacity: 0.6;
    }
  }

  .card-actions {
    position: absolute;
    top: 12px;
    right: 12px;
    opacity: 0;
    transition: all 0.2s ease;

    .action-btn {
      width: 32px;
      height: 32px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: v-bind('gtermThemeVars.textColor');
      background-color: v-bind('gtermThemeVars.cardColor');
      border: 1px solid v-bind('gtermThemeVars.borderColor');
      transition: all 0.2s ease;

      &:hover {
        color: v-bind('gtermThemeVars.primaryColor');
        background-color: v-bind('gtermThemeVars.cardHoverColor');
        border-color: v-bind('gtermThemeVars.primaryColor');
        transform: scale(1.05);
      }

      &:active {
        transform: scale(0.95);
      }

      :deep(.n-icon) {
        font-size: 18px;
      }
    }
  }

  &:hover .card-actions {
    opacity: 1;
  }
}

.text-count {
  font-size: 12px;
  color: v-bind('gtermThemeVars.secondaryText');
  line-height: 1;
  display: flex;
  align-items: center;
}
</style>
