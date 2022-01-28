<template>
  <h2 class="text-center text-3xl">Database stats</h2>
  <ul class="-mx-4 mt-2">
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      songs : <span>{{ stats["songs"] }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      albums : <span>{{ stats["albums"] }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      artists : <span>{{ stats["artists"] }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      Database play time:
      <span>{{ humanizeTime(stats["db_playtime"]) }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      <span>{{ new Date(stats["db_update"] * 1000).toString() }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      Play time : <span>{{ humanizeTime(stats["playtime"]) }}</span>
    </li>
    <li class="p-2 odd:bg-blue-100 dark:odd:bg-gray-800">
      UP time : <span>{{ humanizeTime(stats["uptime"]) }}</span>
    </li>
  </ul>
  <div>
    <button
      aria-label="update-database"
      class="my-2 rounded border-2 border-green-500 p-2 text-green-500 hover:bg-green-500 hover:text-white"
      @click="updateDatabase"
    >
      <font-awesome-icon icon="database"></font-awesome-icon> Update the MPD
      database
    </button>
    <button
      aria-label="delete-cache"
      class="my-2 rounded border-2 border-red-500 p-2 text-red-500 hover:bg-red-500 hover:text-white"
      @click="deleteCache"
    >
      <font-awesome-icon icon="eraser"></font-awesome-icon> Delete the Cover Art
      cache
    </button>
  </div>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { sendCommand, humanizeTime } from "../helpers.js";
export default {
  props: ["stats"],
  components: {
    FontAwesomeIcon,
  },
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
