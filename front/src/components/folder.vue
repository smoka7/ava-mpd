<template>
  <details :open="folders.data != null && folders.data.length > 0">
    <summary
      class="group flex justify-between rounded px-2 py-4 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70 md:p-2"
    >
      <span
        class="w-4/6 overflow-x-hidden text-ellipsis"
        @click.self="listFolders()"
      >
        <FontAwesomeIcon :icon="data.Directory ? 'folder' : 'music'" />
        {{ folderName() }}
      </span>
      <span class="invisible space-x-2 text-sm group-hover:visible">
        <sidebar-button label="add" icon="plus" @click="AddOpen = !AddOpen" />
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
        v-if="folders.data.length > 0"
      >
        <Folder
          class="m-1"
          v-for="(folder, index) in folders.data"
          :key="index"
          :data="folder"
        />
      </div>
    </transition>
    <AddSongs
      :open="AddOpen"
      @close="AddOpen = false"
      @add="FolderCommand('add', $event)"
    />
  </details>
</template>
<script setup>
import { reactive, ref } from "vue";
import { sendCommand } from "../helpers.js";
import endpoints from "../endpoints.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import SidebarButton from "./sidebarButton.vue";
import AddSongs from "./addSongs.vue";

const props = defineProps(["data"]);
const folders = reactive({ data: [] });
const AddOpen = ref(false);

async function listFolders() {
  // if its not a folder do nothing and if folder list exists delete it
  if (!props.data.Directory) return;
  if (folders.data.length > 0) {
    folders.data = [];
    return;
  }
  const response = await sendCommand(endpoints.folders, "list", {
    playlist: props.data.Directory,
  });
  folders.data = [...response.Folders, ...response.Files];
}

function FolderCommand(command, position) {
  const data = {
    playlist: props.data.Directory || props.data.File,
    song: position,
  };
  sendCommand(endpoints.folders, command, data);
}

function folderName() {
  const name = props.data.Directory || props.data.File;
  const index = name.lastIndexOf("/");
  if (index > -1) return name.substring(index + 1);
  return name;
}
</script>
