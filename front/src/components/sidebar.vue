<template>
  <TabGroup
    as="div"
    :selectedIndex="selectedTab"
    @change="changeTab"
    class="relative flex h-screen flex-col rounded bg-white/60 text-primary shadow-md backdrop-blur-3xl dark:bg-gray-700/60 dark:text-white"
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
        <button :class="['normal-tab group', selected ? 'active-tab' : '']">
          <font-awesome-icon
            :icon="['fas', tab.icon]"
            :class="[
              'group-hover:text-secondary',
              selected ? 'text-accent' : 'text-secondary dark:text-lighter',
            ]"
          />
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
<script setup lang="ts">
import { settingsStore } from "../settingsStore";
import { toggleMediaController } from "../helpers";
import { useStore } from "../store";

const store = useStore();
const setting = defineAsyncComponent({
  loader: () => {
    const setStore = settingsStore();
    setStore.getSettings();
    return import("./setting.vue");
  },
  delay: 200,
});

const queueCmp = defineAsyncComponent({
  loader: () => {
    store.getQueue();
    return import("./playlist.vue");
  },
  delay: 200,
});

const storedPlaylist = defineAsyncComponent({
  loader: () => {
    store.getStoredPlaylist();
    return import("./storedPlaylist.vue");
  },
  delay: 200,
});

const folders = defineAsyncComponent({
  loader: () => {
    store.getServerFolders();
    return import("./folders.vue");
  },
  delay: 200,
});

const search = defineAsyncComponent({
  loader: () => import("./search.vue"),
  delay: 200,
});

const selectedTab = computed(() => store.activeTabIndex);

function changeTab(index: number) {
  store.setActiveTab(index);
}

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
    icon: "bookmark",
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
<style lang="postcss">
.fade-enter-active,
.fade-leave-active {
  @apply transition-all duration-200 ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  @apply -translate-y-2 opacity-0;
}

.active-tab {
  @apply bg-white dark:bg-gray-800/70 !important;
}

.normal-tab {
  @apply flex cursor-pointer items-center space-x-1 rounded p-3 duration-300;
  @apply hover:bg-white/60 hover:px-4 focus:outline-none focus:ring focus:ring-lightest;
  @apply dark:hover:bg-gray-800/70;
  @apply md:p-2;
}
</style>
