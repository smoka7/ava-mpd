<template>
  <div class="relative -m-2 flex flex-col overflow-x-hidden pt-6 pb-2">
    <playlist-menu
      v-if="selectedCount > 0"
      @addafter="PlCommand('addafter')"
      @delete="PlCommand('delete')"
      @clear="PlCommand('clear')"
      @removeDuplicate="PlCommand('removeduplicate')"
      @removeInvalid="PlCommand('removeinvalid')"
      @clearSelection="clearSelection()"
      @rename="renameOpen = true"
    />
    <div
      v-for="playlist in playlists"
      :key="playlist.name"
      class="text-primary group mx-2 flex items-center justify-between rounded p-2 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70"
    >
      <span>
        <font-awesome-icon :icon="['fas', playlist.icon]" />
        {{ playlist.name }}
      </span>
      <span class="invisible group-hover:visible md:text-sm">
        <sidebar-button
          :label="'add ' + playlist.name"
          icon="plus"
          @click="sendCommand('api/stored', 'add' + playlist.method)"
        />
        <sidebar-button
          icon="play"
          :label="'play ' + playlist.name"
          @click="sendCommand('api/stored', 'play' + playlist.method)"
        />
      </span>
    </div>
    <details v-for="(playlist, index) in storedPlaylist" :key="index">
      <summary
        class="r group mx-2 flex items-start justify-between rounded py-4 px-2 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70 md:p-2"
      >
        <span
          @click="getSongs(index)"
          class="mr-1 flex w-1/2 cursor-pointer items-center overflow-x-hidden text-ellipsis"
        >
          <FontAwesomeIcon
            icon="angle-right"
            class="mr-1 transform-gpu duration-200"
            :id="'icon-' + playlist.playlist"
          />
          {{ playlist.playlist }}
        </span>
        <span class="flex flex-col items-end space-x-2 text-sm">
          <span>
            {{ playlist.songsCount }} song ({{
              humanizeTime(playlist.duration)
            }})
          </span>
          <span class="invisible flex group-hover:visible">
            <sidebar-button
              label="add"
              icon="plus"
              @click="PlCommand('load', index)"
            />
            <sidebar-button
              label="play"
              icon="play"
              @click="PlCommand('play', index)"
            />
            <sidebar-button
              label="select playlist"
              icon="check-circle"
              :class="playlist.selected ? 'visible text-green-500' : ''"
              @click="toggleSelected(index)"
            />
          </span>
        </span>
      </summary>
      <div
        class="divide-primary border-primary ml-4 flex flex-col divide-y border-l-2 dark:divide-gray-400"
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
import SidebarButton from "./sidebarButton.vue";
export default {
  components: {
    FontAwesomeIcon,
    PlaylistMenu,
    RenamePlaylist,
    SidebarButton,
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
      if (index != null) {
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
      if (el) el.classList.toggle("rotate-90");
    },
  },
  computed: {
    ...mapState({
      storedPlaylist: (state) => state.storedPlaylist,
    }),
  },
};
</script>
