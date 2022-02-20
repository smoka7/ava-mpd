<template>
  <div
    class="z-30 flex h-full w-full flex-col space-y-4 overflow-y-auto rounded bg-lightest/60 p-4 backdrop-blur-3xl dark:text-white md:flex-row md:flex-wrap md:space-x-8"
  >
    <button
      aria-label="close-setting"
      @click="$emit('close')"
      class="fixed top-4 right-6 rounded bg-red-50 p-4 text-red-500 md:bg-transparent"
    >
      <font-awesome-icon icon="times" size="2x"></font-awesome-icon>
    </button>
    <div class="card-class">
      <theme-settings />
    </div>
    <div class="card-class" v-if="connected">
      <database-info />
    </div>
    <div class="card-class" v-if="connected">
      <playback-options/>
    </div>
    <div v-if="connected" class="card-class">
      <outputs />
    </div>
  </div>
</template>
<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import DatabaseInfo from "./databaseInfo.vue";
import Outputs from "./outputs.vue";
import PlaybackOptions from "./playbackOptions.vue";
import ThemeSettings from "./themeSettings.vue";
import { useStore } from "vuex";
import { computed } from "vue";

const emit = defineEmits(["close"]);
const store=useStore();
const connected= computed(()=>store.state.connected);

store.dispatch("getSettings");
</script>
<style>
.card-class {
  @apply w-full rounded bg-white p-4 shadow-md dark:bg-gray-700 md:m-2 md:w-2/5;
}
</style>
