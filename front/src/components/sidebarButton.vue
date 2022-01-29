<template>
  <button
    :aria-label="label"
    class="tooltip rounded p-2 hover:bg-gray-700/10 dark:hover:bg-white/10"
  >
    <span class="tooltiptext">{{ label }}</span>
    <font-awesome-icon :icon="icon" />
  </button>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
export default {
  props: {
    label: String,
    icon: String,
  },
  components: { FontAwesomeIcon },
};
</script>
