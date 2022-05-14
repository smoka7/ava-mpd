<template>
  <Dialog
    :open="open"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <DialogOverlay class="fixed inset-0 bg-black/30 backdrop-blur-sm" />
    <div
      class="absolute top-1/4 right-10 flex max-w-xl flex-col rounded bg-white p-4 dark:bg-gray-800 dark:text-white"
    >
      <DialogTitle class="my-2 text-lg"> Add items to ... </DialogTitle>
      <div class="mt-2 flex w-full flex-col justify-between space-y-2">
        <button
          v-for="(position, index) in postitions"
          :key="position"
          @click="add(position)"
          class="rounded border-2 border-primary p-2 text-lg text-primary hover:bg-primary hover:text-white focus:bg-primary focus:text-white dark:border-accent dark:text-accent hover:dark:bg-accent hover:dark:text-primary focus:dark:bg-accent focus:dark:text-primary"
        >
          {{ texts[index] }}
        </button>
        <button
          @click="$emit('close')"
          class="rounded border-2 border-red-500 p-2 text-lg text-red-500 hover:bg-red-500 hover:text-white focus:bg-red-500 focus:text-white"
        >
          Cancel
        </button>
      </div>
    </div>
  </Dialog>
</template>
<script setup>
import { Dialog, DialogOverlay, DialogTitle } from "@headlessui/vue";

defineProps(["open"]);
const emit = defineEmits(["close", "add"]);
const postitions = ["currentSong", "endOfQueue", "currentAlbum"];
const texts = ["After Current Song", "End Of Queue", "After Current Album"];

function add(position) {
  emit("add", position);
  emit("close");
}
</script>
