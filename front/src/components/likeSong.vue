<template>
  <button aria-label="like-song" :class="btnClass" @click="likeSong">
    <font-awesome-icon :icon="[pLiked ? 'fas' : 'far', 'heart']" size="2x" />
  </button>
</template>
<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import endpoints from "../endpoints.js";
import { sendCommand } from "../helpers";
import { ref } from "vue";

const props = defineProps({
  pLiked: Boolean,
  file: String,
});

const classes = {
  normal: "px-2 md:p-0 text-red-500 transition-all duration-200 transform",
  scale: "scale-105",
  right: "rotate-45",
  left: "-rotate-45",
};

const btnClass = ref([classes.normal]);

function likeSong() {
  sendCommand(endpoints.song, "like", { File: props.file });
  btnClass.value = [classes.normal, classes.scale, classes.right];
  setTimeout(() => {
    btnClass.value = [classes.normal, classes.scale, classes.left];
  }, 150);
  setTimeout(() => {
    btnClass.value = classes.normal;
  }, 300);
}
</script>
