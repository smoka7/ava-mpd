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
      <h2 class="text-lg">themes:</h2>
      <form class="flex space-x-2 items-center" @change="toggleDarkMode">
        <input
          type="radio"
          name="theme"
          value="light"
          id="light"
          v-model="lightOrDark"
        />
        <label for="light">light</label>
        <input
          type="radio"
          name="theme"
          value="dark"
          id="dark"
          v-model="lightOrDark"
        />
        <label for="dark">dark</label>
      </form>
      <h2 class="text-lg">color schemes:</h2>
      <form @change="changeColorScheme" class="flex flex-col space-y-3">
        <div>
          <input
            type="radio"
            name="colorScheme"
            v-model="colorScheme"
            value="auto"
            id="auto"
          />
          <label for="auto" class="ml-2">auto</label>
        </div>
        <div v-for="(scheme, index) in colorSchemes" :key="index" class="">
          <input
            type="radio"
            name="colorScheme"
            :value="index"
            :id="index"
            v-model="colorScheme"
          />
          <label :for="index" class="ml-2">{{ index }} </label>
          <div
            :style="renderScheme(scheme)"
            class="h-8 w-full rounded shadow"
          ></div>
        </div>
      </form>
    </div>
    <div class="card-class">
      <h2 class="text-lg">database stats:</h2>
      <p>
        songs: <span>{{ databaseStats["songs"] }}</span>
      </p>
      <p>
        albums: <span>{{ databaseStats["albums"] }}</span>
      </p>
      <p>
        artists: <span>{{ databaseStats["artists"] }}</span>
      </p>
      <p>
        db_playtime:
        <span>{{ humanizeTime(databaseStats["db_playtime"]) }}</span>
      </p>
      <p>
        last update:
        <span>{{
          new Date(databaseStats["db_update"] * 1000).toString()
        }}</span>
      </p>
      <p>
        playtime: <span>{{ humanizeTime(databaseStats["playtime"]) }}</span>
      </p>
      <p>
        uptime: <span>{{ humanizeTime(databaseStats["uptime"]) }}</span>
      </p>
      <button
        aria-label="update-database"
        class="
          border-2 border-green-500
          text-green-500
          p-2
          my-2
          rounded
          hover:bg-green-500 hover:text-white
        "
        @click="updateDatabase"
      >
        <font-awesome-icon icon="database"></font-awesome-icon> Update the MPD
        database
      </button>
      <button
        aria-label="delete-cache"
        class="
          border-2 border-red-500
          text-red-500
          p-2
          my-2
          rounded
          hover:bg-red-500 hover:text-white
        "
        @click="deleteCache"
      >
        <font-awesome-icon icon="eraser"></font-awesome-icon> Delete the Cover
        Art cache
      </button>
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
          v-if="output['outputenabled'] == 1"
          type="checkbox"
          :name="output['outputname']"
          :aria-label="output['outputname']"
          checked
        />
        <input
          v-else
          type="checkbox"
          :aria-label="output['outputname']"
          :name="output['outputname']"
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
import { sendCommand, humanizeTime } from "../helpers.js";
import { colorSchemes, getRGB, setColorScheme } from "../colors.js";
export default {
  components: {
    FontAwesomeIcon,
  },
  emits: ["close"],
  data() {
    return {
      databaseStats: [],
      outputs: [],
      crossfade: 0,
      replayGainMods: ["off", "track", "album", "auto"],
      replayGain: "off",
      lightOrDark: localStorage.getItem("theme") || "light",
      colorScheme: localStorage.getItem("colorScheme") || "auto",
      colorSchemes: colorSchemes,
    };
  },
  methods: {
    updateDatabase() {
      sendCommand("/api/setting", "updateDatabase");
    },
    deleteCache() {
      sendCommand("/api/setting", "deleteCache");
    },
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
      if (this.lightOrDark === "light") {
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
    humanizeTime: humanizeTime,
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
