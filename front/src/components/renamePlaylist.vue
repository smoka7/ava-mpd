<template>
  <text-input
    :isOpen="InputIsOpen"
    @close="InputIsOpen = false"
    @finish="renamePlaylist"
  >
    <template #title> Enter the new name</template>
    <template #btn>Rename</template>
  </text-input>
</template>
<script setup lang="ts">
import { sendCommand } from "../helpers";
import endpoints from "../endpoints";
import { useStore } from "../store";

const store = useStore();
const props = defineProps<{
  rename: string;
}>();
const InputIsOpen = ref(true);

function renamePlaylist(name: string) {
  sendCommand(endpoints.storedPlaylists, "rename", {
    Playlist: props.rename,
    NewPlaylist: name,
  });
  store.getStoredPlaylist();
}
</script>
