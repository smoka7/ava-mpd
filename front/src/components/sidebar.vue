<template>
  <TabGroup
    as="div"
    class="flex flex-col app-height overflow-y-auto p-2 rounded bg-white m-1 shadow-md dark:bg-gray-700 text-primary dark:text-foreground"
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
import folders from "./folders.vue";
import search from "./search.vue";
import { toggleFolders } from "../helpers.js";
export default {
  components: {
    FontAwesomeIcon,
    storedPlaylist,
    folders,
    search,
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
          "flex space-x-1 items-center md:p-2 p-3 cursor-pointer hover:bg-accent dark:hover:text-primary rounded",
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
  @apply p-2 hover:bg-white dark:hover:bg-gray-700 dark:hover:text-foreground rounded;
}
</style>
