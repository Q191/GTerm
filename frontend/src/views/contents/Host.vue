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
          <n-tree
            block-line
            :data="treeData"
            :pattern="searchText"
            :show-irrelevant-nodes="false"
            :selected-keys="selectedKeys"
            :expanded-keys="expandedKeys"
            :node-props="nodeProps"
            @update:selected-keys="handleSelect"
            @update:expanded-keys="handleExpand"
          />
          <n-dropdown
            trigger="manual"
            placement="bottom-start"
            :show="showDropdown"
            :options="dropdownOptions"
            :x="dropdownX"
            :y="dropdownY"
            @select="handleDropdownSelect"
            @clickoutside="handleClickoutside"
          />
        </div>
        <div class="list-footer">
          <n-input v-model:value="searchText" size="small" clearable placeholder="搜索主机">
            <template #prefix>
              <icon icon="ph:magnifying-glass" />
            </template>
          </n-input>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <div class="content-header">
        <div class="header-left">
          <h2>{{ selectedGroup ? selectedGroup.name : '所有主机' }}</h2>
          <n-badge :value="filteredHosts.length" show-zero type="success" />
        </div>
        <n-button type="primary" size="small" @click="dialogStore.openAddHostDialog()">
          <template #icon>
            <icon icon="ph:plus-bold" />
          </template>
          添加主机
        </n-button>
      </div>

      <!-- 主机网格 -->
      <div class="hosts-grid" v-if="filteredHosts.length > 0">
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
            <div class="protocol-info">
              <icon :icon="getProtocolIcon(host)" />
              <span>{{ getProtocolName(host) }}</span>
            </div>
            <div class="connection-tags" style="margin-left: auto">
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
      <div class="empty-state" v-else>
        <n-result status="404" title="空空如也" :description="`分组「${selectedGroup?.name}」中还没有添加任何主机`" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NTag, NTooltip, NTree, NDropdown, NResult, NBadge, NInput, useMessage, useThemeVars } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import { useDialogStore } from '@/stores/dialog';
import { ListGroup } from '@wailsApp/go/services/GroupSrv';
import { ListHost } from '@wailsApp/go/services/HostSrv';
import { model } from '@wailsApp/go/models';
import { useConnectionStore } from '@/stores/connection';
import { useRouter } from 'vue-router';
import { h, ref, computed, onMounted, watch } from 'vue';

const dialogStore = useDialogStore();
const router = useRouter();
const message = useMessage();
const connectionStore = useConnectionStore();

const groups = ref<model.Group[]>();
const hosts = ref<model.Host[]>();
const selectedGroup = ref<model.Group | null>(null);

const selectedKeys = ref<string[]>([]);
const expandedKeys = ref<string[]>([]);

const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const dropdownOptions = ref<DropdownOption[]>([]);

const searchText = ref('');

const updateDropdownOptions = (type: 'group' | 'host') => {
  if (type === 'group') {
    dropdownOptions.value = [
      {
        label: '编辑分组',
        key: 'edit-group',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      },
      {
        label: '删除分组',
        key: 'delete-group',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      },
    ];
  } else {
    dropdownOptions.value = [
      {
        label: '编辑主机',
        key: 'edit-host',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      },
      {
        label: '删除主机',
        key: 'delete-host',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      },
    ];
  }
};

const currentContextNode = ref<any>(null);

const nodeProps = ({ option }: { option: any }) => {
  return {
    onClick() {
      if (option.key.startsWith('group-')) {
        const groupId = parseInt(option.key.replace('group-', ''));
        const group = groups.value?.find(g => g.id === groupId);
        if (group) {
          const groupHosts = hosts.value?.filter(host => host.group_id === group.id) || [];
          if (groupHosts.length > 0) {
            selectedGroup.value = group;
            selectedKeys.value = [option.key];
          }
        }
      }
    },
    onContextmenu(e: MouseEvent) {
      currentContextNode.value = option;
      dropdownX.value = e.clientX;
      dropdownY.value = e.clientY;
      showDropdown.value = true;

      if (option.key.startsWith('group-')) {
        updateDropdownOptions('group');
      } else if (option.key.startsWith('host-')) {
        updateDropdownOptions('host');
      }

      e.preventDefault();
    },
  };
};

const handleDropdownSelect = (key: string) => {
  showDropdown.value = false;
  if (!currentContextNode.value) return;

  switch (key) {
    case 'edit-group':
      dialogStore.openAddGroupDialog();
      break;
    case 'delete-group':
      // TODO: 实现删除分组功能
      break;
    case 'edit-host':
      const hostId = parseInt(currentContextNode.value.key.replace('host-', ''));
      const host = hosts.value?.find(h => h.id === hostId);
      if (host) {
        dialogStore.openAddHostDialog(true);
      }
      break;
    case 'delete-host':
      // TODO: 实现删除主机功能
      break;
  }
  currentContextNode.value = null;
};

const handleClickoutside = () => {
  showDropdown.value = false;
};

const treeData = computed(() => {
  if (!groups.value) return [];

  const treeNodes = [];

  // 添加"全部主机"节点
  treeNodes.push({
    key: 'all',
    label: '所有主机',
    prefix: () =>
      h(Icon, {
        icon: 'ph:paw-print',
        style: {
          fontSize: '1.2rem',
        },
      }),
  });

  // 添加分组节点
  treeNodes.push(
    ...groups.value.map(group => {
      const groupHosts = hosts.value?.filter(host => host.group_id === group.id) || [];
      return {
        key: `group-${group.id}`,
        label: group.name,
        children:
          groupHosts.length > 0
            ? groupHosts.map(host => ({
                key: `host-${host.id}`,
                label: host.name,
                isLeaf: true,
                prefix: () =>
                  h(Icon, {
                    icon: getDeviceIcon(host),
                    style: {
                      fontSize: '1.2rem',
                    },
                  }),
              }))
            : undefined,
        isLeaf: groupHosts.length === 0,
        prefix: () => {
          const expanded = expandedKeys.value.includes(`group-${group.id}`);
          return h(Icon, {
            icon: expanded
              ? 'ph:folder-open-duotone'
              : groupHosts.length === 0
                ? 'ph:folder-dashed-duotone'
                : 'ph:folders-duotone',
            style: { fontSize: '1.2rem' },
          });
        },
      };
    }),
  );

  // 添加未分组的主机节点
  const ungroupedHosts = hosts.value?.filter(host => !host.group_id) || [];
  if (ungroupedHosts.length > 0) {
    treeNodes.push(
      ...ungroupedHosts.map(host => ({
        key: `host-${host.id}`,
        label: host.name,
        isLeaf: true,
        prefix: () =>
          h(Icon, {
            icon: getDeviceIcon(host),
            style: {
              fontSize: '1rem',
            },
          }),
      })),
    );
  }

  return treeNodes;
});

const handleSelect = (keys: string[]) => {
  if (keys.length > 0) {
    const key = keys[0];
    if (key === 'all') {
      selectedGroup.value = null;
      selectedKeys.value = keys;
    } else if (key.startsWith('host-')) {
      const hostId = parseInt(key.replace('host-', ''));
      const host = hosts.value?.find(h => h.id === hostId);
      if (host) {
        toTerminal(host);
      }
    } else if (key.startsWith('group-')) {
      const groupId = parseInt(key.replace('group-', ''));
      const group = groups.value?.find(g => g.id === groupId);
      if (group) {
        selectedGroup.value = group;
        selectedKeys.value = keys;
      }
    }
  } else {
    selectedGroup.value = null;
    selectedKeys.value = [];
  }
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

const osIconMap = new Map([
  [['cisco ios', 'cisco ios-xe', 'cisco ios-xr', 'cisco nx-os', 'cisco asa', 'cisco'], 'simple-icons:cisco'],
  [['huawei vrp', 'huawei', 'huawei ce', 'huawei ne', 'huawei s', 'huawei ar', 'huawei usg'], 'simple-icons:huawei'],
  [['fortinet', 'fortigate', 'fortios', 'fortimanager', 'fortianalyzer'], 'simple-icons:fortinet'],
  [['mikrotik', 'routeros'], 'simple-icons:mikrotik'],
  [['pfsense', 'opnsense'], 'simple-icons:pfsense'],
  [['juniper', 'junos', 'juniper ex', 'juniper mx', 'juniper srx'], 'simple-icons:juniper'],
  [['hp', 'hpe', 'hp procurve', 'hp comware', 'hpe aruba'], 'simple-icons:hp'],
  [['dell', 'dell emc', 'dell networking'], 'simple-icons:dell'],
  [['red hat enterprise', 'red hat', 'rhel', 'redhat'], 'simple-icons:redhat'],
  [['ubuntu', 'ubuntu server', 'xubuntu', 'kubuntu'], 'simple-icons:ubuntu'],
  [['centos stream', 'centos'], 'simple-icons:centos'],
  [['debian gnu', 'debian'], 'simple-icons:debian'],
  [['opensuse leap', 'opensuse tumbleweed', 'opensuse', 'suse'], 'simple-icons:opensuse'],
  [['fedora server', 'fedora workstation', 'fedora'], 'simple-icons:fedora'],
  [['alma linux', 'almalinux'], 'simple-icons:almalinux'],
  [['kali linux', 'kalilinux', 'kali'], 'simple-icons:kalilinux'],
  [['arch linux', 'archlinux', 'arch', 'manjaro'], 'simple-icons:archlinux'],
  [['rocky linux', 'rockylinux', 'rocky'], 'simple-icons:rockylinux'],
  [['alpine linux', 'alpinelinux', 'alpine'], 'simple-icons:alpinelinux'],
  [['gentoo linux', 'gentoo'], 'simple-icons:gentoo'],
  [['raspberry pi os', 'raspbian'], 'simple-icons:raspberrypi'],
  [['mint', 'linux mint'], 'simple-icons:linuxmint'],
  [['elementary os', 'elementary'], 'simple-icons:elementary'],
  [['zorin os', 'zorin'], 'simple-icons:zorinos'],
  [['pop!_os', 'pop os', 'popos'], 'simple-icons:popos'],
  [['linux'], 'simple-icons:linux'],
]);

const getDeviceIcon = (host: model.Host) => {
  const os = host.metadata?.os?.toLowerCase() || '';
  for (const [keys, icon] of osIconMap) {
    if (keys.some(key => os.includes(key))) {
      return icon;
    }
  }
  return 'ph:computer-tower';
};

const getConnectionCount = (host: model.Host) => {
  return connectionStore.connections.filter(c => c.host === host.host && !c.connectionError).length;
};

const getErrorConnectionCount = (host: model.Host) => {
  return connectionStore.connections.filter(c => c.host === host.host && c.connectionError).length;
};

const getProtocolIcon = (host: model.Host) => {
  const protocol = host.conn_protocol;

  switch (protocol) {
    case 0:
      return 'ph:terminal-duotone';
    case 1:
      return 'ph:broadcast-duotone';
    case 2:
      return 'ph:desktop-duotone';
    case 3:
      return 'ph:monitor-duotone';
    case 4:
      return 'ph:plug-duotone';
    default:
      return 'ph:ghost-duotone';
  }
};

const getProtocolName = (host: model.Host) => {
  const protocol = host.conn_protocol;
  switch (protocol) {
    case 0:
      return 'SSH';
    case 1:
      return 'Telnet';
    case 2:
      return 'RDP';
    case 3:
      return 'VNC';
    case 4:
      return 'Serial';
    default:
      return 'Unknown';
  }
};

const handleExpand = (keys: string[]) => {
  expandedKeys.value = keys;
};

onMounted(async () => {
  await fetchData();
  selectedKeys.value = ['all'];
});

const themeVars = useThemeVars();
</script>

<style lang="less" scoped>
.page-container {
  height: 100%;
  display: flex;
}

.sidebar {
  width: 260px;
  border-right: 1px solid v-bind('themeVars.borderColor');
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
    border-bottom: 1px solid v-bind('themeVars.borderColor');

    .title {
      font-size: 13px;
      font-weight: 600;
      color: v-bind('themeVars.textColorBase');
    }
  }

  .list-content {
    padding: 2px 4px 2px 4px;
    height: calc(100vh - 131px);
    overflow-y: auto;
    scroll-behavior: smooth;
    position: relative;

    &::-webkit-scrollbar {
      display: none;
    }
    -ms-overflow-style: none;
    scrollbar-width: none;

    :deep(.n-tree) {
      .n-tree-node {
        &:hover {
          background: v-bind('`${themeVars.primaryColor}10`');
        }

        &.n-tree-node--selected {
          background: v-bind('`${themeVars.primaryColor}20`');

          .n-tree-node-content {
            color: v-bind('themeVars.primaryColor');
          }
        }

        .n-tree-node-content {
          .n-tree-node-content__text {
            border-bottom: none !important;
            text-decoration: none !important;
          }
        }
      }

      .n-tree__empty {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
      }
    }
  }

  .list-footer {
    padding: 6px;
    border-top: 1px solid v-bind('themeVars.borderColor');
  }
}

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
        color: v-bind('themeVars.textColorBase');
        margin: 0;
      }
    }
  }
}

.hosts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  overflow-y: auto;
  padding-right: 4px;
}

.host-card {
  background: v-bind('themeVars.cardColor');
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  border: 1px solid v-bind('themeVars.borderColor');

  &:hover {
    border-color: v-bind('themeVars.primaryColor');

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
        background: v-bind('`${themeVars.primaryColor}20`');
        color: v-bind('themeVars.primaryColor');
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
          color: v-bind('themeVars.textColorBase');
          margin-bottom: 2px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .host-addr {
          font-size: 12px;
          color: v-bind('themeVars.textColor3');
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
      background: v-bind('`${themeVars.primaryColor}20`');
      width: 32px;
      height: 32px;
      border-radius: 6px !important;

      :deep(.edit-icon) {
        font-size: 16px;
      }

      &:hover {
        color: v-bind('themeVars.primaryColor');
        background: v-bind('`${themeVars.primaryColor}30`');
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
    gap: 8px;
    margin-top: auto;
    border-top: 1px dashed v-bind('themeVars.borderColor');

    .protocol-info {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      color: v-bind('themeVars.textColor3');
      flex-shrink: 0;

      :deep(svg) {
        font-size: 16px;
        color: v-bind('themeVars.primaryColor');
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

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100% - 72px);
  width: 100%;
}
</style>
