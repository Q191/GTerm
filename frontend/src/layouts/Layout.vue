<template>
  <NLayout class="h-screen">
    <NLayoutHeader bordered>
      <Header />
    </NLayoutHeader>

    <NLayout has-sider>
      <NLayoutSider bordered :width="60">
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
import { NLayout, NLayoutContent, NLayoutSider, NScrollbar } from 'naive-ui';
</script>
