<template>
  <div
    v-show="open"
    class="bg-white bg-opacity-30 absolute top-0 bottom-0 left-0 right-0 z-50"
    @click.self="$emit('close')"
  ></div>
  <div
    v-show="open"
    class="
      flex flex-col
      bg-white
      dark:bg-gray-600
      absolute
      right-10
      z-50
      rounded
      border-2 border-primary
      dark:border-lightest
    "
    id="context-menu"
  >
    <div v-for="action in actions" :key="action.title">
      <button :class="btnClass" @click="action.func()">
        {{ action.title }}
      </button>
    </div>
    <div>
      <details>
        <summary @click="getStoredPlaylist" :class="[btnClass, 'flex']">
          add to playlist
        </summary>
        <div
          class="
            border-primary border-t-2
            dark:border-lightest dark:hover:text-black
          "
        >
          <button
            @click="addSongTo(pl.playlist)"
            :class="btnClass"
            v-for="pl in storedPlaylist"
            :key="pl.playlist"
          >
            {{ pl.playlist }}
          </button>
        </div>
      </details>
    </div>
  </div>
</template>
<script>
import { mapState } from "vuex";
import { sendCommand } from "../helpers.js";
export default {
  props: {
    open: Boolean,
    song: {
      position: Number,
      file: String,
    },
  },
  emits: ["close", "showInfo"],
  data() {
    return {
      actions: [
        { title: "delete", func: this.deleteSong },
        { title: "info", func: this.showInfo },
        { title: "remove duplicate songs", func: this.removeDuplicate },
      ],
      btnClass:
        "w-full text-left p-2 hover:bg-blue-100 focus:bg-blue-100 dark:hover:text-black",
    };
  },
  methods: {
    deleteSong(end = -1) {
      let data = {
        Start: Number(this.song.position),
        End: Number(end),
      };
      sendCommand("/api/queue", "delete", data);
      this.$emit("close");
    },
    addSongTo(playlist) {
      sendCommand("api/queue", "addsong", {
        song: this.song.file,
        playlist: playlist,
      });
      this.$emit("close");
    },
    showInfo() {
      this.$emit("showInfo");
      this.$emit("close");
    },
    removeDuplicate() {
      sendCommand("api/queue", "removeduplicate");
      this.$emit("close");
    },
    getStoredPlaylist() {
      if (this.storedPlaylist.length) return;
      this.$store.dispatch("getStoredPlaylist");
    },
  },
  computed: {
    ...mapState({
      storedPlaylist: (state) => state.storedPlaylist,
    }),
  },
};
</script>
