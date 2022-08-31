<template>
  <div
    class="fixed inset-0 flex items-center justify-center p-4 text-center text-7xl text-primary underline decoration-accent md:text-9xl"
  >
    Could't connect to server !!! <br />
    😱 😱 😱
  </div>
  <button
    @click="reConnect"
    class="group fixed bottom-10 right-10 z-20 flex h-20 w-fit items-center justify-center space-x-2 rounded bg-primary p-2 text-4xl text-white hover:bg-accent hover:text-primary"
  >
    <FontAwesomeIcon icon="fa-refresh" :class="btnClass" size="lg" />
    <span> Reconnect </span>
  </button>
</template>
<script setup>
const btnClass = ref(["btn-normal"]);

const store = useStore();

function reConnect() {
  btnClass.value = ["btn-normal", "btn-scale", "btn-right"];

  setTimeout(() => {
    btnClass.value = ["btn-normal", "btn-scale", "btn-left"];
  }, 150);

  setTimeout(() => {
    btnClass.value = "btn-normal btn-animate";
  }, 300);

  store.dispatch("getCurrentSong");
  store.dispatch("getQueue");
  store.dispatch("startCounter");
  store.dispatch("connectToSocket");
}
</script>
<style lang="postcss" scoped>
.btn-normal {
  @apply transform text-accent transition-all duration-200 group-hover:scale-105 group-hover:text-white;
}
.btn-left {
  @apply -rotate-45;
}
.btn-right {
  @apply rotate-45;
}
.btn-scale {
  @apply scale-105;
}
</style>
