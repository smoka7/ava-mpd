<template>
  <TabGroup
    as="div"
    class="
      flex flex-col
      app-height
      overflow-y-auto
      p-2
      pt-1
      rounded
      bg-white
      m-1
      shadow-md
      dark:bg-gray-700 dark:text-white dark:hover:text-primary
    "
  >
    <TabList as="div" class="flex w-full justify-around p-1">
      <Tab v-slot="{ selected }">
        <button :class="selected ? tabClasses[0] : tabClasses[1]">
          <font-awesome-icon icon="list-ul"></font-awesome-icon>
          <span> playlists </span>
        </button>
      </Tab>
      <Tab v-slot="{ selected }">
        <button :class="selected ? tabClasses[0] : tabClasses[1]">
          <font-awesome-icon icon="folder"></font-awesome-icon>
          <span> folders </span>
        </button>
      </Tab>
      <Tab v-slot="{ selected }">
        <button :class="selected ? tabClasses[0] : tabClasses[1]">
          <font-awesome-icon icon="search"></font-awesome-icon>
          <span> search</span>
        </button>
      </Tab>
      <button
        class="p-2 px-6 md:hidden text-primary dark:text-lightest"
        @click="closeFolders"
      >
        <font-awesome-icon icon="times" size="2x"></font-awesome-icon>
      </button>
    </TabList>
    <TabPanels>
      <TabPanel>
        <storedPlaylist></storedPlaylist>
      </TabPanel>
      <TabPanel>
        <folders></folders>
      </TabPanel>
      <TabPanel>
        <search></search>
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
      tabClasses: [
        "md:p-2 p-3 bg-primary text-foreground dark:bg-lightest dark:text-primary rounded",
        "md:p-2 p-3 cursor-pointer text-primary dark:text-foreground hover:bg-blue-200 dark:hover:text-primary rounded",
      ],
    };
  },
  methods: {
    closeFolders: toggleFolders,
  },
};
</script>
