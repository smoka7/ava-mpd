<template>
  <Dialog
    :open="open"
    @close="$emit('close')"
    class="fixed inset-0 z-10 overflow-y-auto"
    id="context-menu"
  >
    <div
      v-show="open"
      class="fixed inset-0 z-0 bg-black/30 backdrop-blur-sm"
    />
    <DialogPanel class="menu">
      <button
        class="menuItem group"
        @click="action.func()"
        v-for="action in actions"
        :key="action.title"
      >
        <FontAwesomeIcon :icon="action.icon" class="menuIcon" />
        <span>{{ action.title }}</span>
      </button>
      <details>
        <summary @click="getStoredPlaylist" class="menuItem group">
          <FontAwesomeIcon icon="plus" class="menuIcon" />
          <span>Add to playlist </span>
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
    </DialogPanel>
  </Dialog>
</template>
<script setup lang="ts">
import { sendCommand } from "../helpers.js";
import endpoints from "../endpoints";

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
  { title: "Delete", icon: "trash", func: deleteSong },
  { title: "Show Song Info", icon: "info-circle", func: showInfo },
  { title: "Shuffle Album", icon: "random", func: shuffleAlbum },
  { title: "Remove Duplicate songs", icon: "eraser", func: removeDuplicate },
  { title: "Clear Selection", icon: "check-circle", func: clearSelection },
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
.menu {
  @apply absolute top-1/4 left-1/4 right-1/4 mx-auto flex w-56 flex-col rounded bg-white p-1 dark:bg-gray-800 dark:text-white md:left-auto md:top-1/3 md:right-12;
}
.menuIcon{
    @apply text-secondary group-hover:text-accent group-focus:text-accent dark:text-lightest;
}
.menuItem {
  @apply flex w-full items-center justify-start space-x-2  overflow-x-hidden text-ellipsis rounded p-2 text-left hover:bg-gray-200 focus:outline-none focus:ring focus:ring-accent dark:hover:text-black;
}
</style>
