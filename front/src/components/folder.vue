<template>
  <div class="flex flex-col">
    <details>
      <summary
        class="
          flex
          text-primary
          dark:text-white
          justify-between
          md:p-2
          px-2
          py-4
          hover:bg-blue-200
          dark:hover:text-primary
          rounded
        "
      >
        <span class="w-5/6 overflow-clip" @click.self="listFolders(index)">
          <FontAwesomeIcon
            :icon="folder.directory ? 'folder' : 'music'"
          ></FontAwesomeIcon>
          {{ folderName }}
        </span>
        <span class="text-sm space-x-2">
          <button aria-label="add" class="sidebar-btn">
            <font-awesome-icon
              @click="addFolder(index)"
              icon="plus"
            ></font-awesome-icon>
          </button>
          <button aria-label="play" class="sidebar-btn">
            <font-awesome-icon
              @click="PlayFolder(index)"
              icon="play"
            ></font-awesome-icon>
          </button>
        </span>
      </summary>
      <div
        class="
          flex flex-col
          border-2 border-primary
          rounded-lg
          m-1
          dark:border-lightest
        "
        v-if="hasChild"
      >
        <Folder
          class="m-1"
          v-for="(folder, index) in folder.children"
          :key="index"
          :data="folder"
        >
        </Folder>
      </div>
    </details>
  </div>
</template>
<script>
import { sendCommand } from "../helpers.js";
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
      folder: this.data,
    };
  },
  methods: {
    async listFolders() {
      // if its not a folder do nothing and if folder list exists delete it
      if (!this.folder.directory) return;
      if (this.hasChild) {
        this.folder.children = [];
        return;
      }
      let request = {
        command: "list",
        data: {
          playlist: this.folder.directory,
        },
      };
      let response = await fetch("/api/folders", {
        method: "POST",
        headers: {
          "Content-Type": "application/json;charset=utf-8",
        },
        body: JSON.stringify(request),
      });
      if (response.ok) {
        let json = await response.json();
        this.folder.children = json;
        return;
      }
      console.log(response.error);
    },
    PlayFolder() {
      let data = {
        playlist: this.folder.directory || this.folder.file,
      };
      sendCommand("api/folders", "play", data);
    },
    addFolder() {
      let data = {
        playlist: this.folder.directory || this.folder.file,
      };
      sendCommand("api/folders", "add", data);
    },
  },
  computed: {
    folderName() {
      let name = this.folder.directory || this.folder.file;
      let index = name.lastIndexOf("/");
      if (index > -1) return name.substring(index + 1);
      return name;
    },
    hasChild() {
      if (this.folder.children && this.folder.children.length > 0) return true;
      return false;
    },
  },
};
</script>
