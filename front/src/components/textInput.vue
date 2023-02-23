<template>
  <Dialog
    :open="isOpen"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <div class="fixed inset-0 bg-black/30 backdrop-blur-sm" />
    <DialogPanel
      class="relative top-1/3 mx-auto flex max-w-xl flex-col rounded bg-white p-4 dark:bg-gray-800 dark:text-white"
    >
      <DialogTitle class="my-2 text-lg">
        <slot name="title"></slot>
      </DialogTitle>
      <input
        type="text"
        v-model="inputText"
        class="w-full rounded bg-gray-200 p-2 text-black outline-none ring-2 ring-primary focus:ring-secondary dark:bg-gray-700 dark:text-white dark:ring-accent dark:focus:ring-secondary"
      />
      <div class="mt-2 flex w-full justify-end space-x-4">
        <button
          @click="
            $emit('finish', inputText);
            $emit('close');
          "
          class="rounded border-2 border-green-500 p-2 text-lg text-green-500 hover:bg-green-500 hover:text-white"
        >
          <slot name="btn"></slot>
        </button>
        <button
          @click="$emit('close')"
          class="rounded border-2 border-red-500 p-2 text-lg text-red-500 hover:bg-red-500 hover:text-white"
        >
          Cancel
        </button>
      </div>
    </DialogPanel>
  </Dialog>
</template>
<script setup lang="ts">
defineEmits(["close", "finish"]);
defineProps(["isOpen"]);
const inputText = ref("");
</script>
