<template>
  <button @click="changeVolume()" >
    <font-awesome-icon :icon="['fas', icon()]" size="lg" />
  </button>
  <progressBar
    class="md:w-40 w-full"
    :data="{ value: volume, max: 100 }"
    @seek="changeVolume"
  />
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import progressBar from "./progressBar.vue";
import { sendCommand } from "../helpers";
export default {
  components: {
    FontAwesomeIcon,
    progressBar,
  },
  props: ["volume"],
  methods: {
    changeVolume(volume) {
      if (volume > 100) volume = 100;
      if (volume < 0) volume = 0;
      if (volume === undefined) volume = this.volume > 0 ? 0 : 100;
      sendCommand("/api/playback", "changeVolume", { start: volume });
    },
    icon() {
      if (this.volume >= 50) return "volume-up";
      else if (this.volume > 0) return "volume-down";
      else return "volume-off";
    },
  },
};
</script>
