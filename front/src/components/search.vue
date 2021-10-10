<template>
  <div class="flex flex-col">
    <form class="p-2 w-full flex" @submit.prevent>
      <Listbox v-model="tag">
        <ListboxButton
          class="
            p-2
            dark:bg-lightest
            bg-primary
            dark:text-primary
            text-foreground
            rounded-l-lg
          "
        >
          <span class="flex items-center">
            {{ tag }}
            <FontAwesomeIcon icon="angle-down" class="ml-2"></FontAwesomeIcon>
          </span>
        </ListboxButton>
        <transition
          enter-active-class="transition duration-100 ease-out"
          enter-from-class="transform scale-95 opacity-0"
          enter-to-class="transform scale-100 opacity-100"
          leave-active-class="transition duration-75 ease-out"
          leave-from-class="transform scale-100 opacity-100"
          leave-to-class="transform scale-95 opacity-0"
        >
          <ListboxOptions
            class="
              fixed
              mt-14
              text-left
              rounded
              shadow-md
              cursor-default
              bg-white
              focus:outline-none focus:ring-2 focus:ring-primary
              dark:bg-gray-700 dark:text-white
            "
          >
            <ListboxOption
              v-for="tag in searchTags"
              :key="tag"
              :value="tag"
              v-slot="{ selected }"
            >
              <li
                :class="{
                  'p-2 focus-within:bg-blue-200 hover:bg-blue-200': !selected,
                  'p-2 bg-lightest': selected,
                }"
              >
                {{ tag }}
                <FontAwesomeIcon
                  v-show="selected"
                  icon="check"
                  class="ml-2 transform-gpu duration-200"
                ></FontAwesomeIcon>
              </li>
            </ListboxOption>
          </ListboxOptions>
        </transition>
      </Listbox>
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
