<template>
  <button
    aria-label="save-queue"
    @click="InputIsOpen = true"
    class="p-2 px-3 text-foreground rounded-full md:p-0 tooltip md:hover:text-blue-200"
  >
    <font-awesome-icon icon="save" size="lg"></font-awesome-icon>
    <span class="tooltiptext">save current queue</span>
  </button>
  <text-input :isOpen="InputIsOpen" @close="InputIsOpen = false" @finish="saveQueue">
    <template #title> Enter the name of playlist </template>
    <template #btn>Save</template>
  </text-input>
</template>
<script>
import TextInput from "./textInput.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { sendCommand } from "../helpers";
export default {
  components: {
    FontAwesomeIcon,
    TextInput,
  },
  data() {
    return {
      InputIsOpen: false,
    };
  },
  methods: {
    saveQueue(name) {
      sendCommand("/api/queue", "save", { playlist: name });
      this.$store.dispatch("getStoredPlaylist");
    },
  },
};
</script>
