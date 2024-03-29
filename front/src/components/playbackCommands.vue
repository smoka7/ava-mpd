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
          'playback-active': command.status,
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
  <div class="flex h-28 w-full justify-around text-white md:h-24">
    <button
      aria-label="previous-song"
      @click="Command('previous')"
      class="duration-200 hover:scale-125 hover:text-lighter"
    >
      <font-awesome-icon icon="step-backward" size="2x" />
    </button>
    <like-song :pLiked="liked" :file="currentSongFile" icon-size="2x" />
    <button
      aria-label="toggle-playback"
      @click="Command('toggle')"
      :class="[
        status.state === 'pause' ? 'border-lightest' : 'border-accent',
        status.state === 'pause' ? 'text-lightest' : 'text-accent',
        'flex h-24 w-24 items-center justify-center rounded-full border-4 duration-200 hover:scale-105 hover:border-8',
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
<script setup lang="ts">
import { toggleMediaController } from "../helpers";
import { useStore } from "../store";

const store = useStore();
const currentSongFile = computed(() => store.currentSong.File);

const liked = computed(() => store.currentSong.Liked);

const status = computed(() => shallowReactive(store.status));

type CommandButton = {
  command: string;
  icon: string;
  status: boolean;
};

const playbackCommands = computed<Array<CommandButton>>(() => [
  {
    command: "consume",
    icon: "minus-square",
    status: store.status.consume,
  },
  {
    command: "single",
    icon: "dice-one",
    status: store.status.single,
  },
  {
    command: "repeat",
    icon: "retweet",
    status: store.status.repeat,
  },
  {
    command: "random",
    icon: "random",
    status: store.status.random,
  },
  {
    command: "clear",
    icon: "eraser",
    status: false,
  },
]);

function Command(command: string) {
  store.sendPlaybackCommand(command);
}
</script>
<style lang="postcss">
.playback-btn {
  @apply scale-125 transform bg-transparent p-2 text-white duration-200 hover:scale-[135%] hover:text-lighter;
}
.playback-active {
  @apply text-accent !important;
}
</style>
