<template>
  <div
    class="absolute top-0 z-50 w-full max-w-full max-h-fit flex justify-between items-center py-1 px-4"
  >
    <Menu>
      <MenuButton class="ml-auto p-2">
        <FontAwesomeIcon icon="ellipsis-h" />
      </MenuButton>
      <transition
        enter-active-class="transition duration-100 ease-out"
        enter-from-class="transform scale-95 opacity-0"
        enter-to-class="transform scale-100 opacity-100"
        leave-active-class="transition duration-100 ease-out"
        leave-from-class="transform scale-100 opacity-100"
        leave-to-class="transform scale-95 opacity-0"
      >
        <MenuItems
          class="fixed md:left-32 md:right-auto right-12 top-20 z-50 flex flex-col text-left rounded shadow-md cursor-default bg-white focus:outline-none focus:ring-2 focus:ring-primary dark:focus:ring-lightest dark:bg-gray-700"
        >
          <MenuItem
            v-slot="{ active }"
            v-for="action in actions"
            :key="action.name"
            @click="$emit(action.method)"
          >
            <div
              :class="{
                'bg-blue-100 text-primary': active,
                'p-2 first:rounded-t last:rounded-b flex items-center space-x-2': true,
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
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";
export default {
  emits: [
    "delete",
    "clearSelection",
    "clear",
    "rename",
    "removeDuplicate",
    "removeInvalid",
  ],
  components: {
    FontAwesomeIcon,
    Menu,
    MenuButton,
    MenuItems,
    MenuItem,
  },
  data() {
    return {
      actions: [
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
      ],
    };
  },
};
</script>
