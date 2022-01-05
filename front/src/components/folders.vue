<template>
  <div class="flex flex-col cursor-pointer">
    <Folder v-for="(folder, index) in folders" :key="index" :data="folder">
    </Folder>
  </div>
</template>
<script>
import Folder from "./folder.vue";
export default {
  components: {
    Folder,
  },
  data() {
    return {
      folders: Array,
    };
  },
  created() {
    this.getFolders();
  },
  methods: {
    async getFolders() {
      let response = await fetch("/api/folders");
      if (response.ok) {
        let json = await response.json();
        this.folders = json;
        return;
      }
      console.log(response.error);
    },
  },
};
</script>
