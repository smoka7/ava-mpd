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
<script setup>
import TextInput from "./textInput.vue";
import { sendCommand } from "../helpers";
import endpoints from "../endpoints.js";
import { ref } from "vue";
import { useStore } from "vuex";

const store=useStore();
const props=defineProps(["rename"]);
const InputIsOpen=ref(true);

function renamePlaylist(name) {
  sendCommand(endpoints.storedPlaylists, "rename", {
    playlist: props.rename,
    song: name,
  });
  store.dispatch("getStoredPlaylist");
};
</script>
