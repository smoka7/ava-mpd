<template>
  <div class="flex cursor-pointer flex-col p-2">
    <Folder v-for="(folder, index) in serverRoot" :key="index" :data="folder" />
  </div>
</template>
<script setup>
import endpoints from '../endpoints';
import { sendCommand } from '../helpers';
import { shallowReactive } from 'vue';
import Folder from './folder.vue';
const resp = await sendCommand(endpoints.folders, 'list', {
  playlist: '',
});

const serverRoot = shallowReactive([...resp.Folders, ...resp.Files]);
</script>
