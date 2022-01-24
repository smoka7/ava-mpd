<template>
  <div class="flex flex-col">
    <form class="p-2 w-full flex" @submit.prevent="search">
      <Listbox v-model="tag">
        <ListboxButton
          class="flex items-center p-2 dark:bg-lightest bg-primary dark:text-primary text-white rounded-l-lg"
        >
          {{ tag }}
          <FontAwesomeIcon icon="angle-right" class="ml-2 rotate-90" />
        </ListboxButton>
        <transition name="fade">
          <ListboxOptions
            class="fixed mt-14 text-left rounded cursor-pointer bg-white/60 dark:bg-gray-700/60 backdrop-blur-3xl"
          >
            <ListboxOption
              v-for="tag in searchTags"
              as="template"
              :key="tag"
              :value="tag"
              v-slot="{ selected }"
            >
              <li
                :class="{
                  'p-2 hover:bg-blue-200 dark:hover:text-primary': true,
                  'bg-lightest dark:text-primary': selected,
                }"
              >
                {{ tag }}
                <FontAwesomeIcon class="ml-2" v-show="selected" icon="check" />
              </li>
            </ListboxOption>
          </ListboxOptions>
        </transition>
      </Listbox>
      <input
        class="p-2 w-full border-primary border-2 border-l-0 outline-none focus:border-lightest dark:text-white rounded-r-lg bg-white/70 dark:bg-gray-700/60 backdrop-blur-3xl"
        type="text"
        name="term"
        v-model.trim="term"
        placeholder="write to search"
      />
    </form>
    <details
      v-for="(term, index) in files"
      :key="index"
      :data="term"
      class="m-1 dark:bg-gray-700/50 bg-white/50 space-y-2 rounded"
    >
      <summary
        class="flex p-2 rounded bg-lightest items-center cursor-pointer dark:text-primary overflow-x-hidden text-ellipsis"
      >
        {{ index }}
      </summary>
      <div>
        <Folder
          v-for="(file, index) in term"
          :key="index"
          :data="file"
          class="last:rounde-b odd:bg-gray-600/10 dark:odd:bg-gray-800/50 dark:hover:odd:bg-gray-800/70"
        />
      </div>
    </details>
  </div>
</template>
<script>
import Folder from "./folder.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
} from "@headlessui/vue";
export default {
  components: {
    Folder,
    Listbox,
    ListboxButton,
    ListboxOptions,
    ListboxOption,
    FontAwesomeIcon,
  },
  data() {
    return {
      searchTags: ["file", "Artist", "Album", "Genre", "Date", "Title"],
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
