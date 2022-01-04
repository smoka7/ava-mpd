<template>
  <h2 class="text-3xl text-center">Database stats</h2>
  <ul class="-mx-4 mt-2">
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      songs : <span>{{ stats["songs"] }}</span>
    </li>
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      albums : <span>{{ stats["albums"] }}</span>
    </li>
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      artists : <span>{{ stats["artists"] }}</span>
    </li>
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      Database play time:
      <span>{{ humanizeTime(stats["db_playtime"]) }}</span>
    </li>
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      <span>{{ new Date(stats["db_update"] * 1000).toString() }}</span>
    </li>
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      Play time : <span>{{ humanizeTime(stats["playtime"]) }}</span>
    </li>
    <li class="odd:bg-blue-100 dark:odd:bg-gray-800 p-2">
      UP time : <span>{{ humanizeTime(stats["uptime"]) }}</span>
    </li>
  </ul>
  <button
    aria-label="update-database"
    class="border-2 border-green-500 text-green-500 p-2 my-2 rounded hover:bg-green-500 hover:text-white"
    @click="updateDatabase"
  >
    <font-awesome-icon icon="database"></font-awesome-icon> Update the MPD
    database
  </button>
  <button
    aria-label="delete-cache"
    class="border-2 border-red-500 text-red-500 p-2 my-2 rounded hover:bg-red-500 hover:text-white"
    @click="deleteCache"
  >
    <font-awesome-icon icon="eraser"></font-awesome-icon> Delete the Cover Art
    cache
  </button>
</template>
<script>
import { sendCommand, humanizeTime } from "../helpers.js";
export default {
  props: ["stats"],
  methods: {
    updateDatabase() {
      sendCommand("/api/setting", "updateDatabase");
    },
    deleteCache() {
      sendCommand("/api/setting", "deleteCache");
    },
    humanizeTime: humanizeTime,
  },
};
</script>
