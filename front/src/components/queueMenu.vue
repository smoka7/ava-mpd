<template>
  <div
    v-show="open"
    class="absolute top-0 bottom-0 left-0 right-0 z-50 bg-white bg-opacity-30"
    @click.self="$emit('close')"
  ></div>
  <div
    v-show="open"
    class="border-primary dark:border-lightest absolute right-10 z-50 flex w-48 flex-col rounded border-2 bg-white text-lg dark:bg-gray-600 dark:text-white"
    id="context-menu"
  >
    <button
      :class="btnClass"
      @click="action.func()"
      v-for="action in actions"
      :key="action.title"
    >
      {{ action.title }}
    </button>
    <details :class="btnClass">
      <summary @click="getStoredPlaylist" class="flex cursor-pointer" >
        add to playlist
      </summary>
        <button
          @click="addSongTo(pl.playlist)"
          :class="btnClass"
          v-for="pl in storedPlaylist"
          :key="pl.playlist"
        >
          {{ pl.playlist }}
        </button>
    </details>
  </div>
</template>
<script setup>
import { computed } from "vue";
import { useStore } from "vuex";
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

const btnClass =
  "w-full text-left p-2 hover:bg-blue-100 focus:bg-blue-100 dark:hover:text-black overflow-x-hidden text-ellipsis even:bg-blue-50 even:dark:bg-gray-800 dark:hover:bg-blue-100";

function deleteSong() {
  props.selectedIds.forEach((id) => {
    sendCommand(endpoints.queue, "delete", { Start: id });
  });

  sendCommand(endpoints.queue, "delete", { Start: props.song.id });
  emit("close");
}

function addSongTo(playlist) {
  props.selectedIds.forEach((id) => {
    sendCommand(endpoints.queue, "addsong", {
      Start: id,
      playlist: playlist,
    });
  });
  sendCommand(endpoints.queue, "addsong", {
    Start: props.song.id,
    playlist: playlist,
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
