<template>
  <div class="flex flex-col relative -m-2 pt-6">
    <playlist-menu
      v-if="selectedCount > 0"
      @delete="PlCommand('delete')"
      @clear="PlCommand('clear')"
      @clearSelection="clearSelection()"
      @rename="renameOpen = true"
    />
    <div
      v-for="playlist in playlists"
      :key="playlist.name"
      class="flex justify-between p-2 mx-2 items-center text-primary dark:text-white hover:bg-white/60 dark:hover:bg-gray-800/70 rounded group"
    >
      <span>
        <font-awesome-icon :icon="['fas', playlist.icon]" />
        {{ playlist.name }}
      </span>
      <span class="md:text-sm invisible group-hover:visible">
        <button
          aria-label="add-most-played-song"
          class="sidebar-btn"
          @click="sendCommand('api/stored', 'add' + playlist.method)"
        >
          <font-awesome-icon icon="plus" />
        </button>
        <button
          aria-label="play-most-played-song"
          class="sidebar-btn"
          @click="sendCommand('api/stored', 'play' + playlist.method)"
        >
          <font-awesome-icon icon="play" />
        </button>
      </span>
    </div>
    <details v-for="(playlist, index) in storedPlaylist" :key="index">
      <summary
        class="flex justify-between items-center md:p-2 py-4 px-2 mx-2 dark:text-white hover:bg-white/60 dark:hover:bg-gray-800/70 r rounded group"
      >
        <span @click="getSongs(index)" class="cursor-pointer">
          <FontAwesomeIcon
            icon="angle-down"
            class="transform-gpu -rotate-90 duration-200"
            :id="'icon-' + playlist.playlist"
          />
          {{ playlist.playlist }}
        </span>
        <span class="text-sm space-x-2">
          <span>
            {{ playlist.songsCount }} song ({{
              humanizeTime(playlist.duration)
            }})
          </span>
          <span class="invisible group-hover:visible">
            <button
              aria-label="select-playlist"
              class="sidebar-btn"
              @click="toggleSelected(index)"
            >
              <font-awesome-icon
                icon="check-circle"
                :class="playlist.selected ? 'text-green-500 visible' : ''"
              />
            </button>
            <button
              aria-label="add"
              class="sidebar-btn"
              @click="PlCommand('load', index)"
            >
              <font-awesome-icon icon="plus" />
            </button>
            <button
              aria-label="play"
              class="sidebar-btn"
              @click="PlCommand('play', index)"
            >
              <font-awesome-icon icon="play" />
            </button>
          </span>
        </span>
      </summary>
      <div
        class="divide-y divide-primary dark:divide-gray-400 flex flex-col ml-4 border-l-2 border-primary"
      >
        <p
          class="px-2 dark:text-white"
          v-for="(song, index) in playlist.songs"
          :key="index"
        >
          {{ song.Title }}<br />
          <span class="text-sm text-gray-800 dark:text-gray-400">
            {{ song.Artist }} -
            {{ song.Album }}
          </span>
        </p>
      </div>
    </details>
    <rename-playlist v-if="renameOpen" :rename="renamePlaylist" />
  </div>
</template>
<script>
import { humanizeTime, sendCommand } from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { mapState } from "vuex";
import PlaylistMenu from "./playlistMenu.vue";
import RenamePlaylist from "./renamePlaylist.vue";
export default {
  components: {
    FontAwesomeIcon,
    PlaylistMenu,
    RenamePlaylist,
  },
  data() {
    return {
      renameOpen: false,
      renamePlaylist: "",
      selectedCount: 0,
      playlists: [
        { name: "Liked Songs", icon: "heart", method: "liked" },
        { name: "Most Played Songs", icon: "chart-area", method: "most" },
      ],
    };
  },
  created() {
    this.$store.dispatch("getStoredPlaylist");
  },
  methods: {
    humanizeTime: humanizeTime,
    sendCommand: sendCommand,
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
    PlCommand(method, index) {
      if (index !== null) {
        this.storedPlaylist[index].selected = true;
      }
      this.storedPlaylist.forEach((p) => {
        if (p.selected) {
          let data = {
            playlist: p.playlist,
          };
          sendCommand("api/stored", method, data);
        }
      });
      this.$store.dispatch("getStoredPlaylist");
    },
    clearSelection() {
      this.storedPlaylist.forEach((p) => {
        p.selected = false;
      });
    },
    toggleSelected(index) {
      this.storedPlaylist[index].selected
        ? this.selectedCount--
        : this.selectedCount++;
      this.renamePlaylist = this.storedPlaylist[index].playlist;
      this.storedPlaylist[index].selected =
        !this.storedPlaylist[index].selected;
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
