<template>
  <div
    class="text-primary bg-gradian flex h-screen w-screen flex-col overflow-hidden text-lg duration-300 dark:text-white md:grid md:grid-cols-4 md:grid-rows-4 md:gap-2 md:p-1"
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
    connectToSocket() {
      let hostname = new URL(window.location.href).host;
      let socket = new WebSocket("ws://" + hostname + "/update");
      socket.onmessage = () => {
        this.updatePlayer();
      };
    },
  },
  created() {
    this.$store.dispatch("getCurrentSong");
    this.$store.dispatch("startCounter");
  },
  mounted() {
    this.connectToSocket();
    setColorScheme();
  },
};
</script>
<style>
.bg-gradian {
  background: radial-gradient(50% 123.47% at 50% 50%, #00ff94 0%, #720059 100%),
    linear-gradient(121.28deg, #669600 0%, #ff0000 100%),
    linear-gradient(360deg, #0029ff 0%, #8fff00 100%),
    radial-gradient(100% 164.72% at 100% 100%, #6100ff 0%, #00ff57 100%),
    radial-gradient(100% 148.07% at 0% 0%, #fff500 0%, #51d500 100%);
  background-blend-mode: screen, color-dodge, overlay, difference, normal;
}
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
