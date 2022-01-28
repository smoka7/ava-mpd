<template>
  <h2 class="mb-2 text-2xl">List of Outputs</h2>
  <div
    v-for="(output, index) in outputs"
    :key="index"
    class="flex cursor-pointer items-center space-x-2 p-1"
    @click="toggleOutput(index)"
  >
    <input
      type="checkbox"
      class="accent-primary dark:accent-lightest"
      :checked="output['outputenabled'] == 1"
      :name="output['outputname']"
      :aria-label="output['outputname']"
    />
    <label :for="output['outputname']">
      {{ output["outputname"] }}
    </label>
  </div>
</template>
<script>
import { sendCommand } from "../helpers.js";
export default {
  props: ["outputs"],
  emits: ["updatesetting"],
  methods: {
    toggleOutput(index) {
      let id = Number(this.outputs[index]["outputid"]);
      if (this.outputs[index]["outputenabled"] == 1) {
        sendCommand("/api/setting", "disableOutput", { Start: id });
        return;
      }
      sendCommand("/api/setting", "enableOutput", { Start: id });
      this.$emit("updatesetting");
    },
  },
};
</script>
