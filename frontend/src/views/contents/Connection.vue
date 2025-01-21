<template>
  <div class="page-container">
    <div class="sidebar" ref="sidebarRef" :style="{ width: sidebarWidth + 'px' }">
      <div class="groups-list">
        <div class="list-header">
          <n-input v-model:value="searchText" size="small" clearable :placeholder="$t('connection.search')">
            <template #prefix>
              <icon icon="ph:magnifying-glass" />
            </template>
          </n-input>
          <n-divider vertical />
          <div class="header-right">
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button text size="large" @click="dialogStore.openGroupDialog">
                  <template #icon>
                    <icon icon="ph:folder-plus" />
                  </template>
                </n-button>
              </template>
              {{ $t('connection.add.group') }}
            </n-tooltip>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button text size="large" @click="dialogStore.openConnDialog()">
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
      </div>
    </div>

    <div class="resize-handle" :style="{ left: sidebarWidth + 'px' }" @mousedown="startResize"></div>

    <div class="main-content">
      <div class="content-header">
        <div class="header-left">
          <h2>{{ selectedGroup ? selectedGroup.name : $t('sider.assets') }}</h2>
          <n-badge :value="filteredConns.length" show-zero type="success" />
        </div>
        <div class="header-right">
          <n-button v-if="selectedGroup" text @click="handleSelect([])">
            <template #icon>
              <icon icon="ph:arrow-left" />
            </template>
            {{ $t('connection.back') }}
          </n-button>
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
            <n-button circle text size="small" class="edit-btn" @click.stop="handleEditConn($event, conn)">
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
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NTag, NTooltip, NTree, NDropdown, NResult, NBadge, NInput, useMessage, useThemeVars } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import { useDialogStore } from '@/stores/dialog';
import { ListGroup } from '@wailsApp/go/services/GroupSrv';
import { ListConnection, DeleteConnection } from '@wailsApp/go/services/ConnectionSrv';
import { enums, model } from '@wailsApp/go/models';
import { useConnectionStore } from '@/stores/connection';
import { useRouter } from 'vue-router';
import { h, ref, computed, onMounted, watch, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';

const dialogStore = useDialogStore();
const router = useRouter();
const message = useMessage();
const connStore = useConnectionStore();
const { t } = useI18n();

const groups = ref<model.Group[]>();
const conns = ref<model.Connection[]>();
const selectedGroup = ref<model.Group | null>(null);

const selectedKeys = ref<string[]>([]);
const expandedKeys = ref<string[]>([]);

const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const dropdownOptions = ref<DropdownOption[]>([]);

const searchText = ref('');

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

const currentContextNode = ref<any>(null);

const nodeProps = ({ option }: { option: any }) => {
  return {
    onClick() {
      if (option.key.startsWith('group-')) {
        const groupId = parseInt(option.key.replace('group-', ''));
        const group = groups.value?.find(g => g.id === groupId);
        if (group) {
          const groupConns = conns.value?.filter(conn => conn.groupID === group.id) || [];
          if (groupConns.length > 0) {
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
      } else if (option.key.startsWith('conn-')) {
        updateDropdownOptions('conn');
      }

      e.preventDefault();
    },
  };
};

const handleDropdownSelect = async (key: string) => {
  showDropdown.value = false;
  if (!currentContextNode.value) return;

  const connId = parseInt(currentContextNode.value.key.replace('conn-', ''));

  switch (key) {
    case 'edit-group':
      dialogStore.openGroupDialog();
      break;
    case 'delete-group':
      // TODO: 实现删除分组功能
      break;
    case 'edit-conn':
      const conn = conns.value?.find(h => h.id === connId);
      if (conn) dialogStore.openConnDialog(true, conn);
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

const treeData = computed(() => {
  if (!groups.value) return [];

  const treeNodes = [];

  treeNodes.push(
    ...groups.value.map(group => {
      const groupConns = conns.value?.filter(conn => conn.groupID === group.id) || [];
      const isEmpty = groupConns.length === 0;
      return {
        key: `group-${group.id}`,
        label: group.name,
        children: isEmpty
          ? undefined
          : groupConns.map(conn => ({
              key: `conn-${conn.id}`,
              label: conn.label,
              isLeaf: true,
              prefix: () =>
                h(Icon, {
                  icon: getDeviceIcon(conn),
                  style: {
                    fontSize: '1.2rem',
                  },
                }),
            })),
        isLeaf: isEmpty,
        prefix: () => {
          const expanded = expandedKeys.value.includes(`group-${group.id}`);
          return h(Icon, {
            icon: isEmpty ? 'ph:folder-dashed-duotone' : expanded ? 'ph:folder-open-duotone' : 'ph:folders-duotone',
            style: { fontSize: '1.2rem' },
          });
        },
      };
    }),
  );

  const ungroupedConns = conns.value?.filter(conn => !conn.groupID) || [];
  if (ungroupedConns.length > 0) {
    treeNodes.push(
      ...ungroupedConns.map(conn => ({
        key: `conn-${conn.id}`,
        label: conn.label,
        isLeaf: true,
        prefix: () =>
          h(Icon, {
            icon: getDeviceIcon(conn),
            style: {
              fontSize: '1.2rem',
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
    if (key.startsWith('conn-')) {
      const connId = parseInt(key.replace('conn-', ''));
      const conn = conns.value?.find(h => h.id === connId);
      if (conn) {
        toTerminal(conn);
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

const filteredConns = computed(() => {
  if (!selectedGroup.value) return conns.value || [];
  return conns.value?.filter(conn => conn.groupID === selectedGroup.value?.id) || [];
});

watch([() => dialogStore.connDialogVisible, () => dialogStore.groupDialogVisible], ([connVisible, groupVisible]) => {
  if (!connVisible || !groupVisible) {
    fetchData();
  }
});

const toTerminal = (conn: model.Connection) => {
  const connection = {
    id: Date.now(),
    connId: conn.id,
    label: `${conn.label} (${connStore.connections.filter(c => c.host === conn.host).length + 1})`,
    host: conn.host,
    username: conn.credential?.username || '',
  };
  connStore.addConnection(connection);
  router.push({ name: 'Terminal' });
};

const handleEditConn = (event: MouseEvent, conn: model.Connection) => {
  event.preventDefault();
  dialogStore.openConnDialog(true, conn);
};

const handleEditGroup = (event: MouseEvent) => {
  event.preventDefault();
  dialogStore.openGroupDialog();
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

const handleExpand = (keys: string[]) => {
  expandedKeys.value = keys;
};

onMounted(async () => {
  await fetchData();
});

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

onMounted(() => {
  window.addEventListener('sidebar-width-change', ((e: CustomEvent) => {
    sidebarWidth.value = e.detail;
  }) as EventListener);
});

onUnmounted(() => {
  window.removeEventListener('sidebar-width-change', ((e: CustomEvent) => {
    sidebarWidth.value = e.detail;
  }) as EventListener);
});
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
    }
  }

  .list-content {
    padding: 2px 4px 2px 4px;
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

    .header-right {
      :deep(.n-button) {
        font-size: 13px;
        color: v-bind('themeVars.textColor3');
        transition: color 0.2s ease;

        &:hover {
          color: v-bind('themeVars.primaryColor');
        }
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
</style>
