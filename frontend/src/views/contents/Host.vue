<template>
  <div class="page-container">
    <!-- 左侧边栏 -->
    <div class="sidebar">
      <!-- 分组列表 -->
      <div class="groups-list">
        <div class="list-header">
          <span class="title">分组</span>
          <n-button text size="tiny" @click="dialogStore.openAddGroupDialog">
            <template #icon>
              <icon icon="ph:plus-bold" />
            </template>
          </n-button>
        </div>
        <div class="list-content">
          <div
            v-for="group in groups"
            :key="group.id"
            class="group-item"
            :class="{ active: selectedGroup?.id === group.id }"
            @click="selectGroup(group)"
          >
            <icon icon="ph:folder-simple-duotone" />
            <span class="name">{{ group.name }}</span>
            <span class="count">{{ getHostCountInGroup(group) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <div class="content-header">
        <div class="header-left">
          <h2>{{ selectedGroup ? selectedGroup.name : '所有主机' }}</h2>
          <n-tag size="small">{{ filteredHosts.length }}</n-tag>
        </div>
        <n-button type="primary" size="small" @click="dialogStore.openAddHostDialog()">
          <template #icon>
            <icon icon="ph:plus-bold" />
          </template>
          添加主机
        </n-button>
      </div>

      <!-- 主机网格 -->
      <div class="hosts-grid">
        <div v-for="host in filteredHosts" :key="host.id" class="host-card" @click="toTerminal(host)">
          <div class="card-header">
            <div class="card-left">
              <div class="os-icon">
                <icon :icon="getDeviceIcon(host)" />
              </div>
              <div class="card-info">
                <div class="host-name">{{ host.name }}</div>
                <div class="host-addr">{{ host.credential?.username }}@{{ host.host }}</div>
              </div>
            </div>
            <n-button circle text size="small" class="edit-btn" @click.stop="handleEditHost">
              <template #icon>
                <icon icon="ph:pencil-simple" class="edit-icon" />
              </template>
            </n-button>
          </div>

          <div class="card-body">
            <div class="host-name">{{ host.name }}</div>
            <div class="host-addr">{{ host.credential?.username }}@{{ host.host }}</div>
          </div>

          <div class="card-footer">
            <div class="last-connected">
              <icon icon="ph:clock-counter-clockwise-duotone" />
              <span>{{ getLastConnectedTime(host) }}</span>
            </div>
            <div class="connection-tags">
              <n-tooltip trigger="hover" v-if="getConnectionCount(host) > 0">
                <template #trigger>
                  <div>
                    <n-tag size="tiny" type="success">
                      {{ getConnectionCount(host) }}
                    </n-tag>
                  </div>
                </template>
                活跃连接
              </n-tooltip>
              <n-tooltip trigger="hover" v-if="getErrorConnectionCount(host) > 0">
                <template #trigger>
                  <div>
                    <n-tag size="tiny" type="error">
                      {{ getErrorConnectionCount(host) }}
                    </n-tag>
                  </div>
                </template>
                断开连接
              </n-tooltip>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NTag, NTooltip, useMessage } from 'naive-ui';
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
const selectedGroup = ref<model.Group | null>(null);

const selectGroup = (group: model.Group) => {
  selectedGroup.value = group;
};

const getHostCountInGroup = (group: model.Group) => {
  return hosts.value?.filter(host => host.group_id === group.id).length || 0;
};

const filteredHosts = computed(() => {
  if (!selectedGroup.value) return hosts.value || [];
  return hosts.value?.filter(host => host.group_id === selectedGroup.value?.id) || [];
});

watch(
  () => dialogStore.hostDialogVisible,
  newVal => {
    if (!newVal) {
      fetchData();
    }
  },
);

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

const handleEditGroup = (event: MouseEvent) => {
  event.preventDefault();
  dialogStore.openAddGroupDialog();
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

const fetchData = async () => {
  const [groupsData, hostsData] = await Promise.all([fetchGroups(), fetchHosts()]);
  groups.value = groupsData;
  hosts.value = hostsData;
};

const getDeviceIcon = (host: model.Host) => {
  const os = host.metadata?.os?.toLowerCase() || '';

  if (os.includes('cisco')) {
    return 'simple-icons:cisco';
  } else if (os.includes('huawei')) {
    return 'simple-icons:huawei';
  } else if (os.includes('redhat') || os.includes('rhel')) {
    return 'simple-icons:redhat';
  } else if (os.includes('ubuntu')) {
    return 'simple-icons:ubuntu';
  } else if (os.includes('centos')) {
    return 'simple-icons:centos';
  } else if (os.includes('debian')) {
    return 'simple-icons:debian';
  } else if (os.includes('suse')) {
    return 'simple-icons:opensuse';
  } else if (os.includes('fedora')) {
    return 'simple-icons:fedora';
  } else if (os.includes('linux')) {
    return 'simple-icons:linux';
  }

  return 'ph:computer-tower-duotone';
};

const getLastConnectedTime = (host: model.Host) => {
  // 这里可以根据实际情况返回上次连接时间
  // 暂时返回一个示例时间
  return '10分钟前';
};

const getConnectionCount = (host: model.Host) => {
  return connectionStore.connections.filter(c => c.host === host.host && !c.connectionError).length;
};

const getErrorConnectionCount = (host: model.Host) => {
  return connectionStore.connections.filter(c => c.host === host.host && c.connectionError).length;
};

onMounted(async () => {
  await fetchData();
});
</script>

<style lang="less" scoped>
.page-container {
  height: 100%;
  display: flex;
  background: v-bind('gtermThemeVars.cardColor');
}

// 左侧边栏样式
.sidebar {
  width: 260px;
  border-right: 1px solid v-bind('`${gtermThemeVars.primaryColor}10`');
  display: flex;
  flex-direction: column;
}

.groups-list {
  flex: 1;
  display: flex;
  flex-direction: column;

  .list-header {
    padding: 16px 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid v-bind('`${gtermThemeVars.primaryColor}10`');

    .title {
      font-size: 13px;
      font-weight: 600;
      color: v-bind('gtermThemeVars.textColor');
      opacity: 0.6;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
  }

  .list-content {
    flex: 1;
    overflow-y: auto;
    padding: 0 12px;
  }

  .group-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px;
    border-radius: 6px;
    cursor: pointer;
    margin: 2px 0;
    color: v-bind('gtermThemeVars.textColor');
    transition: all 0.2s ease;

    &:hover {
      background: v-bind('`${gtermThemeVars.primaryColor}08`');
    }

    &.active {
      background: v-bind('`${gtermThemeVars.primaryColor}15`');
      color: v-bind('gtermThemeVars.primaryColor');
    }

    .name {
      flex: 1;
      font-size: 14px;
    }

    .count {
      font-size: 12px;
      opacity: 0.5;
      min-width: 20px;
      height: 20px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}

// 主内容区样式
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  padding: 24px;

  .content-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;

    .header-left {
      display: flex;
      align-items: center;
      gap: 12px;

      h2 {
        font-size: 20px;
        font-weight: 600;
        color: v-bind('gtermThemeVars.textColor');
        margin: 0;
      }
    }
  }
}

// 主机网格样式
.hosts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  overflow-y: auto;
  padding-right: 4px;
}

.host-card {
  background: v-bind('gtermThemeVars.cardHoverColor');
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  border: 1px solid transparent;

  &:hover {
    border-color: v-bind('`${gtermThemeVars.primaryColor}30`');

    .card-header {
      .edit-btn {
        opacity: 1;
      }
    }
  }

  .card-header {
    padding: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .card-left {
      display: flex;
      align-items: center;
      gap: 10px;
      flex: 1;
      min-width: 0;

      .os-icon {
        width: 38px;
        height: 38px;
        border-radius: 6px;
        background: v-bind('`${gtermThemeVars.primaryColor}10`');
        color: v-bind('gtermThemeVars.primaryColor');
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 22px;
        flex-shrink: 0;
      }

      .card-info {
        min-width: 0;

        .host-name {
          font-size: 14px;
          font-weight: 600;
          color: v-bind('gtermThemeVars.textColor');
          margin-bottom: 2px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .host-addr {
          font-size: 12px;
          color: v-bind('gtermThemeVars.secondaryText');
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }

    .edit-btn {
      opacity: 0;
      transition: all 0.2s ease;
      margin-left: 8px;
      background: v-bind('`${gtermThemeVars.primaryColor}10`');
      width: 32px;
      height: 32px;
      border-radius: 6px !important;

      :deep(.edit-icon) {
        font-size: 16px;
      }

      &:hover {
        color: v-bind('gtermThemeVars.primaryColor');
        background: v-bind('`${gtermThemeVars.primaryColor}20`');
      }
    }
  }

  .card-body {
    display: none;
  }

  .card-footer {
    padding: 8px 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: auto;
    background: v-bind('`${gtermThemeVars.primaryColor}05`');

    .last-connected {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      color: v-bind('gtermThemeVars.textColor');
      opacity: 0.5;

      :deep(svg) {
        font-size: 14px;
      }
    }

    .connection-tags {
      display: flex;
      gap: 4px;
      align-items: center;
    }

    :deep(.n-tag) {
      padding: 0;
      width: 18px;
      height: 18px;
      line-height: 18px;
      text-align: center;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>
