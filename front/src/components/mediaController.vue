<template>
  <div
    class="flex h-full max-h-full flex-col justify-between space-y-2 bg-primary p-2 text-white md:rounded-l"
  >
    <album-art
      :url="AlbumArtUrl"
      :altText="currentSong.Title + 'cover'"
      id="albumArt"
      class="-m-2 flex h-[60dvh] w-[104%] flex-col justify-end space-y-1 self-center rounded-tl"
    >
      <div
        v-if="currentSong.Title"
        class="absolute my-2 flex h-[50%] w-full flex-col justify-end self-start overflow-hidden text-ellipsis bg-gradient-to-b from-transparent via-primary/40 to-primary p-4 text-left text-3xl backdrop-opacity-70 md:text-xl"
      >
        <span>{{ currentSong.Title }}</span>
        <span> {{ currentSong.Album }} </span>
        <span> {{ currentSong.Artist }} ({{ currentSong.Date }}) </span>
      </div>
    </album-art>
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
