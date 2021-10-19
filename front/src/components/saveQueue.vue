<template>
  <button
    aria-label="save-queue"
    @click="toggleDialog(true)"
    class="
      p-2
      px-3
      text-foreground
      rounded-full
      md:p-0
      tooltip
      md:hover:text-blue-200
    "
  >
    <font-awesome-icon icon="save" size="lg"></font-awesome-icon>
    <span class="tooltiptext">save current queue</span>
  </button>
  <Dialog
    :open="isOpen"
    @close="toggleDialog"
    class="fixed inset-0 z-10 overflow-y-auto rounded"
  >
    <DialogOverlay class="fixed inset-0 bg-black opacity-30" />
    <div
      class="
        relative
        flex flex-col
        max-w-xl
        p-4
        top-1/3
        mx-auto
        bg-white
        rounded
      "
    >
      <DialogTitle class="my-2 text-lg">Enter the name of playlist</DialogTitle>
      <input
        type="text"
        v-model="playListName"
        class="border-2 border-primary rounded p-2 text-black w-full"
        placeholder="playlist name"
      />
      <div class="flex w-full justify-end mt-2 space-x-4">
        <button
          @click="saveQueue"
          class="
            p-2
            border-2 border-green-500
            text-green-500 text-lg
            hover:bg-green-500 hover:text-white
            rounded
          "
        >
          save
        </button>
        <button
          @click="toggleDialog(false)"
          class="
            p-2
            border-2 border-red-500
            text-red-500 text-lg
            hover:bg-red-500 hover:text-white
            rounded
          "
        >
          cancel
        </button>
      </div>
    </div>
  </Dialog>
</template>
<script>
import { Dialog, DialogOverlay, DialogTitle } from "@headlessui/vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { sendCommand } from "../helpers";
export default {
  components: {
    Dialog,
    DialogOverlay,
    DialogTitle,
    FontAwesomeIcon,
  },
  data() {
    return {
      playListName: "",
      isOpen: false,
    };
  },
  methods: {
    toggleDialog(open) {
      this.isOpen = open;
    },
    saveQueue() {
      sendCommand("/api/queue", "save", { playlist: this.playListName });
      this.toggleDialog();
    },
  },
};
</script>
