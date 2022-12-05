<template>
  <Dialog
    :open="open"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <div class="fixed inset-0 bg-black/30 backdrop-blur-sm" />
    <DialogPanel
      class="menu"
    >
      <DialogTitle class="m-2 text-lg"> Add items to ... </DialogTitle>
      <div class="mt-2 flex w-full flex-col justify-between space-y-2">
        <button
          v-for="(position, index) in postitions"
          :key="position"
          @click="add('add', position)"
          class="menuItem"
        >
          {{ texts[index] }}
        </button>
        <details v-if="storedPlaylist != null">
          <summary class="menuItem cursor-pointer">
            <FontAwesomeIcon
              icon="angle-right"
              class="open-rotate transform-gpu p-1 text-accent duration-200"
            />
            Stored Playlists
          </summary>
          <div class="flex flex-col space-y-1">
            <button
              v-for="(playlist, index) in storedPlaylist"
              :key="index"
              @click="add('addToPlaylist', playlist.Name)"
              class="menuItem"
            >
              {{ playlist.Name }}
            </button>
          </div>
        </details>
      </div>
    </DialogPanel>
  </Dialog>
</template>
<script setup>
defineProps(["open", "storedPlaylist"]);
const emit = defineEmits(["close", "add", "addToPlaylist"]);
const postitions = ["currentSong", "endOfQueue", "currentAlbum"];
const texts = ["After Current Song", "End Of Queue", "After Current Album"];

function add(command, position) {
  emit(command, position);
  emit("close");
}
</script>
