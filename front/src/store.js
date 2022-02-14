import { createStore } from "vuex";
import { fetchOrFail } from "./helpers";
import endpoints from "./endpoints";
const store = createStore({
  state: {
    currentSong: {
      Info: {},
      Liked: false,
    },
    durationInterval: Number,
    status: {},
    albumArt: "default",
    storedPlaylist: [],
    queue: {
      Albums: [],
      Length: 0,
      Duration: 0,
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
    setCounter(state) {
      clearInterval(state.durationInterval);
      state.durationInterval = setInterval(() => {
        if (state.status.state == "play") {
          state.status.elapsed = Number(state.status.elapsed) + 0.1;
        }
      }, 100);
    },
  },
  actions: {
    async getCurrentSong(store) {
      let response = await fetchOrFail(endpoints.status);
      store.commit("setCurrentSong", response.CurrentSong);
      store.commit("setStatus", response.Status);
      store.commit("setAlbumArt", response.AlbumArt);
    },
    async getStoredPlaylist() {
      let response = await fetchOrFail(endpoints.storedPlaylists);
      store.commit("setStoredPlaylist", response);
    },
    async getQueue() {
      let response = await fetchOrFail(endpoints.queue);
      store.commit("setQueue", response);
    },
    startCounter() {
      store.commit("setCounter");
    },
  },
});
export default store;
