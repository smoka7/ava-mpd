<template>
      <div
      class="flex flex-col absolute inset-0 z-auto p-8 space-y-2 overflow-y-scroll overflow-x-hidden bg-lightest"
    >
      <div class="flex flex-row items-center justify-between">
        <h1
          class="flex space-x-4 items-center text-bold underline decoration-2 text-4xl text-ellipsis font-bold mr-2"
        >
          <span>
            {{ info.Title }}
          </span>
          <likeSong :pLiked="liked" :file="info['file']" />
        </h1>
        <button
          aria-label="close-info"
          @click="$emit('close')"
          class="px-4 py-2 rounded-full"
        >
          <font-awesome-icon class="text-red-500" icon="times" size="2x" />
        </button>
      </div>
      <div
        class="flex flex-col md:flex-row justify-around md:space-x-2 md:space-y-0 space-y-2"
      >
        <album-art
          :url="albumArt"
          class="md:w-1/2 h-fit md:sticky top-0"
        ></album-art>
        <ul class="flex flex-col md:w-1/2 rounded bg-white dark:bg-gray-700">
          <li
            v-for="(value, index) in info"
            :key="index"
            class="even:bg-blue-50 dark:even:bg-gray-800 last:rounded-b first:rounded-t p-2"
          >
            {{ index }} : {{ value }}
          </li>
          <li
            v-for="(sticker, index) in stickers"
            :key="index"
            class="even:bg-blue-50 dark:even:bg-gray-800 last:rounded-b p-2"
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
        data: { song: this.song },
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
