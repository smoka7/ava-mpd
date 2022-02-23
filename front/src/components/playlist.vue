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
      class="bg-secondary text-primary h-content absolute bottom-0 z-10 flex flex-wrap-reverse w-full items-center justify-between space-x-4 space-y-2 px-4 text-base md:h-6 md:text-sm"
    >
      <span>{{ queue.Length }} Tracks </span>
      <span v-if="selectedIds.length > 0">
        {{ selectedIds.length }} selected
      </span>
      <span>duration: {{ humanizeTime(queue.Duration) }} </span>
      <button
        aria-label="goto-current-song"
        @click="scrollToCurrentSong()"
        class="hover:bg-primary hover:text-secondary px-2"
      >
        <font-awesome-icon icon="compact-disc" />
        Current Song
      </button>
      <button
        aria-label="close-playlist"
        @click="closePlaylist"
        class="hover:bg-primary hover:text-secondary px-2 md:hidden"
      >
        <font-awesome-icon icon="times" class="scale-125"/>
      </button>
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
      selected: { id: -1, pos: -1 },
      selectedIds: [],
    };
  },
  async mounted() {
    await this.$store.dispatch("getQueue");
    this.changeCurrentSongTo(this.currentSongPos);
    this.scrollToCurrentSong();
  },
  methods: {
    humanizeTime: humanizeTime,
    closePlaylist: toggleMediaController,
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
    changeCurrentSongTo(id) {
      const el = document.getElementById("song" + id);
      if (el !== null) {
        el.id = "currentSong";
      }
    },
    scrollToCurrentSong() {
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
    currentSongPos: function(newSong, oldSong) {
      const oldId = "song" + oldSong;
      const el = document.getElementById("currentSong");
      if (el) el.id = oldId;
      this.changeCurrentSongTo(newSong);
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
