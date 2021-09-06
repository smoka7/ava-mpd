<template>
  <div
    class="bg-transparent absolute top-0 bottom-0 left-0 right-0"
    @click.self="$emit('close')"
  >
    <div
      class="
        flex flex-col
        absolute
        top-0
        left-0
        bottom-0
        right-0
        md:w-3/4 md:left-1/4 md:top-1
        app-height
        z-20
        p-8
        space-y-2
        border-2 border-primary
        rounded
        overflow-auto
        bg-lightest
      "
    >
      <div class="flex flex-row align-baseline justify-between">
        <h1 class="text-primary text-4xl">song Info</h1>
        <button
          aria-label="close-info"
          @click="$emit('close')"
          class="px-4 py-2 rounded-full"
        >
          <font-awesome-icon
            class="text-red-500"
            icon="times"
            size="2x"
          ></font-awesome-icon>
        </button>
      </div>
      <div
        class="
          flex flex-col
          rounded
          p-4
          bg-white
          dark:text-white dark:bg-gray-700
        "
      >
        <p v-for="(value, index) in info" :key="index">
          {{ index }} : {{ value }}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
export default {
  components: { FontAwesomeIcon },
  props: ["song"],
  data() {
    return {
      info: {},
    };
  },
  methods: {},
  async mounted() {
    let request = {
      command: "info",
      data: { song: this.song },
    };
    let response = await fetch("/api/song", {
      method: "POST",
      headers: {
        "Content-Type": "application/json;charset=utf-8",
      },
      body: JSON.stringify(request),
    });
    if (response.ok) {
      this.info = await response.json();
    }
  },
};
</script>
