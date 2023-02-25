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
            'open-rotate mr-2 transform-gpu duration-200',
          ]"
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
          class="hide-unfocused text-primary dark:text-white"
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
        <button
          class="group relative cursor-pointer"
          @click.stop="play(song.Id)"
        >
          <FontAwesomeIcon
            :id="song.Pos"
            :class="[
              isCurrSong(song.Pos)
                ? ' group-hover:invisible group-focus:invisible'
                : 'hide-unfocused',
              'text-green-500',
            ]"
            :icon="['fas', 'play']"
          />
          <span
            :class="[
              isCurrSong(song.Pos)
                ? 'hide-unfocused'
                : ' group-focus-within:invisible group-hover:invisible group-focus:invisible',
              'absolute inset-0',
            ]"
          >
            {{ song.Track }}
          </span>
        </button>
        <span>
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
              : 'hide-unfocused text-primary dark:text-white'
          "
        >
          <FontAwesomeIcon icon="check-circle" />
        </button>
        <button
          aria-label="show-menu"
          @click="$emit('showMenu', song.Pos, song.Id, $event)"
          class="hide-unfocused mr-2 text-primary dark:text-white"
        >
          <FontAwesomeIcon icon="ellipsis-h" />
        </button>
        <span>
          {{ humanizeTime(Number(song.Duration)) }}
        </span>
      </div>
    </div>
  </details>
</template>
<script setup lang="ts">
import type { Album } from "../store";
import endpoints from "../endpoints";
import { sendCommand, humanizeTime } from "../helpers";

defineEmits(["showMenu", "select", "selectAlbum"]);
const props = defineProps<{
  album: Album;
  currentSongPos: number;
  currentAlbum: boolean;
  selectedIds: Array<number>;
}>();

function play(id: number) {
  const data = {
    ID: Number(id),
  };
  sendCommand(endpoints.queue, "play", data);
}

function isCurrSong(pos: number): boolean {
  return pos === props.currentSongPos;
}

function isSelected(id: number): boolean {
  return props.selectedIds.indexOf(id) > -1;
}

function startMoveSong(event: DragEvent, id: number) {
  if (event.dataTransfer == null) {
    return;
  }
  event.dataTransfer.setData("id", id.toString());
  event.dataTransfer.dropEffect = "move";
  event.dataTransfer.effectAllowed = "move";
}

function moveSong(event: DragEvent, position: number) {
  if (event.dataTransfer == null) {
    return;
  }
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
  @apply hover:bg-white dark:text-white hover:dark:bg-gray-800/70;
}

.current-album-header {
  @apply album-header;
  @apply bg-white dark:bg-gray-800/70 dark:text-white !important;
}

.song {
  @apply my-1 flex w-full justify-between p-1;
  @apply hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70;
  @apply md:rounded md:px-3;
}

#currentSong {
  @apply bg-white py-1.5 dark:bg-gray-800/70;
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
