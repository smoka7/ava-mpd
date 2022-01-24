<template>
  <div
    class="flex flex-col w-full h-full rounded overflow-hidden backdrop-blur-3xl bg-white/60 dark:bg-gray-700/60 fbg"
  >
    <div
      v-if="queue.Length == 0"
      class="flex items-center justify-center p-4 w-full h-full text-9xl decoration-accent underline"
    >
      Queue is empty!
    </div>
    <div
      v-if="!songInfo && queue.Length != 0"
      class="flex space-x-4 justify-between items-center w-full absolute bottom-0 h-10 md:h-6 text-base md:text-sm text-primary px-4 bg-secondary z-10"
    >
      <span>{{ queue.Length }} Tracks </span>
      <span>duration: {{ humanizeTime(queue.Duration) }} </span>
      <button
        aria-label="goto-current-song"
        @click="ScrollToCurrentSong()"
        class="px-2 hover:text-secondary hover:bg-primary"
      >
        <font-awesome-icon icon="compact-disc" />
        Current Song
      </button>
      <button
        aria-label="close-playlist"
        @click="closePlaylist"
        class="px-2 md:hidden hover:text-secondary hover:bg-primary"
      >
        <font-awesome-icon icon="times" size="2x" />
      </button>
    </div>
    <div class="overflow-y-auto space-y-1 md:px-2 mb-10 md:mb-6">
      <album
        v-for="(album, index) in queue.Albums"
        :key="index"
        :album="album"
        :currentAlbum="currentAlbum(album)"
        :currentSongPos="currentSongPos"
        @showMenu="showMenu"
      />
    </div>
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
import { humanizeTime, toggleMediaController } from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { defineAsyncComponent, shallowReactive } from "vue";
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
    this.ScrollToCurrentSong();
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
        el.id = "currentSong";
      }
    },
    ScrollToCurrentSong() {
      let el = document.getElementById("currentSong");
      if (el) el.scrollIntoView({ block: "center", behavior: "smooth" });
    },
    currentAlbum(album) {
      return (
        Number(this.currentSongPos) >= Number(album.Songs[0].Pos) &&
        Number(this.currentSongPos) <=
          Number(album.Songs[album.Songs.length - 1].Pos)
      );
    },
  },
  computed: {
    ...mapState({
      queue: (state) => shallowReactive(state.queue),
      currentSongPos: (state) => state.currentSong.Pos,
    }),
  },
  watch: {
    currentSongPos: function (newSong, oldSong) {
      let oldId = "song" + oldSong;
      let el = document.getElementById("currentSong");
      if (el) el.id = oldId;
      this.ChangeCurrentSongTo(newSong);
      this.ScrollToCurrentSong();
    },
  },
};
</script>
<style>
#currentSong {
  @apply bg-red-300 dark:text-primary !important;
}
</style>
