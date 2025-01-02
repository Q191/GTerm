<template>
  <div class="inventory-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <NInput placeholder="搜索主机或用户@主机名...">
        <template #prefix>
          <Icon icon="ph:magnifying-glass-duotone" />
        </template>
      </NInput>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar">
      <NButton size="small" type="primary" ghost @click="() => dialogStore.openAddHostDialog()">
        <template #icon>
          <Icon icon="ph:plus-circle-duotone" />
        </template>
        添加主机
      </NButton>
      <NButton size="small" ghost @click="dialogStore.openAddGroupDialog">
        <template #icon>
          <Icon icon="ph:folder-simple-plus-duotone" />
        </template>
        新建分组
      </NButton>
    </div>

    <!-- 分组列表 -->
    <div class="section">
      <div class="section-header">
        <NH3 prefix="bar">分组</NH3>
        <span class="text-count">共 {{ groups?.length }} 个分组</span>
      </div>

      <NGrid x-gap="12" y-gap="12" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
        <NGi v-for="v in groups" :key="v.id">
          <div class="card group">
            <div class="card-content">
              <div class="card-icon">
                <Icon icon="ph:folders-duotone" />
              </div>
              <div class="card-info">
                <div class="card-title">{{ v.name }}</div>
                <div class="card-subtitle">1 台主机</div>
              </div>
            </div>
            <div class="card-actions">
              <NButton text circle class="action-btn">
                <template #icon>
                  <Icon icon="ph:pencil-simple-duotone" />
                </template>
              </NButton>
            </div>
          </div>
        </NGi>
      </NGrid>
    </div>

    <!-- 主机列表 -->
    <div class="section">
      <div class="section-header">
        <NH3 prefix="bar">主机</NH3>
        <span class="text-count">共 {{ hosts?.length }} 台主机</span>
      </div>

      <NGrid x-gap="12" y-gap="12" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
        <NGi v-for="v in hosts" :key="v.id">
          <div class="card group" @click="toTerminal(v)">
            <div class="card-content">
              <div class="card-icon">
                <Icon icon="ph:linux-logo-duotone" />
              </div>
              <div class="card-info">
                <div class="card-title">{{ v.name }}</div>
                <div class="card-subtitle">ssh://{{ v.credential?.username }}@{{ v.host }}</div>
              </div>
            </div>
            <div class="card-actions">
              <NButton text circle class="action-btn" @click.stop="handleEditHost">
                <template #icon>
                  <Icon icon="ph:pencil-simple-duotone" />
                </template>
              </NButton>
            </div>
          </div>
        </NGi>
      </NGrid>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NGi, NGrid, NH3, NInput, useMessage } from 'naive-ui';
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';
import { ListGroup, ListHost } from '@wailsApp/go/services/GroupSrv';
import { model } from '@wailsApp/go/models';
import { useConnectionStore } from '@/stores/connection';

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
  overflow-y: auto;
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
    }
  }
}

.card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-radius: 8px;
  background-color: v-bind('gtermThemeVars.cardColor');
  border: 1px solid v-bind('gtermThemeVars.borderColor');
  transition: all 0.2s ease;

  &:hover {
    background-color: v-bind('gtermThemeVars.cardHoverColor');
    border-color: v-bind('gtermThemeVars.splitColor');
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
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
    width: 40px;
    height: 40px;
    border-radius: 8px;
    background-color: v-bind('gtermThemeVars.primaryColor');
    color: white;
    font-size: 24px;

    :deep(svg) {
      width: 24px;
      height: 24px;
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

.text-count {
  font-size: 12px;
  color: v-bind('gtermThemeVars.secondaryText');
  line-height: 1;
  display: flex;
  align-items: center;
}
</style>
