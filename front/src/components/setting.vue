<template>
  <div
    class="masonry grid h-full w-full flex-col rounded p-2 dark:text-white md:flex-wrap"
  >
    <div class="card-class">
      <theme-settings />
    </div>
    <div class="card-class">
      <database-info />
    </div>
    <div class="card-class">
      <playback-options />
    </div>
    <div class="card-class space-y-1">
      <outputs />
    </div>
    <DownloadCover />
    <div class="card-class">
      <MappingHelp />
    </div>
  </div>
</template>
<style lang="postcss">
.card-class {
  @apply w-full rounded bg-white/60 p-4 backdrop-blur-3xl dark:bg-gray-700/60 md:m-1;
}

.card-header {
  @apply mb-2 text-2xl  underline decoration-secondary  decoration-2;
}
.masonry {
  @apply grid-cols-1 md:grid-cols-2;
  gap: 1rem;
  grid-template-rows: masonry;
}
</style>
