<template>
  <div
    class="flex h-full max-h-full flex-col justify-around space-y-4 bg-primary p-2 text-white md:rounded"
  >
    <div
      class="my-2 flex w-full flex-col self-center overflow-hidden text-ellipsis text-left text-xl md:max-h-[30%]"
      v-if="currentSong.Title"
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
      class="aspect-square w-2/3 flex-shrink-0 self-center md:w-1/2"
    />
    <volume-control :v-if="status != null" :volume="status.volume" />
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

<script setup>
import endpoints from "../endpoints.js";
import { sendCommand, humanizeTime } from "../helpers";

defineEmits(["openSetting"]);

const store = useStore();
const AlbumArtUrl = computed(() => store.state.albumArt);
const status = computed(() => shallowReactive(store.state.status));
const currentSong = computed(() => shallowReactive(store.state.currentSong));
const hoverProgress = ref(0);

function seek(time) {
  sendCommand(endpoints.playback, "seek", { start: Number(time) });
}

function sethoverProgress(e) {
  hoverProgress.value = e;
}
</script>
