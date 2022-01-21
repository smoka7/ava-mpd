<template>
  <div class="flex flex-col cursor-pointer">
    <Folder v-for="(folder, index) in folders" :key="index" :data="folder" />
  </div>
</template>
<script>
import { getFolders } from "../helpers";
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
    this.folders = await getFolders("");
  },
};
</script>
