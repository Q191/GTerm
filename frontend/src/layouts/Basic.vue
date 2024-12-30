<template>
  <NLayout class="h-screen">
    <Header />
    <router-view />
  </NLayout>
</template>

<script lang="ts" setup>
import Header from '@/layouts/Header.vue';
import { NLayout } from 'naive-ui';
</script>
