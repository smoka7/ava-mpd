<template>
  <h2 class="card-header">Database stats</h2>
  <ul class="-mx-4 list">
    <li>
      songs : <span>{{ stats.songs }}</span>
    </li>
    <li>
      albums : <span>{{ stats.albums }}</span>
    </li>
    <li>
      artists : <span>{{ stats.artists }}</span>
    </li>
    <li>
      Database play time:
      <span>{{ humanizeTime(stats.db_playtime) }}</span>
    </li>
    <li>
      Last database update:
      <span>{{ new Date(stats.db_update * 1000).toString() }}</span>
    </li>
    <li>
      Play time : <span>{{ humanizeTime(stats.playtime) }}</span>
    </li>
    <li>
      Up time : <span>{{ humanizeTime(stats.uptime) }}</span>
    </li>
  </ul>
  <div class="mt-2 flex flex-col items-start justify-start space-y-2">
    <button
      aria-label="update-database"
      class="rounded border-2 border-green-500 p-2 text-green-500 hover:bg-green-500 hover:text-white"
      @click="updateDatabase"
    >
      <font-awesome-icon icon="database" /> Update the MPD database
      <span v-if="updating != null" class="animate-pulse">
        updating the database...
      </span>
    </button>
    <button
      aria-label="delete-cache"
      class="rounded border-2 border-red-500 p-2 text-red-500 hover:bg-red-500 hover:text-white"
      @click="deleteCache"
    >
      <font-awesome-icon icon="eraser" /> Delete the Cover Art cache
    </button>
  </div>
</template>
<script setup lang="ts">
import { humanizeTime } from "../helpers.js";

const store = useStore();
const stats = computed(() => store.state.settings.DatabaseStats);
const updating = computed(() => store.state.status.updating_db);

function updateDatabase() {
  store.dispatch("sendCommandToSetting", {
    command: "updateDatabase",
  });
}

function deleteCache() {
  store.dispatch("sendCommandToSetting", {
    command: "deleteCache",
  });
}
</script>
