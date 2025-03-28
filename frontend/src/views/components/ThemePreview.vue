<template>
  <div class="theme-preview-container">
    <div class="preview-content" :style="previewStyle">
      <div class="line" :style="{ color: theme?.foreground }">$ ls -la</div>
      <div class="line" :style="{ color: theme?.brightGreen }">drwxr-xr-x 1 user group 48 Mar 28 01:51 .</div>
      <div class="line" :style="{ color: theme?.brightGreen }">drwxr-xr-x 3 user group 96 Mar 28 01:51 ..</div>
      <div class="line" :style="{ color: theme?.brightBlue }">drwxr-xr-x 2 user group 64 Mar 28 01:51 themes</div>
      <div class="line" :style="{ color: theme?.brightYellow }">-rw-r--r-- 1 user group 128 Mar 28 01:51 README.md</div>
      <div class="line" :style="{ color: theme?.brightRed }">$ git status</div>
      <div class="line" :style="{ color: theme?.brightMagenta }">On branch main</div>
      <div class="line" :style="{ color: theme?.brightCyan }">Your branch is up to date with 'origin/main'.</div>
      <div class="line" :style="{ color: theme?.foreground }">$</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, toRef, toValue } from 'vue';

const props = defineProps<{
  theme: any;
}>();

const themeRef = toRef(props, 'theme');

const previewStyle = computed(() => ({
  backgroundColor: toValue(themeRef)?.background || '#000000',
}));
</script>

<style lang="less" scoped>
.theme-preview-container {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  overflow: hidden;
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
}

.preview-content {
  padding: 8px;
  height: 100%;
  overflow: hidden;
}

.line {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 2px;
}
</style>
