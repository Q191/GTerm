<template>
  <div class="file-transfer-container">
    <div class="integrated-header">
      <div class="left-section">
        <div class="connection-dropdown">
          <n-select
            v-model:value="selectedConnectionId"
            :options="connectionOptions"
            :placeholder="t('frontend.file_transfer.select_connection')"
            :disabled="connecting"
            class="connection-select"
            @update:value="handleConnectionChange"
          />

          <n-tooltip placement="bottom">
            <template #trigger>
              <n-button
                quaternary
                circle
                :disabled="!selectedConnectionId"
                @click="isConnected ? disconnectSftp() : connectSftp()"
                :loading="connecting"
                class="connection-action-btn"
              >
                <template #icon>
                  <n-icon>
                    <icon :icon="isConnected ? 'ph:link-break-bold' : 'ph:link-bold'" />
                  </n-icon>
                </template>
              </n-button>
            </template>
            {{ isConnected ? t('frontend.file_transfer.disconnect') : t('frontend.file_transfer.connect') }}
          </n-tooltip>
        </div>

        <n-divider v-if="isConnected" vertical style="height: 20px; margin: 0 8px" />

        <div v-if="isConnected" class="path-navigator">
          <div class="path-control">
            <n-tooltip placement="bottom">
              <template #trigger>
                <n-button quaternary circle class="path-nav-btn" @click="navigateHome">
                  <template #icon>
                    <n-icon><icon icon="ph:house-bold" /></n-icon>
                  </template>
                </n-button>
              </template>
              {{ t('frontend.file_transfer.navigate_home') }}
            </n-tooltip>

            <n-tooltip placement="bottom" :disabled="!canNavigateUp">
              <template #trigger>
                <n-button quaternary circle class="path-nav-btn" @click="navigateUp" :disabled="!canNavigateUp">
                  <template #icon>
                    <n-icon><icon icon="ph:arrow-up-bold" /></n-icon>
                  </template>
                </n-button>
              </template>
              {{ t('frontend.file_transfer.navigate_up') }}
            </n-tooltip>

            <n-input class="path-input" v-model:value="remotePathInput" @keydown.enter="navigateToPath">
              <template #prefix>
                <n-icon><icon icon="ph:folder-open-bold" /></n-icon>
              </template>
            </n-input>

            <n-tooltip placement="bottom">
              <template #trigger>
                <n-button quaternary class="path-nav-go-btn" @click="navigateToPath">
                  <template #icon>
                    <n-icon><icon icon="ph:arrow-circle-right-bold" /></n-icon>
                  </template>
                  {{ t('frontend.file_transfer.goto') }}
                </n-button>
              </template>
              {{ t('frontend.file_transfer.navigate_to_path') }}
            </n-tooltip>
          </div>
        </div>
      </div>

      <div v-if="isConnected" class="right-section">
        <n-button-group size="small">
          <n-tooltip placement="bottom">
            <template #trigger>
              <n-button quaternary @click="refreshRemoteFiles">
                <template #icon>
                  <n-icon><icon icon="ph:arrows-clockwise-bold" /></n-icon>
                </template>
              </n-button>
            </template>
            {{ t('frontend.file_transfer.refresh') }}
          </n-tooltip>

          <n-tooltip placement="bottom">
            <template #trigger>
              <n-button quaternary @click="uploadToRemote">
                <template #icon>
                  <n-icon><icon icon="ph:arrow-up-bold" /></n-icon>
                </template>
              </n-button>
            </template>
            {{ t('frontend.file_transfer.upload') }}
          </n-tooltip>

          <n-tooltip placement="bottom" :disabled="!remoteSelectedFiles.length">
            <template #trigger>
              <n-button quaternary @click="downloadFromRemote" :disabled="!remoteSelectedFiles.length">
                <template #icon>
                  <n-icon><icon icon="ph:arrow-down-bold" /></n-icon>
                </template>
              </n-button>
            </template>
            {{ t('frontend.file_transfer.download') }}
          </n-tooltip>
        </n-button-group>

        <n-tooltip placement="bottom">
          <template #trigger>
            <n-popover
              trigger="click"
              placement="bottom"
              :width="380"
              :show="showTasksPopover"
              @update:show="showTasksPopover = $event"
            >
              <template #trigger>
                <n-button quaternary circle @click="toggleTasksPopover">
                  <template #icon>
                    <n-icon>
                      <icon icon="ph:arrows-down-up-bold" />
                    </n-icon>
                    <n-badge type="info" :value="activeTasks.length" :show="activeTasks.length > 0" :max="99" />
                  </template>
                </n-button>
              </template>

              <div class="tasks-popover-content">
                <div v-if="transferTasks.length === 0" class="no-tasks">
                  <span>{{ $t('frontend.file_transfer.no_tasks') }}</span>
                </div>
                <div v-else class="tasks-list">
                  <div
                    v-for="task in transferTasks"
                    :key="task.id"
                    class="task-item"
                    :class="{ 'is-completed': task.status === FileTransferTaskState.COMPLETED }"
                  >
                    <div class="task-info">
                      <div class="task-title">
                        <n-icon class="task-type-icon">
                          <icon :icon="task.isUpload ? 'ph:arrow-up-bold' : 'ph:arrow-down-bold'" />
                        </n-icon>
                        <span class="task-name">{{ getTaskFilename(task) }}</span>
                      </div>
                      <div class="task-details">
                        <span class="task-progress-text">
                          {{ formatFileSize(task.transferred) }} / {{ formatFileSize(task.size) }}
                        </span>
                        <span class="task-status-text">{{ getTaskStatusText(task) }}</span>
                      </div>
                    </div>
                    <n-progress
                      type="line"
                      :percentage="getTaskPercentage(task)"
                      :processing="task.status === FileTransferTaskState.PROGRESS"
                      :status="getTaskProgressStatus(task)"
                      :show-indicator="false"
                      :height="4"
                      class="task-progress-bar"
                    />
                  </div>
                </div>
                <div class="tasks-actions">
                  <n-button v-if="hasCompletedTasks" size="small" @click="clearCompletedTasks">
                    {{ $t('frontend.file_transfer.clear_completed') }}
                  </n-button>
                  <n-button v-if="transferTasks.length > 0" size="small" @click="clearAllTasks">
                    {{ $t('frontend.file_transfer.clear_all') }}
                  </n-button>
                </div>
              </div>
            </n-popover>
          </template>
          {{ t('frontend.file_transfer.transfer_tasks') }}
        </n-tooltip>
      </div>
    </div>

    <div v-if="!isConnected" class="empty-state">
      <div class="empty-content">
        <n-icon size="64" class="empty-icon">
          <icon icon="ph:cloud-slash-bold" />
        </n-icon>
        <div class="empty-text">{{ $t('frontend.file_transfer.no_connection') }}</div>
      </div>
    </div>

    <div v-else class="file-browser">
      <div class="file-list-container">
        <div class="file-list-header">
          <div class="file-grid-cell checkbox-cell">
            <n-checkbox
              :checked="allRemoteSelected"
              :indeterminate="someRemoteSelected && !allRemoteSelected"
              @update:checked="handleCheckAll"
            />
          </div>
          <div class="file-grid-cell name-cell">{{ t('frontend.file_transfer.name') }}</div>
          <div class="file-grid-cell size-cell">{{ t('frontend.file_transfer.size') }}</div>
          <div class="file-grid-cell time-cell">{{ t('frontend.file_transfer.modified') }}</div>
          <div class="file-grid-cell owner-cell">{{ t('frontend.file_transfer.owner') }}</div>
          <div class="file-grid-cell group-cell">{{ t('frontend.file_transfer.group') }}</div>
          <div class="file-grid-cell perm-cell">{{ t('frontend.file_transfer.permissions') }}</div>
        </div>

        <div class="file-list-body">
          <n-scrollbar>
            <div class="file-grid">
              <div
                v-for="file in remoteFiles"
                :key="file.name"
                class="file-grid-row"
                :class="{ 'is-selected': remoteSelectedFiles.includes(file.name) }"
                @click="handleRowClick(file)"
              >
                <div class="file-grid-cell checkbox-cell" @click.stop>
                  <n-checkbox
                    :checked="remoteSelectedFiles.includes(file.name)"
                    @update:checked="checked => handleFileSelect(file, checked)"
                  />
                </div>
                <div class="file-grid-cell name-cell">
                  <div class="file-name">
                    <n-icon class="file-icon" :class="{ 'is-folder': file.isDir }">
                      <icon :icon="file.isDir ? 'ph:folder-simple-bold' : 'ph:file-bold'" />
                    </n-icon>
                    <span class="file-label">{{ file.name }}</span>
                  </div>
                </div>
                <div class="file-grid-cell size-cell">
                  {{ file.isDir ? '-' : formatFileSize(file.size) }}
                </div>
                <div class="file-grid-cell time-cell">
                  {{ formatDate(file.modTime) }}
                </div>
                <div class="file-grid-cell owner-cell">
                  {{ file.owner }}
                </div>
                <div class="file-grid-cell group-cell">
                  {{ file.group }}
                </div>
                <div class="file-grid-cell perm-cell">
                  {{ file.permissions }}
                </div>
              </div>
            </div>
          </n-scrollbar>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { Icon } from '@iconify/vue';
import { useI18n } from 'vue-i18n';
import {
  NSelect,
  NButton,
  NIcon,
  NProgress,
  NScrollbar,
  NCheckbox,
  NButtonGroup,
  NInput,
  NDivider,
  NPopover,
  NBadge,
  useThemeVars,
  NTooltip,
} from 'naive-ui';
import {
  ConnectSFTP,
  DisconnectSFTP,
  ListRemoteFiles,
  UploadFiles,
  DownloadFiles,
  SelectDownloadDirectory,
  SelectUploadFiles,
} from '@wailsApp/go/services/FileTransferSrv';
import { ListConnection } from '@wailsApp/go/services/ConnectionSrv';
import { EventsOn, EventsOff } from '@wailsApp/runtime';
import { useCall } from '@/utils/call';
import { v4 as uuidv4 } from 'uuid';
import { enums, types } from '@wailsApp/go/models';

const { FileTransferTaskState } = enums;
const { t } = useI18n();
const { call } = useCall();
const themeVars = useThemeVars();

const selectedConnectionId = ref<number | null>(null);
const connecting = ref(false);
const isConnected = ref(false);

const remoteAbsolutePath = ref('');
const remotePathInput = ref('');
const remoteFiles = ref<types.FileTransferItemInfo[]>([]);
const remoteSelectedFiles = ref<string[]>([]);
const showTasksPopover = ref(false);

const transferTasks = ref<types.FileTransferTask[]>([]);
const connections = ref<any[]>([]);

const pathHistory = ref<string[]>([]);
const currentHistoryIndex = ref(-1);

const connectionOptions = computed(() => {
  if (!connections.value || !Array.isArray(connections.value)) {
    return [];
  }
  return connections.value.map(conn => ({
    label: conn.label,
    value: conn.id,
  }));
});

const fetchConns = async () => {
  const result = await call(ListConnection);
  connections.value = result.data || [];
};

const allRemoteSelected = computed(() => {
  return remoteFiles.value.length > 0 && remoteSelectedFiles.value.length === remoteFiles.value.length;
});

const someRemoteSelected = computed(() => {
  return remoteSelectedFiles.value.length > 0;
});

const activeTasks = computed(() => {
  return transferTasks.value.filter(
    task => task.status === FileTransferTaskState.PENDING || task.status === FileTransferTaskState.PROGRESS,
  );
});

const hasCompletedTasks = computed(() => {
  return transferTasks.value.some(
    task => task.status === FileTransferTaskState.COMPLETED || FileTransferTaskState.ERROR,
  );
});

const toggleTasksPopover = () => {
  showTasksPopover.value = !showTasksPopover.value;
};

const getTaskFilename = (task: types.FileTransferTask): string => {
  const pathParts = task.source.split('/');
  return pathParts[pathParts.length - 1];
};

const getTaskPercentage = (task: types.FileTransferTask): number => {
  if (task.size === 0) return 0;
  return Math.min(Math.floor((task.transferred / task.size) * 100), 100);
};

const getTaskProgressStatus = (task: types.FileTransferTask): 'success' | 'error' | 'warning' | undefined => {
  switch (task.status) {
    case FileTransferTaskState.COMPLETED:
      return 'success';
    case FileTransferTaskState.ERROR:
      return 'error';
    case FileTransferTaskState.PENDING:
      return 'warning';
    default:
      return undefined;
  }
};

const getTaskStatusText = (task: types.FileTransferTask): string => {
  switch (task.status) {
    case FileTransferTaskState.COMPLETED:
      return t('frontend.file_transfer.completed');
    case FileTransferTaskState.ERROR:
      return task.error || t('frontend.file_transfer.failed');
    case FileTransferTaskState.PENDING:
      return t('frontend.file_transfer.pending');
    case FileTransferTaskState.PROGRESS:
      return `${getTaskPercentage(task)}%`;
    default:
      return '';
  }
};

const clearCompletedTasks = () => {
  transferTasks.value = transferTasks.value.filter(
    task => task.status === FileTransferTaskState.PENDING || task.status === FileTransferTaskState.PROGRESS,
  );
};

const clearAllTasks = () => {
  transferTasks.value = [];
};

const connectSftp = async () => {
  if (!selectedConnectionId.value) return;

  connecting.value = true;
  const resp = await call(ConnectSFTP, {
    args: [selectedConnectionId.value],
  });

  if (resp.ok) {
    isConnected.value = true;
    pathHistory.value = [];
    currentHistoryIndex.value = -1;
    await navigateTo('', false);
  }
  connecting.value = false;
};

const refreshRemoteFiles = async () => {
  const resp = await call(ListRemoteFiles, {
    args: [remoteAbsolutePath.value],
  });

  if (resp.ok) {
    const response = resp.data as types.FileList;
    remoteFiles.value = response.files || [];
    remoteAbsolutePath.value = response.absolutePath || '';
    remotePathInput.value = remoteAbsolutePath.value;
    remoteSelectedFiles.value = [];
  }
};

const canNavigateUp = computed(() => {
  return remoteAbsolutePath.value !== '' && remoteAbsolutePath.value !== '/';
});

const addToHistory = (path: string) => {
  if (pathHistory.value[pathHistory.value.length - 1] !== path) {
    pathHistory.value.push(path);
    currentHistoryIndex.value = pathHistory.value.length - 1;
  }
};

const navigateTo = async (path: string, addHistory = true) => {
  remoteAbsolutePath.value = path;

  if (addHistory) {
    addToHistory(path);
  }

  await refreshRemoteFiles();
};

const navigateHome = () => {
  navigateTo('');
};

const navigateUp = () => {
  if (!canNavigateUp.value) return;

  const path = remoteAbsolutePath.value;
  const lastSlashIndex = path.lastIndexOf('/');

  if (lastSlashIndex <= 0) {
    navigateTo('/');
  } else {
    let parentPath = path.substring(0, lastSlashIndex);
    if (parentPath === '') {
      parentPath = '/';
    }
    navigateTo(parentPath);
  }
};

const navigateToPath = () => {
  if (!remotePathInput.value) return;
  navigateTo(remotePathInput.value);
};

const handleFileClick = (file: types.FileTransferItemInfo) => {
  if (file.isDir) {
    const path = remoteAbsolutePath.value;
    let newPath;

    if (path === '' || path === '/') {
      newPath = `/${file.name}`;
    } else {
      newPath = `${path}/${file.name}`;
    }
    navigateTo(newPath);
  }
};

const handleRowClick = (file: types.FileTransferItemInfo) => {
  if (file.isDir) {
    handleFileClick(file);
  } else {
    handleFileSelect(file, !remoteSelectedFiles.value.includes(file.name));
  }
};

const handleFileSelect = (file: types.FileTransferItemInfo, checked: boolean) => {
  if (checked) {
    if (!remoteSelectedFiles.value.includes(file.name)) {
      remoteSelectedFiles.value.push(file.name);
    }
  } else {
    const index = remoteSelectedFiles.value.indexOf(file.name);
    if (index !== -1) {
      remoteSelectedFiles.value.splice(index, 1);
    }
  }
};

const handleCheckAll = (checked: boolean) => {
  if (checked) {
    remoteSelectedFiles.value = remoteFiles.value.map(file => file.name);
  } else {
    remoteSelectedFiles.value = [];
  }
};

const disconnectSftp = async () => {
  const resp = await call(DisconnectSFTP);
  if (resp.ok) {
    isConnected.value = false;
    remoteFiles.value = [];
  }
};

const uploadToRemote = async () => {
  const filesResponse = await call(SelectUploadFiles, {
    args: [t('frontend.file_transfer.select_files_title')],
  });

  if (!filesResponse.ok) {
    transferTasks.value = transferTasks.value.map(task => {
      if (task.isUpload && task.status === FileTransferTaskState.PENDING) {
        return { ...task, status: FileTransferTaskState.ERROR, error: filesResponse.msg };
      }
      return task;
    });

    return;
  }

  const sources = filesResponse.data || [];
  if (sources.length === 0) {
    return;
  }

  for (const source of sources) {
    const filename = source.split('/').pop() || t('frontend.file_transfer.unknown_file');
    const fileSize = 0;

    const task: types.FileTransferTask = {
      id: uuidv4(),
      source: source,
      destination: `${remoteAbsolutePath.value === '' ? '/' : remoteAbsolutePath.value}/${filename}`,
      size: fileSize,
      transferred: 0,
      isUpload: true,
      status: FileTransferTaskState.PENDING,
    };

    transferTasks.value.push(task);
  }

  if (!showTasksPopover.value && sources.length > 0) {
    showTasksPopover.value = true;
  }

  const response = await call(UploadFiles, {
    args: [sources, remoteAbsolutePath.value],
  });

  if (response.ok) {
    await refreshRemoteFiles();
  } else {
    transferTasks.value = transferTasks.value.map(task => {
      if (task.isUpload && task.status === FileTransferTaskState.PENDING) {
        return { ...task, status: FileTransferTaskState.ERROR, error: response.msg };
      }
      return task;
    });
  }
};

const downloadFromRemote = async () => {
  if (!remoteSelectedFiles.value.length) return;

  const files = remoteFiles.value.filter(f => remoteSelectedFiles.value.includes(f.name) && !f.isDir);
  const sources = files.map(f => `${remoteAbsolutePath.value === '' ? '/' : remoteAbsolutePath.value}/${f.name}`);

  const dirResponse = await call(SelectDownloadDirectory, {
    args: [t('frontend.file_transfer.select_directory_title')],
  });

  if (!dirResponse.ok) {
    transferTasks.value = transferTasks.value.map(task => {
      if (!task.isUpload && task.status === FileTransferTaskState.PENDING) {
        return { ...task, status: FileTransferTaskState.ERROR, error: dirResponse.msg };
      }
      return task;
    });

    return;
  }

  const downloadPath = dirResponse.data;

  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    const source = sources[i];
    const destination = `${downloadPath}/${file.name}`;

    const task: types.FileTransferTask = {
      id: uuidv4(),
      source: source,
      destination: destination,
      size: file.size,
      transferred: 0,
      isUpload: false,
      status: FileTransferTaskState.PENDING,
    };

    transferTasks.value.push(task);
  }

  if (!showTasksPopover.value && files.length > 0) {
    showTasksPopover.value = true;
  }

  const response = await call(DownloadFiles, {
    args: [sources, downloadPath],
  });

  if (!response.ok) {
    transferTasks.value = transferTasks.value.map(task => {
      if (!task.isUpload && task.status === FileTransferTaskState.PENDING) {
        return { ...task, status: FileTransferTaskState.ERROR, error: response.msg };
      }
      return task;
    });
  }
};

const formatFileSize = (size: number): string => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
  if (size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`;
  return `${(size / 1024 / 1024 / 1024).toFixed(2)} GB`;
};

const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleString();
};

const handleTransferProgress = (data: any) => {
  const taskType = data.type === 'upload';
  const fileName = data.fileName || '';
  let task = transferTasks.value.find(t => t.isUpload === taskType && t.source.endsWith(fileName));
  if (!task && fileName) {
    const newTask: types.FileTransferTask = {
      id: uuidv4(),
      source: taskType ? fileName : `${remoteAbsolutePath.value}/${fileName}`,
      destination: taskType ? `${remoteAbsolutePath.value}/${fileName}` : fileName,
      size: data.total || 0,
      transferred: data.transferred || 0,
      isUpload: taskType,
      status: FileTransferTaskState.PROGRESS,
    };

    transferTasks.value.push(newTask);
    task = newTask;
  }

  if (task) {
    task.status = FileTransferTaskState.PROGRESS;
    task.transferred = data.transferred || 0;
    task.size = data.total || task.size;

    if (data.progress >= 100) {
      setTimeout(() => {
        task!.status = FileTransferTaskState.COMPLETED;
      }, 500);
    }
  }
};

const handleConnectionChange = () => {
  if (isConnected.value) {
    disconnectSftp();
  }
};

onMounted(() => {
  fetchConns().then(() => {
    if (connections.value.length > 0) {
      selectedConnectionId.value = connections.value[0].id;
    }
  });

  EventsOn('transfer:progress', handleTransferProgress);
});

onBeforeUnmount(() => {
  if (isConnected.value) {
    disconnectSftp();
  }

  EventsOff('transfer:progress');
});
</script>

<style lang="less" scoped>
.file-transfer-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 0;
  background-color: v-bind('themeVars.bodyColor');
}

.integrated-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: v-bind('themeVars.cardColor');
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  height: 48px;
  flex-shrink: 0;
}

.left-section {
  display: flex;
  align-items: center;
  flex: 1;
  overflow: hidden;
}

.right-section {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.connection-dropdown {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.connection-select {
  width: 180px;
}

.connection-action-btn {
  margin-left: 4px;
}

.path-navigator {
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.path-control {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.path-nav-btn {
  margin: 0;
  flex-shrink: 0;
}

.path-input {
  flex: 1;
  min-width: 100px;
}

.path-nav-go-btn {
  margin: 0;
  flex-shrink: 0;
}

.tasks-popover-content {
  max-height: 300px;
  overflow-y: auto;
}

.no-tasks {
  padding: 16px;
  text-align: center;
  color: v-bind('themeVars.textColor3');
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  background-color: v-bind('themeVars.tableColorStriped');
  border-radius: 4px;
  padding: 8px;

  &.is-completed {
    opacity: 0.8;
  }
}

.task-info {
  margin-bottom: 6px;
}

.task-title {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
}

.task-type-icon {
  font-size: 14px;
  color: v-bind('themeVars.primaryColor');
}

.task-name {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-details {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: v-bind('themeVars.textColor3');
}

.task-progress-bar {
  margin-top: 2px;
}

.tasks-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 8px;
  border-top: 1px solid v-bind('themeVars.borderColor');
  margin-top: 8px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: v-bind('themeVars.cardColor');
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.empty-icon {
  color: v-bind('themeVars.textColor3');
}

.empty-text {
  font-size: 16px;
  color: v-bind('themeVars.textColor3');
}

.file-browser {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: v-bind('themeVars.cardColor');
  height: 0;
  min-height: 0;
}

.file-list-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.file-list-header {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background-color: v-bind('themeVars.tableColorStriped');
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  font-weight: 500;
  flex-shrink: 0;
}

.file-list-body {
  flex: 1;
  overflow: hidden;
  min-height: 0;
  position: relative;
}

.file-list-body .n-scrollbar {
  height: 100%;
}

.file-grid-row {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: v-bind('themeVars.tableColorHover');
  }

  &.is-selected {
    background-color: v-bind('themeVars.primaryColorSuppl');
  }
}

.file-grid-cell {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 4px;
}

.checkbox-cell {
  width: 40px;
  flex-shrink: 0;
}

.name-cell {
  flex: 3;
  min-width: 200px;
}

.size-cell {
  flex: 1;
  min-width: 100px;
}

.time-cell {
  flex: 2;
  min-width: 180px;
}

.owner-cell {
  flex: 1;
  min-width: 100px;
}

.group-cell {
  flex: 1;
  min-width: 100px;
}

.perm-cell {
  flex: 1;
  min-width: 100px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 18px;

  &.is-folder {
    color: v-bind('themeVars.warningColor');
  }
}

.file-label {
  font-weight: 500;
}
</style>
