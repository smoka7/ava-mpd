<template>
  <div class="flex flex-col">
    <div
      class="
        flex
        justify-between
        p-2
        text-primary
        dark:text-white
        hover:bg-blue-200
        rounded
      "
    >
      <span>
        <font-awesome-icon icon="heart"></font-awesome-icon>
        liked Songs
      </span>
      <span class="space-x-2 text-sm">
        <button aria-label="add-liked-song">
          <font-awesome-icon @click="addLiked" icon="plus"></font-awesome-icon>
        </button>
        <button aria-label="play-liked-song">
          <font-awesome-icon @click="playLiked" icon="play"></font-awesome-icon>
        </button>
      </span>
    </div>
    <div
      class="
        flex
        justify-between
        p-2
        text-primary
        dark:text-white
        hover:bg-blue-200
        rounded
      "
    >
      <span>
        <font-awesome-icon icon="chart-area"></font-awesome-icon>
        most Played Songs
      </span>
      <span class="space-x-2 text-sm">
        <button aria-label="add-most-played-song">
          <font-awesome-icon
            @click="addMostPlayed"
            icon="plus"
          ></font-awesome-icon>
        </button>
        <button aria-label="play-most-played-song">
          <font-awesome-icon
            @click="playMostPlayed"
            icon="play"
          ></font-awesome-icon>
        </button>
      </span>
    </div>
    <details v-for="(playlist, index) in storedPlaylist" :key="index">
      <summary
        class="
          flex
          justify-between
          md:p-2
          py-4
          px-2
          text-primary
          cursor-pointer
          border-primary
          dark:text-white dark:border-gray-800
          hover:bg-blue-200
          rounded
        "
      >
        <span @click="getSongs(index)">
          <FontAwesomeIcon
            icon="angle-down"
            class="transform-gpu -rotate-90 duration-200"
            :id="'icon-' + playlist.playlist"
          ></FontAwesomeIcon>
          {{ playlist.playlist }}
        </span>
        <span class="text-sm space-x-2">
          {{ playlist.songsCount }} song ({{ humanizeTime(playlist.duration) }})
          <button aria-label="add">
            <font-awesome-icon
              @click="addPlaylist(index)"
              icon="plus"
            ></font-awesome-icon>
          </button>
          <button aria-label="play">
            <font-awesome-icon
              @click="PlayPlaylist(index)"
              icon="play"
            ></font-awesome-icon>
          </button>
        </span>
      </summary>
      <div class="space-y-1 divide-y dark:divide-gray-800 flex flex-col">
        <p class="px-2" v-for="(song, index) in playlist.songs" :key="index">
          {{ song.Title }}<br />
          <span class="text-sm text-gray-800 dark:text-gray-100">
            {{ song.Artist }}
            {{ song.Album }}
          </span>
        </p>
      </div>
    </details>
  </div>
</template>
<script>
import { humanizeTime, sendCommand } from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { mapState } from "vuex";
export default {
  components: {
    FontAwesomeIcon,
  },
  created() {
    this.$store.dispatch("getStoredPlaylist");
  },
  methods: {
    humanizeTime: humanizeTime,
    async getSongs(index) {
      this.animate(this.storedPlaylist[index].playlist);
      if (
        this.storedPlaylist[index].songs &&
        this.storedPlaylist[index].songs.length
      ) {
        this.storedPlaylist[index].songs = [];
        return;
      }
      let request = {
        command: "list",
        data: {
          playlist: this.storedPlaylist[index].playlist,
        },
      };
      let response = await fetch("/api/stored", {
        method: "POST",
        headers: {
          "Content-Type": "application/json;charset=utf-8",
        },
        body: JSON.stringify(request),
      });
      if (response.ok) {
        let json = await response.json();
        this.storedPlaylist[index].songs = json;
        return;
      }
      console.log(response.error);
    },
    PlayPlaylist(index) {
      let data = {
        playlist: this.storedPlaylist[index].playlist,
      };
      sendCommand("api/stored", "play", data);
    },
    playLiked() {
      sendCommand("api/stored", "playliked");
    },
    addLiked() {
      sendCommand("api/stored", "addliked");
    },
    playMostPlayed() {
      sendCommand("api/stored", "playmost");
    },
    addMostPlayed() {
      sendCommand("api/stored", "addmost");
    },
    addPlaylist(index) {
      let data = {
        playlist: this.storedPlaylist[index].playlist,
      };
      sendCommand("api/stored", "load", data);
    },
    animate(id) {
      let el = document.getElementById("icon-" + id);
      if (el) el.classList.toggle("-rotate-90");
    },
  },
  computed: {
    ...mapState({
      storedPlaylist: (state) => state.storedPlaylist,
    }),
  },
};
</script>
