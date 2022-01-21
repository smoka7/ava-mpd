<template>
  <details :open="currentAlbum" @click="loadSongs()">
    <summary
      class="flex px-8 py-2 md:rounded sticky top-0 bg-lightest items-center cursor-pointer dark:text-primary"
    >
      <FontAwesomeIcon
        :icon="['fas', currentAlbum ? 'compact-disc' : 'angle-right']"
        class="mr-2 transform-gpu duration-200"
        :id="'icon-' + album.Songs[0].Pos"
      />
      {{ album.Artist }} - {{ album.Album }} ({{ album.Date }})
    </summary>
    <div v-if="currentAlbum || open">
      <div
        v-for="song in album.Songs"
        :key="song.Pos"
        draggable="true"
        @dragover.prevent
        @dragenter.prevent
        @drop="moveSong($event, song.Pos)"
        @dragstart="startMoveSong($event, song.Pos)"
        class="md:py-1 py-2 px-4 grid grid-cols-12 items-center md:rounded md:m-1 text-black dark:text-white odd:bg-gray-600/10 dark:odd:bg-gray-800/50 dark:hover:odd:bg-gray-800/70 hover:bg-white/60 dark:hover:bg-gray-800/70 group"
        :id="'song' + song.Pos"
      >
        <span
          class="col-start-1 col-end-3 md:col-end-2 px-2 cursor-pointer"
          @click.stop="play(song.Pos)"
        >
          <FontAwesomeIcon
            :id="song.Pos"
            :class="
              song.Pos !== currentSongPos
                ? 'invisible group-hover:visible text-green-500 mr-2'
                : 'text-red-500 mr-2'
            "
            :icon="[
              'fas',
              song.Pos !== currentSongPos ? 'play' : 'compact-disc',
            ]"
          />
          {{ song.Track }}
        </span>
        <span
          class="col-start-3 md:col-start-2 col-end-11 overflow-x-hidden text-ellipsis"
        >
          {{ song.Title }}
        </span>
        <span
          class="col-start-11 md:col-start-12 col-end-13 cursor-pointer"
          @click="$emit('showMenu', song.Pos, $event)"
        >
          <FontAwesomeIcon
            :id="'delete' + song.Pos"
            class="mr-2 invisible group-hover:visible"
            icon="ellipsis-h"
          />
          {{ humanizeTime(song.Duration) }}
        </span>
      </div>
    </div>
  </details>
</template>
<script>
import { sendCommand, humanizeTime } from "../helpers";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
export default {
  props: ["album", "currentSongPos", "currentAlbum"],
  components: { FontAwesomeIcon },
  emits: ["showMenu"],
  data() {
    return {
      open: this.currentAlbum,
    };
  },
  methods: {
    humanizeTime: humanizeTime,
    play(id) {
      let data = {
        Start: Number(id),
      };
      sendCommand("/api/queue", "play", data);
    },
    loadSongs() {
      this.open = true;
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
  },
};
</script>
<style>
details[open] > summary > svg {
  @apply rotate-90;
}
details[open] summary ~ * {
  animation: sweep 0.2s ease-in-out;
}
@keyframes sweep {
  0% {
    @apply opacity-0 mt-2;
  }
  100% {
    @apply opacity-100 mt-0;
  }
}
</style>
