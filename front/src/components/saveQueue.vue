<template>
  <button
    aria-label="save-queue"
    @click="InputIsOpen = true"
    class="tooltip playback-btn"
  >
    <font-awesome-icon icon="save" size="lg" />
    <span class="tooltiptext">save current queue</span>
  </button>
  <text-input
    :isOpen="InputIsOpen"
    @close="InputIsOpen = false"
    @finish="saveQueue"
  >
    <template #title> Enter the name of playlist </template>
    <template #btn>Save</template>
  </text-input>
</template>
<script setup lang="ts">
import { sendCommand } from "../helpers";
import endpoints from "../endpoints";

const store = useStore();
const InputIsOpen = ref(false);

function saveQueue(name) {
  sendCommand(endpoints.queue, "save", { Playlist: name });
  store.dispatch("getStoredPlaylist");
}
</script>
