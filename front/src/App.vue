<template>
  <div
    class="flex h-screen w-screen flex-col overflow-hidden bg-gradient-to-br from-lightest via-lighter to-accent text-primary duration-300 dark:text-white md:flex-row md:space-x-2 md:p-2"
    v-if="store.connected"
  >
    <div
      class="fixed inset-0 z-10 h-[100dvh] flex-shrink-0 md:static md:h-full md:w-[30%]"
      id="mediaController"
      v-if="store.queue.Length > 0"
    >
      <media-controller />
    </div>
    <sidebar class="h-[100dvh] w-full md:h-full" id="queue" />
    <songInfo v-if="showSongInfo" @close="closeInfo" />
  </div>
  <reconnectButton v-else />
</template>
<script setup lang="ts">
import { loadTheme, setColorScheme } from "./colors";
import { handleKey } from "./keymap";
import { useStore } from "./store";

const store = useStore();

store.connectToSocket();
store.getCurrentSong();
store.startCounter();
loadTheme();
setColorScheme();

document.addEventListener("keydown", listenKeyEvents);

onUnmounted(() => {
  document.removeEventListener("keydown", listenKeyEvents);
});

function listenKeyEvents(event: KeyboardEvent) {
  if (!event.target) {
    return;
  }
  const pressedModifire = event.ctrlKey || event.altKey || event.shiftKey;
  const target = event.target as HTMLTextAreaElement;
  const pressdInInput = target.tagName.toUpperCase() == "INPUT";
  if (pressdInInput || pressedModifire) {
    return;
  }
  if (event) {
    handleKey(event.key);
  }
}

const songInfo = defineAsyncComponent(
  () => import("./components/songInfo.vue")
);
const showSongInfo = computed(() => store.song.show);

async function closeInfo() {
  store.clearSongInfo();
}
</script>
<style lang="postcss">
* {
  @apply focus:outline-none focus:ring-2 focus:ring-accent  !important;
  @apply focus:rounded;
}

.tooltip {
  @apply relative inline-block;
}
.tooltip .tooltiptext {
  @apply invisible absolute bottom-full left-1/2 z-10 -ml-10 flex w-fit scale-0 rounded bg-lightest/90 p-1 text-center text-sm text-primary duration-100;
}
.tooltip:hover .tooltiptext {
  @apply visible scale-100;
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

.hide-unfocused {
  @apply invisible group-focus-within:visible group-hover:visible group-focus:visible;
}
</style>
