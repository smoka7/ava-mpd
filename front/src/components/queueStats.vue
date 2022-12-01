<template>
  <div
    class="h-content absolute inset-0 top-auto flex w-full flex-wrap-reverse items-center justify-between space-y-2 bg-secondary px-4 text-base text-white md:space-x-4 md:space-y-0 md:rounded-b md:py-1 md:text-sm"
  >
    <span
      class="flex w-full justify-between duration-200 hover:text-base md:w-auto md:divide-x-2 md:divide-white"
    >
      <span class="px-1 font-bold"> {{ length }} Songs </span>
      <span class="px-1 font-bold"> duration: {{ humanizeTime(duration) }} </span>
    </span>
    <span v-if="selectedLength > 0"> {{ selectedLength }} selected </span>
    <queuePagination @goToCurrent="emits('scrollToCurrentSong')" />
  </div>
</template>
<script setup>
import { humanizeTime } from "../helpers.js";
defineProps(["length", "duration", "selectedLength"]);
const emits = defineEmits(["scrollToCurrentSong"]);
</script>
