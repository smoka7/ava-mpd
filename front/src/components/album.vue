<template>
  <details :open="currentAlbum">
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
      <div
        v-for="song in album.Songs"
        :key="song.Pos"
        draggable="true"
        @dragover.prevent
        @dragenter.prevent
        @drop="moveSong($event, song.Pos)"
        @dragstart="startMoveSong($event, song.Id)"
        class="group grid grid-cols-12 items-center py-2 px-4 text-black even:bg-gray-600/10 hover:bg-white/60 dark:text-white dark:even:bg-gray-800/50 dark:hover:bg-gray-800/70 dark:hover:even:bg-gray-800/70 md:m-1 md:rounded md:py-1"
        :id="song.Pos !== currentSongPos ? '' : 'currentSong'"
      >
        <span
          class="col-start-1 col-end-3 cursor-pointer px-2 md:col-end-2"
          @click.stop="play(song.Id)"
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
          class="col-start-3 col-end-9 overflow-x-hidden text-ellipsis md:col-start-2 justify-start"
        >
          {{ song.Title }}
        </span>
        <button
          class="col-start-9 col-end-13 flex cursor-pointer items-center justify-end space-x-2 md:col-start-11"
        >
          <FontAwesomeIcon
            label="select song"
            icon="check-circle"
            :class="
              isSelected(song.Id)
                ? 'visible text-green-500 dark:text-green-400'
                : 'text-primary invisible group-hover:visible dark:text-white'
            "
            @click="$emit('select', song.Id)"
          />
          <FontAwesomeIcon
            @click="$emit('showMenu', song.Pos, song.Id, $event)"
            class="text-primary invisible mr-2 group-hover:visible dark:text-white"
            icon="ellipsis-h"
          />
          <span>
            {{ humanizeTime(song.Duration) }}
          </span>
        </button>
      </div>
  </details>
</template>
<script setup>
import { ref } from "vue";
import endpoints from "../endpoints.js";
import { sendCommand, humanizeTime } from "../helpers";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

defineEmits(["showMenu", "select"]);

const props = defineProps([
  "album",
  "currentSongPos",
  "currentAlbum",
  "selectedIds",
]);

function play(id) {
  const data = {
    ID: Number(id),
  };
  sendCommand(endpoints.queue, "play", data);
}

function isSelected(id) {
  return props.selectedIds.indexOf(id) > -1;
}

function startMoveSong(event, id) {
  event.dataTransfer.setData("id", id);
  event.dataTransfer.dropEfect = "move";
  event.dataTransfer.effectAllowed = "move";
}

function moveSong(event, position) {
  const start = event.dataTransfer.getData("id");
  const data = {
    Start: Number(start),
    Finish: Number(position),
  };
  sendCommand(endpoints.queue, "move", data);
}
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
