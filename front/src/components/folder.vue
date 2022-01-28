<template>
  <details>
    <summary
      class="group flex justify-between rounded px-2 py-4 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70 md:p-2"
    >
      <span
        class="w-4/6 overflow-x-hidden text-ellipsis"
        @click.self="listFolders()"
      >
        <FontAwesomeIcon :icon="data.Directory ? 'folder' : 'music'" />
        {{ folderName }}
      </span>
      <span class="invisible space-x-2 text-sm group-hover:visible">
        <button
          aria-label="add"
          class="sidebar-btn"
          @click="FolderCommand('add')"
        >
          <font-awesome-icon icon="plus" />
        </button>
        <button
          aria-label="play"
          class="sidebar-btn"
          @click="FolderCommand('play')"
        >
          <font-awesome-icon icon="play" />
        </button>
      </span>
    </summary>
    <transition name="fade">
      <div
        class="border-primary dark:border-lightest m-1 flex flex-col border-l-2"
        v-if="opened"
      >
        <Folder
          class="m-1"
          v-for="(folder, index) in folders"
          :key="index"
          :data="folder"
        />
      </div>
    </transition>
  </details>
</template>
<script>
import { shallowReactive } from "vue";
import { sendCommand, getFolders } from "../helpers.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
export default {
  components: {
    FontAwesomeIcon,
  },
  props: {
    data: Object,
  },
  data() {
    return {
      folders: shallowReactive([]),
      opened: false,
    };
  },
  methods: {
    async listFolders() {
      // if its not a folder do nothing and if folder list exists delete it
      if (!this.data.Directory) return;
      if (this.opened) {
        this.folders = [];
        this.opened = false;
        return;
      }
      this.opened = true;
      this.folders = await getFolders(this.data.Directory);
    },
    FolderCommand(command) {
      let data = {
        playlist: this.data.Directory || this.data.File,
      };
      sendCommand("api/folders", command, data);
    },
  },
  computed: {
    folderName() {
      let name = this.data.Directory || this.data.File;
      let index = name.lastIndexOf("/");
      if (index > -1) return name.substring(index + 1);
      return name;
    },
  },
};
</script>
