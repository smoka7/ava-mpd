import { fetchOrFail, sendCommand } from "./helpers";
import endpoints from "./endpoints";
type StatusResponse = {
  Status: Status;
  CurrentSong: CurrentSong;
};
type SongCommandData = {
  ID: number;
  File: string;
};
type CurrentSong = {
  Album: string;
  Artist: string;
  Date: string;
  Genre: string;
  Id: string;
  Pos: string;
  Title: string;
  File: string;
  Liked: boolean;
};
type Status = {
  state: "pause" | "play" | "stop";
  duration: number;
  mixrampdb?: number;
  xfade?: number;
  updatingDB?: boolean;
  volume: number;
  elapsed: number;
  consume: boolean;
  random: boolean;
  repeat: boolean;
  single: boolean;
};
type Queue = {
  Length: number;
  CurrentSongPage: number;
  CurrentPage: number;
  LastPage: number;
  Duration: number;
  Albums: Array<Album>;
};
type Album = {
  Songs: Array<Song>;
  Album: string;
  Artist: string;
  Date: string;
};
type Song = {
  Album: string;
  Artist: string;
  Title: string;
  Pos: string;
  Id: string;
  Track: string;
  Duration: string;
};
type StoredPlaylist = {
  Name: string;
  SongsCount: number;
  Duration: number;
};
type DatabaseStats = {
  albums: string;
  artists: string;
  db_playtime: string;
  db_update: string;
  playtime: string;
  songs: string;
  uptime: string;
};

type SettingsResponse = {
  Outputs: Array<Array<string>>;
  DatabaseStats: DatabaseStats;
  ReplayGain: string;
  DownloadCoverArt: boolean;
};
export const useStore = defineStore("main", {
  state: () => {
    return {
      connected: true,
      currentSong: {} as CurrentSong,
      activeTabIndex: 0,
      durationInterval: 0 as number,
      status: {} as Status,
      albumArt: "default",
      storedPlaylist: [] as Array<StoredPlaylist>,
      queue: {} as Queue,
      song: {
        info: {},
        stickers: {},
        albumArt: "",
        liked: false,
        show: false,
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
    };
  },
  actions: {
    setCurrentSong(song: CurrentSong) {
      this.currentSong = song;
    },
    setActiveTab(tab: number) {
      this.activeTabIndex = tab;
    },
    setStatus(status: Status) {
      this.status = status;
      this.connected = true;
    },
    setQueue(playlist: Queue) {
      this.queue = playlist;
    },
    setStoredPlaylist(playlist: Array<StoredPlaylist>) {
      this.storedPlaylist = playlist;
    },
    setAlbumArt(albumArt) {
      this.albumArt = albumArt;
    },
    setConnected(connected: boolean) {
      this.connected = connected;
      if (!connected) {
        this.queue = {} as Queue;
      }
    },
    setServerFolders(response) {
      this.serverFolders = [...response.Folders, ...response.Files];
    },
    setSettings(settings: SettingsResponse) {
      this.settings = settings;
    },
    setCounter() {
      clearInterval(this.durationInterval);
      if (this.status.state == "play") {
        this.durationInterval = setInterval(() => {
          this.status.elapsed = Number(this.status.elapsed) + 1;
        }, 1000);
      }
    },
    async getCurrentSong() {
      const response = await fetchOrFail<StatusResponse>(endpoints.status);
      if (response.Status != null || response.CurrentSong != null) {
        this.setCurrentSong(response.CurrentSong);
        this.setStatus(response.Status);
        this.setCounter();
        this.getCurrentSongAlbumart();
        return;
      }
      this.getCurrentSong();
    },
    async getCurrentSongAlbumart() {
      const albumArt = await sendCommand(endpoints.song, "albumArt", {
        ID: Number(this.currentSong.Id),
      });
      this.setAlbumArt(albumArt.Url);
    },
    async getServerFolders() {
      const response = await sendCommand(endpoints.folders, "list", {
        File: "",
      });
      this.setServerFolders(response);
    },
    async getStoredPlaylist() {
      const response = await fetchOrFail<Array<StoredPlaylist>>(
        endpoints.storedPlaylists
      );
      this.setStoredPlaylist(response);
    },
    async getQueue(page = -1) {
      const response = await fetchOrFail<Queue>(
        endpoints.queue + "?page=" + page
      );
      if (response) this.setQueue(response);
    },
    async sendPlaybackCommand(command: string) {
      sendCommand(endpoints.playback, command);
    },
    startCounter() {
      this.setCounter();
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
          event == "sticker" ||
          event == "mixer" ||
          event == "update"
        ) {
          this.getCurrentSong();
          this.setCounter();
        }

        if (event == "playlist") {
          this.getQueue();
        }
      };

      socket.onerror = (err) => {
        this.setConnected(false);
        console.log(err);
      };

      socket.onclose = () => {
        this.setConnected(false);
      };
    },
    async getSongInfo(songId: number) {
      const song = await sendCommand(endpoints.song, "info", {
        ID: songId,
      });
      const albumArt = await sendCommand(endpoints.song, "albumArt", {
        ID: songId,
      });

      if (song.Info == null && song.Stickers == null) return;

      this.song.info = song.Info;
      this.song.stickers = song.Stickers;
      this.song.albumArt = albumArt.Url;
      this.song.show = true;

      const index = song.Stickers.findIndex((stick) => {
        return stick.Name == "liked" && stick.Value == "true";
      });
      this.song.liked = index > -1;
    },
    clearSongInfo() {
      this.song.show = false;
    },
    async getSettings() {
      const response = await fetchOrFail<SettingsResponse>(endpoints.setting);
      this.setSettings(response);
    },
    async sendCommandToSetting(payload: {
      command: string;
      data?: { Value: number };
    }) {
      sendCommand(endpoints.setting, payload.command, payload.data);
      this.getSettings();
    },
    async toggleLike(data: SongCommandData) {
      sendCommand(endpoints.song, "like", data);
    },
  },
});
