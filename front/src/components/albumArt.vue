<template>
  <div
    @click.self="zoom"
    class="from-lightest to-accent via-lighter rounded bg-gradient-to-br"
  >
    <img
      v-if="!defaultAlbumArt"
      @click="zoom"
      loading="lazy"
      :alt="altText"
      :class="albumArtClass"
      :src="defaultAlbumArt ? '' : url"
    />
  </div>
</template>
<script>
export default {
  props: ["url", "altText"],
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
<style></style>
