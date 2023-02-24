<template>
  <div class="flex flex-col space-y-2">
    <h2 class="card-header">Playback Options</h2>
    <div class="flex items-center justify-between space-x-2">
      <label for="crossfade">Crossfade:</label>
      <input
        class="rounded border border-accent p-2 dark:bg-gray-700"
        name="crossfade"
        type="number"
        min="0"
        aria-label="crossfade"
        v-model="crossFade"
      />
    </div>
    <div class="flex items-center justify-between space-x-2">
      <label for="mixrampdb" class="w-full">MixRampdb: </label>
      <input
        class="w-1/2 rounded border border-accent p-2 dark:bg-gray-700"
        name="mixrampdb"
        type="number"
        aria-label="mixrampdb"
        v-model="mixRampDB"
      />
    </div>
    <div class="flex items-center justify-between space-x-2">
      <label for="replaygain"> Replay gain:</label>
      <select
        name="replaygain"
        aria-label="replayGain"
        class="mt-2 rounded border border-accent bg-white p-2 dark:bg-gray-700"
        v-model="replayGain"
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
<script setup lang="ts">
import { useStore } from "../store";
const store = useStore();
const replayGainMods = ["off", "track", "album", "auto"];

const replayGain = computed({
  get() {
    return store.settings.ReplayGain;
  },
  set(value: string) {
    const index = replayGainMods.indexOf(value);
    store.sendCommandToSetting({
      command: "setGain",
      data: { Value: Number(index) },
    });
  },
});

const mixRampDB = computed({
  get() {
    return store.status.mixrampdb || 0;
  },
  set(value) {
    store.sendCommandToSetting({
      command: "mixrampdb",
      data: {
        Value: Number(value),
      },
    });
  },
});

const crossFade = computed({
  get() {
    return store.status.xfade || 0;
  },
  set(value: number) {
    store.sendCommandToSetting({
      command: "crossfade",
      data: { Value: value < 1 ? 0 : value },
    });
  },
});
</script>
