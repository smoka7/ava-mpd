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
const store = useStore();
const outputs = computed(() => store.state.settings.Outputs);

function toggleOutput(index) {
  const id = Number(outputs.value[index]["outputid"]);

  let command = "enableOutput";
  if (outputs.value[index]["outputenabled"] == 1) {
    command = "disableOutput";
  }

  store.dispatch("sendCommandToSetting", {
    command: command,
    data: { Value: id },
  });
}
</script>
