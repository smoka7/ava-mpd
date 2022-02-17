<template>
  <div
    class="grid-rows-8 grid grid-cols-1 bg-primary p-4 text-white md:static md:grid-cols-12 md:grid-rows-3 md:rounded"
    id="mediaController"
  >
    <div
      class="row-end-9 row-start-7 flex h-28 w-full justify-center space-x-10 text-lightest md:col-start-1 md:col-end-4 md:row-start-1 md:row-end-3 md:h-auto md:w-auto md:space-x-6"
    >
      <button
        aria-label="previous-song"
        @click="playbackCommand('previous')"
        class="scale-110 hover:text-accent md:scale-75"
      >
        <font-awesome-icon icon="step-backward" size="2x" />
      </button>
      <like-song :pLiked="liked" :file="currentSong.file" />
      <button
        aria-label="toggle-playback"
        @click="playbackCommand('toggle')"
        :class="[
          status.state === 'pause' ? 'bg-green-400' : 'bg-accent',
          'flex aspect-square items-center justify-center rounded-full text-white hover:text-primary',
        ]"
      >
        <font-awesome-icon
          :icon="['fas', status.state === 'play' ? 'pause' : 'play']"
          size="lg"
        />
      </button>
      <button
        aria-label="stop-song"
        @click="playbackCommand('stop')"
        class="scale-110 hover:text-accent md:scale-75"
      >
        <font-awesome-icon icon="stop" size="2x" />
      </button>
      <button
        aria-label="next-song"
        @click="playbackCommand('next')"
        class="scale-110 hover:text-accent md:scale-75"
      >
        <font-awesome-icon icon="step-forward" size="2x" />
      </button>
    </div>
    <div
      class="row-start-1 row-end-3 flex flex-col-reverse justify-around md:col-start-4 md:col-end-11 md:row-start-1 md:row-end-3 md:flex-row md:justify-start"
    >
      <album-art
        :url="albumArt"
        :altText="currentSong.Title + 'cover'"
        id="albumArt"
        class="aspect-square w-full flex-shrink-0 md:mx-2 md:h-24 md:w-24"
      />
      <div
        class="my-2 flex w-full flex-col self-center text-left text-2xl md:w-auto md:text-base"
        v-if="currentSong.Title"
      >
        <span class="text-ellipsis">
          {{ currentSong.Title }}
        </span>
        <span>
          {{ currentSong.Artist }} - {{ currentSong.Album }} ({{
            currentSong.Date
          }})
        </span>
      </div>
    </div>
    <div
      class="row-start-3 row-end-5 flex w-full flex-col-reverse justify-start md:col-start-11 md:col-end-13 md:row-start-1 md:row-end-3 md:flex-col"
    >
      <div class="flex h-10 justify-around space-x-2 md:mb-1 md:justify-end">
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
          @click="$emit('openSetting')"
          :class="[btnClass.noramal, btnClass.base]"
        >
          <span class="tooltiptext">settings</span>
          <font-awesome-icon icon="cog" size="lg" />
        </button>
        <button
          aria-label="open-playlist"
          @click="openPlaylist"
          :class="[
            btnClass.noramal,
            btnClass.base,
            'rounded-full bg-white p-2 text-primary md:hidden',
          ]"
        >
          <font-awesome-icon icon="list-ul" size="lg" />
        </button>
        <button
          aria-label="open-folders"
          @click="openFolders"
          :class="[
            btnClass.noramal,
            btnClass.base,
            'rounded-full bg-white p-2 text-primary md:hidden',
          ]"
        >
          <font-awesome-icon icon="folder" size="lg" />
        </button>
      </div>
      <div
        class="mb-2 flex w-full items-center space-x-2 md:m-0 md:justify-end"
      >
        <volume-control :volume="status.volume" />
      </div>
    </div>
    <div
      class="row-start-6 row-end-6 flex w-full flex-col md:col-start-1 md:col-end-13 md:row-start-3 md:row-end-4 md:w-auto"
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
      />
    </div>
  </div>
</template>

<script>
import progressBar from "./progressBar.vue";
import albumArt from "./albumArt.vue";
import saveQueue from "./saveQueue.vue";
import likeSong from "./likeSong.vue";
import volumeControl from "./volumeControl.vue";
import { shallowReactive } from "vue";
import endpoints from "../endpoints.js";
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
    likeSong,
    volumeControl,
  },
  emits: ["openSetting"],
  setup() {
    const store = useStore();
    const playbackCommands = reactive([
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
      normal: "bg-transparent",
      base: "p-2 md:p-0 tooltip md:hover:text-accent",
    };
    return { playbackCommands, btnClass };
  },
  methods: {
    seek(time) {
      sendCommand(endpoints.playback, "seek", { start: Number(time) });
    },
    playbackCommand(command) {
      sendCommand(endpoints.playback, command);
    },
    openPlaylist: toggleMediaController,
    openFolders: toggleFolders,
    humanizeTime: humanizeTime,
  },
  computed: {
    ...mapState({
      albumArt: (state) => state.albumArt,
      liked: (state) => state.currentSong.Liked,
      currentSong: (state) => shallowReactive(state.currentSong.Info),
      status: (state) => shallowReactive(state.status),
    }),
  },
};
</script>
