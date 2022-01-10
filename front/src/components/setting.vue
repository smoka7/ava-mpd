<template>
  <div
    class="fixed md:static flex flex-col md:flex-row md:flex-wrap space-y-4 md:space-x-8 bottom-0 top-0 left-0 z-30 w-full app-height p-4 rounded overflow-y-auto bg-lightest/60 backdrop-blur-3xl dark:text-white"
  >
    <button
      aria-label="close-setting"
      @click="$emit('close')"
      class="fixed top-4 right-6 p-4 bg-red-50 rounded md:bg-transparent text-red-500"
    >
      <font-awesome-icon icon="times" size="2x"></font-awesome-icon>
    </button>
    <div class="card-class">
      <theme-settings />
    </div>
    <div class="card-class">
      <database-info :stats="databaseStats" />
    </div>
    <div class="card-class">
      <playback-options
        :crossfade="crossfade"
        :replayGain="replayGain"
        @updatesetting="getSettings()"
      />
    </div>
    <div class="card-class">
      <outputs :outputs="outputs" @updatesetting="getSettings()" />
    </div>
  </div>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import DatabaseInfo from "./databaseInfo.vue";
import Outputs from "./outputs.vue";
import PlaybackOptions from "./playbackOptions.vue";
import ThemeSettings from "./themeSettings.vue";
export default {
  components: {
    FontAwesomeIcon,
    DatabaseInfo,
    ThemeSettings,
    Outputs,
    PlaybackOptions,
  },
  emits: ["close"],
  data() {
    return {
      databaseStats: [],
      outputs: [],
      crossfade: 0,
      replayGain: "off",
    };
  },
  methods: {
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
