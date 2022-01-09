<template>
  <Dialog
    :open="isOpen"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <DialogOverlay class="fixed inset-0 bg-black opacity-30" />
    <div
      class="relative flex flex-col max-w-xl p-4 top-1/3 mx-auto bg-foreground dark:bg-gray-800 dark:text-foreground rounded"
    >
      <DialogTitle class="my-2 text-lg">
        <slot name="title"></slot>
      </DialogTitle>
      <input
        type="text"
        v-model="inputText"
        class="w-full p-2 rounded outline-none ring-2 ring-primary focus:ring-secondary text-black bg-gray-200 dark:bg-gray-700 dark:text-foreground dark:ring-lighter dark:focus:ring-lightest"
      />
      <div class="flex w-full justify-end mt-2 space-x-4">
        <button
          @click="
            $emit('finish', inputText);
            $emit('close');
          "
          class="p-2 border-2 border-green-500 text-green-500 text-lg hover:bg-green-500 hover:text-foreground rounded"
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
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
export default {
  props: ["isOpen"],
  emits: ["close", "finish"],
  components: {
    Dialog,
    DialogOverlay,
    DialogTitle,
    FontAwesomeIcon,
  },
  data() {
    return {
      inputText: "",
    };
  },
};
</script>
