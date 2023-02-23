<template>
  <div>
    <SwitchGroup as="div" class="flex items-center justify-between p-2">
      <SwitchLabel class="mr-4 w-4/6 text-lg md:mr-2"
        >Download missing covers from
        <a
          href="https://musicbrainz.org"
          target="_blank"
          class="underline decoration-secondary"
          >musicBrainz
        </a>
      </SwitchLabel>
      <Switch
        v-model="download"
        class="relative ml-2 inline-flex h-10 w-16 items-center rounded bg-lightest shadow-inner md:h-7 md:w-14"
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
    <span class="text-base"
      >Downloaded covers will be saved in users cache directory.</span
    >
  </div>
</template>
<script setup lang="ts">
const store = useStore();
const download = computed({
  get() {
    return store.state.settings.DownloadCoverArt;
  },
  set() {
    store.dispatch("sendCommandToSetting", {
      command: "download",
    });
  },
});
</script>
