<template>
  <progress
    class="h-2 w-full cursor-pointer rounded bg-lightest transition-all duration-200 ease-in-out hover:h-4"
    max="100"
    ref="progress"
    @click="seek"
    @mousemove="sendHover"
    :value="progressedPercentage"
  />
</template>
<script setup lang="ts">
const progress = ref<HTMLInputElement | null>(null);

const props = defineProps<{
  data: { value: number; max: number };
}>();
const emit = defineEmits(["seek", "sendhover"]);

function seek(e: MouseEvent) {
  const seek = findSeekValue(e);
  emit("seek", seek);
}

function sendHover(e: MouseEvent) {
  const seek = findSeekValue(e);
  emit("sendhover", seek);
}

function findSeekValue(e: MouseEvent): number {
  const percent = e.offsetX / (progress.value?.offsetWidth || 0);
  return Math.floor(percent * props.data.max + 0.01 * props.data.max);
}

const progressedPercentage = computed(
  () => (props.data.value / props.data.max) * 100 || 0,
);
</script>
<style lang="postcss">
progress {
  appearance: none;
}
::-webkit-progress-value {
  @apply rounded bg-accent;
}
::-webkit-progress-bar {
  @apply rounded bg-lightest;
}
::-moz-progress-bar {
  @apply bg-accent;
}
</style>
