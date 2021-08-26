import { createStore } from "vuex";
const store = createStore({
  state: {
    currentSong: Object,
    durationInterval: Number,
    status: Object,
    albumArt: "default",
    storedPlaylist: Array,
    Queue: Array,
  },
  mutations: {
    setCurrentSong(state, song) {
      state.currentSong = song;
    },
    setStatus(state, status) {
      state.status = status;
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
    startCounter() {
      store.commit("setCounter");
    },
  },
});
export default store;
