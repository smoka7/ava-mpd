<template>
  <div
    class="flex h-full w-full flex-col justify-between gap-2 rounded p-2 dark:text-white md:flex-row md:flex-wrap"
  >
    <div class="card-class">
      <theme-settings />
    </div>
    <div class="card-class" v-if="connected">
      <database-info />
    </div>
    <div class="card-class" v-if="connected">
      <playback-options />
    </div>
    <div class="card-class space-y-1" v-if="connected">
      <outputs />
    </div>
    <div class="card-class" v-if="connected">
      <DownloadCover />
    </div>
  </div>
</template>
<script setup>
import DatabaseInfo from "./databaseInfo.vue";
import Outputs from "./outputs.vue";
import PlaybackOptions from "./playbackOptions.vue";
import ThemeSettings from "./themeSettings.vue";
import DownloadCover from "./downloadCover.vue";
import { useStore } from "vuex";
import { computed } from "vue";

const store = useStore();
const connected = computed(() => store.state.connected);

store.dispatch("getSettings");
</script>
<style lang="postcss">
.card-class {
  @apply w-full rounded bg-white p-4 dark:bg-gray-700 md:m-1 md:w-[48%];
}

.card-header {
  @apply mb-2 text-2xl  underline decoration-secondary  decoration-2;
}
</style>
