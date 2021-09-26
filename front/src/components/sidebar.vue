<template>
  <div
    class="
      flex flex-col
      app-height
      overflow-y-auto
      p-2
      rounded
      bg-white
      m-1
      shadow-md
      dark:bg-gray-700 dark:text-white dark:hover:text-primary
    "
  >
    <div class="flex w-full my-1 justify-around">
      <button
        :class="activeTab == 0 ? tabClasses[0] : tabClasses[1]"
        @click="changeTab(0)"
      >
        <font-awesome-icon icon="list-ul"></font-awesome-icon>
        <span> playlists </span>
      </button>
      <button
        :class="activeTab == 1 ? tabClasses[0] : tabClasses[1]"
        @click="changeTab(1)"
      >
        <font-awesome-icon icon="folder"></font-awesome-icon>
        <span> folders </span>
      </button>
      <button
        :class="activeTab == 2 ? tabClasses[0] : tabClasses[1]"
        @click="changeTab(2)"
      >
        <font-awesome-icon icon="search"></font-awesome-icon>
        <span> search</span>
      </button>
      <button
        class="p-2 px-6 md:hidden text-primary dark:text-lightest"
        @click="closeFolders"
      >
        <font-awesome-icon icon="times" size="2x"></font-awesome-icon>
      </button>
    </div>
    <storedPlaylist v-if="activeTab == 0"></storedPlaylist>
    <folders v-if="activeTab == 1"></folders>
    <search v-if="activeTab == 2"></search>
  </div>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import storedPlaylist from "./storedPlaylist.vue";
import folders from "./folders.vue";
import search from "./search.vue";
import { toggleFolders } from "../helpers.js";
export default {
  components: {
    FontAwesomeIcon,
    storedPlaylist,
    folders,
    search,
  },
  data() {
    return {
      activeTab: 0,
      tabClasses: [
        "md:p-2 p-3 bg-primary text-foreground dark:bg-lightest dark:text-primary rounded",
        "md:p-2 p-3 cursor-pointer text-primary dark:text-foreground hover:bg-blue-200 dark:hover:text-primary rounded",
      ],
    };
  },
  methods: {
    changeTab(id) {
      this.activeTab = id;
    },
    closeFolders: toggleFolders,
  },
};
</script>
