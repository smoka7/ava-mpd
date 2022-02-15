<template>
  <div class="flex cursor-pointer flex-col">
    <Folder v-for="(folder, index) in folders" :key="index" :data="folder" />
  </div>
</template>
<script>
import endpoints from "../endpoints";
import { sendCommand } from "../helpers";
import { shallowReactive } from "vue";
import Folder from "./folder.vue";
export default {
  components: {
    Folder,
  },
  data() {
    return {
      folders: shallowReactive([]),
    };
  },
  async created() {
    const resp = await sendCommand(endpoints.folders, "list", {
      playlist: "",
    });
    this.folders = [...resp.Folders, ...resp.Files];
  },
};
</script>
