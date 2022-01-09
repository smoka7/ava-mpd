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
<script>
import TextInput from "./textInput.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { sendCommand } from "../helpers";
export default {
  props: ["rename"],
  components: {
    FontAwesomeIcon,
    TextInput,
  },
  data() {
    return {
      InputIsOpen: true,
    };
  },
  methods: {
    renamePlaylist(name) {
      sendCommand("/api/stored", "rename", {
        playlist: this.rename,
        song: name,
      });
      this.$store.dispatch("getStoredPlaylist");
    },
  },
};
</script>
