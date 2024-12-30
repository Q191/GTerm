<template>
  <div class="inventory-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <NInput placeholder="搜索主机或用户@主机名..." round>
        <template #prefix>
          <Icon icon="ph:magnifying-glass-duotone" />
        </template>
        <template #suffix>
          <NButton size="tiny" type="primary" ghost>连接</NButton>
        </template>
      </NInput>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar">
      <NButton size="small" type="primary" ghost @click="dialogStore.openAddServerDialog">
        <template #icon>
          <Icon icon="ph:plus-circle-duotone" />
        </template>
        添加服务器
      </NButton>
      <NButton size="small" ghost>
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
        <span class="text-count">共 3 个分组</span>
      </div>

      <NGrid x-gap="12" y-gap="12" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
        <NGi v-for="index in 3" :key="index">
          <div class="card group" @click="toTerminal">
            <div class="card-content">
              <div class="card-icon">
                <Icon icon="ph:folders-duotone" />
              </div>
              <div class="card-info">
                <div class="card-title">测试分组</div>
                <div class="card-subtitle">7 台主机</div>
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
        <NH3 prefix="bar">主机列表</NH3>
        <span class="text-count">共 100 台主机</span>
      </div>

      <NGrid x-gap="12" y-gap="12" cols="2 s:2 m:3 l:4 xl:6" responsive="screen">
        <NGi v-for="index in 100" :key="index">
          <div class="card group" @click="toTerminal">
            <div class="card-content">
              <div class="card-icon">
                <Icon icon="ph:linux-logo-duotone" />
              </div>
              <div class="card-info">
                <div class="card-title">测试服务器</div>
                <div class="card-subtitle">ssh://root@192.168.1.100</div>
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
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, NGi, NGrid, NH3, NInput } from 'naive-ui';
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';
import { gtermTheme } from '@/themes/gterm-theme';

const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const router = useRouter();

const toTerminal = () => {
  router.push({ name: 'Terminal' });
};

const gtermThemeVars = computed(() => {
  return gtermTheme(prefStore.isDark);
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
  }

  &:hover .card-actions {
    opacity: 1;
  }

  .action-btn {
    color: v-bind('gtermThemeVars.textColor');
    opacity: 0.7;

    &:hover {
      opacity: 1;
      background-color: v-bind('gtermThemeVars.splitColor');
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
