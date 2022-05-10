<template>
  <TabGroup
    as="div"
    class="fbg relative flex h-screen flex-col rounded bg-white/60 text-primary shadow-md backdrop-blur-3xl dark:bg-gray-700/60 dark:text-white"
  >
    <TabList
      as="div"
      class="sticky top-0 mb-1 flex w-full justify-around rounded-t bg-white/30 p-1 dark:bg-gray-700/60"
    >
      <Tab
        v-for="(tab, name) in tabs"
        :key="name"
        v-slot="{ selected }"
        as="template"
      >
        <button :class="[tabClasses.normal, selected ? tabClasses.active : '']">
          <font-awesome-icon :icon="['fas', tab.icon]" />
          <span class="hidden md:block">
            {{ name }}
          </span>
        </button>
      </Tab>
      <button
        aria-label="close-sidebar"
        class="p-2 px-6 dark:text-lightest md:hidden"
        @click="toggleMediaController"
      >
        <font-awesome-icon icon="times" size="2x" />
      </button>
    </TabList>
    <TabPanels class="h-full overflow-y-auto">
      <TabPanel v-for="(tab, name) in tabs" :key="name" as="div">
        <suspense>
          <component :is="tab.component" />
          <template #fallback>
            <div
              class="flex h-full w-full items-center justify-center p-4 text-7xl underline decoration-accent md:text-9xl"
            >
              Loading...
            </div>
          </template>
        </suspense>
      </TabPanel>
    </TabPanels>
  </TabGroup>
</template>
<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { TabGroup, TabList, Tab, TabPanels, TabPanel } from "@headlessui/vue";
import { toggleMediaController } from "../helpers.js";
import { defineAsyncComponent } from "vue";
import { useStore } from "vuex";

const store = useStore();

const setting = defineAsyncComponent({
  loader: () => import("./setting.vue"),
  delay: 200,
});

const queueCmp = defineAsyncComponent({
  loader: () => {
    store.dispatch("getQueue");
    return import("./playlist.vue");
  },
  delay: 200,
});

const storedPlaylist = defineAsyncComponent({
  loader: () => {
    store.dispatch("getStoredPlaylist");
    return import("./storedPlaylist.vue");
  },
  delay: 200,
});

const folders = defineAsyncComponent({
  loader: () => {
    store.dispatch("getServerFolders");
    return import("./folders.vue");
  },
  delay: 200,
});

const search = defineAsyncComponent({
  loader: () => import("./search.vue"),
  delay: 200,
});

const tabClasses = {
  active: "bg-accent text-primary",
  normal:
    "flex space-x-1 items-center md:p-2 p-3 cursor-pointer hover:bg-white/60 dark:hover:bg-gray-800/70 dark:hover:text-white rounded",
};

const tabs = {
  Queue: {
    component: queueCmp,
    icon: "list-ul",
  },
  Settings: {
    component: setting,
    icon: "cog",
  },
  Playlists: {
    component: storedPlaylist,
    icon: "list-ul",
  },
  Folders: {
    component: folders,
    icon: "folder",
  },
  Search: {
    component: search,
    icon: "search",
  },
};
</script>
<style>
.fade-enter-active,
.fade-leave-active {
  @apply transition-all duration-200 ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  @apply -translate-y-2 opacity-0;
}
</style>
