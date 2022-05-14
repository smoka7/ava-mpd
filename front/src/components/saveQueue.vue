<template>
  <button
    aria-label="save-queue"
    @click="InputIsOpen = true"
    class="tooltip scale-125 transform bg-transparent p-2 hover:text-accent"
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
<script setup>
import TextInput from "./textInput.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { sendCommand } from "../helpers";
import endpoints from "../endpoints.js";
import { useStore } from "vuex";
import { ref } from "vue";

const store = useStore();
const InputIsOpen = ref(false);

function saveQueue(name) {
  sendCommand(endpoints.queue, "save", { playlist: name });
  store.dispatch("getStoredPlaylist");
}
</script>
