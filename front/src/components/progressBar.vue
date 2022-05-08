<template>
  <progress
    class="h-2 w-full cursor-pointer rounded bg-lightest transition-all duration-100 ease-in hover:h-3"
    max="100"
    ref="progress"
    @click="seek"
    @mousedown.prevent="MouseDown"
    @mouseup.prevent="MouseUp"
    :value="progressedPercentage"
  />
</template>
<script>
export default {
  props: {
    data: Object,
  },
  emits: ["seek"],
  methods: {
    MouseDown() {
      this.$refs.progress.addEventListener("mousemove", this.seek);
    },
    MouseUp() {
      this.$refs.progress.removeEventListener("mousemove", this.seek);
    },
    seek(e) {
      const percent = e.offsetX / this.$refs.progress.offsetWidth;
      const seek = Math.floor(percent * this.data.max);
      this.$emit("seek", seek);
    },
  },
  computed: {
    progressedPercentage() {
      return (this.data.value / this.data.max) * 100 || 0;
    },
  },
};
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
