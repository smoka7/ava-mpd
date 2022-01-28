<template>
  <div
    class="bg-lightest absolute inset-0 z-auto flex flex-col space-y-2 overflow-x-hidden overflow-y-scroll p-8"
  >
    <div class="flex flex-row items-center justify-between">
      <h1
        class="text-bold mr-2 flex items-center space-x-4 text-ellipsis text-4xl font-bold underline decoration-2"
      >
        <span>
          {{ info.Title }}
        </span>
        <likeSong :pLiked="liked" :file="info['file']" />
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
        v-if="albumArt != 'default'"
        :url="albumArt"
        :altText="info['Title'] + ' cover'"
        class="top-0 h-fit md:sticky md:w-1/2"
      />
      <ul
        :class="{
          'flex w-full flex-col rounded bg-white dark:bg-gray-700': true,
          'md:w-1/2': albumArt != 'default',
        }"
      >
        <li
          v-for="(value, index) in info"
          :key="index"
          class="p-2 first:rounded-t last:rounded-b even:bg-blue-50 dark:even:bg-gray-800"
        >
          {{ index }} : {{ value }}
        </li>
        <li
          v-for="(sticker, index) in stickers"
          :key="index"
          class="p-2 last:rounded-b even:bg-blue-50 dark:even:bg-gray-800"
        >
          {{ sticker.Name }} : {{ sticker.Value }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import likeSong from "./likeSong.vue";
import albumArt from "./albumArt.vue";
export default {
  components: {
    FontAwesomeIcon,
    albumArt,
    likeSong,
  },
  props: ["song"],
  data() {
    return {
      info: {},
      stickers: {},
      albumArt: "",
      liked: false,
    };
  },
  methods: {
    isItLiked() {
      let index = this.stickers.findIndex((stick) => {
        if (stick.Name == "liked" && stick.Value == "true") return true;
      });
      if (index > -1) {
        this.liked = true;
      }
    },
    async getInfo() {
      let request = {
        command: "info",
        data: { start: this.song },
      };
      let response = await fetch("/api/song", {
        method: "POST",
        headers: {
          "Content-Type": "application/json;charset=utf-8",
        },
        body: JSON.stringify(request),
      });
      if (response.ok) {
        let song = await response.json();
        this.info = song.Info;
        this.stickers = song.Stickers;
        this.albumArt = song.AlbumArt;
        this.isItLiked();
      }
    },
  },
  async mounted() {
    this.getInfo();
  },
};
</script>
