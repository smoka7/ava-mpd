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
        class="col-span-12 col-start-1 row-start-1 row-end-1 mr-2 flex space-x-4 overflow-x-clip text-ellipsis md:col-span-10"
      >
        <button
          class="group relative cursor-pointer"
          @click.stop="play(song.Id)"
        >
          <FontAwesomeIcon
            :id="song.Pos"
            :class="[
              isCurrSong(song.Pos)
                ? ' text-green-500 group-hover:invisible group-focus:invisible'
                : 'hide-unfocused text-primary',
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
        class="md;pr-0 col-span-5 col-start-8 row-start-2 row-end-2 flex cursor-pointer items-center justify-between space-x-2 pr-2 md:col-span-2 md:col-start-11 md:row-start-1 md:row-end-1"
      >
        <button
          aria-label="show-menu"
          @click="$emit('showMenu', song.Pos, song.Id, $event)"
          class="hide-unfocused mr-2 text-primary dark:text-white"
        >
          <FontAwesomeIcon icon="ellipsis-h" />
        </button>
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

        <likeSong
          :p-liked="song.Liked"
          :file="song.File"
          icon-size="lg"
          color="text-primary/80 dark:text-white/80"
        />

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
  return Number(pos) === props.currentSongPos;
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
  @apply my-1 grid w-full grid-cols-12 grid-rows-2 justify-between p-1 md:grid-rows-1;
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
