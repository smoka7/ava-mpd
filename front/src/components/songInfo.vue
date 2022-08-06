<template>
  <div
    class="absolute inset-0 z-auto flex flex-col space-y-2 overflow-x-hidden overflow-y-scroll bg-lightest p-8"
  >
    <div class="flex flex-row items-center justify-between">
      <h1
        class="mr-2 flex items-center space-x-4 text-ellipsis text-4xl font-bold text-primary underline decoration-accent"
      >
        <span>
          {{ state.info.Title }}
        </span>
        <likeSong :pLiked="state.liked" :file="state.info['file']" />
      </h1>
      <button
        aria-label="close-info"
        @click="$emit('close')"
        class="rounded-full px-4 py-2"
      >
        <font-awesome-icon class="text-red-500" icon="times" size="2x" />
      </button>
    </div>
    <div
      class="flex flex-col justify-around space-y-2 md:flex-row md:space-x-2 md:space-y-0"
    >
      <album-art
        v-if="state.albumArt != 'default'"
        :url="state.albumArt"
        :altText="state.info['Title'] + ' cover'"
        class="top-0 h-fit md:sticky md:w-1/2"
      />
      <ul
        :class="{
          'flex w-full flex-col rounded bg-white dark:bg-gray-700': true,
          'md:w-1/2': state.albumArt != 'default',
        }"
      >
        <li
          v-for="(value, index) in state.info"
          :key="index"
          class="p-2 first:rounded-t last:rounded-b even:bg-blue-50 dark:even:bg-gray-800"
        >
          {{ index }} : {{ value }}
        </li>
        <li
          v-for="(sticker, index) in state.stickers"
          :key="index"
          class="p-2 last:rounded-b even:bg-blue-50 dark:even:bg-gray-800"
        >
          {{ sticker.Name }} : {{ sticker.Value }}
        </li>
      </ul>
    </div>
  </div>
</template>
<script setup>
import endpoints from "../endpoints.js";
import { sendCommand } from "../helpers";
import { reactive, onMounted } from "vue";

const props = defineProps(["song"]);
const state = reactive({
  info: {},
  stickers: {},
  albumArt: "",
  liked: false,
});

function isItLiked() {
  const index = state.stickers.findIndex((stick) => {
    if (stick.Name == "liked" && stick.Value == "true") return true;
  });
  if (index > -1) {
    state.liked = true;
  }
}

async function getInfo() {
  const song = await sendCommand(endpoints.song, "info", {
    ID: props.song,
  });
  const albumArt = await sendCommand(endpoints.song, "albumArt", {
    ID: props.song,
  });
  state.info = song.Info;
  state.stickers = song.Stickers;
  state.albumArt = albumArt.Url;
  isItLiked();
}

onMounted(() => {
  getInfo();
});
</script>
