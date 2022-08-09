import store from "./store";

const PlaybackKeyMaps = {
  p: { name: "toggle", action: "toggle" },
  Backspace: { name: "play", action: "play" },
  ",": { name: "previous", action: "previous" },
  ".": { name: "next", action: "next" },
  s: { name: "stop", action: "stop" },
  u: { name: "consume", action: "consume" },
  y: { name: "single", action: "single" },
  r: { name: "repeat", action: "repeat" },
  z: { name: "random", action: "random" },
  c: { name: "clear", action: "clear" },
  f: { name: "seek forward", action: "seekForward" },
  b: { name: "seek backward", action: "seekBackward" },
  9: { name: "Decrease Volume by 5%", action: "volumeDown" },
  0: { name: "Increase Volume by 5%", action: "volumeUp" },
};
const otherMaps = {
  l: { name: "like playing song", func: like },
  i: { name: "show playing song info", func: showInfo },
};

const tabMappings = {
  1: { name: "show queue tab" },
  2: { name: "show setting Tab" },
  3: { name: "show playlists Tab" },
  4: { name: "show server folders Tab" },
  5: { name: "show folder tab" },
};

export function handleKey(key) {
  if (PlaybackKeyMaps[key]) {
    store.dispatch("sendPlaybackCommand", PlaybackKeyMaps[key].action);
    return;
  }

  if (otherMaps[key]) {
    otherMaps[key].func();
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