<template>
  <div class="flex w-full items-center justify-between space-x-1 md:w-auto">
    <button
      v-if="queue.CurrentPage > 1"
      aria-label="Previous-page"
      @click="loadPage(queue.CurrentPage - 1)"
      class="flex items-center justify-between space-x-1 px-2 hover:bg-primary hover:text-secondary"
    >
      <font-awesome-icon icon="angle-right" flip="horizontal" />
      <span> previous </span>
    </button>
    <button
      v-if="queue.CurrentPage !== queue.LastPage"
      aria-label="next-page"
      @click="loadPage(queue.CurrentPage + 1)"
      class="flex items-center justify-between space-x-1 px-2 hover:bg-primary hover:text-secondary"
    >
      <span> next </span>
      <font-awesome-icon icon="angle-right" />
    </button>
    <button
      aria-label="goto-current-song"
      @click="$emit('goToCurrent')"
      class="flex items-center justify-between space-x-1 px-2 hover:bg-primary hover:text-secondary"
    >
      <font-awesome-icon icon="compact-disc" />
      <span> Current Song </span>
    </button>
  </div>
</template>
<script setup>
import { useStore } from "vuex";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { computed } from "vue";

const store = useStore();
const queue = computed(() => store.state.queue);
defineEmits(["goToCurrent"]);

function loadPage(page) {
  store.dispatch("getQueue", page);
}
</script>
