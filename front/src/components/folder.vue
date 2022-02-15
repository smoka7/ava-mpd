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
        <sidebar-button label="add" icon="plus" @click="FolderCommand('add')" />
        <sidebar-button
          label="play"
          icon="play"
          @click="FolderCommand('play')"
        />
      </span>
    </summary>
    <transition name="fade">
      <div
        class="m-1 flex flex-col border-l-2 border-primary dark:border-lightest"
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
import { sendCommand} from "../helpers.js";
import endpoints from "../endpoints.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import SidebarButton from "./sidebarButton.vue";
export default {
  components: {
    FontAwesomeIcon,
    SidebarButton,
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
      const response = await sendCommand(endpoints.folders, "list", {
        playlist: this.data.Directory,
      });
      this.folders = [...response.Folders, ...response.Files];
    },
    FolderCommand(command) {
      const data = {
        playlist: this.data.Directory || this.data.File,
      };
      sendCommand(endpoints.folders, command, data);
    },
  },
  computed: {
    folderName() {
      const name = this.data.Directory || this.data.File;
      const index = name.lastIndexOf("/");
      if (index > -1) return name.substring(index + 1);
      return name;
    },
  },
};
</script>
