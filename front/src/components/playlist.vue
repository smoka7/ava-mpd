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
    <queue-menu
      :open="state.menu"
      :song="state.selected"
      :selected-ids="state.selectedIds"
      @close="hideMenu"
      @clearSelection="state.selectedIds = []"
      @showInfo="showInfo"
    />
    <div
      v-if="queue.Length == 0"
      class="flex h-full w-full items-center justify-center p-4 text-7xl underline decoration-accent md:text-9xl"
    >
      Queue is empty!
    </div>
    <QueueStats
      v-else
      @scrollToCurrentSong="scrollToCurrentSong"
      :length="queue.Length"
      :duration="queue.Duration"
      :selectedLength="state.selectedIds.length"
    />
  </div>
</template>
<script setup lang="ts">
const queueMenu = defineAsyncComponent(() => import("./queueMenu.vue"));

const store = useStore();

const state = reactive({
  menu: false,
  selected: { id: -1, pos: -1 },
  selectedIds: [],
});

const queue = computed(() => shallowReactive(store.state.queue));

const currentSongPos = computed(() => store.state.currentSong.Pos);

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

async function showInfo() {
  await store.dispatch("getSongInfo", state.selected.id);
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
