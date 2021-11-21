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
      z-10
    "
    id="mediaController"
  >
    <div
      class="
        md:col-start-1 md:col-end-3 md:row-start-1 md:row-end-3
        row-start-7 row-end-9
        w-full
        h-28
        md:h-auto
        md:w-auto
        flex
        justify-around
      "
    >
      <button
        aria-label="previous-song"
        @click="playbackCommand('previous')"
        class="text-lightest py-7 px-8 rounded-full md:p-3"
      >
        <font-awesome-icon icon="step-backward" size="lg" />
      </button>
      <button
        aria-label="toggle-playback"
        @click="playbackCommand('toggle')"
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
        @click="playbackCommand('stop')"
        class="text-lightest py-7 px-8 rounded-full md:pl-3 md:p-0"
      >
        <font-awesome-icon icon="stop" size="lg" />
      </button>
      <button
        aria-label="next-song"
        @click="playbackCommand('next')"
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
      <album-art
        :url="albumArt"
        id="albumArt"
        class="w-full md:w-24 md:mx-4"
      ></album-art>
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
            aria-label="like-song"
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
      <div class="flex space-x-2 h-10 md:mb-1 md:justify-end justify-around">
        <save-queue />
        <button
          v-for="command in playbackCommands"
          :key="command.icon"
          :aria-label="command.command"
          @click="playbackCommand(command.command)"
          :class="[
            command.status == '1' ? btnClass.active : btnClass.normal,
            btnClass.base,
          ]"
        >
          <span class="tooltiptext">{{ command.command }}</span>
          <font-awesome-icon :icon="command.icon" size="lg" />
        </button>
        <button
          aria-label="setting"
          @click="openSetting"
          :class="[btnClass.noramal, btnClass.base]"
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
          v-else-if="status.volume < 50 && status.volume > 0"
          icon="volume-down"
          size="lg"
        />
        <font-awesome-icon v-else icon="volume-off" size="lg" />
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
  </div>
</template>

<script>
import progressBar from "./progressBar.vue";
import albumArt from "./albumArt.vue";
import saveQueue from "./saveQueue.vue";
import {
  sendCommand,
  humanizeTime,
  toggleMediaController,
  toggleFolders,
} from "../helpers";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { mapState } from "vuex";
import { useStore } from "vuex";
import { reactive, computed } from "vue";
export default {
  components: {
    progressBar,
    FontAwesomeIcon,
    albumArt,
    saveQueue,
  },
  emits: ["openSetting"],
  setup() {
    const store = useStore();
    let playbackCommands = reactive([
      {
        command: "consume",
        icon: "minus-square",
        status: computed(() => store.state.status.consume),
      },
      {
        command: "single",
        icon: "dice-one",
        status: computed(() => store.state.status.single),
      },
      {
        command: "repeat",
        icon: "retweet",
        status: computed(() => store.state.status.repeat),
      },
      {
        command: "random",
        icon: "random",
        status: computed(() => store.state.status.random),
      },
      {
        command: "clear",
        icon: "eraser",
        status: 0,
      },
    ]);
    const btnClass = {
      active: "text-green-500 bg-transparent",
      normal: "text-foreground bg-transparent",
      base: "p-2 md:p-0 rounded-full tooltip md:hover:text-blue-200",
    };
    return { playbackCommands, btnClass };
  },
  methods: {
    seek(time) {
      sendCommand("/api/playback", "seek", { start: Number(time) });
    },
    playbackCommand(command) {
      sendCommand("/api/playback", command);
    },
    changeVolume(volume) {
      if (volume > 100) volume = 100;
      if (volume < 0) volume = 0;
      sendCommand("/api/playback", "changeVolume", { start: volume });
    },
    likeSong() {
      sendCommand("/api/song", "like", { song: this.currentSong.file });
      let el = document.getElementById("like-btn");
      el.classList.add("scale-105", "rotate-45");
      setTimeout(() => {
        el.classList.remove("rotate-45");
        el.classList.add("-rotate-45");
      }, 150);
      setTimeout(() => {
        el.classList.remove("scale-105", "-rotate-45");
      }, 300);
    },
    openSetting() {
      this.$emit("openSetting");
    },
    clearQueue() {
      sendCommand("/api/playback", "clear");
    },
    openPlaylist: toggleMediaController,
    openFolders: toggleFolders,
    humanizeTime: humanizeTime,
  },
  created() {},
  computed: {
    ...mapState(["currentSong", "albumArt", "status"]),
  },
};
</script>
