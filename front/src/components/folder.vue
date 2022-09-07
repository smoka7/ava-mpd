<template>
  <details :open="folders.data != null && folders.data.length > 0">
    <summary
      class="group flex justify-between rounded px-2 py-4 duration-200 hover:bg-white/60 hover:py-[1.1rem] dark:text-white dark:hover:bg-gray-800/70 md:p-2"
      @click="listFolders()"
    >
      <span class="w-4/6 overflow-x-hidden text-ellipsis">
        <FontAwesomeIcon
          :icon="data.Directory ? 'folder' : 'music'"
          class="group-hover:text-secondary"
        />
        {{ folderName() }}
      </span>
      <span class="hide-unfocused space-x-2 text-sm">
        <sidebar-button label="add" icon="plus" @click.stop="openAdd" />
        <sidebar-button
          label="play"
          icon="play"
          @click.stop="FolderCommand('play')"
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
      :storedPlaylist="storedPlaylist"
      @close="AddOpen = false"
      @add="FolderCommand('add', $event)"
      @addPlaylist="FolderCommand('addToPlaylist', $event)"
    />
  </details>
</template>
<script setup>
import { sendCommand } from "../helpers.js";
import endpoints from "../endpoints.js";

const props = defineProps(["data"]);
const folders = reactive({ data: [] });
const AddOpen = ref(false);

const store = useStore();
const storedPlaylist = computed(() => store.state.storedPlaylist);

async function listFolders() {
  // if its not a folder do nothing and if folder list exists delete it
  if (!props.data.Directory) return;
  if (folders.data.length > 0) {
    folders.data = [];
    return;
  }
  const response = await sendCommand(endpoints.folders, "list", {
    File: props.data.Directory,
  });
  folders.data = [...response.Folders, ...response.Files];
}

function FolderCommand(command, position) {
  const data = {
    File: props.data.Directory || props.data.File,
    Pos: position,
  };
  sendCommand(endpoints.folders, command, data);
}

function folderName() {
  const name = props.data.Directory || props.data.File;
  const index = name.lastIndexOf("/");
  if (index > -1) return name.substring(index + 1);
  return name;
}
function openAdd() {
  AddOpen.value = !AddOpen.value;
  if (storedPlaylist.length) return;
  store.dispatch("getStoredPlaylist");
}
</script>
