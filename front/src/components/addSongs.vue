<template>
  <Dialog
    :open="open"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <div class="fixed inset-0 bg-black/30 backdrop-blur-sm" />
    <DialogPanel class="menu">
      <DialogTitle class="m-2 text-lg"> Add items to ... </DialogTitle>
      <div class="mt-2 flex w-full flex-col justify-between space-y-2">
        <button
          v-for="(position, index) in postitions"
          :key="position"
          @click="
            emit('add', position);
            emit('close');
          "
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
              @click="
                emit('addToPlaylist', playlist.Name);
                emit('close');
              "
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
<script setup lang="ts">
import type { StoredPlaylist } from "../store";
defineProps<{ open: boolean; storedPlaylist?: Array<StoredPlaylist> }>();
type SongPosition = "currentSong" | "endOfQueue" | "currentAlbum";

const emit = defineEmits<{
  (e: "close"): void;
  (e: "add", value: SongPosition): void;
  (e: "addToPlaylist", value: string): void;
}>();

const postitions: Array<SongPosition> = [
  "currentSong",
  "endOfQueue",
  "currentAlbum",
];

const texts = ["After Current Song", "End Of Queue", "After Current Album"];
</script>
