<template>
  <div
    class="
      fixed
      bottom-0
      grid grid-rows-8 grid-cols-1
      w-full
      h-screen
      p-4
      md:h-36 md:grid-rows-3 md:grid-cols-12 md:p-4
      bg-primary
      text-foreground
      z-20
    "
    id="mediaController"
  >
    <div
      class="
        md:col-start-1 md:col-end-3 md:row-start-1 md:row-end-3
        row-start-7 row-end-9
        w-full
        md:w-auto
        flex
        justify-around
      "
    >
      <button
        aria-label="previous-song"
        @click="sendCommand('/api/playback', 'previous')"
        class="text-lightest py-7 px-8 rounded-full md:p-3"
      >
        <font-awesome-icon icon="step-backward" size="lg" />
      </button>
      <button
        aria-label="toggle-playback"
        @click="sendCommand('/api/playback', 'toggle')"
        :class="[
          status.state === 'pause' ? 'bg-green-400' : 'bg-accent',
          'text-white md:py-7 md:px-7 rounded-full px-12 flex items-center',
        ]"
      >
        <font-awesome-icon
          v-if="status.state === 'play'"
          icon="pause"
          size="lg"
        />
        <font-awesome-icon v-else icon="play" size="lg" />
      </button>
      <button
        aria-label="stop-song"
        @click="sendCommand('/api/playback', 'stop')"
        class="text-lightest py-7 px-8 rounded-full md:pl-3 md:p-0"
      >
        <font-awesome-icon icon="stop" size="lg" />
      </button>
      <button
        aria-label="next-song"
        @click="sendCommand('/api/playback', 'next')"
        class="text-lightest py-7 px-8 rounded-full md:p-3"
      >
        <font-awesome-icon icon="step-forward" size="lg" />
      </button>
    </div>
    <div
      class="
        md:col-start-3 md:col-end-11 md:row-start-1 md:row-end-3
        row-start-1 row-end-4
        w-full
        md:w-auto
        flex
        md:flex-row
        flex-col-reverse
        justify-start
      "
    >
      <div
        @click.self="zoom"
        :class="
          albumArtZoomed
            ? 'fixed top-0 left-0 w-screen h-screen flex justify-center z-50 md:py-10 py-44'
            : ''
        "
      >
        <div
          v-if="defaultAlbumArt"
          class="rounded-xl w-full md:w-20 md:h-20 md:mx-4 avatar"
        ></div>
        <img
          v-show="!defaultAlbumArt"
          id="albumArt"
          @click="zoom"
          :class="
            albumArtZoomed
              ? 'rounded-xl w-auto h-auto'
              : 'rounded-xl w-full md:w-20 md:h-20 md:mx-4 cursor-pointer'
          "
          height="80"
          width="80"
          loading="lazy"
          :src="defaultAlbumArt ? '' : albumArt"
          :alt="
            'album art of ' + currentSong.Title + ' by ' + currentSong.Artist
          "
        />
      </div>
      <div
        class="
          flex flex-col
          self-center
          w-full
          md:w-auto
          my-2
          md:m-0
          text-xl
          md:text-base
          text-left
        "
        v-if="currentSong.Title"
      >
        <span>
          {{ currentSong.Title }}
          <button
            id="like-btn"
            class="rating text-red-500 transition-all duration-200 transform"
            @click="likeSong"
          >
            <font-awesome-icon
              v-if="currentSong.liked == 'true'"
              :icon="['fas', 'heart']"
              size="lg"
            />
            <font-awesome-icon v-else :icon="['far', 'heart']" size="lg" />
          </button>
        </span>
        <span>
          {{ currentSong.Artist }} - {{ currentSong.Album }} ({{
            currentSong.Date
          }})
        </span>
      </div>
    </div>
    <div
      class="
        md:col-start-11 md:col-end-13 md:row-start-1 md:row-end-3
        row-start-5 row-end-5
        w-full
        flex flex-col-reverse
        md:flex-col
        justify-start
        my-2
        md:m-0
        p-2
        md:p-0
      "
    >
      <div class="flex space-x-2 h-10 md:mb-1 md:justify-end justify-between">
        <button
          aria-label="save-queue"
          @click="togglePlaylistSave"
          class="p-2 px-3 text-foreground rounded-full md:p-0 tooltip"
        >
          <font-awesome-icon icon="save" size="lg"></font-awesome-icon>
          <span class="tooltiptext">save current queue</span>
        </button>
        <button
          aria-label="clear-queue"
          @click="clearQueue"
          class="p-2 px-3 text-foreground rounded-full md:p-0 tooltip"
        >
          <font-awesome-icon icon="times" size="lg"></font-awesome-icon>
          <span class="tooltiptext">clear queue</span>
        </button>
        <button
          aria-label="repeat"
          @click="sendCommand('/api/playback', 'repeat')"
          :class="[
            status.repeat == '1'
              ? 'md:text-green-500 bg-green-500 md:bg-transparent'
              : 'md:text-foreground md:dark:text-gray-100 md:bg-transparent bg-gray-700',
            'p-2 md:p-0 text-white rounded-full tooltip',
          ]"
        >
          <span class="tooltiptext">repeat</span>
          <font-awesome-icon icon="retweet" size="lg" />
        </button>
        <button
          aria-label="single"
          @click="sendCommand('/api/playback', 'single')"
          :class="[
            status.single == '1'
              ? 'md:text-green-500 bg-green-500 md:bg-transparent'
              : 'md:text-foreground md:dark:text-gray-100 md:bg-transparent bg-gray-700',
            'px-2.5 md:p-0 text-white rounded-full tooltip',
          ]"
        >
          <span class="tooltiptext">single</span>
          <font-awesome-icon icon="dice-one" size="lg" />
        </button>
        <button
          aria-label="random"
          @click="sendCommand('/api/playback', 'random')"
          :class="[
            status.random == '1'
              ? 'md:text-green-500 bg-green-500 md:bg-transparent'
              : 'md:text-foreground md:dark:text-gray-100 md:bg-transparent bg-gray-700',
            'p-2 md:p-0 text-white rounded-full tooltip',
          ]"
        >
          <span class="tooltiptext">random</span>
          <font-awesome-icon icon="random" size="lg" />
        </button>
        <button
          aria-label="consume"
          @click="sendCommand('/api/playback', 'consume')"
          :class="[
            status.consume == '1'
              ? 'md:text-green-500 bg-green-500 md:bg-transparent'
              : 'md:text-foreground md:dark:text-gray-100 md:bg-transparent bg-gray-700',
            'p-2 md:p-0 text-white rounded-full tooltip',
          ]"
        >
          <span class="tooltiptext">consume</span>
          <font-awesome-icon icon="eraser" size="lg" />
        </button>
        <button
          aria-label="setting"
          @click="openSetting"
          class="
            md:bg-transparent md:text-foreground
            bg-foreground
            text-primary
            rounded-full
            p-2
            md:px-2
            tooltip
          "
        >
          <span class="tooltiptext">settings</span>
          <font-awesome-icon icon="cog" size="lg"></font-awesome-icon>
        </button>
        <button
          aria-label="open-playlist"
          @click="openPlaylist"
          class="bg-foreground text-primary rounded-full p-2 md:hidden"
        >
          <font-awesome-icon icon="list-ul" size="lg"></font-awesome-icon>
        </button>
        <button
          aria-label="open-folders"
          @click="openFolders"
          class="bg-foreground text-primary rounded-full p-2 md:hidden"
        >
          <font-awesome-icon icon="folder" size="lg"></font-awesome-icon>
        </button>
      </div>
      <div
        class="
          flex
          items-center
          space-x-2
          w-full
          mb-2
          md:m-0 md:justify-end
          text-foreground
        "
      >
        <font-awesome-icon
          v-if="status.volume >= 50"
          icon="volume-up"
          size="lg"
        />
        <font-awesome-icon
          v-if="status.volume < 50 && status.volume > 0"
          icon="volume-down"
          size="lg"
        />
        <font-awesome-icon
          v-if="status.volume == 0"
          icon="volume-off"
          size="lg"
        />
        <progressBar
          class="md:w-32 w-full"
          :data="{ value: status.volume, max: 100 }"
          @seek="changeVolume"
        ></progressBar>
      </div>
    </div>
    <div
      class="
        md:col-start-1 md:col-end-13 md:row-start-3 md:row-end-4
        row-start-6 row-end-6
        w-full
        md:w-auto
        flex flex-col
      "
      v-if="status.elapsed"
    >
      <p class="flex justify-between">
        <span>
          {{ humanizeTime(status.elapsed) }}
        </span>
        <span>
          {{ humanizeTime(status.duration) }}
        </span>
      </p>
      <progress-bar
        v-if="status.elapsed"
        :data="{ value: status.elapsed, max: status.duration }"
        v-on:seek="seek"
      ></progress-bar>
    </div>
    <form
      v-if="playListSave"
      @submit.prevent="saveQueue"
      class="
        fixed
        top-0
        left-0
        w-screen
        h-screen
        flex
        justify-center
        items-center
        z-50
        bg-transparent bg-blue-200 bg-opacity-40
      "
      @click.self="togglePlaylistSave"
    >
      <input
        type="text"
        v-model="playListName"
        class="
          border-2 border-primary
          rounded
          mx-2
          px-2
          h-12
          text-black
          w-full
          md:w-1/3
        "
        placeholder="playlist name"
        autofocus
      />
    </form>
  </div>
</template>

<script>
import progressBar from "./progressBar.vue";
import {
  sendCommand,
  humanizeTime,
  toggleMediaPlayer,
  toggleFolders,
} from "../helpers";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { mapState } from "vuex";
import { writeColors } from "../colors";
export default {
  components: {
    progressBar,
    FontAwesomeIcon,
  },
  emits: ["openSetting"],
  data() {
    return {
      albumArtZoomed: false,
      playListName: "",
      playListSave: false,
    };
  },
  methods: {
    seek(time) {
      sendCommand("/api/playback", "seek", { start: Number(time) });
    },
    changeVolume(volume) {
      if (volume > 100) volume = 100;
      if (volume < 0) volume = 0;
      sendCommand("/api/playback", "changeVolume", { start: volume });
    },
    zoom() {
      this.albumArtZoomed = !this.albumArtZoomed;
    },
    likeSong() {
      let el = document.getElementById("like-btn");
      el.classList.add("scale-105", "rotate-45");
      setTimeout(() => {
        el.classList.remove("rotate-45");
        el.classList.add("-rotate-45");
      }, 150);
      setTimeout(() => {
        el.classList.remove("scale-105", "-rotate-45");
      }, 300);
      sendCommand("/api/song", "like", { song: this.currentSong.file });
    },
    openSetting() {
      this.$emit("openSetting");
    },
    clearQueue() {
      sendCommand("/api/playback", "clear");
    },
    togglePlaylistSave() {
      this.playListSave = !this.playListSave;
    },
    saveQueue() {
      sendCommand("/api/queue", "save", { playlist: this.playListName });
      this.togglePlaylistSave();
    },
    openPlaylist: toggleMediaPlayer,
    openFolders: toggleFolders,
    humanizeTime: humanizeTime,
    sendCommand: sendCommand,
  },
  computed: {
    ...mapState(["currentSong", "albumArt", "status"]),
    defaultAlbumArt() {
      writeColors();
      return this.albumArt === "default";
    },
  },
};
</script>
<style>
.avatar {
  background: radial-gradient(50% 123.47% at 50% 50%, #00ff94 0%, #720059 100%),
    linear-gradient(121.28deg, #669600 0%, #ff0000 100%),
    linear-gradient(360deg, #0029ff 0%, #8fff00 100%),
    radial-gradient(100% 164.72% at 100% 100%, #6100ff 0%, #00ff57 100%),
    radial-gradient(100% 148.07% at 0% 0%, #fff500 0%, #51d500 100%);
  background-blend-mode: screen, color-dodge, overlay, difference, normal;
}
</style>
