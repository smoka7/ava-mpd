<template>
  <details :open="currentAlbum" @click="loadSongs()">
    <summary
      class="bg-lightest dark:text-primary sticky top-0 flex cursor-pointer items-center px-8 py-2 md:rounded"
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
        class="group grid grid-cols-12 items-center py-2 px-4 text-black odd:bg-gray-600/10 hover:bg-white/60 dark:text-white dark:odd:bg-gray-800/50 dark:hover:bg-gray-800/70 dark:hover:odd:bg-gray-800/70 md:m-1 md:rounded md:py-1"
        :id="'song' + song.Pos"
      >
        <span
          class="col-start-1 col-end-3 cursor-pointer px-2 md:col-end-2"
          @click.stop="play(song.Pos)"
        >
          <FontAwesomeIcon
            :id="song.Pos"
            :class="
              song.Pos !== currentSongPos
                ? 'invisible mr-2 text-green-500 group-hover:visible'
                : 'mr-2 text-red-500'
            "
            :icon="[
              'fas',
              song.Pos !== currentSongPos ? 'play' : 'compact-disc',
            ]"
          />
          {{ song.Track }}
        </span>
        <span
          class="col-start-3 col-end-11 overflow-x-hidden text-ellipsis md:col-start-2"
        >
          {{ song.Title }}
        </span>
        <span
          class="col-start-11 col-end-13 cursor-pointer md:col-start-12"
          @click="$emit('showMenu', song.Pos, $event)"
        >
          <FontAwesomeIcon
            :id="'delete' + song.Pos"
            class="invisible mr-2 group-hover:visible"
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
<style scoped>
details[open] > summary > svg {
  @apply rotate-90;
}
details[open] summary ~ * {
  animation: sweep 0.2s ease-in-out;
}
@keyframes sweep {
  0% {
    @apply mt-2 opacity-0;
  }
  100% {
    @apply mt-0 opacity-100;
  }
}
</style>
