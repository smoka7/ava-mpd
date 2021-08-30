<template>
  <div
    class="
      fixed
      bottom-0
      top-0
      left-0
      flex flex-col
      bg-white
      dark:text-white
      w-full
      h-full
      rounded
      md:static
      overflow-hidden
      app-height
      dark:bg-gray-700
    "
    id="queue"
  >
    <div class="overflow-y-auto snap-container">
      <details v-for="(album, index) in queue" :key="index" open>
        <summary
          id="album-title"
          class="
            flex
            px-8
            py-2
            sticky
            top-0
            text-primary
            bg-lightest
            items-center
          "
          @click="animate(album[0].Pos)"
        >
          <FontAwesomeIcon
            icon="angle-down"
            class="mr-2 transform-gpu duration-200"
            :id="'icon-' + album[0].Pos"
          ></FontAwesomeIcon>
          {{ album[0].Artist }} - {{ album[0].Album }} ({{ album[0].Date }})
        </summary>
        <div
          @mouseenter="show(song.Pos)"
          @mouseleave="show(song.Pos)"
          @click="showMenu(song.Pos, song.file, $event)"
          v-for="song in album"
          :key="song.Pos"
          class="
            py-1
            px-4
            grid grid-cols-12
            snap-child
            hover:bg-blue-100
            dark:hover:text-black
          "
          :id="'song' + song.Pos"
        >
          <span
            class="
              col-start-1 col-end-3
              md:col-end-2
              justify-self-start
              pr-2
              w-full
              flex
              items-center
            "
            @click.stop="play(song.Pos)"
          >
            <FontAwesomeIcon
              :id="song.Pos"
              v-if="song.Pos !== currentSong.Pos"
              class="invisible text-green-500 mr-2 cursor-pointer"
              icon="play"
            ></FontAwesomeIcon>
            <FontAwesomeIcon
              v-else
              class="text-red-500 mr-2"
              icon="compact-disc"
            ></FontAwesomeIcon>
            {{ song.Track }}
          </span>
          <span class="col-start-3 md:col-start-2 col-end-11 self-start">
            {{ song.Title }}
          </span>
          <span
            class="col-start-11 md:col-start-12 col-end-13 justify-self-center"
          >
            <FontAwesomeIcon
              :id="'delete' + song.Pos"
              class="text-primary mr-2 invisible"
              icon="ellipsis-h"
            ></FontAwesomeIcon>
            {{ humanizeTime(song.duration) }}
          </span>
        </div>
      </details>
    </div>
    <button
      aria-label="close-playlist"
      @click="closePlaylist"
      class="
        fixed
        right-4
        bottom-4
        bg-red-500
        text-white
        rounded-full
        p-2
        md:hidden
        w-20
        h-20
      "
    >
      <font-awesome-icon icon="arrow-right" size="2x"></font-awesome-icon>
    </button>
    <div
      v-show="menu"
      class="bg-transparent absolute top-0 bottom-0 left-0 right-0"
      @click.self="hideMenu"
    >
      <div
        id="context-menu"
        class="
          flex flex-col
          bg-white
          dark:bg-gray-600
          absolute
          z-50
          rounded
          border-2 border-primary
          cursor-pointer
        "
      >
        <div
          class="p-2 hover:bg-blue-100 dark:hover:text-black"
          @click="deleteSong(selected)"
        >
          delete
        </div>
        <details>
          <summary
            @click="getStoredPlaylist"
            class="flex p-2 hover:bg-blue-100 dark:hover:text-black"
          >
            add to playlist
          </summary>
          <div class="border-primary border-t-2 border-b-2">
            <p
              @click="addSongTo(pl.playlist)"
              class="px-4 py-1 hover:bg-blue-100 dark:border-gray-800"
              v-for="pl in storedPlaylist"
              :key="pl.playlist"
            >
              {{ pl.playlist }}
            </p>
          </div>
        </details>
        <div class="p-2 hover:bg-blue-100 dark:hover:text-black">info</div>
      </div>
    </div>
  </div>
</template>
<script>
import { mapState } from "vuex";
import { sendCommand, humanizeTime, toggleMediaController } from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
export default {
  components: {
    FontAwesomeIcon,
  },
  data() {
    return {
      menu: false,
      selected: -1,
      uri: "",
    };
  },
  async mounted() {
    setTimeout(() => {
      let el = document.getElementById("song" + this.currentSong.Pos);
      if (el !== null) {
        el.classList.add("bg-red-100");
        el.classList.add("dark:text-black");
        el.scrollIntoView({ block: "center", behavior: "smooth" });
      }
    }, 1000);
  },
  methods: {
    humanizeTime: humanizeTime,
    closePlaylist: toggleMediaController,
    getStoredPlaylist() {
      if (this.storedPlaylist.length) return;
      this.$store.dispatch("getStoredPlaylist");
    },
    addSongTo(playlist) {
      sendCommand("api/queue", "addsong", {
        song: this.uri,
        playlist: playlist,
      });
    },
    show(id) {
      let el = document.getElementById(id);
      if (el) el.classList.toggle("invisible");
      el = document.getElementById("delete" + id);
      if (el) el.classList.toggle("invisible");
    },
    showMenu(pos, uri, e) {
      (this.selected = pos), (this.uri = uri);
      let el = document.getElementById("context-menu");
      el.style.left = e.pageX + "px";
      el.style.top = e.pageY + "px";
      this.menu = true;
    },
    hideMenu() {
      this.selected = { position: -1, file: "" };
      this.menu = false;
    },
    animate(id) {
      let el = document.getElementById("icon-" + id);
      if (el) el.classList.toggle("-rotate-90");
    },
    play(id) {
      let data = {
        Start: Number(id),
      };
      sendCommand("/api/queue", "play", data);
    },
    deleteSong(start = -1, end = -1) {
      let data = {
        Start: Number(start),
        End: Number(end),
      };
      sendCommand("/api/queue", "delete", data);
    },
  },
  computed: {
    ...mapState({
      queue: (state) => state.queue,
      storedPlaylist: (state) => state.storedPlaylist,
      currentSong: (state) => state.currentSong,
    }),
  },
  watch: {
    currentSong: function (newSong, oldSong) {
      let oldId = "song" + oldSong.Pos;
      let newId = "song" + newSong.Pos;
      let el = document.getElementById(oldId);
      if (el) {
        el.classList.remove("bg-red-100");
        el.classList.remove("dark:text-black");
      }
      el = document.getElementById(newId);
      if (el) {
        el.classList.toggle("dark:text-black");
        el.classList.toggle("bg-red-100");
        el.scrollIntoView({ block: "center", behavior: "smooth" });
      }
    },
  },
};
</script>
