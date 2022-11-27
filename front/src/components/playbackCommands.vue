<template>
  <div class="flex w-full flex-col justify-end">
    <div class="flex h-10 justify-around space-x-2">
      <save-queue />
      <button
        v-for="command in playbackCommands"
        :key="command.icon"
        :aria-label="command.command"
        @click="Command(command.command)"
        :class="{
          'playback-active': command.status == '1',
          'playback-btn tooltip': true,
        }"
      >
        <span class="tooltiptext">{{ command.command }}</span>
        <font-awesome-icon :icon="command.icon" size="lg" />
      </button>
      <button
        aria-label="setting"
        @click="toggleMediaController"
        class="playback-btn tooltip md:hidden"
      >
        <span class="tooltiptext">queue...</span>
        <font-awesome-icon icon="list-ul" size="lg" />
      </button>
    </div>
  </div>
  <div class="flex h-28 w-full justify-around text-lightest md:h-24">
    <button
      aria-label="previous-song"
      @click="Command('previous')"
      class="duration-200 hover:scale-125 hover:text-lighter"
    >
      <font-awesome-icon icon="step-backward" size="2x" />
    </button>
    <like-song :pLiked="liked" :file="currentSong.File" />
    <button
      aria-label="toggle-playback"
      @click="Command('toggle')"
      :class="[
        status.state === 'pause' ? 'bg-secondary' : 'bg-accent',
        'flex aspect-square items-center justify-center rounded-full text-white duration-200 hover:scale-105 hover:text-lightest',
      ]"
    >
      <font-awesome-icon
        :icon="['fas', status.state === 'play' ? 'pause' : 'play']"
        size="2x"
      />
    </button>
    <button
      aria-label="stop-song"
      @click="Command('stop')"
      class="duration-200 hover:scale-125 hover:text-lighter"
    >
      <font-awesome-icon icon="stop" size="2x" />
    </button>
    <button
      aria-label="next-song"
      @click="Command('next')"
      class="duration-200 hover:scale-125 hover:text-lighter"
    >
      <font-awesome-icon icon="step-forward" size="2x" />
    </button>
  </div>
</template>
<script setup>
import { toggleMediaController } from "../helpers.js";
const store = useStore();
const currentSong = computed(() => shallowReactive(store.state.currentSong));

const liked = computed(() => store.state.currentSong.Liked);

const status = computed(() => shallowReactive(store.state.status));

const playbackCommands = computed(() => [
  {
    command: "consume",
    icon: "minus-square",
    status: store.state.status.consume,
  },
  {
    command: "single",
    icon: "dice-one",
    status: store.state.status.single,
  },
  {
    command: "repeat",
    icon: "retweet",
    status: store.state.status.repeat,
  },
  {
    command: "random",
    icon: "random",
    status: store.state.status.random,
  },
  {
    command: "clear",
    icon: "eraser",
    status: 0,
  },
]);

function Command(command) {
  store.dispatch("sendPlaybackCommand", command);
}
</script>
<style lang="postcss">
.playback-btn {
  @apply scale-125 transform bg-transparent p-2 text-lightest duration-200 hover:scale-[135%] hover:text-lighter;
}
.playback-active {
  @apply text-accent !important;
}
</style>
