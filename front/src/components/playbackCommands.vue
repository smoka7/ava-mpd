<template>
  <div class="flex w-full flex-col justify-end">
    <div class="flex h-10 justify-around space-x-2">
      <save-queue />
      <button
        v-for="command in playbackCommands"
        :key="command.icon"
        :aria-label="command.command"
        @click="Command(command.command)"
        :class="[
          command.status == '1' ? btnClass.active : btnClass.normal,
          btnClass.base,
        ]"
      >
        <span class="tooltiptext">{{ command.command }}</span>
        <font-awesome-icon :icon="command.icon" size="lg" />
      </button>
      <button
        aria-label="setting"
        @click="toggleMediaController"
        :class="[btnClass.normal, btnClass.base, 'md:hidden']"
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
      class="hover:scale-125 hover:text-accent"
    >
      <font-awesome-icon icon="step-backward" size="2x" />
    </button>
    <like-song :pLiked="liked" :file="currentSong.File" />
    <button
      aria-label="toggle-playback"
      @click="Command('toggle')"
      :class="[
        status.state === 'pause' ? 'bg-green-400' : 'bg-accent',
          'flex items-center justify-center rounded-full px-10 text-white hover:text-primary md:px-8 hover:scale-105',
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
      class="hover:scale-125 hover:text-accent"
    >
      <font-awesome-icon icon="stop" size="2x" />
    </button>
    <button
      aria-label="next-song"
      @click="Command('next')"
      class="hover:scale-125 hover:text-accent"
    >
      <font-awesome-icon icon="step-forward" size="2x" />
    </button>
  </div>
</template>
<script setup>
const store = useStore();
const currentSong = computed(() =>
  shallowReactive(store.state.currentSong),
);

const liked = computed(() => store.state.currentSong.Liked);

const status = computed(() => shallowReactive(store.state.status));

const btnClass = {
  active: "text-green-500 bg-transparent",
  normal: "bg-transparent hover:text-accent",
  base: "p-2 tooltip transform scale-125 hover:scale-[135%]",
};

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
