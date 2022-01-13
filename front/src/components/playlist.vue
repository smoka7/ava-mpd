<template>
  <div
    class="flex flex-col w-full h-full rounded overflow-hidden backdrop-blur-3xl bg-white/60 dark:bg-gray-700/60"
  >
    <div class="overflow-y-auto snap-y">
      <div
        v-if="!queue"
        class="flex items-center justify-center p-4 w-full h-full text-9xl decoration-accent underline"
      >
        Queue is empty!
      </div>
      <details v-for="(album, index) in queue" :key="index" open>
        <summary
          class="flex px-8 md:py-2 py-4 md:mx-2 my-1 md:rounded sticky top-0 bg-lightest items-center cursor-pointer snap-both dark:text-primary"
          @click="animate(album[0].Pos)"
        >
          <FontAwesomeIcon
            icon="angle-down"
            class="mr-2 transform-gpu duration-200"
            :id="'icon-' + album[0].Pos"
          ></FontAwesomeIcon>
          {{ album[0].Artist }} - {{ album[0].Album }} ({{ album[0].Date }})
        </summary>
        <div
          draggable="true"
          @dragover.prevent
          @dragenter.prevent
          @drop="moveSong($event, song.Pos)"
          @dragstart="startMoveSong($event, song.Pos)"
          v-for="song in album"
          :key="song.Pos"
          class="md:py-1 py-2 px-4 grid grid-cols-12 md:rounded md:mx-2 text-black dark:text-white cursor-pointer group odd:bg-gray-600/10 dark:odd:bg-gray-800/50 dark:hover:odd:bg-gray-800/70 hover:bg-white/60 dark:hover:bg-gray-800/70"
          :id="'song' + song.Pos"
        >
          <span
            class="col-start-1 col-end-3 md:col-end-2 justify-self-start pr-2 w-full flex items-center cursor-pointer"
            @click.stop="play(song.Pos)"
          >
            <FontAwesomeIcon
              :id="song.Pos"
              :class="
                song.Pos !== currentSong.Pos
                  ? 'invisible group-hover:visible text-green-500 mr-2'
                  : 'text-red-500 mr-2'
              "
              :icon="[
                'fas',
                song.Pos !== currentSong.Pos ? 'play' : 'compact-disc',
              ]"
            />
            {{ song.Track }}
          </span>
          <span class="col-start-3 md:col-start-2 col-end-11 self-start">
            {{ song.Title }}
          </span>
          <span
            class="col-start-11 md:col-start-12 col-end-13 justify-self-center"
          >
            <FontAwesomeIcon
              @click="showMenu(song.Pos, $event)"
              :id="'delete' + song.Pos"
              class="mr-2 invisible group-hover:visible"
              icon="ellipsis-h"
            />
            {{ humanizeTime(song.duration) }}
          </span>
        </div>
      </details>
    </div>
    <button
      aria-label="close-playlist"
      @click="closePlaylist"
      class="fixed right-4 bottom-4 bg-red-500 text-white rounded-full p-2 md:hidden w-20 h-20"
    >
      <font-awesome-icon icon="arrow-right" size="2x"></font-awesome-icon>
    </button>
    <teleport to="#app">
      <queue-menu
        :open="menu"
        @close="hideMenu"
        :song="selected"
        @showInfo="songInfo = true"
      />
    </teleport>
    <songInfo v-if="songInfo" :song="selected" @close="songInfo = false" />
  </div>
</template>
<script>
import { mapState } from "vuex";
import {
  sendCommand,
  humanizeTime,
  toggleMediaController,
} from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import songInfo from "./songInfo.vue";
import queueMenu from "./queueMenu.vue";
export default {
  components: {
    FontAwesomeIcon,
    songInfo,
    queueMenu,
  },
  data() {
    return {
      menu: false,
      songInfo: false,
      selected: -1,
    };
  },
  async mounted() {
    await this.$store.dispatch("getQueue");
    let el = document.getElementById("song" + this.currentSong.Pos);
    if (el !== null) {
      el.classList.add("current-song");
      el.scrollIntoView({ block: "center", behavior: "smooth" });
    }
  },
  methods: {
    humanizeTime: humanizeTime,
    closePlaylist: toggleMediaController,
    showMenu(pos, e) {
      this.selected = Number(pos);
      let el = document.getElementById("context-menu");
      el.style.top = e.pageY + "px";
      this.menu = true;
    },
    startMoveSong(event, position) {
      event.dataTransfer.setData("position", position);
      event.dataTransfer.dropEfect = "move";
      event.dataTransfer.effectAllowed = "move";
    },
    moveSong(event, position) {
      const start = event.dataTransfer.getData("position");
      let data = {
        Start: Number(start),
        Finish: Number(position),
      };
      sendCommand("/api/queue", "move", data);
    },
    hideMenu() {
      this.menu = false;
    },
    animate(id) {
      let el = document.getElementById("icon-" + id);
      if (el) el.classList.toggle("-rotate-90");
    },
    play(id) {
      let data = {
        Start: Number(id),
      };
      sendCommand("/api/queue", "play", data);
    },
  },
  computed: {
    ...mapState({
      queue: (state) => state.queue,
      currentSong: (state) => state.currentSong,
    }),
  },
  watch: {
    currentSong: function (newSong, oldSong) {
      let oldId = "song" + oldSong.Pos;
      let newId = "song" + newSong.Pos;
      let el = document.getElementById(oldId);
      if (el) {
        el.classList.remove("current-song");
      }
      el = document.getElementById(newId);
      if (el) {
        el.classList.add("current-song");
        el.scrollIntoView({ block: "center", behavior: "smooth" });
      }
    },
  },
};
</script>
<style>
.current-song {
  @apply bg-red-200 dark:text-primary !important;
}
</style>
