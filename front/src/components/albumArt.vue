<template>
  <div
    @click.self="zoom"
    :class="
      isZoomed
        ? 'fixed top-12 left-0 md:left-1/3 w-screen h-auto md:w-auto z-50'
        : ''
    "
  >
    <div v-if="defaultAlbumArt" class="rounded w-full h-full avatar"></div>
    <img
      v-show="!defaultAlbumArt"
      @click="zoom"
      loading="lazy"
      class="rounded cursor-pointer shadow-md"
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
