<template>
  <TabGroup
    as="div"
    class="flex flex-col overflow-y-auto p-2 rounded bg-white/60 backdrop-blur-3xl shadow-md dark:bg-gray-700/60 darK:backdrop-blur-3xl text-primary dark:text-white fbg"
  >
    <TabList as="div" class="flex w-full justify-around p-1">
      <Tab
        v-for="(tab, name) in tabs"
        :key="name"
        v-slot="{ selected }"
        as="template"
      >
        <button :class="[tabClasses.normal, selected ? tabClasses.active : '']">
          <font-awesome-icon :icon="['fas', tab.icon]"></font-awesome-icon>
          <span>
            {{ name }}
          </span>
        </button>
      </Tab>
      <button
        aria-label="close-sidebar"
        class="p-2 px-6 md:hidden dark:text-lightest"
        @click="closeFolders"
      >
        <font-awesome-icon icon="times" size="2x"></font-awesome-icon>
      </button>
    </TabList>
    <TabPanels as="template">
      <TabPanel v-for="tab in tabs" :key="tab.icon" as="template">
        <component :is="tab.component" />
      </TabPanel>
    </TabPanels>
  </TabGroup>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { TabGroup, TabList, Tab, TabPanels, TabPanel } from "@headlessui/vue";
import storedPlaylist from "./storedPlaylist.vue";
import { toggleFolders } from "../helpers.js";
import { defineAsyncComponent } from "vue";
export default {
  components: {
    FontAwesomeIcon,
    storedPlaylist,
    folders: defineAsyncComponent(() => import("./folders.vue")),
    search: defineAsyncComponent(() => import("./search.vue")),
    Tab,
    TabList,
    TabPanel,
    TabGroup,
    TabPanels,
  },
  data() {
    return {
      tabClasses: {
        active: "bg-accent text-primary",
        normal:
          "flex space-x-1 items-center md:p-2 p-3 cursor-pointer hover:bg-white/60 dark:hover:bg-gray-800/70 dark:hover:text-white rounded",
      },
      tabs: {
        Playlists: {
          component: "storedPlaylist",
          icon: "list-ul",
        },
        Folders: {
          component: "folders",
          icon: "folder",
        },
        Search: {
          component: "search",
          icon: "search",
        },
      },
    };
  },
  methods: {
    closeFolders: toggleFolders,
  },
};
</script>
<style>
.sidebar-btn {
  @apply p-2 dark:hover:bg-white/10 hover:bg-gray-700/10 rounded;
}
.fade-enter-active,
.fade-leave-active {
  @apply transition-all duration-200 ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  @apply opacity-0 -translate-y-2;
}
</style>
