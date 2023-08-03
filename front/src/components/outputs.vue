<template>
  <h2 class="card-header">List of Outputs</h2>
  <div
    v-for="(output, index) in outputs"
    :key="index"
    class="flex cursor-pointer items-center space-x-2 p-1"
    @click="toggleOutput(index)"
  >
    <input
      type="checkbox"
      class="accent-secondary"
      :checked="output.Enabled"
      :name="output.Name"
      :aria-label="output.Name"
    />
    <label :for="output.Name">
      {{ output.Name }}
    </label>
  </div>
</template>
<script setup lang="ts">
import { settingsStore } from "../settingsStore";

const store = settingsStore();
const outputs = computed(() => store.settings.Outputs);

function toggleOutput(index: number) {
  const id = Number(outputs.value[index].ID);

  store.sendCommandToSetting({
    command: "toggleOutput",
    data: { Value: id },
  });
}
</script>
