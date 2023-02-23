<template>
  <div class="flex w-full items-center justify-between space-x-1 md:w-auto">
    <button
      v-if="queue.CurrentPage > 1"
      aria-label="Previous-page"
      @click="loadPage(queue.CurrentPage - 1)"
      class="pagination-button"
    >
      <font-awesome-icon icon="angle-right" flip="horizontal" />
      <span> previous </span>
    </button>
    <button
      v-if="queue.CurrentPage !== queue.LastPage"
      aria-label="next-page"
      @click="loadPage(queue.CurrentPage + 1)"
      class="pagination-button"
    >
      <span> next </span>
      <font-awesome-icon icon="angle-right" />
    </button>
    <button
      aria-label="goto-current-song"
      @click="$emit('goToCurrent')"
      class="pagination-button"
    >
      <font-awesome-icon icon="compact-disc" />
      <span> Current Song </span>
    </button>
  </div>
</template>
<script setup lang="ts">
import { useStore } from "../store";
const store = useStore();
const queue = computed(() => store.queue);
defineEmits(["goToCurrent"]);

function loadPage(page) {
  store.getQueue(page);
}
</script>
<style lang="postcss">
.pagination-button {
  @apply flex items-center justify-between space-x-1 rounded px-2 duration-200 hover:bg-primary hover:py-1.5 hover:text-lightest;
}
</style>
