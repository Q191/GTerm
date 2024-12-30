<template>
  <div class="flex items-center relative">
    <div class="tab-container-wrapper">
      <div ref="tabContainerRef" class="tab-container">
        <NTag
          v-for="tab in visibleTabs"
          :key="tab"
          size="large"
          :bordered="false"
          closable
          style="margin-right: 10px; cursor: pointer"
          :theme-overrides="{
            color: selectedTab === tab ? 'rgba(242, 201, 125, 0.16)' : 'rgba(255, 255, 255, .08)',
            textColor: selectedTab === tab ? '#f2c97d' : '',
          }"
          @close="closeTab(tab)"
          @click="selectTab(tab)"
        >
          {{ tab }}
          <template #icon>
            <Debian size="20" />
          </template>
        </NTag>
      </div>
    </div>
    <NDropdown v-if="hasHiddenTabs" trigger="click" :options="dropdownOptions" @select="handleSelect">
      <NButton class="pl-2" quaternary size="small"> 展开 </NButton>
    </NDropdown>
    <NButton class="pl-2" quaternary size="small" @click="addTab">
      <template #icon>
        <NIcon>
          <Add />
        </NIcon>
      </template>
    </NButton>
  </div>
</template>

<script setup lang="ts">
import Debian from '@/assets/icons/Debian.vue';
import { Add } from '@vicons/ionicons5';
import { NButton, NDropdown, NIcon, NTag } from 'naive-ui';
import { computed, nextTick, ref, watch } from 'vue';

const props = defineProps<{
  tabWidth: number;
}>();
const tabs = ref<string[]>(['测试服务器', 'web 服务器', '后端服务器']);
const selectedTab = ref<string | null>(null);
const tabContainerRef = ref<HTMLElement | null>(null);

const tabContentWidth = computed(() => props.tabWidth - 430);

const tabContentWidthStyle = computed(() => ({ width: `${tabContentWidth.value}px` }));

const visibleTabs = computed(() => {
  const containerWidth = tabContentWidth.value;
  let totalWidth = 0;
  return tabs.value.filter(tab => {
    const tabWidth = tab.length * 14 + 60; // 估算每个标签的宽度
    if (totalWidth + tabWidth <= containerWidth) {
      totalWidth += tabWidth;
      return true;
    }
    return false;
  });
});

const dropdownOptions = computed(() => {
  return tabs.value.map(tab => ({
    label: tab,
    key: tab,
    disabled: visibleTabs.value.includes(tab),
  }));
});

const hasHiddenTabs = computed(() => visibleTabs.value.length < tabs.value.length);

const addTab = () => {
  const newTab = `服务器 ${tabs.value.length + 1}`;
  tabs.value.push(newTab);
};

const closeTab = (tab: string) => {
  tabs.value = tabs.value.filter(t => t !== tab);
};

const selectTab = (tab: string) => {
  selectedTab.value = tab;
};

const handleSelect = (key: string) => {
  selectTab(key);
};

watch(tabs, () => {
  nextTick(() => {
    if (tabContainerRef.value) {
      tabContainerRef.value.scrollLeft = 0;
    }
  });
});
</script>

<style lang="less" scoped>
.tab-container-wrapper {
  max-width: v-bind(tabContentWidthStyle);
  overflow: hidden;
}

.tab-container {
  display: flex;
  white-space: nowrap;
  flex-wrap: nowrap;
  padding: 10px 0;
}
</style>
