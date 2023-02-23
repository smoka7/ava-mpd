<template>
  <div class="flex flex-col space-y-4">
    <h2 class="card-header">Color Schemes</h2>
    <SwitchGroup as="div" class="flex items-center p-2">
      <SwitchLabel class="mr-4 text-lg md:mr-2">Dark Mode</SwitchLabel>
      <Switch
        isDark
        v-model="isDark"
        class="relative inline-flex h-10 w-20 items-center rounded bg-lightest shadow-inner md:h-7 md:w-14"
      >
        <span class="sr-only">Enable Dark Mode</span>
        <span
          class="inline-block h-8 w-8 translate-x-1 transform rounded bg-primary duration-300 dark:translate-x-10 md:h-5 md:w-6 md:dark:translate-x-7"
          isDark
        />
      </Switch>
    </SwitchGroup>
    <RadioGroup
      v-model="colorScheme"
      class="flex flex-col space-y-3 rounded bg-white p-2 dark:bg-gray-600"
    >
      <RadioGroupOption
        v-for="(scheme, index) in colorSchemes"
        :key="index"
        as="template"
        v-slot="{ checked }"
        :value="index"
      >
        <div
          :class="[
            { 'rounded bg-gray-200 text-primary': checked },
            'hover:lighter group flex w-full cursor-pointer flex-col space-y-2 rounded p-2 duration-200 dark:hover:bg-lightest',
          ]"
        >
          <span
            class="flex w-full items-center justify-between group-hover:text-primary"
          >
            <RadioGroupLabel>{{ index }}</RadioGroupLabel>
            <font-awesome-icon
              v-if="checked"
              class="text-primary"
              icon="check-circle"
            />
          </span>
          <span :style="renderScheme(scheme)" class="h-10 w-full rounded">
          </span>
        </div>
      </RadioGroupOption>
    </RadioGroup>
  </div>
</template>
<script setup lang="ts">
import {
  colorSchemes,
  getRGB,
  setColorScheme,
  type colorScheme,
} from "../colors";

const colorScheme = ref(localStorage.getItem("colorScheme") || "first");
const isDark = ref(localStorage.getItem("theme") == "dark");

function toggleDarkMode() {
  if (!isDark.value) {
    document.documentElement.classList.remove("dark");
    localStorage.setItem("theme", "light");
    return;
  }
  localStorage.setItem("theme", "dark");
  document.documentElement.classList.add("dark");
}

function changeColorScheme() {
  localStorage.setItem("colorScheme", colorScheme.value);
  setColorScheme();
}

function renderScheme(scheme: colorScheme): string {
  let style = "background-image: linear-gradient(90deg,";
  let percent = 0;

  scheme.forEach((element) => {
    const rgb = getRGB(element);
    style += rgb + " " + percent + "%,";
    percent += 20;
    style += rgb + " " + percent + "%,";
  });

  return style + getRGB(scheme[4]) + " 100%);";
}

watch(colorScheme, () => {
  changeColorScheme();
});

watch(isDark, () => {
  toggleDarkMode();
});
</script>
