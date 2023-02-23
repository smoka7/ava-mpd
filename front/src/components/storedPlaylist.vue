<template>
  <div class="relative flex flex-col overflow-x-hidden pt-6 pb-2">
    <PlaylistMenu
      v-if="state.selectedCount > 0"
      @delete="playlistCommand('delete')"
      @clear="playlistCommand('clear')"
      @removeDuplicate="playlistCommand('removeduplicate')"
      @removeInvalid="playlistCommand('removeinvalid')"
      @clearSelection="clearSelection()"
      @rename="state.renameOpen = true"
    />
    <button
      v-for="playlist in playlists"
      :key="playlist.name"
      class="group mx-2 flex items-center justify-between rounded p-2 text-primary duration-300 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70"
    >
      <span>
        <font-awesome-icon :icon="['fas', playlist.icon]" class="menuIcon" />
        {{ playlist.name }}
      </span>
      <span class="hide-unfocused md:text-sm">
        <sidebar-button
          :label="'add ' + playlist.name"
          icon="plus"
          @click="
            AddCmp.playlist = playlist.method;
            AddCmp.open = true;
          "
        />
        <sidebar-button
          icon="play"
          :label="'play ' + playlist.name"
          @click="
            sendCommand(endpoints.storedPlaylists, 'play' + playlist.method)
          "
        />
      </span>
    </button>
    <details v-for="(playlist, index) in storedPlaylist" :key="index">
      <summary
        class="group mx-2 flex items-start justify-between rounded py-4 px-2 duration-300 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70 md:p-2"
        @click.self="getSongs(index)"
      >
        <span
          @click="getSongs(index)"
          class="mr-1 flex w-1/2 cursor-pointer items-center overflow-x-hidden text-ellipsis py-2"
        >
          <FontAwesomeIcon
            icon="angle-right"
            class="open-rotate transform-gpu p-1 menuIcon duration-200"
            :id="'icon-' + playlist.Name"
          />
          {{ playlist.Name }}
        </span>
        <span class="flex flex-col items-end space-x-2 text-sm">
          <span @click="getSongs(index)">
            {{ playlist.SongsCount }} song ({{
              humanizeTime(playlist.Duration)
            }})
          </span>
          <span class="hide-unfocused flex">
            <sidebar-button
              label="add"
              icon="plus"
              @click="
                AddCmp.index = index;
                AddCmp.open = true;
              "
            />
            <sidebar-button
              label="play"
              icon="play"
              @click="playlistCommand('play', index)"
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
      <div class="mx-4 my-1 flex flex-col rounded bg-gray-600/20">
        <p
          class="flex items-center justify-between p-2 dark:text-white md:flex-row"
          v-for="(song, index) in playlist.songs"
          :key="index"
        >
          <span>{{ song.Title }}</span>
          <span class="text-sm text-gray-800 dark:text-gray-200">
            {{ song.Artist }} -
            {{ song.Album }}
          </span>
        </p>
      </div>
    </details>
    <rename-playlist v-if="state.renameOpen" :rename="state.renamePlaylist" />
    <AddSongs
      :open="AddCmp.open"
      @close="AddCmp.open = false"
      @add="addPlaylist($event)"
    />
  </div>
</template>
<script setup lang="ts">
import endpoints from "../endpoints";
import { humanizeTime, sendCommand } from "../helpers.js";

const store = useStore();

const storedPlaylist = computed(() =>
  shallowReactive(store.state.storedPlaylist)
);

const state = reactive({
  renameOpen: false,
  renamePlaylist: "",
  selectedCount: 0,
});

const playlists = [
  { name: "Liked Songs", icon: "heart", method: "liked" },
  { name: "Most Played Songs", icon: "chart-area", method: "most" },
];

const AddCmp = reactive({
  open: false,
  playlist: "",
  index: -1,
});

function addPlaylist(position) {
  const data = {
    Pos: position,
  };
  if (AddCmp.index > -1) {
    playlistCommand("load", AddCmp.index, position);
    AddCmp.index = -1;
    return;
  }
  const command = "add" + AddCmp.playlist;
  sendCommand(endpoints.storedPlaylists, command, data);
}

async function getSongs(index) {
  if (
    storedPlaylist.value[index].songs &&
    storedPlaylist.value[index].songs.length
  ) {
    storedPlaylist.value[index].songs = [];
    return;
  }
  storedPlaylist.value[index].songs = await sendCommand(
    endpoints.storedPlaylists,
    "list",
    { Playlist: storedPlaylist.value[index].Name }
  );
}

function playlistCommand(method, index, position) {
  if (index != null) {
    storedPlaylist.value[index].selected = true;
  }
  storedPlaylist.value.forEach((p) => {
    if (p.selected) {
      const data = {
        Playlist: p.Name,
        Pos: position,
      };
      sendCommand(endpoints.storedPlaylists, method, data);
    }
  });
  store.dispatch("getStoredPlaylist");
}

function clearSelection() {
  storedPlaylist.value.forEach((p) => {
    p.selected = false;
  });
  state.selectedCount = 0;
}

function toggleSelected(index) {
  storedPlaylist.value[index].selected
    ? state.selectedCount--
    : state.selectedCount++;
  state.renamePlaylist = storedPlaylist.value[index].Name;
  storedPlaylist.value[index].selected = !storedPlaylist.value[index].selected;
}
</script>
