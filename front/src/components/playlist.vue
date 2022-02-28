<template>
  <div
    class="fbg flex h-full w-full flex-col overflow-hidden rounded bg-white/60 backdrop-blur-3xl dark:bg-gray-700/60"
  >
    <div
      v-if="queue.Length == 0"
      class="decoration-accent flex h-full w-full items-center justify-center p-4 text-7xl underline md:text-9xl"
    >
      {{ message }}
    </div>
    <div
      v-if="!songInfo"
      class="bg-secondary text-primary h-content absolute bottom-0 z-10 flex w-full flex-wrap-reverse items-center justify-between space-x-4 space-y-2 px-4 text-base md:h-6 md:text-sm"
    >
      <span>
        {{ queue.Length }} Tracks / duration:
        {{ humanizeTime(queue.Duration) }}
      </span>
      <span v-if="selectedIds.length > 0">
        {{ selectedIds.length }} selected
      </span>
      <queuePagination @goToCurrent="scrollToCurrentSong" />
    </div>
    <div class="mb-10 space-y-1 overflow-y-auto md:mb-6 md:px-2">
      <album
        v-for="(album, index) in queue.Albums"
        :key="index"
        :album="album"
        :currentAlbum="currentAlbum(album)"
        :selectedIds="selectedIds"
        :currentSongPos="currentSongPos"
        @showMenu="showMenu"
        @select="selectSong"
      />
    </div>
    <teleport to="#app">
      <queue-menu
        :open="menu"
        :song="selected"
        :selected-ids="selectedIds"
        @close="hideMenu"
        @clearSelection="selectedIds = []"
        @showInfo="songInfo = true"
      />
    </teleport>
    <songInfo v-if="songInfo" :song="selected.id" @close="songInfo = false" />
  </div>
</template>
<script>
import { mapState } from "vuex";
import { humanizeTime } from "../helpers.js";
import { defineAsyncComponent, shallowReactive } from "vue";
import album from "./album.vue";
import queuePagination from "./queuePagination.vue";
export default {
  components: {
    album,
    queuePagination,
    songInfo: defineAsyncComponent(() => import("./songInfo.vue")),
    queueMenu: defineAsyncComponent(() => import("./queueMenu.vue")),
  },
  data() {
    return {
      menu: false,
      songInfo: false,
      selected: { id: -1, pos: -1 },
      selectedIds: [],
    };
  },
  async mounted() {
    await this.$store.dispatch("getQueue");
    this.scrollToCurrentSong();
  },
  methods: {
    humanizeTime: humanizeTime,
    showMenu(pos, id) {
      this.selected.pos = Number(pos);
      this.selected.id = Number(id);
      this.menu = true;
    },
    hideMenu() {
      this.menu = false;
    },
    selectSong(id) {
      const index = this.selectedIds.indexOf(id);
      if (index > -1) {
        this.selectedIds.splice(index, 1);
        return;
      }
      this.selectedIds.push(id);
    },
    async scrollToCurrentSong() {
      if (this.queue.CurrentPage != this.queue.currentSongPage) {
        await this.$store.dispatch("getQueue", this.queue.currentSongPage);
      }
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
      message: (state) => {
        if (!state.connected) {
          return "Couldn't Connect to server!!";
        }
        if (state.queue.Length == 0) {
          return "Queue is empty!";
        }
      },
    }),
  },
  watch: {
    currentSongPos: function() {
      this.scrollToCurrentSong();
    },
  },
};
</script>
<style>
#currentSong {
  @apply dark:text-primary bg-red-300 !important;
}
</style>
