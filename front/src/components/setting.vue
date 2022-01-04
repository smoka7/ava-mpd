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
    <div class="card-class space-y-1">
      <SwitchGroup>
        <SwitchLabel class="mr-2">Dark Mode</SwitchLabel>
        <Switch
          v-model="isDark"
          :class="isDark ? 'bg-teal-900' : 'bg-teal-700'"
          isDark
          class="
            relative
            bg-lighter
            inline-flex
            items-center
            h-6
            rounded
            w-11
          "
        >
          <span class="sr-only">Enable notifications</span>
          <span
            :class="isDark ? 'translate-x-6' : 'translate-x-1'"
            class="inline-block w-4 h-4 transform bg-primary rounded"
            isDark
          />
        </Switch>
      </SwitchGroup>
      <h2 class="text-lg">color schemes:</h2>
      <RadioGroup v-model="colorScheme" class="flex flex-col space-y-3">
        <template v-for="(scheme, index) in colorSchemes" :key="index">
          <RadioGroupOption v-slot="{ checked }" :value="index">
            <div
              :class="[
                { 'bg-blue-200 text-primary rounded': checked },
                'p-2 cursor-pointer',
              ]"
            >
              <div class="flex justify-between">
                <RadioGroupLabel>{{ index }}</RadioGroupLabel>
                <font-awesome-icon
                  v-show="checked"
                  class="text-primary ml-auto"
                  icon="check-circle"
                ></font-awesome-icon>
              </div>
              <div
                :style="renderScheme(scheme)"
                class="h-8 w-full rounded shadow"
              ></div>
            </div>
          </RadioGroupOption>
        </template>
      </RadioGroup>
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
    <div class="card-class">
      <h2 class="text-lg mb-2">list of outputs:</h2>
      <div
        v-for="(output, index) in outputs"
        :key="index"
        class="space-x-2 cursor-pointer"
        @click="toggleOutput(index)"
      >
        <input
          type="checkbox"
          class="accent-primary"
          :checked="output['outputenabled'] == 1"
          :name="output['outputname']"
          :aria-label="output['outputname']"
        />
        <label :for="output['outputname']">
          {{ output["outputname"] }}
        </label>
      </div>
    </div>
  </div>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import {
  RadioGroup,
  RadioGroupLabel,
  RadioGroupOption,
  Switch,
  SwitchGroup,
  SwitchLabel,
} from "@headlessui/vue";
import DatabaseInfo from "./databaseInfo.vue";
import { sendCommand } from "../helpers.js";
import { colorSchemes, getRGB, setColorScheme } from "../colors.js";
export default {
  components: {
    FontAwesomeIcon,
    RadioGroup,
    RadioGroupLabel,
    RadioGroupOption,
    Switch,
    SwitchGroup,
    SwitchLabel,
    DatabaseInfo,
  },
  emits: ["close"],
  data() {
    return {
      databaseStats: [],
      outputs: [],
      crossfade: 0,
      replayGainMods: ["off", "track", "album", "auto"],
      replayGain: "off",
      isDark: localStorage.getItem("theme") == "dark",
      colorScheme: localStorage.getItem("colorScheme") || "first",
      colorSchemes: colorSchemes,
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
    toggleDarkMode() {
      if (!this.isDark) {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("theme", "light");
        return;
      }
      localStorage.setItem("theme", "dark");
      document.documentElement.classList.add("dark");
    },
    changeColorScheme() {
      localStorage.setItem("colorScheme", this.colorScheme);
      setColorScheme();
    },
    renderScheme(scheme) {
      return (
        "background-image: linear-gradient(90deg," +
        getRGB(scheme[0]) +
        " 0%," +
        getRGB(scheme[0]) +
        " 20%," +
        getRGB(scheme[1]) +
        " 20%," +
        getRGB(scheme[1]) +
        " 40%," +
        getRGB(scheme[2]) +
        " 40%," +
        getRGB(scheme[2]) +
        " 60%," +
        getRGB(scheme[3]) +
        " 60%," +
        getRGB(scheme[3]) +
        " 80%," +
        getRGB(scheme[4]) +
        " 80%," +
        getRGB(scheme[4]) +
        " 100%);"
      );
    },
  },
  async created() {
    this.crossfade = Number(this.$store.state.status.xfade || 0);
    await this.getSettings();
  },
  watch: {
    colorScheme(newCS) {
      this.changeColorScheme(newCS);
    },
    isDark() {
      this.toggleDarkMode();
    },
  },
};
</script>
<style>
.card-class {
  @apply w-full md:w-2/5 md:m-2 p-4 shadow-md bg-white dark:bg-gray-700 rounded;
}
</style>
