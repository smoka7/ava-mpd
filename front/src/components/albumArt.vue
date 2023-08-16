<template>
  <div @click="zoom" class="relative -mx-8">
    <img
      v-if="!defaultAlbumArt"
      loading="lazy"
      :alt="altText"
      :class="albumArtClass"
      :src="url"
    />
    <slot></slot>
  </div>
</template>
<script setup lang="ts">
const props = defineProps(["url", "altText"]);
const isZoomed = ref(false);

function zoom() {
  isZoomed.value = !isZoomed.value;
}

const defaultAlbumArt = computed(() => props.url === "default");
const albumArtClass = computed(() => {
  return {
    "fixed z-50 top-2 bottom-2 w-screen md:w-1/2 inset-0 md:left-1/4 sm:w-1/2 sm:right-2 sm:left-auto":
      isZoomed.value,
    "w-full h-full absolute": !isZoomed.value,
    "rounded-tl cursor-pointer aspect-square": true,
  };
});
</script>
