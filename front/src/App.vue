<template>
  <div
    class="from-lightest to-accent via-lighter text-primary flex h-screen w-screen flex-col overflow-hidden bg-gradient-to-br text-lg duration-300 dark:text-white md:grid md:grid-cols-4 md:grid-rows-4 md:gap-2 md:p-1"
  >
    <sidebar
      id="sidebar"
      class="fixed inset-0 h-full md:static md:col-span-1 md:col-start-1 md:col-end-1 md:row-start-1 md:row-end-4"
    />
    <div
      class="fixed inset-0 md:relative md:col-span-3 md:col-start-2 md:col-end-5 md:row-start-1 md:row-end-4"
      id="queue"
    >
      <setting v-if="settingIsOpen" @close="toggleSetting()" />
      <playlist v-else />
    </div>
    <media-controller
      @openSetting="toggleSetting()"
      class="fixed inset-0 z-10 md:col-span-4 md:col-start-1 md:col-end-5 md:row-span-1 md:row-start-4 md:row-end-4"
    />
  </div>
</template>
<script>
import MediaController from "./components/mediaController.vue";
import Playlist from "./components/playlist.vue";
import Sidebar from "./components/sidebar.vue";
import { setColorScheme } from "./colors.js";
import { toggleMediaController } from "./helpers.js";
import { defineAsyncComponent } from "vue";
export default {
  components: {
    MediaController,
    Playlist,
    Sidebar,
    setting: defineAsyncComponent(() => import("./components/setting.vue")),
  },
  data() {
    return {
      settingIsOpen: false,
    };
  },
  methods: {
    updatePlayer() {
      this.$store.dispatch("getCurrentSong");
      this.$store.dispatch("getQueue");
    },
    toggleSetting() {
      this.settingIsOpen = !this.settingIsOpen;
      toggleMediaController();
    },
  },
  created() {
    this.$store.dispatch("connectToSocket");
    this.$store.dispatch("getCurrentSong");
    this.$store.dispatch("startCounter");
  },
  mounted() {
    setColorScheme();
  },
};
</script>
<style>
.tooltip {
  @apply relative inline-block;
}
.tooltip .tooltiptext {
  @apply bg-lightest/90 text-primary invisible absolute bottom-full left-1/2 z-10 -ml-10 inline w-auto rounded p-2 text-center;
}
.tooltip:hover .tooltiptext {
  @apply visible;
}
::-webkit-scrollbar {
  @apply w-2;
}
::-webkit-scrollbar-track {
  @apply bg-lightest;
}
::-webkit-scrollbar-thumb {
  @apply bg-primary rounded-xl;
}
::-webkit-scrollbar-thumb:hover {
  @apply bg-secondary rounded-xl;
}
@supports (-moz-appearance: none) {
  .fbg {
    @apply bg-white dark:bg-gray-700 !important;
  }
}
</style>
