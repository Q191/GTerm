<template>
  <n-layout class="h-screen">
    <n-layout-header bordered>
      <Header ref="headerRef" />
    </n-layout-header>

    <n-layout :has-sider="!isTerminal">
      <n-layout-sider v-if="!isTerminal" bordered :width="60">
        <Sider />
      </n-layout-sider>
      <n-layout-content>
        <n-scrollbar>
          <router-view />
        </n-scrollbar>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script lang="ts" setup>
import Header from '@/layouts/Header.vue';
import Sider from '@/layouts/Sider.vue';
import { NLayout, NLayoutContent, NLayoutHeader, NLayoutSider, NScrollbar } from 'naive-ui';
import { useRoute } from 'vue-router';
import { computed, ref, provide } from 'vue';

const route = useRoute();
const isTerminal = computed(() => route.name === 'Terminal');
const headerRef = ref();

provide(
  'connectionTabs',
  computed(() => headerRef.value?.connectionTabsRef),
);
</script>
