<template>
  <div class="tooltip flex w-full items-center justify-center space-x-2 px-2">
    <button
      @click="changeVolume()"
      :aria-label="volume > 0 ? 'mute' : 'max-the-volume'"
      class="flex items-center"
    >
      <font-awesome-icon :icon="['fas', icon()]" size="lg" />
    </button>
    <progressBar :data="{ value: volume, max: 100 }" @seek="changeVolume" />
    <span class="tooltiptext">volume {{ volume }}%</span>
  </div>
</template>
<script setup lang="ts">
import { sendCommand } from "../helpers";
import endpoints from "../endpoints";

const props = defineProps(["volume"]);

function changeVolume(volume) {
  if (volume > 100) volume = 100;
  if (volume < 0) volume = 0;
  if (volume === undefined) volume = props.volume > 0 ? 0 : 100;
  sendCommand(endpoints.playback, "changeVolume", { start: volume });
}

function icon() {
  if (props.volume >= 50) return "volume-up";
  else if (props.volume > 0) return "volume-down";
  else return "volume-off";
}
</script>
