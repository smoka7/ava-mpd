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
<script setup>
import { sendCommand } from "../helpers.js";
import { computed } from "vue";
import { useStore } from "vuex";
import endpoints from "../endpoints.js";

const store = useStore();
const outputs = computed(() => store.state.settings.Outputs);

function toggleOutput(index) {
  const id = Number(outputs.value[index]["outputid"]);
  if (outputs.value[index]["outputenabled"] == 1) {
    sendCommand(endpoints.setting, "disableOutput", { Value: id });
    store.dispatch("getSettings");
    return;
  }
  sendCommand(endpoints.setting, "enableOutput", { Value: id });
  store.dispatch("getSettings");
}
</script>
