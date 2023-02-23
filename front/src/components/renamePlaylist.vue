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
import endpoints from "../endpoints.js";

const store=useStore();
const props=defineProps(["rename"]);
const InputIsOpen=ref(true);

function renamePlaylist(name) {
  sendCommand(endpoints.storedPlaylists, "rename", {
    Playlist: props.rename,
    NewPlaylist: name,
  });
  store.dispatch("getStoredPlaylist");
};
</script>
