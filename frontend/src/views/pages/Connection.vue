<template>
  <div class="page-container">
    <div class="sidebar" ref="sidebarRef" :style="{ width: sidebarWidth + 'px' }">
      <div class="groups-list">
        <div class="list-header">
          <n-input
            v-model:value="searchText"
            size="small"
            clearable
            :placeholder="$t('connection.search')"
            :allow-input="value => !/\s/.test(value)"
          >
            <template #prefix>
              <icon icon="ph:magnifying-glass" />
            </template>
          </n-input>
          <n-divider vertical />
          <div class="header-right">
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button text size="large" @click="handleAddGroup">
                  <template #icon>
                    <icon icon="ph:folder-plus" />
                  </template>
                </n-button>
              </template>
              {{ $t('connection.add.group') }}
            </n-tooltip>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button text size="large" @click="handleAddConn">
                  <template #icon>
                    <icon icon="ph:plus" />
                  </template>
                </n-button>
              </template>
              {{ $t('connection.add.conn') }}
            </n-tooltip>
          </div>
        </div>
        <div class="list-content">
          <!-- 资产列表部分 -->
          <div class="section-header">
            <div class="section-title">资产列表</div>
            <n-button text size="tiny" @click="assetListCollapsed = !assetListCollapsed">
              <template #icon>
                <icon :icon="assetListCollapsed ? 'ph:caret-right' : 'ph:caret-down'" />
              </template>
            </n-button>
          </div>
          <div class="asset-list" v-show="!assetListCollapsed">
            <div
              v-for="conn in filteredAssets"
              :key="conn.id"
              class="asset-item"
              @click="handleSelectConn(conn)"
              @contextmenu="handleConnContextMenu($event, conn)"
            >
              <div class="asset-icon">
                <icon :icon="getDeviceIcon(conn)" />
              </div>
              <div class="asset-info">
                <div class="asset-name">{{ conn.label }}</div>
              </div>
            </div>
          </div>

          <!-- 分组列表部分 -->
          <div class="section-header">
            <div class="section-title">资产分组</div>
            <n-button text size="tiny" @click="groupListCollapsed = !groupListCollapsed">
              <template #icon>
                <icon :icon="groupListCollapsed ? 'ph:caret-right' : 'ph:caret-down'" />
              </template>
            </n-button>
          </div>
          <div class="group-list" v-show="!groupListCollapsed">
            <div
              v-for="group in filteredGroups"
              :key="group.id"
              class="group-item"
              :class="{ active: selectedGroup?.id === group.id }"
              @click="handleSelectGroup(group)"
              @contextmenu="handleGroupContextMenu($event, group)"
            >
              <div class="group-icon">
                <icon :icon="getGroupConnCount(group) > 0 ? 'ph:folders-duotone' : 'ph:folder-dashed-duotone'" />
              </div>
              <div class="group-info">
                <div class="group-name">{{ group.name }}</div>
              </div>
              <div class="group-count">{{ getGroupConnCount(group) }}</div>
            </div>
          </div>

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
      </div>
    </div>

    <div class="resize-handle" :style="{ left: sidebarWidth + 'px' }" @mousedown="startResize"></div>

    <div class="main-content">
      <div class="content-header">
        <div class="header-left">
          <h2>{{ selectedGroup ? selectedGroup.name : $t('sider.assets') }}</h2>
          <n-badge :value="filteredConns.length" show-zero type="success" />
        </div>
      </div>

      <div class="conns-grid" v-if="filteredConns.length > 0">
        <div v-for="conn in filteredConns" :key="conn.id" class="conn-card" @click="toTerminal(conn)">
          <div class="card-header">
            <div class="card-left">
              <div class="os-icon">
                <icon :icon="getDeviceIcon(conn)" />
              </div>
              <div class="card-info">
                <div class="conn-name">{{ conn.label }}</div>
                <div v-if="conn.connProtocol === enums.ConnProtocol.SSH" class="conn-info">
                  {{ conn.credential?.username }}@{{ conn.host }}
                </div>
                <div v-if="conn.connProtocol === enums.ConnProtocol.SERIAL" class="conn-info">
                  {{ conn.serialPort }}
                </div>
              </div>
            </div>
            <n-button circle text size="small" class="edit-btn" @click.stop="handleEditConn(conn)">
              <template #icon>
                <icon icon="ph:pencil-simple" class="edit-icon" />
              </template>
            </n-button>
          </div>

          <div class="card-footer">
            <div class="protocol-info">
              <icon :icon="getProtocolIcon(conn)" />
              <span>{{ conn.connProtocol }}</span>
            </div>
            <div class="connection-tags" style="margin-left: auto">
              <n-tooltip trigger="hover" v-if="getConnCount(conn) > 0">
                <template #trigger>
                  <div>
                    <n-tag size="tiny" type="success">
                      {{ getConnCount(conn) }}
                    </n-tag>
                  </div>
                </template>
                {{ $t('connection.connection.active') }}
              </n-tooltip>
              <n-tooltip trigger="hover" v-if="getErrorConnCount(conn) > 0">
                <template #trigger>
                  <div>
                    <n-tag size="tiny" type="error">
                      {{ getErrorConnCount(conn) }}
                    </n-tag>
                  </div>
                </template>
                {{ $t('connection.connection.disconnected') }}
              </n-tooltip>
            </div>
          </div>
        </div>
      </div>
      <div class="empty-state" v-else>
        <n-result
          status="404"
          :title="$t('connection.empty.title')"
          :description="
            selectedGroup
              ? $t('connection.empty.group_desc', { name: selectedGroup.name })
              : $t('connection.empty.all_desc')
          "
        />
      </div>
    </div>

    <connection-modal
      v-model:show="showConnModal"
      :is-edit="isEditConn"
      :connection="editConnection"
      @success="handleConnSuccess"
    />
    <group-modal
      v-model:show="showGroupModal"
      :is-edit="isEditGroup"
      :group="editGroup"
      @success="handleGroupSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import {
  NButton,
  NTag,
  NTooltip,
  NDropdown,
  NResult,
  NBadge,
  NInput,
  useMessage,
  useThemeVars,
  useDialog,
} from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import { ListGroup, DeleteGroup } from '@wailsApp/go/services/GroupSrv';
import { ListConnection, DeleteConnection } from '@wailsApp/go/services/ConnectionSrv';
import { enums, model } from '@wailsApp/go/models';
import { useConnectionStore } from '@/stores/connection';
import { useRouter } from 'vue-router';
import { h, ref, computed, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import ConnectionModal from '@/views/modals/ConnectionModal.vue';
import GroupModal from '@/views/modals/GroupModal.vue';

const router = useRouter();
const message = useMessage();
const connStore = useConnectionStore();
const { t } = useI18n();

const groups = ref<model.Group[]>();
const conns = ref<model.Connection[]>();
const selectedGroup = ref<model.Group | null>(null);

const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const dropdownOptions = ref<DropdownOption[]>([]);

const searchText = ref('');

const assetListCollapsed = ref(false);
const groupListCollapsed = ref(false);

// 对话框状态
const showConnModal = ref(false);
const showGroupModal = ref(false);
const isEditConn = ref(false);
const isEditGroup = ref(false);
const editConnection = ref<model.Connection | undefined>(undefined);
const editGroup = ref<model.Group | undefined>(undefined);

const currentContextNode = ref<any>(null);

const dialog = useDialog();

const updateDropdownOptions = (type: 'group' | 'conn') => {
  if (type === 'group') {
    dropdownOptions.value = [
      {
        label: t('connection.menu.edit_group'),
        key: 'edit-group',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      },
      {
        label: t('connection.menu.delete_group'),
        key: 'delete-group',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      },
    ];
  } else {
    dropdownOptions.value = [
      {
        label: t('connection.menu.edit_conn'),
        key: 'edit-conn',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      },
      {
        label: t('connection.menu.delete_conn'),
        key: 'delete-conn',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      },
    ];
  }
};

const handleDropdownSelect = async (key: string) => {
  showDropdown.value = false;
  if (!currentContextNode.value) return;

  const connId = parseInt(currentContextNode.value.key.replace('conn-', ''));

  switch (key) {
    case 'edit-group':
      handleEditGroup();
      break;
    case 'delete-group':
      await handleDeleteGroup();
      break;
    case 'edit-conn':
      const conn = conns.value?.find(h => h.id === connId);
      if (conn) handleEditConn(conn);
      break;
    case 'delete-conn':
      const resp = await DeleteConnection(connId);
      if (!resp.ok) {
        message.error(resp.msg);
        return;
      }
      message.success(t('message.deleteSuccess'));
      await fetchData();
      break;
  }
  currentContextNode.value = null;
};

const handleClickoutside = () => {
  showDropdown.value = false;
};

const filteredConns = computed(() => {
  if (!selectedGroup.value) return conns.value || [];
  return conns.value?.filter(conn => conn.groupID === selectedGroup.value?.id) || [];
});

const toTerminal = (conn: model.Connection) => {
  const connection = {
    id: Date.now(),
    connId: conn.id,
    label: `${conn.label} (${connStore.connections.filter(c => c.host === conn.host).length + 1})`,
    host: conn.host,
    username: conn.credential?.username || '',
    theme: conn.theme || 'Default',
  };
  connStore.addConnection(connection);
  router.push({ name: 'Terminal' });
};

const handleEditConn = (conn: model.Connection) => {
  isEditConn.value = true;
  editConnection.value = conn;
  showConnModal.value = true;
};

const handleConnSuccess = () => {
  fetchData();
};

const handleAddConn = () => {
  isEditConn.value = false;
  editConnection.value = undefined;
  showConnModal.value = true;
};

const handleAddGroup = () => {
  isEditGroup.value = false;
  editGroup.value = undefined;
  showGroupModal.value = true;
};

const handleEditGroup = () => {
  if (!currentContextNode.value) return;
  const groupId = parseInt(currentContextNode.value.key.replace('group-', ''));
  const group = groups.value?.find(g => g.id === groupId);
  if (group) {
    isEditGroup.value = true;
    editGroup.value = group;
    showGroupModal.value = true;
  }
};

const handleGroupSuccess = () => {
  fetchData();
};

const fetchGroups = async () => {
  const resp = await ListGroup();
  if (!resp.ok) {
    message.error(resp.msg);
  }
  return resp.data;
};

const fetchConns = async () => {
  const resp = await ListConnection();
  if (!resp.ok) {
    message.error(resp.msg);
  }
  return resp.data;
};

const fetchData = async () => {
  const [groupsData, connsData] = await Promise.all([fetchGroups(), fetchConns()]);
  groups.value = groupsData;
  conns.value = connsData;
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

const getDeviceIcon = (conn: model.Connection) => {
  const os = conn.metadata?.os?.toLowerCase() || '';
  for (const [keys, icon] of osIconMap) {
    if (keys.some(key => os.includes(key))) {
      return icon;
    }
  }
  return 'ph:computer-tower';
};

const getConnCount = (conn: model.Connection) => {
  return connStore.connections.filter(c => c.host === conn.host && !c.connectionError).length;
};

const getErrorConnCount = (conn: model.Connection) => {
  return connStore.connections.filter(c => c.host === conn.host && c.connectionError).length;
};

const getProtocolIcon = (conn: model.Connection) => {
  const protocol = conn.connProtocol;
  switch (protocol) {
    case enums.ConnProtocol.SSH:
      return 'ph:terminal-duotone';
    case enums.ConnProtocol.TELNET:
      return 'ph:broadcast-duotone';
    case enums.ConnProtocol.RDP:
      return 'ph:desktop-duotone';
    case enums.ConnProtocol.VNC:
      return 'ph:monitor-duotone';
    case enums.ConnProtocol.SERIAL:
      return 'ph:plug-duotone';
    default:
      return 'ph:gconn-duotone';
  }
};

const themeVars = useThemeVars();

const sidebarRef = ref<HTMLElement | null>(null);
const sidebarWidth = ref(Number(localStorage.getItem('sidebarWidth')) || 260);
const minWidth = 260;
const maxWidth = 380;
const isResizing = ref(false);

const startResize = (e: MouseEvent) => {
  e.preventDefault();
  isResizing.value = true;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';

  const startX = e.clientX;
  const startWidth = sidebarWidth.value;

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing.value) return;

    const delta = e.clientX - startX;
    const newWidth = Math.min(Math.max(startWidth + delta, minWidth), maxWidth);
    sidebarWidth.value = newWidth;
    localStorage.setItem('sidebarWidth', newWidth.toString());
  };

  const handleMouseUp = () => {
    isResizing.value = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
  };

  document.addEventListener('mousemove', handleMouseMove);
  document.addEventListener('mouseup', handleMouseUp);
};

onMounted(async () => {
  await fetchData();

  window.addEventListener('sidebar-width-change', ((e: CustomEvent) => {
    sidebarWidth.value = e.detail;
  }) as EventListener);
});

onUnmounted(() => {
  window.removeEventListener('sidebar-width-change', ((e: CustomEvent) => {
    sidebarWidth.value = e.detail;
  }) as EventListener);
});

const filteredAssets = computed(() => {
  if (!searchText.value) return conns.value || [];
  return (
    conns.value?.filter(
      conn =>
        conn.label.toLowerCase().includes(searchText.value.toLowerCase()) ||
        conn.host?.toLowerCase().includes(searchText.value.toLowerCase()),
    ) || []
  );
});

const filteredGroups = computed(() => {
  if (!searchText.value) return groups.value || [];
  return groups.value?.filter(group => group.name.toLowerCase().includes(searchText.value.toLowerCase())) || [];
});

const getGroupConnCount = (group: model.Group) => {
  return conns.value?.filter(conn => conn.groupID === group.id).length || 0;
};

const handleSelectConn = (conn: model.Connection) => {
  toTerminal(conn);
};

const handleSelectGroup = (group: model.Group | null) => {
  if (selectedGroup.value?.id === group?.id) {
    selectedGroup.value = null;
  } else {
    selectedGroup.value = group;
  }
};

const handleConnContextMenu = (event: MouseEvent, conn: model.Connection) => {
  event.preventDefault();
  dropdownX.value = event.clientX;
  dropdownY.value = event.clientY;
  showDropdown.value = true;
  currentContextNode.value = { key: `conn-${conn.id}` };
  updateDropdownOptions('conn');
};

const handleGroupContextMenu = (event: MouseEvent, group: model.Group) => {
  event.preventDefault();
  dropdownX.value = event.clientX;
  dropdownY.value = event.clientY;
  showDropdown.value = true;
  currentContextNode.value = { key: `group-${group.id}` };
  updateDropdownOptions('group');
};

const handleDeleteGroup = async () => {
  if (!currentContextNode.value) return;
  const groupId = parseInt(currentContextNode.value.key.replace('group-', ''));
  const group = groups.value?.find(g => g.id === groupId);
  if (!group) return;

  dialog.warning({
    title: t('connection.delete.group.title'),
    content: t('connection.delete.group.content', { name: group.name }),
    positiveText: t('connection.delete.group.confirm'),
    negativeText: t('connection.delete.group.cancel'),
    onPositiveClick: async () => {
      const resp = await DeleteGroup(groupId);
      if (!resp.ok) {
        message.error(resp.msg);
        return;
      }
      message.success(t('message.deleteSuccess'));
      await fetchData();
    },
  });
};
</script>

<style lang="less" scoped>
.page-container {
  height: 100%;
  display: flex;
  position: relative;
}

.sidebar {
  border-right: none;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.groups-list {
  flex: 1;
  display: flex;
  flex-direction: column;

  .list-header {
    padding: 6px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid v-bind('themeVars.borderColor');

    .title {
      font-size: 13px;
      font-weight: 600;
      color: v-bind('themeVars.textColorBase');
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 4px;
      margin-left: auto;
    }
  }

  .list-content {
    padding: 2px 4px;
    height: calc(100vh - 80px);
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
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  padding: 24px;
  position: relative;

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

.conns-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  overflow-y: auto;
  padding-right: 4px;
}

.conn-card {
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

        .conn-name {
          font-size: 14px;
          font-weight: 600;
          color: v-bind('themeVars.textColorBase');
          margin-bottom: 2px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .conn-info {
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
      margin-left: auto;
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

.resize-handle {
  width: 4px;
  cursor: col-resize;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 10;
  margin-left: -2px;
  border-right: 1px solid v-bind('themeVars.borderColor');
  background: transparent;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 8px 4px;

  .section-title {
    padding: 0;
    font-size: 12px;
    font-weight: 600;
    color: v-bind('themeVars.textColor3');
  }

  :deep(.n-button) {
    width: 16px;
    height: 16px;
    color: v-bind('themeVars.textColor3');

    .n-button__icon {
      font-size: 14px;
    }
  }
}

.asset-list,
.group-list {
  .asset-item,
  .group-item {
    display: flex;
    align-items: center;
    padding: 4px 8px;
    margin: 0 4px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      background: v-bind('`${themeVars.primaryColor}10`');
    }

    .asset-icon,
    .group-icon {
      width: 20px;
      height: 20px;
      border-radius: 50%;
      background: v-bind('`${themeVars.primaryColor}20`');
      color: v-bind('themeVars.primaryColor');
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 16px;
      flex-shrink: 0;
    }

    .asset-info,
    .group-info {
      flex: 1;
      min-width: 0;
      margin-left: 6px;

      .asset-name,
      .group-name {
        font-size: 14px;
        color: v-bind('themeVars.textColorBase');
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .group-count {
      font-size: 12px;
      color: v-bind('themeVars.textColor3');
      margin-left: 6px;
      min-width: 16px;
      text-align: right;
    }
  }

  .group-item {
    &.active {
      background: v-bind('`${themeVars.primaryColor}20`');

      .group-name {
        color: v-bind('themeVars.primaryColor');
      }
    }
  }
}
</style>
