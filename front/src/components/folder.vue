<template>
  <details>
    <summary
      class="flex justify-between md:p-2 px-2 py-4 dark:text-white hover:bg-white/60 dark:hover:bg-gray-800/70 rounded group"
    >
      <span
        class="w-4/6 overflow-x-hidden text-ellipsis"
        @click.self="listFolders()"
      >
        <FontAwesomeIcon :icon="data.Directory ? 'folder' : 'music'" />
        {{ folderName }}
      </span>
      <span class="text-sm space-x-2 invisible group-hover:visible">
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
    <div
      class="flex flex-col border-2 border-primary rounded-lg m-1 dark:border-lightest"
      v-if="opened"
    >
      <Folder
        class="m-1"
        v-for="(folder, index) in folders"
        :key="index"
        :data="folder"
      />
    </div>
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
