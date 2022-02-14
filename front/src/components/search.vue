<template>
  <div class="flex flex-col">
    <form class="flex w-full p-2" @submit.prevent="search">
      <Listbox v-model="tag">
        <ListboxButton
          class="dark:bg-lightest bg-primary dark:text-primary flex items-center rounded-l-lg p-2 text-white"
        >
          {{ tag }}
          <FontAwesomeIcon icon="angle-right" class="ml-2 rotate-90" />
        </ListboxButton>
        <transition name="fade">
          <ListboxOptions
            class="fixed mt-14 cursor-pointer rounded bg-white/60 text-left backdrop-blur-3xl dark:bg-gray-700/60"
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
                  'dark:hover:text-primary p-2 hover:bg-blue-200': true,
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
        class="border-primary focus:border-lightest w-full rounded-r-lg border-2 border-l-0 bg-white/70 p-2 outline-none backdrop-blur-3xl dark:bg-gray-700/60 dark:text-white"
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
      class="m-1 space-y-2 rounded bg-white/50 dark:bg-gray-700/50"
    >
      <summary
        class="group flex cursor-pointer items-center justify-between overflow-x-hidden text-ellipsis rounded p-2 dark:text-white"
      >
        <span>
          {{ index }}
        </span>
        <span class="invisible space-x-2 text-sm group-hover:visible">
          <sidebar-button
            label="add"
            icon="plus"
            @click="command('searchadd', index)"
          />
          <sidebar-button
            label="play"
            icon="play"
            @click="command('searchplay', index)"
          />
        </span>
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
import SidebarButton from "./sidebarButton.vue";
import endpoints from "../endpoints.js";
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
    SidebarButton,
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
      const response = await this.command("search", this.term);
      if (response.ok) {
        const json = await response.json();
        if (this.files.length) this.files = [];
        this.files = json.Songs;
        return;
      }
    },
    async command(cm, title) {
      const request = {
        terms: [this.tag, title],
        command: cm,
      };
      return await fetch(endpoints.search, {
        method: "POST",
        headers: {
          "Content-Type": "application/json;charset=utf-8",
        },
        body: JSON.stringify(request),
      });
    },
  },
};
</script>
