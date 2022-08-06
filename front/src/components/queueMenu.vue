<template>
  <Dialog
    :open="open"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto"
    id="context-menu"
  >
    <DialogOverlay
      v-show="open"
      class="fixed inset-0 z-0 bg-black/30 backdrop-blur-sm"
    />
    <div
      class="drop-blur-3xl absolute top-1/4 left-1/4 right-1/4 mx-auto flex w-52 flex-col divide-y rounded bg-white p-1 hover:divide-transparent dark:bg-gray-800 dark:text-white md:left-auto md:top-1/3 md:right-12"
    >
      <button
        class="menuItem"
        @click="action.func()"
        v-for="action in actions"
        :key="action.title"
      >
        {{ action.title }}
      </button>
      <details>
        <summary @click="getStoredPlaylist" class="menuItem flex">
          add to playlist
        </summary>
        <button
          @click="addSongTo(pl.Name)"
          class="menuItem"
          v-for="pl in storedPlaylist"
          :key="pl.Name"
        >
          {{ pl.Name }}
        </button>
      </details>
    </div>
  </Dialog>
</template>
<script setup>
import { sendCommand } from "../helpers.js";
import endpoints from "../endpoints.js";

const store = useStore();
const emit = defineEmits(["close", "showInfo", "clearSelection"]);
const storedPlaylist = computed(() => store.state.storedPlaylist);

const props = defineProps({
  open: Boolean,
  song: {
    id: Number,
    pos: Number,
  },
  selectedIds: Array,
});

const actions = [
  { title: "Delete", func: deleteSong },
  { title: "Show Song Info", func: showInfo },
  { title: "Shuffle Album", func: shuffleAlbum },
  { title: "Remove Duplicate songs", func: removeDuplicate },
  { title: "Clear Selection", func: clearSelection },
];

function deleteSong() {
  props.selectedIds.forEach((id) => {
    sendCommand(endpoints.queue, "delete", { ID: Number(id) });
  });

  sendCommand(endpoints.queue, "delete", { ID: props.song.id });
  emit("clearSelection");
  emit("close");
}

function addSongTo(playlist) {
  props.selectedIds.forEach((id) => {
    sendCommand(endpoints.queue, "addsong", {
      ID: Number(id),
      Playlist: playlist,
    });
  });
  sendCommand(endpoints.queue, "addsong", {
    ID: props.song.id,
    Playlist: playlist,
  });
  emit("close");
}

function shuffleAlbum() {
  sendCommand(endpoints.queue, "shuffleAlbum", {
    Start: props.song.pos,
  });
  emit("close");
}

function showInfo() {
  emit("showInfo");
  emit("close");
}
function clearSelection() {
  emit("clearSelection");
  emit("close");
}

function removeDuplicate() {
  sendCommand(endpoints.queue, "removeduplicate");
  emit("close");
}

function getStoredPlaylist() {
  if (storedPlaylist.length) return;
  store.dispatch("getStoredPlaylist");
}
</script>
<style>
.menuItem {
  @apply w-full overflow-x-hidden text-ellipsis rounded p-2  text-left hover:bg-blue-100 dark:hover:bg-blue-100 dark:hover:text-black;
}
</style>
