<template>
  <div
    class="
      fixed
      md:static
      flex flex-col
      md:flex-row md:flex-wrap
      space-y-4
      bottom-0
      top-0
      left-0
      z-30
      w-full
      app-height
      p-4
      rounded
      overflow-y-auto
      bg-lightest
      dark:text-white
    "
  >
    <button
      aria-label="close-setting"
      @click="$emit('close')"
      class="fixed top-4 right-6 bg-primary px-4 py-2 rounded-full md:hidden"
    >
      <font-awesome-icon
        class="text-red-500"
        icon="times"
        size="2x"
      ></font-awesome-icon>
    </button>
    <div class="card-class">
      <theme-settings />
    </div>
    <div class="card-class">
      <database-info :stats="databaseStats" />
    </div>
    <div class="card-class">
      <h2 class="text-lg mb-2">playback options:</h2>
      <label for="crossfade"> crossfade </label>
      <input
        class="border border-blue-500 p-2 rounded dark:bg-gray-700"
        name="crossfade"
        type="number"
        min="0"
        aria-label="crossfade"
        @input="setCrossfade"
        v-model="crossfade"
      />
      <label for="replayain"> replayGain </label>
      <select
        name="replaygain"
        aria-label="replaygain"
        class="
          border border-blue-500
          p-2
          rounded
          bg-white
          dark:bg-gray-700
          mt-2
        "
        @change="setReplayGain"
        v-model.lazy="replayGain"
      >
        <option
          :value="mode"
          v-for="(mode, index) in replayGainMods"
          :key="index"
        >
          {{ mode }}
        </option>
      </select>
    </div>
    <div class="card-class"></div>
  </div>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import DatabaseInfo from "./databaseInfo.vue";
import ThemeSettings from "./themeSettings.vue";
import { sendCommand } from "../helpers.js";
export default {
  components: {
    FontAwesomeIcon,
    DatabaseInfo,
    ThemeSettings,
  },
  emits: ["close"],
  data() {
    return {
      databaseStats: [],
      outputs: [],
      crossfade: 0,
      replayGainMods: ["off", "track", "album", "auto"],
      replayGain: "off",
    };
  },
  methods: {
    toggleOutput(index) {
      let id = Number(this.outputs[index]["outputid"]);
      if (this.outputs[index]["outputenabled"] == 1) {
        sendCommand("/api/setting", "disableOutput", { Start: id });
        return;
      }
      sendCommand("/api/setting", "enableOutput", { Start: id });
    },
    setCrossfade() {
      let second = Number(this.crossfade);
      if (second < 0) {
        second = 0;
        this.crossfade = 0;
      }
      sendCommand("/api/setting", "crossfade", { Start: second });
    },
    setReplayGain() {
      let index = this.replayGainMods.indexOf(this.replayGain);
      sendCommand("/api/setting", "setGain", { Start: Number(index) });
    },
    async getSettings() {
      let response = await fetch("/api/setting");
      if (response.ok) {
        let json = await response.json();
        this.outputs = json.outputs;
        this.databaseStats = json.databaseStats;
        this.replayGain = json.replayGain.replay_gain_mode;
        return;
      }
      console.log(response.error);
    },
  },
  async created() {
    this.crossfade = Number(this.$store.state.status.xfade || 0);
    await this.getSettings();
  },
};
</script>
<style>
.card-class {
  @apply w-full md:w-2/5 md:m-2 p-4 shadow-md bg-white dark:bg-gray-700 rounded;
}
</style>
