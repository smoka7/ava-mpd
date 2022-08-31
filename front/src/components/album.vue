<template>
  <details :open="currentAlbum" class="album">
    <summary
      :class="[currentAlbum ? 'current-album-header' : 'album-header', 'group']"
    >
      <span>
        <FontAwesomeIcon
          :icon="['fas', currentAlbum ? 'compact-disc' : 'angle-right']"
          :class="[
            currentAlbum ? 'text-accent' : '',
            'mr-2 transform-gpu duration-200 open-rotate',
          ]"
          size="lg"
        />
        {{ album.Artist }} | {{ album.Album }}
      </span>
      <button
        class="w-fit shrink-0 cursor-pointer"
        aria-label="select-the-album"
      >
        <FontAwesomeIcon
          label="select Album"
          icon="check-circle"
          size="lg"
          class="invisible text-primary group-hover:visible"
          @click.stop="$emit('selectAlbum', album.Album)"
        />
        {{ album.Date }}
      </button>
    </summary>
    <div
      v-for="song in album.Songs"
      :key="song.Pos"
      draggable="true"
      @dragover.prevent
      @dragenter.prevent
      @drop="moveSong($event, song.Pos)"
      @dragstart="startMoveSong($event, song.Id)"
      class="song group"
      tabindex="0"
      :id="isCurrSong(song.Pos) ? 'currentSong' : ''"
    >
      <div
        class="mr-2 flex flex-grow-0 space-x-4 overflow-x-clip text-ellipsis"
      >
        <span class="group relative cursor-pointer" @click.stop="play(song.Id)">
          <FontAwesomeIcon
            :id="song.Pos"
            :class="
              isCurrSong(song.Pos)
                ? 'text-accent'
                : 'invisible text-green-500 group-hover:visible'
            "
            :icon="['fas', isCurrSong(song.Pos) ? 'compact-disc' : 'play']"
          />
          <span
            :class="[
              isCurrSong(song.Pos)
                ? 'hidden'
                : 'absolute inset-0 group-hover:invisible',
            ]"
          >
            {{ song.Track }}
          </span>
        </span>
        <span class="">
          {{ song.Title }}
        </span>
      </div>
      <div
        class="flex flex-shrink-0 cursor-pointer items-center justify-between space-x-2 md:w-2/12"
      >
        <button
          @click="$emit('select', song.Id)"
          aria-label="select-the-song"
          :class="
            isSelected(song.Id)
              ? 'visible text-green-500 dark:text-green-400'
              : 'invisible text-primary group-hover:visible dark:text-white'
          "
        >
          <FontAwesomeIcon icon="check-circle" />
        </button>
        <button
          aria-label="show-menu"
          @click="$emit('showMenu', song.Pos, song.Id, $event)"
          class="invisible mr-2 text-primary group-hover:visible dark:text-white"
        >
          <FontAwesomeIcon icon="ellipsis-h" />
        </button>
        <span>
          {{ humanizeTime(song.Duration) }}
        </span>
      </div>
    </div>
  </details>
</template>
<script setup>
import endpoints from "../endpoints.js";
import { sendCommand, humanizeTime } from "../helpers";

defineEmits(["showMenu", "select", "selectAlbum"]);

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

function isCurrSong(pos) {
  return pos === props.currentSongPos;
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
<style lang="postcss">
details[open] > summary > span > .open-rotate {
  @apply rotate-90;
}
details[open] summary ~ * {
  animation: sweep 0.2s ease-in-out;
}

.album[open] {
  @apply rounded bg-white/40 p-1 dark:bg-gray-600/50;
}

.album-header {
  @apply sticky top-0 flex cursor-pointer items-center justify-between rounded px-4 py-2;
  @apply transition-transform duration-200;
  @apply hover:bg-lightest hover:py-3 hover:text-primary dark:text-white dark:hover:text-primary;
}

.current-album-header {
  @apply album-header;
  @apply bg-lightest text-primary hover:py-2 !important;
}

.song {
  @apply my-1 flex w-full justify-between p-1 duration-300;
  @apply hover:bg-white/60 hover:py-1.5 dark:text-white dark:hover:bg-gray-800/70;
  @apply md:rounded md:px-3;
}

#currentSong {
  @apply bg-white/60 dark:bg-gray-800/70;
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
