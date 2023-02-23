import store from "./store";

export const playbackMappings = {
  p: { name: "toggle playback (play/pause)", action: "toggle" },
  Backspace: { name: "play", action: "play" },
  ",": { name: "previous song", action: "previous" },
  ".": { name: "next song", action: "next" },
  s: { name: "stop", action: "stop" },
  u: { name: "consume", action: "consume" },
  y: { name: "single", action: "single" },
  r: { name: "repeat", action: "repeat" },
  z: { name: "random", action: "random" },
  c: { name: "clear queue", action: "clear" },
  f: { name: "seek forward 15s", action: "seekForward" },
  b: { name: "seek backward 15s", action: "seekBackward" },
  9: { name: "Decrease Volume by 5%", action: "volumeDown" },
  0: { name: "Increase Volume by 5%", action: "volumeUp" },
};
export const otherMappings = {
  l: { name: "like playing song", func: like },
  i: { name: "show playing song info", func: showInfo },
};

export const tabMappings = {
  1: { name: "show queue tab" },
  2: { name: "show setting tab" },
  3: { name: "show playlists tab" },
  4: { name: "show server folders tab" },
  5: { name: "show folder tab" },
};

export function handleKey(key) {
  if (playbackMappings[key]) {
    store.dispatch("sendPlaybackCommand", playbackMappings[key].action);
    return;
  }

  if (otherMappings[key]) {
    otherMappings[key].func();
    return;
  }

  if (tabMappings[key]) {
    changeActiveTab(key);
  }
}

function like() {
  store.dispatch("toggleLike", { File: "" });
}

function changeActiveTab(tab) {
  store.commit("setActiveTab", tab - 1);
}

function showInfo() {
  if (store.state.song.show) {
    store.dispatch("clearSongInfo", Number(store.state.currentSong.Id));
    return;
  }
  store.dispatch("getSongInfo", Number(store.state.currentSong.Id));
}
