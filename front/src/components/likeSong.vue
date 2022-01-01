<template>
  <button aria-label="like-song" :class="btnClass" @click="likeSong">
    <font-awesome-icon :icon="[pLiked ? 'fas' : 'far', 'heart']" size="lg" />
  </button>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { sendCommand } from "../helpers";
export default {
  props: {
    pLiked: Boolean,
    file: "",
  },
  components: {
    FontAwesomeIcon,
  },
  data() {
    return {
      btnClass:
        "px-2 md:p-0 text-red-500 transition-all duration-200 transform",
    };
  },
  methods: {
    likeSong() {
      const classes = {
        normal:
          "px-2 md:p-0 text-red-500 transition-all duration-200 transform",
        scale: "scale-105",
        right: "rotate-45",
        left: "-rotate-45",
      };
      sendCommand("/api/song", "like", { song: this.file });
      this.btnClass = [classes.normal, classes.scale, classes.right];
      setTimeout(() => {
        this.btnClass = [classes.normal, classes.scale, classes.left];
      }, 150);
      setTimeout(() => {
        this.btnClass = classes.normal;
      }, 300);
    },
  },
};
</script>
