<template>
  <div
    v-show="open"
    class="absolute top-0 bottom-0 left-0 right-0 z-50 bg-white bg-opacity-30"
    @click.self="$emit('close')"
  ></div>
  <div
    v-show="open"
    class="border-primary dark:border-lightest absolute right-10 z-50 flex w-48 flex-col rounded border-2 bg-white text-lg dark:bg-gray-600 dark:text-white"
    id="context-menu"
  >
    <div v-for="action in actions" :key="action.title">
      <button :class="btnClass" @click="action.func()">
        {{ action.title }}
      </button>
    </div>
    <details class="w-full">
      <summary @click="getStoredPlaylist" :class="[btnClass, 'flex']">
        add to playlist
      </summary>
      <div
        class="border-primary dark:border-lightest border-t-2 dark:hover:text-black"
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
</template>
<script>
import { mapState } from "vuex";
import { sendCommand } from "../helpers.js";
import endpoints from "../endpoints.js";
export default {
  props: {
    open: Boolean,
    song: Number,
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
        "w-full text-left p-2 hover:bg-blue-100 focus:bg-blue-100 dark:hover:text-black overflow-x-hidden text-ellipsis",
    };
  },
  methods: {
    deleteSong(end = -1) {
      let data = {
        Start: this.song,
        End: Number(end),
      };
      sendCommand(endpoints.queue, "delete", data);
      this.$emit("close");
    },
    addSongTo(playlist) {
      sendCommand(endpoints.queue, "addsong", {
        Start: this.song,
        playlist: playlist,
      });
      this.$emit("close");
    },
    showInfo() {
      this.$emit("showInfo");
      this.$emit("close");
    },
    removeDuplicate() {
      sendCommand(endpoints.queue, "removeduplicate");
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
