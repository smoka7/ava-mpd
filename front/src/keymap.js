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
};

export function handleKey(key) {
  if (PlaybackKeyMaps[key]) {
    store.dispatch("sendPlaybackCommand", PlaybackKeyMaps[key].action);
  }
}
