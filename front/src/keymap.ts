import { useStore } from "./store";
type PlaybackKey = {
  name: string;
  action: string;
  key: string;
};
export const playbackMappings: Array<PlaybackKey> = [
  { key: "p", name: "toggle playback (play/pause)", action: "toggle" },
  { key: "Backspace", name: "play", action: "play" },
  { key: ",", name: "previous song", action: "previous" },
  { key: ".", name: "next song", action: "next" },
  { key: "s", name: "stop", action: "stop" },
  { key: "u", name: "consume", action: "consume" },
  { key: "y", name: "single", action: "single" },
  { key: "r", name: "repeat", action: "repeat" },
  { key: "z", name: "random", action: "random" },
  { key: "c", name: "clear queue", action: "clear" },
  { key: "f", name: "seek forward 15s", action: "seekForward" },
  { key: "b", name: "seek backward 15s", action: "seekBackward" },
  { key: "9", name: "Decrease Volume by 5%", action: "volumeDown" },
  { key: "0", name: "Increase Volume by 5%", action: "volumeUp" },
];

type OtherKey = {
  name: string;
  func: Function;
  key: string;
};
export const otherMappings: Array<OtherKey> = [
  { key: "l", name: "like playing song", func: like },
  { key: "i", name: "show playing song info", func: showInfo },
];

type TabKey = {
  name: string;
  key: string;
};

export const tabMappings: Array<TabKey> = [
  { key: "1", name: "show queue tab" },
  { key: "2", name: "show setting tab" },
  { key: "3", name: "show playlists tab" },
  { key: "4", name: "show server folders tab" },
  { key: "5", name: "show folder tab" },
];

export function handleKey(key: string) {
  const store = useStore();
  const playback = playbackMappings.find((mapping) => mapping.key == key);
  if (playback) {
    store.sendPlaybackCommand(playback.action);
    return;
  }

  const func = otherMappings.find((mapping) => mapping.key == key);
  if (func) {
    func.func();
    return;
  }

  const tab = tabMappings.find((mapping) => mapping.key == key);
  if (tab) {
    changeActiveTab(Number(key));
  }
}

function like() {
  const store = useStore();
  store.toggleLike({ File: "" });
}

function changeActiveTab(tab: number) {
  const store = useStore();
  store.setActiveTab(tab - 1);
}

function showInfo() {
  const store = useStore();
  if (store.song.show) {
    store.clearSongInfo();
    return;
  }
  store.getSongInfo(Number(store.currentSong.Id));
}
