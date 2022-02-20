<template>
  <div
    class="absolute top-0 z-50 flex max-h-fit w-full max-w-full items-center justify-between py-1 px-4"
  >
    <Menu>
      <MenuButton class="ml-auto p-2">
        <FontAwesomeIcon icon="ellipsis-h" />
      </MenuButton>
      <transition name="fade">
        <MenuItems
          class="absolute right-12 top-10 z-50 flex cursor-default flex-col rounded bg-white/60 text-left shadow-md backdrop-blur-3xl focus:outline-none dark:bg-gray-700/60 md:left-32 md:right-2"
        >
          <MenuItem
            v-slot="{ active }"
            v-for="action in actions"
            :key="action.name"
            @click="$emit(action.method)"
          >
            <div
              :class="{
                'text-primary bg-blue-100': active,
                'flex items-center space-x-2 p-2 first:rounded-t last:rounded-b': true,
              }"
            >
              <FontAwesomeIcon :icon="['fas', action.icon]" />
              <span>{{ action.name }}</span>
            </div>
          </MenuItem>
        </MenuItems>
      </transition>
    </Menu>
  </div>
</template>
<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";
const emit = defineEmits([
  "addafter",
  "delete",
  "clearSelection",
  "clear",
  "rename",
  "removeDuplicate",
  "removeInvalid",
]);
const actions = [
  {
    name: "Add after current",
    icon: "plus",
    method: "addafter",
  },
  {
    name: "Delete Playlists",
    icon: "trash",
    method: "delete",
  },
  {
    name: "Clear Playlists",
    icon: "eraser",
    method: "clear",
  },
  {
    name: "Rename Playlists",
    icon: "edit",
    method: "rename",
  },
  {
    name: "Remove Duplicate songs",
    icon: "eraser",
    method: "removeDuplicate",
  },
  {
    name: "Remove Invalid songs",
    icon: "eraser",
    method: "removeInvalid",
  },
  {
    name: "Clear Selection",
    icon: "check-circle",
    method: "clearSelection",
  },
];
</script>
