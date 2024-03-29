import { fetchOrFail, sendCommand } from "./helpers";
import endpoints from "./endpoints";
import { settingsStore } from "./settingsStore";

type StatusResponse = {
  Status: Status;
  CurrentSong: CurrentSong;
};

type SongCommandData = {
  ID?: number;
  File?: string;
};

type CurrentSong = {
  Album: string;
  Artist: string;
  Date: string;
  Genre: string;
  Id: number;
  Pos: number;
  Title: string;
  File: string;
  Liked: boolean;
};

type Status = {
  state: "pause" | "play" | "stop";
  duration: number;
  mixrampdb?: number;
  xfade?: number;
  updating_db?: boolean;
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

export type Album = {
  Songs: Songs;
  Album: string;
  Artist: string;
  Date: string;
};

type Song = {
  Album: string;
  Artist: string;
  Title: string;
  File: string;
  Liked: boolean;
  Pos: number;
  Id: number;
  Track: number;
  Duration: number;
};

export type Songs = Array<Song>;

export type File = { File: string };

export type Directory = { Directory: string };

export type FoldersResponse = {
  Directories: Array<Directory>;
  Files: Array<File>;
};

type SongInfo = {
  album: string;
  albumArtist: string;
  albumArtistSort: string;
  artist: string;
  artistSort: string;
  composer: string;
  date: string;
  disc: string;
  format: string;
  genre: string;
  label: string;
  "Last-Modified": string;
  MUSICBRAINZ_ALBUMARTISTID: string;
  MUSICBRAINZ_ALBUMID: string;
  MUSICBRAINZ_ARTISTID: string;
  MUSICBRAINZ_RELEASETRACKID: string;
  MUSICBRAINZ_TRACKID: string;
  MUSICBRAINZ_WORKID: string;
  OriginalDate: string;
  time: string;
  title: string;
  track: string;
  duration: string;
  file: string;
};

type AlbumArtResponse = { Url: string };

type SongInfoResponse = {
  Info: SongInfo;
  Stickers: Array<Stickers>;
  liked: boolean;
  albumArt: string;
  show: boolean;
};

type Stickers = {
  Name: string;
  Value: string;
};

export type StoredPlaylist = {
  Name: string;
  Selected?: boolean;
  SongsCount: number;
  Duration: number;
  Songs: Array<{ Album: string; Title: string; Artist: string }>;
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
      song: {} as SongInfoResponse,
      serverFolders: {} as FoldersResponse,
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
    setAlbumArt(albumArt: string) {
      this.albumArt = albumArt;
    },
    setConnected(connected: boolean) {
      this.connected = connected;
      if (!connected) {
        this.queue = {} as Queue;
      }
    },
    setServerFolders(response: FoldersResponse) {
      this.serverFolders = response;
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
      const albumArt = await sendCommand<AlbumArtResponse>(
        endpoints.song,
        "albumArt",
        {
          ID: Number(this.currentSong.Id),
        }
      );
      if (albumArt) this.setAlbumArt(albumArt.Url);
    },
    async getServerFolders() {
      const response = await sendCommand<FoldersResponse>(
        endpoints.folders,
        "list",
        {
          File: "",
        }
      );
      if (response) this.setServerFolders(response);
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
          this.getQueue(this.queue.CurrentPage);
          this.setCounter();
        }

        if (event == "output" || event == "database" || event == "options") {
          settingsStore().getSettings();
        }

        if (event == "stored_playlist") {
          this.getStoredPlaylist();
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
      const song = await sendCommand<SongInfoResponse>(endpoints.song, "info", {
        ID: songId,
      });
      const albumArt = await sendCommand<AlbumArtResponse>(
        endpoints.song,
        "albumArt",
        {
          ID: songId,
        }
      );
      if (song == null || albumArt == null) return;
      if (song.Info == null && song.Stickers == null) return;

      this.song.Info = song.Info;
      this.song.Stickers = song.Stickers;
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
    async toggleLike(data: SongCommandData) {
      sendCommand(endpoints.song, "like", data);
    },
  },
});
