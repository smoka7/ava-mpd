<template>
  <div
    class="flex h-full max-h-full flex-col justify-between space-y-2 bg-primary p-2 text-white md:rounded"
  >
    <div
      class="flex w-full flex-col space-y-1 self-center sm:flex-row sm:justify-between md:max-h-[60dvh] md:flex-col md:justify-end"
      v-if="currentSong.Title"
    >
      <div
        class="my-2 flex w-full flex-col self-start overflow-hidden text-ellipsis text-left text-xl sm:w-2/3"
      >
        <span class="text-ellipsis">
          {{ currentSong.Title }}
        </span>
        <span> {{ currentSong.Album }} </span>
        <span> {{ currentSong.Artist }} ({{ currentSong.Date }}) </span>
      </div>
      <album-art
        :url="AlbumArtUrl"
        :altText="currentSong.Title + 'cover'"
        id="albumArt"
        class="aspect-square w-2/3 flex-shrink-0 self-center sm:w-1/5 md:w-3/5"
      />
    </div>
    <volume-control :v-if="status != null" :volume="Number(status.volume)" />
    <PlaybackCommands />
    <div class="flex w-full flex-col justify-start" v-if="status.elapsed">
      <p class="flex justify-between">
        <span>
          {{ humanizeTime(status.elapsed) }}
        </span>
        <span>
          {{ humanizeTime(status.duration) }}
        </span>
      </p>
      <div class="tooltip">
        <progress-bar
          v-if="status.elapsed"
          :data="{ value: status.elapsed, max: status.duration }"
          v-on:seek="seek"
          @sendhover="sethoverProgress"
        />
        <span class="tooltiptext">{{ humanizeTime(hoverProgress) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import endpoints from "../endpoints";
import { sendCommand, humanizeTime } from "../helpers";
import { useStore } from "../store";
defineEmits(["openSetting"]);

const store = useStore();
const AlbumArtUrl = computed(() => store.albumArt);
const status = computed(() => shallowReactive(store.status));
const currentSong = computed(() => shallowReactive(store.currentSong));
const hoverProgress = ref(0);

function seek(time: number) {
  sendCommand(endpoints.playback, "seek", { start: time });
}

function sethoverProgress(e: number) {
  hoverProgress.value = e;
}
</script>
