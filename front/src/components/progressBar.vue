<template>
  <progress
    class="h-2 w-full cursor-pointer rounded bg-lightest transition-all duration-100 ease-in hover:h-3"
    max="100"
    ref="progress"
    @click="seek"
    @mousemove="sendHover"
    :value="progressedPercentage"
  />
</template>
<script setup>
import { computed, ref } from "vue";

const progress = ref(null);

const props = defineProps({
  data: Object,
});
const emit = defineEmits(["seek", "sendhover"]);

function seek(e) {
  const seek = findSeekValue(e);
  emit("seek", seek);
}

function sendHover(e) {
  const seek = findSeekValue(e);
  emit("sendhover", seek);
}

function findSeekValue(e) {
  const percent = e.offsetX / progress.value.offsetWidth;
  return Math.floor(percent * props.data.max);
}

const progressedPercentage = computed(
  () => (props.data.value / props.data.max) * 100 || 0,
);
</script>
<style>
progress {
  -webkit-appearance: none;
}
::-webkit-progress-value {
  @apply rounded bg-gradient-to-r from-lighter via-accent to-secondary;
}
::-webkit-progress-bar {
  @apply rounded bg-lightest;
}
::-moz-progress-bar {
  @apply rounded bg-gradient-to-r from-lighter via-accent to-secondary;
}
</style>
