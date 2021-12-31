<template>
  <div @click.self="zoom" class="avatar rounded">
    <img
      v-show="!defaultAlbumArt"
      @click="zoom"
      loading="lazy"
      :class="albumArtClass"
      :src="defaultAlbumArt ? '' : url"
    />
  </div>
</template>
<script>
export default {
  props: ["url"],
  data() {
    return { isZoomed: false };
  },
  methods: {
    zoom() {
      this.isZoomed = !this.isZoomed;
    },
  },
  computed: {
    defaultAlbumArt() {
      return this.url === "default";
    },
    albumArtClass() {
      return {
        "fixed z-50 top-2 bottom-2 w-screen h-auto md:w-1/2 aspect-square inset-0 md:left-1/4":
          this.isZoomed,
        "w-full aspect-square": !this.isZoomed,
        "rounded cursor-pointer": true,
      };
    },
  },
};
</script>
<style>
.avatar {
  background: radial-gradient(50% 123.47% at 50% 50%, #00ff94 0%, #720059 100%),
    linear-gradient(121.28deg, #669600 0%, #ff0000 100%),
    linear-gradient(360deg, #0029ff 0%, #8fff00 100%),
    radial-gradient(100% 164.72% at 100% 100%, #6100ff 0%, #00ff57 100%),
    radial-gradient(100% 148.07% at 0% 0%, #fff500 0%, #51d500 100%);
  background-blend-mode: screen, color-dodge, overlay, difference, normal;
}
</style>
