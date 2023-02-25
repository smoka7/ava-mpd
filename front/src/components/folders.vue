<template>
  <div class="flex cursor-pointer flex-col p-2">
    <Folder
      v-for="(directory, index) in serverRoot.Directories"
      :key="index"
      :directory="directory"
    />
    <Folder
      v-for="(file, index) in serverRoot.Files"
      :key="index"
      :file="file"
    />
  </div>
</template>
<script setup lang="ts">
import { useStore } from "../store";

const store = useStore();
const serverRoot = computed(() => shallowReactive(store.serverFolders));
</script>
