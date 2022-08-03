<template>
  <div class="flex flex-col space-y-2">
    <h2 class="card-header">Playback Options</h2>
    <div class="flex flex-row items-center justify-between space-x-2">
      <label for="crossfade"> Crossfade : {{ storeCrossfade }} </label>
      <input
        class="rounded border border-blue-500 p-2 dark:bg-gray-700"
        name="crossfade"
        type="number"
        min="0"
        aria-label="crossfade"
        @input="setCrossfade"
        v-model="options.crossfade"
      />
    </div>
    <div class="flex flex-row items-center justify-between space-x-2">
      <label for="mixrampdb" class="w-full"
        >MixRampdb: {{ storeMixrampdb }}
      </label>
      <input
        class="w-1/2 rounded border border-blue-500 p-2 dark:bg-gray-700"
        name="mixrampdb"
        type="number"
        aria-label="mixrampdb"
        @input="setMixrampdb"
        v-model="options.mixrampdb"
      />
    </div>
    <div class="flex flex-row items-center justify-between space-x-2">
      <label for="replaygain"> Replay gain : {{ storeGain }}</label>
      <select
        name="replaygain"
        aria-label="replayGain"
        class="mt-2 rounded border border-blue-500 bg-white p-2 dark:bg-gray-700"
        @change="setReplayGain"
        v-model="options.gain"
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

const store = useStore();
const storeGain = computed(() => store.state.settings.ReplayGain);
const storeCrossfade = computed(() => store.state.status.xfade || 0);
const storeMixrampdb = computed(() => store.state.status.mixrampdb || 0);
const options = reactive({ gain: storeGain.value, mixrampdb: storeMixrampdb.value, crossfade: storeCrossfade.value });
const replayGainMods = ["off", "track", "album", "auto"];

function setCrossfade() {
  let second = Number(options.crossfade);
  if (second < 0) {
    second = 0;
    options.crossfade = 0;
  }
  sendCommand(endpoints.setting, "crossfade", { Value: second });
  store.dispatch("getSettings");
}

function setMixrampdb() {
  sendCommand(endpoints.setting, "mixrampdb", { Value: Number(options.mixrampdb) });
  store.dispatch("getSettings");
}

function setReplayGain() {
  const index = replayGainMods.indexOf(options.gain);
  sendCommand(endpoints.setting, "setGain", { Value: Number(index) });
  store.dispatch("getSettings");
}

function sendComm(command, data) {
  sendCommand(endpoints.setting, command, data);
  store.dispatch("getSettings");
}
</script>
