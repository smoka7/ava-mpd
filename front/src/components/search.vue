<template>
  <div class="flex flex-col">
    <form class="p-2 w-full flex" @submit.prevent>
      <select
        class="
          p-2
          dark:bg-lightest
          bg-primary
          dark:text-primary
          text-foreground
          rounded-l-lg
        "
        name="tag"
        v-model.lazy="tag"
      >
        <option v-for="tag in searchTags" :key="tag" :value="tag">
          {{ tag }}
        </option>
      </select>
      <input
        class="
          p-2
          w-full
          border-primary border-2
          dark:text-primary dark:placeholder-gray-900 dark:border-lightest
          rounded-r-lg
        "
        type="text"
        name="term"
        @input="search"
        v-model.trim="term"
        placeholder="write to search"
      />
    </form>
    <div class="flex flex-col divide-y dark:divide-gray-800">
      <Folder v-for="(file, index) in files" :key="index" :data="file">
      </Folder>
    </div>
  </div>
</template>
<script>
import Folder from "./folder.vue";
export default {
  components: {
    Folder,
  },
  data() {
    return {
      searchTags: ["file", "artist", "album", "genre", "date", "title"],
      tag: "file",
      term: "",
      files: [],
    };
  },
  methods: {
    async search() {
      if (this.term === "") {
        this.files = [];
        return;
      }
      let request = {
        terms: [this.tag, this.term],
      };
      let response = await fetch("/api/search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json;charset=utf-8",
        },
        body: JSON.stringify(request),
      });
      if (response.ok) {
        let json = await response.json();
        if (this.files.length) this.files = [];
        this.files = json;
        return;
      }
    },
  },
};
</script>
