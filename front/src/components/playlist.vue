<template>
  <div
    class="flex flex-col w-full h-full rounded overflow-hidden backdrop-blur-3xl bg-white/60 dark:bg-gray-700/60"
  >
    <div
      v-if="queue.Length == 0"
      class="flex items-center justify-center p-4 w-full h-full text-9xl decoration-accent underline"
    >
      Queue is empty!
    </div>
    <div
      v-else
      class="flex space-x-4 justify-between items-center w-full sticky top-0 h-8 md:h-6 text-base md:text-sm text-primary px-4 bg-secondary z-10"
    >
      <span>{{ queue.Length }} Tracks </span>
      <span>duration: {{ humanizeTime(queue.Duration) }} </span>
    </div>
    <div class="overflow-y-auto space-y-1 px-2">
      <album
        v-for="(album, index) in queue.Albums"
        :key="index"
        :album="album"
        :currentSongPos="currentSongPos"
        @showMenu="showMenu"
      />
    </div>
    <button
      aria-label="close-playlist"
      @click="closePlaylist"
      class="fixed right-4 bottom-4 bg-red-500 text-white rounded-full p-2 md:hidden w-20 h-20"
    >
      <font-awesome-icon icon="arrow-right" size="2x" />
    </button>
    <teleport to="#app">
      <queue-menu
        :open="menu"
        @close="hideMenu"
        :song="selected"
        @showInfo="songInfo = true"
      />
    </teleport>
    <songInfo v-if="songInfo" :song="selected" @close="songInfo = false" />
  </div>
</template>
<script>
import { mapState } from "vuex";
import {
  sendCommand,
  humanizeTime,
  toggleMediaController,
} from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { defineAsyncComponent } from "vue";
import album from "./album.vue";
export default {
  components: {
    album,
    FontAwesomeIcon,
    songInfo: defineAsyncComponent(() => import("./songInfo.vue")),
    queueMenu: defineAsyncComponent(() => import("./queueMenu.vue")),
  },
  data() {
    return {
      menu: false,
      songInfo: false,
      selected: -1,
    };
  },
  async mounted() {
    await this.$store.dispatch("getQueue");
    this.ChangeCurrentSongTo(this.currentSongPos);
  },
  methods: {
    humanizeTime: humanizeTime,
    closePlaylist: toggleMediaController,
    showMenu(pos, e) {
      this.selected = Number(pos);
      let el = document.getElementById("context-menu");
      el.style.top = e.pageY + "px";
      this.menu = true;
    },
    hideMenu() {
      this.menu = false;
    },
    ChangeCurrentSongTo(id) {
      let el = document.getElementById("song" + id);
      if (el !== null) {
        el.classList.add("current-song");
        el.scrollIntoView({ block: "center", behavior: "smooth" });
      }
    },
  },
  computed: {
    ...mapState({
      queue: (state) => state.queue,
      currentSongPos: (state) => state.currentSong.Pos,
    }),
  },
  watch: {
    currentSongPos: function (newSong, oldSong) {
      let oldId = "song" + oldSong;
      let newId = "song" + newSong;
      let el = document.getElementById(oldId);
      if (el) {
        el.classList.remove("current-song");
      }
      this.ChangeCurrentSongTo(newId);
    },
  },
};
</script>
<style>
.current-song {
  @apply bg-red-300 dark:text-primary !important;
}
</style>
