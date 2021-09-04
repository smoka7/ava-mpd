<template>
  <div class="flex flex-col">
    <details>
      <summary class="flex text-primary dark:text-white justify-between p-2">
        <span class="w-5/6" @click.self="listFolders(index)">
          <FontAwesomeIcon
            :icon="folder.directory ? 'folder' : 'music'"
          ></FontAwesomeIcon>
          {{ folderName }}
        </span>
        <span class="text-sm space-x-2">
          <button aria-label="add">
            <font-awesome-icon
              @click="addFolder(index)"
              icon="plus"
            ></font-awesome-icon>
          </button>
          <button aria-label="play">
            <font-awesome-icon
              @click="PlayFolder(index)"
              icon="play"
            ></font-awesome-icon>
          </button>
        </span>
      </summary>
      <div class="flex flex-col border-2 border-primary rounded-lg m-1">
        <Folder
          class="ml-2"
          v-for="(folder, index) in folder.folders"
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
      // if its not a folder do nothing
      if (!this.folder.directory) {
        return;
      }
      //if folder list exists delete it
      if (this.folder.folders && this.folder.folders.length) {
        this.folder.folders = [];
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
        this.folder.folders = json;
        return;
      }
      console.log(response.error);
    },
    PlayFolder(index) {
      let data = {
        playlist: this.folder.directory || this.folder.file,
      };
      sendCommand("api/folders", "play", data);
    },
    addFolder(index) {
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
  },
};
</script>
