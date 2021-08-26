<template>
  <progress
    class="w-full bg-lightest rounded h-4 cursor-pointer"
    max="100"
    ref="progress"
    @click="seek"
    @mousedown.prevent="MouseDown"
    @mouseup.prevent="MouseUp"
    :value="progressedPercentage"
  ></progress>
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
      let percent = e.offsetX / this.$refs.progress.offsetWidth;
      let seek = Math.floor(percent * this.data.max);
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
  @apply bg-gradient-to-r from-lighter via-accent to-secondary rounded;
}
::-webkit-progress-bar {
  @apply bg-lightest rounded;
}
::-moz-progress-bar {
  @apply bg-gradient-to-r from-lighter via-accent to-secondary rounded;
}
</style>
