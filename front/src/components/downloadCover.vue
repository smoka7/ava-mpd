<template>
  <div>
    <SwitchGroup as="div" class="flex items-center p-2">
      <SwitchLabel class="mr-4 text-lg md:mr-2"
        >Download missing cover from musicBrainz</SwitchLabel
      >
      <Switch
        v-model="download"
        class="relative inline-flex h-10 w-16 items-center rounded bg-lightest shadow-inner md:h-7 md:w-14"
      >
        <span class="sr-only">Download missing cover from musicBrainz</span>
        <span
          :class="[
            'inline-block h-8 w-8 transform rounded bg-primary duration-300 md:h-5 md:w-6',
            download ? ' translate-x-7' : 'translate-x-1',
          ]"
        />
      </Switch>
    </SwitchGroup>
  </div>
</template>
<script setup>
import { Switch, SwitchGroup, SwitchLabel } from "@headlessui/vue";
import { computed } from "vue";
import { useStore } from "vuex";
import { sendCommand } from "../helpers";
import endpoints from "../endpoints";

const store = useStore();
const download = computed({
  get() {
    return store.state.settings.DownloadCoverArt;
  },
  set() {
    sendCommand(endpoints.setting, "download");
    store.dispatch("getSettings");
  },
});
</script>
