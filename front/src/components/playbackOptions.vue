<template>
  <h2 class="text-2xl mb-2">Playback Options</h2>
  <div class="flex flex-col space-y-2">
    <div class="space-x-2">
      <label for="crossfade"> Crossfade </label>
      <input
        class="border border-blue-500 p-2 rounded dark:bg-gray-700"
        name="crossfade"
        type="number"
        min="0"
        aria-label="crossfade"
        @input="setCrossfade"
        v-model="crossfade"
      />
    </div>
    <div class="space-x-2">
      <label for="replaygain"> Replay gain </label>
      <select
        name="replaygain"
        aria-label="replayGain"
        class="border border-blue-500 p-2 rounded bg-white dark:bg-gray-700 mt-2"
        @change="setReplayGain"
        v-model.lazy="Gain"
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
  </div>
</template>
<script>
import { sendCommand } from "../helpers.js";
export default {
  emits: ["updatesetting"],
  props: ["crossfade", "replayGain"],
  data() {
    return {
      replayGainMods: ["off", "track", "album", "auto"],
      Gain: this.replayGain,
    };
  },
  methods: {
    setCrossfade() {
      let second = Number(this.crossfade);
      if (second < 0) {
        second = 0;
        this.crossfade = 0;
      }
      sendCommand("/api/setting", "crossfade", { Start: second });
      this.$emit("updatesetting");
    },
    setReplayGain() {
      let index = this.replayGainMods.indexOf(this.Gain);
      sendCommand("/api/setting", "setGain", { Start: Number(index) });
      this.$emit("updatesetting");
    },
  },
};
</script>
