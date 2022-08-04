import { createStore } from "vuex";
import { fetchOrFail, sendCommand } from "./helpers";
import endpoints from "./endpoints";

const store = createStore({
  state: {
    connected: true,
    currentSong: {},
    durationInterval: Number,
    status: {},
    albumArt: "default",
    storedPlaylist: [],
    queue: {
      Albums: [],
      Length: 0,
      Duration: 0,
    },
    serverFolders: [],
    settings: {
      Outputs: {},
      DatabaseStats: {
        albums: "",
        artists: "",
        db_playtime: "",
        db_update: "",
        playtime: "",
        songs: "",
        uptime: "",
      },
    },
  },
  mutations: {
    setCurrentSong(state, song) {
      state.currentSong = song;
    },
    setStatus(state, status) {
      state.status = status;
    },
    setQueue(state, playlist) {
      state.queue = playlist;
    },
    setStoredPlaylist(state, playlist) {
      state.storedPlaylist = playlist;
    },
    setAlbumArt(state, albumArt) {
      state.albumArt = albumArt;
    },
    setConnected(state, connected) {
      state.connected = connected;
    },
    setServerFolders(state, response) {
      state.serverFolders = [...response.Folders, ...response.Files];
    },
    setSettings(state, settings) {
      state.settings = settings;
    },
    setCounter(state) {
      clearInterval(state.durationInterval);
      if (state.status.state == "play") {
        state.durationInterval = setInterval(() => {
          state.status.elapsed = Number(state.status.elapsed) + 1;
        }, 1000);
      }
    },
  },
  actions: {
    async getCurrentSong(store) {
      const response = await fetchOrFail(endpoints.status);
      if (response.Status != null || response.CurrentSong != null) {
        store.commit("setCurrentSong", response.CurrentSong);
        store.commit("setStatus", response.Status);
        store.commit("setCounter");
        store.dispatch("getCurrentSongAlbumart");
        return;
      }
      store.dispatch("getCurrentSong");
    },
    async getCurrentSongAlbumart() {
      const albumArt = await sendCommand(endpoints.song, "albumArt", {
        ID: Number(store.state.currentSong.Id),
      });
      store.commit("setAlbumArt", albumArt.Url);
    },
    async getServerFolders() {
      const response = await sendCommand(endpoints.folders, "list", {
        File: "",
      });
      store.commit("setServerFolders", response);
    },
    async getStoredPlaylist() {
      const response = await fetchOrFail(endpoints.storedPlaylists);
      store.commit("setStoredPlaylist", response);
    },
    async getQueue(_, page = -1) {
      const response = await fetchOrFail(endpoints.queue + "?page=" + page);
      store.commit("setQueue", response);
    },
    startCounter() {
      store.commit("setCounter");
    },
    connectToSocket() {
      let hostname = new URL(window.location.href).host;

      if (import.meta.env.DEV && import.meta.env["VITE_SERVER_URL"] != "") {
        hostname = import.meta.env["VITE_SERVER_URL"];
      }

      const socket = new WebSocket("ws://" + hostname + "/update");

      socket.onmessage = (message) => {
        const data = JSON.parse(message.data);
        const event = data.Subsystem;
        if (
          event == "player" ||
          event == "options" ||
          event == "stickers" ||
          event == "update"
        ) {
          store.dispatch("getCurrentSong");
          store.commit("setCounter");
        }

        if (event == "playlist") {
          store.dispatch("getQueue");
        }
      };
      socket.onerror = (err) => {
        store.commit("setConnected", false);
        console.log(err);
      };
    },
    async getSettings() {
      const response = await fetchOrFail(endpoints.setting);
      store.commit("setSettings", response);
    },
  },
});
export default store;
