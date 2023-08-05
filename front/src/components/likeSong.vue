<template>
  <button aria-label="like-song" :class="btnClass" @click="likeSong">
    <font-awesome-icon
      :icon="[pLiked ? 'fas' : 'far', 'heart']"
      :class="[pLiked ? 'text-rose-500' : color || 'text-rose-500']"
      :size="iconSize"
    />
  </button>
</template>
<script setup lang="ts">
import { useStore } from "../store";

const props = defineProps<{
  pLiked?: boolean;
  color?: string;
  file?: string;
  iconSize: "2x" | "lg";
}>();

const classes = {
  normal:
    "px-2 md:p-0 text-rose-500 transition-all duration-200 transform hover:scale-125",
  scale: "scale-105",
  right: "rotate-45",
  left: "-rotate-45",
};

const btnClass = ref([classes.normal]);
const store = useStore();

function likeSong() {
  store.toggleLike({ File: props.file });

  btnClass.value = [classes.normal, classes.scale, classes.right];

  setTimeout(() => {
    btnClass.value = [classes.normal, classes.scale, classes.left];
  }, 150);

  setTimeout(() => {
    btnClass.value = [classes.normal];
  }, 300);
}
</script>
