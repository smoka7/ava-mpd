<template>
  <h2 class="card-header">Database stats</h2>
  <ul class="-mx-4 mt-2">
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      songs : <span>{{ stats.songs }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      albums : <span>{{ stats.albums }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      artists : <span>{{ stats.artists }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      Database play time:
      <span>{{ humanizeTime(stats.db_playtime) }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      Last database update:
      <span>{{ new Date(stats.db_update * 1000).toString() }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      Play time : <span>{{ humanizeTime(stats.playtime) }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      Up time : <span>{{ humanizeTime(stats.uptime) }}</span>
    </li>
  </ul>
  <div class="flex flex-col">
    <button
      aria-label="update-database"
      class="my-2 rounded border-2 border-green-500 p-2 text-green-500 hover:bg-green-500 hover:text-white"
      @click="updateDatabase"
    >
      <font-awesome-icon icon="database" /> Update the MPD database
    </button>
    <p v-if="updating != null" class="animate-pulse">
      updating the database...
    </p>
    <button
      aria-label="delete-cache"
      class="my-2 rounded border-2 border-red-500 p-2 text-red-500 hover:bg-red-500 hover:text-white"
      @click="deleteCache"
    >
      <font-awesome-icon icon="eraser" /> Delete the Cover Art cache
    </button>
  </div>
</template>
<script setup>
import { sendCommand, humanizeTime } from "../helpers.js";
import endpoints from "../endpoints.js";
import { useStore } from "vuex";
import { computed } from "vue";

const store = useStore();
const stats = computed(() => store.state.settings.DatabaseStats);
const updating = computed(() => store.state.status.updating_db);

function updateDatabase() {
  sendCommand(endpoints.setting, "updateDatabase");
}

function deleteCache() {
  sendCommand(endpoints.setting, "deleteCache");
}
</script>
