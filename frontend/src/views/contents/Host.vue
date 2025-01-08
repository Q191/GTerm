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
    <div class="section" v-if="groups?.length">
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
    <div class="section" v-if="hosts?.length">
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
              <div class="os-icon">
                <icon :icon="getDeviceIcon(v)" />
              </div>
              <div class="host-info">
                <div class="host-name">{{ v.name }}</div>
                <div class="host-addr">{{ v.credential?.username }}@{{ v.host }}</div>
              </div>
            </div>
            <div class="card-actions">
              <n-button text circle class="action-btn" @click.stop="handleEditHost">
                <template #icon>
                  <icon icon="ph:pencil-simple-duotone" />
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

onMounted(async () => {
  await fetchData();
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
  padding: 10px 12px;
  border-radius: 12px;
  background-color: v-bind('gtermThemeVars.cardColor');
  border: 1px solid v-bind('gtermThemeVars.borderColor');
  transition: all 0.2s ease;
  cursor: pointer;

  &:hover {
    border-color: v-bind('gtermThemeVars.primaryColor');
    background-color: v-bind('`${gtermThemeVars.primaryColor}08`');

    .card-icon {
      color: white;
      background: v-bind('gtermThemeVars.primaryColor');
    }

    .card-actions {
      opacity: 1;
    }
  }

  .card-content {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    min-width: 0;
  }

  .card-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 44px;
    height: 44px;
    border-radius: 10px;
    background: v-bind('`${gtermThemeVars.primaryColor}15`');
    color: v-bind('gtermThemeVars.primaryColor');
    font-size: 22px;
    transition: all 0.2s ease;

    :deep(svg) {
      filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
    }
  }

  .card-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;

    .card-title {
      font-size: 14px;
      font-weight: 600;
      color: v-bind('gtermThemeVars.textColor');
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 1.2;
    }

    .card-subtitle {
      font-size: 12px;
      color: v-bind('gtermThemeVars.secondaryText');
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 1.2;
    }
  }

  .card-actions {
    opacity: 0;
    transition: opacity 0.2s ease;
    width: 32px;
    height: 32px;
    border-radius: 8px;

    .action-btn {
      width: 100%;
      height: 100%;
      border-radius: inherit;
      display: flex;
      align-items: center;
      justify-content: center;
      color: v-bind('gtermThemeVars.textColor');
      background: v-bind('gtermThemeVars.cardHoverColor');
      transition: all 0.2s ease;

      &:hover {
        background-color: v-bind('gtermThemeVars.primaryColor');
        color: white;
      }
    }
  }
}

// 主机卡片样式
.card.host {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 12px;
  background-color: v-bind('gtermThemeVars.cardColor');
  border: 1px solid v-bind('gtermThemeVars.borderColor');
  transition: all 0.2s ease;
  cursor: pointer;

  &:hover {
    border-color: v-bind('gtermThemeVars.primaryColor');
    background-color: v-bind('`${gtermThemeVars.primaryColor}08`');

    .os-icon {
      color: white;
      background: v-bind('gtermThemeVars.primaryColor');
    }

    .card-actions {
      opacity: 1;
    }
  }

  .card-content {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    min-width: 0;
  }

  .os-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 44px;
    height: 44px;
    border-radius: 10px;
    background: v-bind('`${gtermThemeVars.primaryColor}15`');
    color: v-bind('gtermThemeVars.primaryColor');
    font-size: 22px;
    transition: all 0.2s ease;

    :deep(svg) {
      filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
    }
  }

  .host-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;

    .host-name {
      font-size: 14px;
      font-weight: 600;
      color: v-bind('gtermThemeVars.textColor');
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 1.2;
    }

    .host-addr {
      font-size: 12px;
      color: v-bind('gtermThemeVars.secondaryText');
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 1.2;
    }
  }

  .card-actions {
    opacity: 0;
    transition: opacity 0.2s ease;
    width: 32px;
    height: 32px;
    border-radius: 8px;

    .action-btn {
      width: 100%;
      height: 100%;
      border-radius: inherit;
      display: flex;
      align-items: center;
      justify-content: center;
      color: v-bind('gtermThemeVars.textColor');
      background: v-bind('gtermThemeVars.cardHoverColor');
      transition: all 0.2s ease;

      &:hover {
        background-color: v-bind('gtermThemeVars.primaryColor');
        color: white;
      }
    }
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
