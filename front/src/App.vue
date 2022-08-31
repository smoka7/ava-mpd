<template>
  <div
    class="flex h-screen w-screen flex-col overflow-hidden bg-gradient-to-br from-lightest via-lighter to-accent text-lg text-primary duration-300 dark:text-white md:flex-row md:space-x-2 md:p-2"
  >
    <div
      class="fixed inset-0 z-10 h-screen md:static md:h-full md:w-1/4"
      id="mediaController"
    >
      <media-controller />
    </div>
    <sidebar class="h-screen md:h-full md:w-3/4" id="queue" />
    <songInfo v-if="showSongInfo" @close="closeInfo" />
  </div>
</template>
<script setup>
import { loadTheme, setColorScheme } from "./colors.js";
import { handleKey } from "./keymap";

const store = useStore();

store.dispatch("connectToSocket");
store.dispatch("getCurrentSong");
store.dispatch("startCounter");
loadTheme();
setColorScheme();

document.addEventListener("keydown", listenKeyEvents);

onUnmounted(() => {
  document.removeEventListener("keydown", listenKeyEvents);
});

function listenKeyEvents(event) {
  if (event && event.target.tagName.toUpperCase() !== "INPUT") {
    handleKey(event.key);
  }
}

const songInfo = defineAsyncComponent(() => import("./components/songInfo.vue"));
const showSongInfo = computed(() => store.state.song.show);

async function closeInfo() {
  await store.dispatch("clearSongInfo");
}

</script>
<style lang="postcss">
* {
  @apply focus:outline-none focus:outline-accent;
}

.tooltip {
  @apply relative inline-block;
}
.tooltip .tooltiptext {
  @apply invisible absolute bottom-full left-1/2 z-10 -ml-10 inline w-auto rounded bg-lightest/90 p-2 text-center text-primary;
}
.tooltip:hover .tooltiptext {
  @apply visible;
}
::-webkit-scrollbar {
  @apply w-2;
}
::-webkit-scrollbar-track {
  @apply bg-lightest;
}
::-webkit-scrollbar-thumb {
  @apply bg-primary;
}
::-webkit-scrollbar-thumb:hover {
  @apply bg-secondary;
}
@supports (-moz-appearance: none) {
  .fbg {
    @apply bg-white/70 dark:bg-gray-700/70 !important;
  }
}
</style>
