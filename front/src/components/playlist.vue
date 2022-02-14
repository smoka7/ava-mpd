<template>
  <div
    class="fbg flex h-full w-full flex-col overflow-hidden rounded bg-white/60 backdrop-blur-3xl dark:bg-gray-700/60"
  >
    <div
      v-if="queue.Length == 0"
      class="flex h-full w-full items-center justify-center p-4 text-9xl underline decoration-accent"
    >
      Queue is empty!
    </div>
    <div
      v-if="!songInfo && queue.Length != 0"
      class="absolute bottom-0 z-10 flex h-10 w-full items-center justify-between space-x-4 bg-secondary px-4 text-base text-primary md:h-6 md:text-sm"
    >
      <span>{{ queue.Length }} Tracks </span>
      <span>duration: {{ humanizeTime(queue.Duration) }} </span>
      <button
        aria-label="goto-current-song"
        @click="ScrollToCurrentSong()"
        class="px-2 hover:bg-primary hover:text-secondary"
      >
        <font-awesome-icon icon="compact-disc" />
        Current Song
      </button>
      <button
        aria-label="close-playlist"
        @click="closePlaylist"
        class="px-2 hover:bg-primary hover:text-secondary md:hidden"
      >
        <font-awesome-icon icon="times" size="2x" />
      </button>
    </div>
    <div class="mb-10 space-y-1 overflow-y-auto md:mb-6 md:px-2">
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
      const el = document.getElementById("context-menu");
      el.style.top = e.pageY + "px";
      this.menu = true;
    },
    hideMenu() {
      this.menu = false;
    },
    ChangeCurrentSongTo(id) {
      const el = document.getElementById("song" + id);
      if (el !== null) {
        el.id = "currentSong";
      }
    },
    ScrollToCurrentSong() {
      const el = document.getElementById("currentSong");
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
      currentSongPos: (state) => state.currentSong.Info.Pos,
    }),
  },
  watch: {
    currentSongPos: function (newSong, oldSong) {
      const oldId = "song" + oldSong;
      const el = document.getElementById("currentSong");
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
