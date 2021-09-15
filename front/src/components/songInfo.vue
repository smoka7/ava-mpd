<template>
  <div
    class="bg-transparent absolute top-0 bottom-0 left-0 right-0"
    @click.self="$emit('close')"
  >
    <div
      class="
        flex flex-col
        absolute
        top-0
        left-0
        bottom-0
        right-0
        md:w-3/4 md:left-1/4 md:top-1
        app-height
        z-auto
        p-8
        space-y-2
        rounded
        overflow-auto
        bg-lightest
      "
    >
      <div class="flex flex-row align-baseline justify-between">
        <h1 class="text-primary text-bold text-4xl">song Info</h1>
        <button
          aria-label="close-info"
          @click="$emit('close')"
          class="px-4 py-2 rounded-full"
        >
          <font-awesome-icon
            class="text-red-500"
            icon="times"
            size="2x"
          ></font-awesome-icon>
        </button>
      </div>
      <div
        class="
          flex flex-col
          md:flex-row
          justify-around
          md:space-x-2 md:space-y-0
          space-y-2
        "
      >
        <album-art :url="albumArt" class="md:w-1/2"></album-art>
        <dl
          class="
            flex flex-col
            rounded
            p-4
            bg-white
            dark:text-white dark:bg-gray-700
          "
        >
          <p v-for="(value, index) in info" :key="index" class="flex space-x-2">
            <dt>{{ index }}:</dt>
            <dd>
              {{ value }}
            </dd>
          </p>
          <p v-for="(sticker, index) in stickers" :key="index">
            {{ sticker.Name }} : {{ sticker.Value }}
          </p>
        </dl>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import albumArt from "./albumArt.vue";
export default {
  components: { FontAwesomeIcon, albumArt },
  props: ["song"],
  data() {
    return {
      info: {},
      stickers: {},
      albumArt: "",
    };
  },
  methods: {},
  async mounted() {
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
    }
  },
};
</script>
