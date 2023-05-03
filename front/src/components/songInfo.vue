<template>
  <Dialog
    :open="true"
    class="absolute inset-0 z-20 flex flex-col space-y-2 overflow-x-hidden overflow-y-scroll bg-lightest p-4 md:inset-4 md:rounded md:p-8"
  >
    <div class="flex flex-col-reverse items-center justify-start md:flex-row">
      <h1
        class="flex items-center space-x-4 text-ellipsis text-3xl font-bold text-primary underline decoration-accent md:mr-2"
      >
        <span>
          {{ state.Info["Title"] }}
        </span>
        <likeSong :pLiked="state.liked" :file="state.Info['file']" icon-size="2x"/>
      </h1>
      <button
        aria-label="close-info"
        @click="$emit('close')"
        class="ml-auto rounded-full px-4 py-2"
      >
        <font-awesome-icon class="text-red-500" icon="times" size="2x" />
      </button>
    </div>
    <div
      class="flex flex-col justify-around space-y-2 md:flex-row md:space-x-2 md:space-y-0"
    >
      <album-art
        v-if="state.albumArt != 'default'"
        :url="state.albumArt"
        :altText="state.Info['Title'] + ' cover'"
        class="top-0 h-fit md:sticky md:w-1/2"
      />
      <ul
        :class="{
          'list flex w-full flex-col rounded bg-white dark:bg-gray-700 dark:text-white': true,
          'md:w-1/2': state.albumArt != 'default',
        }"
      >
        <li
          v-for="(value, index) in state.Info"
          :key="index"
          class="p-2 first:rounded-t last:rounded-b"
        >
          {{ index }} : {{ value }}
        </li>
        <li
          v-for="(sticker, index) in state.Stickers"
          :key="index"
          class="p-2 last:rounded-b"
        >
          {{ sticker.Name }} : {{ sticker.Value }}
        </li>
      </ul>
    </div>
  </Dialog>
</template>
<script setup lang="ts">
import { useStore } from "../store";
const state = computed(() => useStore().song);
</script>
