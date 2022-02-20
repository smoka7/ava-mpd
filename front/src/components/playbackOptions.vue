<template>
  <h2 class="mb-2 text-2xl">Playback Options</h2>
  <div class="flex flex-col space-y-2">
    <div class="space-x-2">
      <label for="crossfade"> Crossfade </label>
      <input
        class="rounded border border-blue-500 p-2 dark:bg-gray-700"
        name="crossfade"
        type="number"
        min="0"
        aria-label="crossfade"
        @input="setCrossfade"
        v-model="crossfade"
      />
    </div>
    <div class="space-x-2">
      <label for="replaygain"> Replay gain :{{ storeGain }}</label>
      <select
        name="replaygain"
        aria-label="replayGain"
        class="mt-2 rounded border border-blue-500 bg-white p-2 dark:bg-gray-700"
        @change="setReplayGain"
        v-model="gain.state"
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
<script setup>
import { sendCommand } from "../helpers.js";
import endpoints from "../endpoints.js";
import { useStore } from "vuex";
import { computed, reactive } from "vue";

const gain = reactive({ state: "off" });
const store = useStore();
const storeGain = computed(() => store.state.settings.ReplayGain);
const crossfade = computed(() => store.state.status.xfade || 0);
const replayGainMods = ["off", "track", "album", "auto"];

function setCrossfade() {
  let second = Number(crossfade);
  if (second < 0) {
    second = 0;
    crossfade = 0;
  }
  sendCommand(endpoints.setting, "crossfade", { Start: second });
  store.dispatch("getSettings");
}

function setReplayGain() {
  const index = replayGainMods.indexOf(gain.state);
  sendCommand(endpoints.setting, "setGain", { Start: Number(index) });
  store.dispatch("getSettings");
}
</script>
