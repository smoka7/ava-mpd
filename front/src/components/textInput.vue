<template>
  <Dialog
    :open="isOpen"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <DialogOverlay class="fixed inset-0 bg-black/30 backdrop-blur-sm" />
    <div
      class="relative flex flex-col max-w-xl p-4 top-1/3 mx-auto bg-white dark:bg-gray-800 dark:text-white rounded"
    >
      <DialogTitle class="my-2 text-lg">
        <slot name="title"></slot>
      </DialogTitle>
      <input
        type="text"
        v-model="inputText"
        class="w-full p-2 rounded outline-none ring-2 ring-primary focus:ring-secondary text-black bg-gray-200 dark:bg-gray-700 dark:text-white dark:ring-accent dark:focus:ring-secondary"
      />
      <div class="flex w-full justify-end mt-2 space-x-4">
        <button
          @click="
            $emit('finish', inputText);
            $emit('close');
          "
          class="p-2 border-2 border-green-500 text-green-500 text-lg hover:bg-green-500 hover:text-white rounded"
        >
          <slot name="btn"></slot>
        </button>
        <button
          @click="$emit('close')"
          class="p-2 border-2 border-red-500 text-red-500 text-lg hover:bg-red-500 hover:text-white rounded"
        >
          Cancel
        </button>
      </div>
    </div>
  </Dialog>
</template>
<script>
import { Dialog, DialogOverlay, DialogTitle } from "@headlessui/vue";
export default {
  props: ["isOpen"],
  emits: ["close", "finish"],
  components: {
    Dialog,
    DialogOverlay,
    DialogTitle,
  },
  data() {
    return {
      inputText: "",
    };
  },
};
</script>
