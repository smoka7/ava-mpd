<template>
  <button
    @click="open = true"
    class="rounded border-2 border-lightest p-2 hover:border-lighter hover:bg-lighter"
  >
    show mappings list
  </button>
  <Dialog
    :open="open"
    @close="open = false"
    class="fixed inset-0 z-10 overflow-y-auto"
  >
    <DialogOverlay
      v-show="open"
      class="fixed inset-0 z-0 bg-black/30 backdrop-blur-sm"
    />

    <transition name="fade">
      <DialogPanel>
        <div
          class="absolute flex h-screen w-full flex-col overflow-y-auto bg-white p-1 dark:bg-gray-800 dark:text-white md:top-auto md:left-1/4 md:w-1/2 md:p-4"
        >
          <table class="table-auto rounded">
            <thead class="bg-lightest text-left text-primary">
              <tr>
                <th class="px-2">key</th>
                <th class="px-2">action</th>
              </tr>
            </thead>
            <tbody class="">
              <tr v-for="(map, index) in playbackMappings" :key="index">
                <td class="px-2">
                  {{ index }}
                </td>
                <td class="px-2">
                  {{ map.name }}
                </td>
              </tr>
              <tr v-for="(map, index) in otherMappings" :key="index">
                <td class="px-2">
                  {{ index }}
                </td>
                <td class="px-2">
                  {{ map.name }}
                </td>
              </tr>
              <tr v-for="(map, index) in tabMappings" :key="index">
                <td class="px-2">
                  {{ index }}
                </td>
                <td class="px-2">
                  {{ map.name }}
                </td>
              </tr>
            </tbody>
          </table>
          <button
            @click="open = false"
            class="self-end rounded border-2 border-red-500 p-2 text-red-500 hover:bg-red-500 hover:text-white"
          >
            Got it!
          </button>
        </div>
      </DialogPanel>
    </transition>
  </Dialog>
</template>
<script setup>
import { otherMappings, tabMappings, playbackMappings } from "../keymap";
const open = ref(false);
</script>
<style lang="postcss">
tbody tr {
  @apply even:bg-gray-200 even:dark:bg-gray-700;
}
</style>
