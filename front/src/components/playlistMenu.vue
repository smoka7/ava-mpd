<template>
  <button class="ml-auto p-2" @click="open = !open">
    <FontAwesomeIcon icon="ellipsis-h" />
  </button>
  <Dialog
    :open="open"
    @close="open = false"
    class="fixed inset-0 z-10 overflow-y-auto"
  >
    <DialogOverlay
      v-show="open"
      class="fixed inset-0 z-0 bg-black/30 backdrop-blur-sm"
    />

    <transition name="fade">
      <div
                class="drop-blur-3xl absolute top-1/4 left-1/4 right-1/4 mx-auto flex w-52 flex-col divide-y rounded bg-white p-1 hover:divide-transparent dark:bg-gray-800 dark:text-white md:left-auto md:w-52 md:top-20 md:right-12"
      >
        <button
          v-for="action in actions"
          :key="action.name"
          @click="
            $emit(action.method);
            open = false;
          "
          class="menuItem flex items-center justify-start space-x-2"
        >
          <FontAwesomeIcon :icon="action.icon" class="text-accent" />
          <span>{{ action.name }}</span>
        </button>
      </div>
    </transition>
  </Dialog>
</template>
<script setup>
const open = ref(false);

defineEmits([
  "delete",
  "clearSelection",
  "clear",
  "rename",
  "removeDuplicate",
  "removeInvalid",
]);
const actions = [
  {
    name: "Delete Playlists",
    icon: "trash",
    method: "delete",
  },
  {
    name: "Clear Playlists",
    icon: "eraser",
    method: "clear",
  },
  {
    name: "Rename Playlists",
    icon: "edit",
    method: "rename",
  },
  {
    name: "Remove Duplicate songs",
    icon: "eraser",
    method: "removeDuplicate",
  },
  {
    name: "Remove Invalid songs",
    icon: "eraser",
    method: "removeInvalid",
  },
  {
    name: "Clear Selection",
    icon: "check-circle",
    method: "clearSelection",
  },
];
</script>
