import { createStore } from "vuex";
const store = createStore({
  state: {
    currentSong: Object,
    durationInterval: Number,
    status: Object,
    albumArt: "default",
    storedPlaylist: [],
    queue: {
            Albums: [],
      Length: 0,
            Duration:0,
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
      let response = await fetch("/api/status");
      if (response.ok) {
        let json = await response.json();
        store.commit("setCurrentSong", json.CurrentSong);
        store.commit("setStatus", json.Status);
        store.commit("setAlbumArt", json.AlbumArt);
        return;
      }
      console.log(response.error);
    },
    async getStoredPlaylist() {
      let response = await fetch("/api/stored");
      if (response.ok) {
        let json = await response.json();
        store.commit("setStoredPlaylist", json);
        return;
      }
      console.log(response.error);
    },
    async getQueue() {
      let response = await fetch("/api/queue");
      if (response.ok) {
        let json = await response.json();
        store.commit("setQueue", json);
        return;
      }
      console.log(response.error);
    },
    startCounter() {
      store.commit("setCounter");
    },
  },
});
export default store;
