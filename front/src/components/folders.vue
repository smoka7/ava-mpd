<template>
  <div class="flex cursor-pointer flex-col p-2">
    <Folder v-for="(folder, index) in serverRoot" :key="index" :data="folder" />
  </div>
</template>
<script setup lang="ts">
import { useStore } from "../store";

const store = useStore();
const serverRoot = computed(() => shallowReactive(store.serverFolders));
</script>
