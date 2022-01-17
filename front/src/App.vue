<template>
  <div
    class="flex flex-col md:grid md:grid-rows-4 md:grid-cols-4 md:gap-2 md:p-1 text-lg w-screen h-screen duration-300 text-primary dark:text-white bg-gradian overflow-hidden"
  >
    <sidebar
      id="sidebar"
      class="md:col-start-1 md:col-end-1 md:row-start-1 md:row-end-4 md:col-span-1 fixed inset-0 md:static h-full"
    />
    <div
      class="md:col-start-2 md:col-end-5 md:col-span-3 md:row-start-1 md:row-end-4 fixed inset-0 md:relative"
      id="queue"
    >
      <setting v-if="settingIsOpen" @close="toggleSetting()" />
      <playlist v-else />
    </div>
    <media-controller
      @openSetting="toggleSetting()"
      class="md:col-span-4 md:col-start-1 md:col-end-5 md:row-start-4 md:row-end-4 md:row-span-1 fixed inset-0 z-10"
    />
  </div>
</template>
<script>
import MediaController from "./components/mediaController.vue";
import Playlist from "./components/playlist.vue";
import Sidebar from "./components/sidebar.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { setColorScheme } from "./colors.js";
import { toggleMediaController } from "./helpers.js";
import { defineAsyncComponent } from "vue";
export default {
  components: {
    MediaController,
    Playlist,
    Sidebar,
    FontAwesomeIcon,
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
  @apply invisible bg-lightest/90 text-primary w-auto text-center p-2 rounded bottom-full left-1/2 -ml-10 absolute z-10 inline;
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
</style>
