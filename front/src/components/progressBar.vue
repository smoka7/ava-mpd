<template>
  <progress
    class="bg-lightest h-4 w-full cursor-pointer rounded"
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
  @apply from-lighter via-accent to-secondary rounded bg-gradient-to-r;
}
::-webkit-progress-bar {
  @apply bg-lightest h-4 rounded;
}
::-moz-progress-bar {
  @apply from-lighter via-accent to-secondary rounded bg-gradient-to-r;
}
</style>
