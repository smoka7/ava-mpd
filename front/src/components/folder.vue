<template>
  <details :open="IsOpen">
    <summary
      class="group flex justify-between rounded px-2 py-4 duration-200 hover:bg-white/60 dark:text-white dark:hover:bg-gray-800/70 md:p-2"
      @click="listFolders()"
    >
      <span class="w-4/6 overflow-x-hidden text-ellipsis">
        <FontAwesomeIcon
          :icon="directory?.Directory !== undefined ? 'folder' : 'music'"
          class="menuIcon"
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
        v-if="IsOpen"
      >
        <Folder
          class="m-1"
          v-for="(directory, index) in folders.Directories"
          :key="index"
          :directory="directory"
        />
        <Folder
          class="m-1"
          v-for="(file, index) in folders.Files"
          :key="index"
          :file="file"
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
<script setup lang="ts">
import { sendCommand } from "../helpers";
import {
  useStore,
  type Directory,
  type File,
  type FoldersResponse,
} from "../store";
import endpoints from "../endpoints";

const props = defineProps<{
  directory?: Directory;
  file?: File;
}>();

const folders = reactive({
  Files: [] as Array<File>,
  Directories: [] as Array<Directory>,
});

const AddOpen = ref(false);
const IsOpen = computed(
  () => folders.Directories.length > 0 || folders.Files.length > 0
);

const store = useStore();
const storedPlaylist = computed(() => store.storedPlaylist);

async function listFolders() {
  if (props.file) return;
  if (folders.Directories.length > 0 || folders.Files.length > 0) {
    folders.Directories = [];
    folders.Files = [];
    return;
  }
  const response = await sendCommand<FoldersResponse>(
    endpoints.folders,
    "list",
    {
      File: props.directory?.Directory,
    }
  );
  if (response) {
    folders.Directories = response.Directories;
    folders.Files = response.Files;
  }
}

function FolderCommand(command: string, position?: string) {
  const data = {
    File: props.directory?.Directory || props.file?.File,
    Pos: position,
  };
  sendCommand(endpoints.folders, command, data);
}

function folderName() {
  const name = props.directory?.Directory || props.file?.File || "Empty";
  const index = name.lastIndexOf("/");
  if (index > -1) return name.substring(index + 1);
  return name;
}
function openAdd() {
  AddOpen.value = !AddOpen.value;
  if (storedPlaylist.value.length) return;
  store.getStoredPlaylist();
}
</script>
