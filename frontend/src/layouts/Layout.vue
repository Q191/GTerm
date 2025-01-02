<template>
  <NLayout class="h-screen">
    <NLayoutHeader bordered>
      <Header ref="headerRef" />
    </NLayoutHeader>

    <NLayout :has-sider="!isTerminal">
      <NLayoutSider v-if="!isTerminal" bordered :width="60">
        <Sider />
      </NLayoutSider>
      <NLayoutContent>
        <NScrollbar>
          <router-view />
        </NScrollbar>
      </NLayoutContent>
    </NLayout>
  </NLayout>
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
