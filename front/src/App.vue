<template>
  <div class="flex flex-col-reverse">
    <main class="h-full w-full mx-auto flex">
      <sidebar id="sidebar" class="md:w-1/4 w-full fixed md:static h-full">
      </sidebar>
      <div class="w-full md:w-3/4 md:ml-2 md:m-1">
        <setting v-if="settingIsOpen" @close="toggleSetting()"></setting>
        <playlist v-else></playlist>
      </div>
    </main>
    <media-controller @openSetting="toggleSetting()"> </media-controller>
  </div>
</template>
<script>
import MediaController from "./components/mediaController.vue";
import Playlist from "./components/playlist.vue";
import Sidebar from "./components/sidebar.vue";
import Setting from "./components/setting.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { setColorScheme } from "./colors.js";
export default {
  components: {
    MediaController,
    Playlist,
    Sidebar,
    FontAwesomeIcon,
    Setting,
  },
  data() {
    return {
      settingIsOpen: false,
    };
  },
  methods: {
    updatePlayer() {
      this.$store.dispatch("getCurrentSong");
    },
    toggleSetting() {
      this.settingIsOpen = !this.settingIsOpen;
    },
    connectToSocket() {
      let hostname = new URL(window.location.href).hostname;
      let socket = new WebSocket("ws://" + hostname + ":3001/update");
      socket.onmessage = (event) => {
        let some = JSON.parse(event.data);
        if (some.Subsystem == "playlist") {
          this.$store.dispatch("getCurrentPlaylist");
        }
        if (
          some.Subsystem == "player" ||
          some.Subsystem == "mixer" ||
          some.Subsystem == "sticker" ||
          some.Subsystem == "options"
        ) {
          this.updatePlayer();
        }
      };
    },
  },
  created() {
    this.$store.dispatch("getCurrentSong");
    this.$store.dispatch("getCurrentPlaylist");
    this.$store.dispatch("startCounter");
  },
  mounted() {
    this.connectToSocket();
    setColorScheme();
  },
};
</script>
