<template>
  <button
    :aria-label="label"
    class="tooltip rounded p-2 hover:bg-gray-700/10 dark:hover:bg-white/10"
  >
    <span class="tooltiptext">{{ label }}</span>
    <font-awesome-icon :icon="icon" size="lg" />
  </button>
</template>
<script setup lang="ts">
defineProps<{
  label: string;
  icon: string;
}>();
</script>
