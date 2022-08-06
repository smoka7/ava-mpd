<template>
  <div class="flex h-full w-full flex-col">
    <div
      class="mb-10 space-y-1 overflow-y-auto overflow-x-hidden md:mb-6 md:px-2"
    >
      <Album
        v-for="(album, index) in queue.Albums"
        :key="index"
        :album="album"
        :currentAlbum="currentAlbum(album)"
        :selectedIds="state.selectedIds"
        :currentSongPos="currentSongPos"
        @showMenu="showMenu"
        @selectAlbum="selectAlbum"
        @select="selectSong"
      />
    </div>
    <teleport to="#app">
      <queue-menu
        :open="state.menu"
        :song="state.selected"
        :selected-ids="state.selectedIds"
        @close="hideMenu"
        @clearSelection="state.selectedIds = []"
        @showInfo="state.songInfo = true"
      />
    </teleport>
    <songInfo
      v-if="state.songInfo"
      :song="state.selected.id"
      @close="state.songInfo = false"
    />
    <div
      v-if="queue.Length == 0"
      class="flex h-full w-full items-center justify-center p-4 text-7xl underline decoration-accent md:text-9xl"
    >
      {{ message }}
    </div>
    <div
      v-if="!state.songInfo"
      class="h-content absolute right-0 bottom-0 left-0 flex w-full flex-wrap-reverse items-center justify-between space-y-2 rounded-b bg-secondary px-4 text-base text-primary md:space-x-4 md:space-y-0 md:py-1 md:text-sm"
    >
      <span>
        {{ queue.Length }} Tracks / duration:
        {{ humanizeTime(queue.Duration) }}
      </span>
      <span v-if="state.selectedIds.length > 0">
        {{ state.selectedIds.length }} selected
      </span>
      <queuePagination @goToCurrent="scrollToCurrentSong" />
    </div>
  </div>
</template>
<script setup>
import { useStore } from "vuex";
import { humanizeTime } from "../helpers.js";
import {
  defineAsyncComponent,
  shallowReactive,
  reactive,
  watch,
  computed,
  onMounted,
} from "vue";

const songInfo = defineAsyncComponent(() => import("./songInfo.vue"));
const queueMenu = defineAsyncComponent(() => import("./queueMenu.vue"));

const store = useStore();

const state = reactive({
  menu: false,
  songInfo: false,
  selected: { id: -1, pos: -1 },
  selectedIds: [],
});

const queue = computed(() => shallowReactive(store.state.queue));

const currentSongPos = computed(() => store.state.currentSong.Pos);

const message = computed(() => {
  if (!store.state.connected) {
    return "Couldn't Connect to server!!";
  }
  if (queue.value.Length == 0) {
    return "Queue is empty!";
  }
  return "";
});

onMounted(() => {
  scrollToCurrentSong();
});

function showMenu(pos, id) {
  state.selected.pos = Number(pos);
  state.selected.id = Number(id);
  state.menu = true;
}

function hideMenu() {
  state.menu = false;
}

function selectAlbum(name) {
  const album = queue.value.Albums.find((album) => {
    return album.Album == name;
  });

  if (album !== null) {
    album.Songs.forEach((song) => {
      selectSong(song.Id);
    });
  }
}

function selectSong(id) {
  const index = state.selectedIds.indexOf(id);
  if (index > -1) {
    state.selectedIds.splice(index, 1);
    return;
  }
  state.selectedIds.push(id);
}

async function scrollToCurrentSong() {
  if (queue.value.CurrentPage != queue.value.CurrentSongPage) {
    await store.dispatch("getQueue", queue.value.CurrentSongPage);
  }
  const el = document.getElementById("currentSong");
  if (el) el.scrollIntoView({ block: "center", behavior: "smooth" });
}

function currentAlbum(album) {
  return (
    Number(currentSongPos.value) >= Number(album.Songs[0].Pos) &&
    Number(currentSongPos.value) <=
      Number(album.Songs[album.Songs.length - 1].Pos)
  );
}

watch(currentSongPos, () => {
  scrollToCurrentSong();
});
</script>
<style>
#currentSong {
  @apply bg-red-300 dark:text-primary !important;
}
</style>
