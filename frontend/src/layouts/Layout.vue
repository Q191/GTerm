<template>
  <NLayout class="layout-container">
    <NLayoutHeader bordered class="header">
      <Header ref="headerRef" />
    </NLayoutHeader>

    <NLayout :has-sider="!isTerminal" class="content">
      <NLayoutSider v-if="!isTerminal" bordered :width="50">
        <Sider />
      </NLayoutSider>
      <NLayoutContent>
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>

<script lang="ts" setup>
import { NLayout, NLayoutContent, NLayoutHeader, NLayoutSider } from 'naive-ui';
import { computed, provide, ref } from 'vue';
import { useRoute } from 'vue-router';
import Header from '@/layouts/Header.vue';
import Sider from '@/layouts/Sider.vue';

const route = useRoute();
const isTerminal = computed(() => route.name === 'Terminal');
const headerRef = ref();

provide(
  'connectionTabs',
  computed(() => headerRef.value?.connectionTabsRef),
);
</script>

<style lang="less" scoped>
.layout-container {
  height: 100vh;
  .header {
    height: 38px;
  }
  .content {
    height: calc(100vh - 38px);
  }
}
</style>
