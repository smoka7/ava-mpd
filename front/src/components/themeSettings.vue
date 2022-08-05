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
    <RadioGroup v-model="colorScheme" class="flex flex-col space-y-3">
      <template v-for="(scheme, index) in colorSchemes" :key="index">
        <RadioGroupOption v-slot="{ checked }" :value="index">
          <button
            :class="[
              { 'rounded bg-blue-200 text-primary': checked },
              'hover:lighter group flex w-full cursor-pointer flex-col space-y-2 rounded p-2 duration-300 dark:hover:bg-lightest',
            ]"
          >
            <span
              class="flex w-full items-center justify-between group-hover:text-primary"
            >
              <RadioGroupLabel>{{ index }}</RadioGroupLabel>
              <font-awesome-icon
                v-show="checked"
                class="text-primary"
                icon="check-circle"
              />
            </span>
            <span :style="renderScheme(scheme)" class="h-10 w-full rounded">
            </span>
          </button>
        </RadioGroupOption>
      </template>
    </RadioGroup>
  </div>
</template>
<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import {
  RadioGroup,
  RadioGroupLabel,
  RadioGroupOption,
  Switch,
  SwitchGroup,
  SwitchLabel,
} from "@headlessui/vue";
import { colorSchemes, getRGB, setColorScheme } from "../colors.js";
export default {
  components: {
    FontAwesomeIcon,
    RadioGroup,
    RadioGroupLabel,
    RadioGroupOption,
    Switch,
    SwitchGroup,
    SwitchLabel,
  },
  data() {
    return {
      isDark: localStorage.getItem("theme") == "dark",
      colorScheme: localStorage.getItem("colorScheme") || "first",
      colorSchemes: colorSchemes,
    };
  },
  methods: {
    toggleDarkMode() {
      if (!this.isDark) {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("theme", "light");
        return;
      }
      localStorage.setItem("theme", "dark");
      document.documentElement.classList.add("dark");
    },
    changeColorScheme() {
      localStorage.setItem("colorScheme", this.colorScheme);
      setColorScheme();
    },
    renderScheme(scheme) {
      let style = "background-image: linear-gradient(90deg,";
      let percent = 0;

      scheme.forEach((element) => {
        const rgb = getRGB(element);
        style += rgb + " " + percent + "%,";
        percent += 20;
        style += rgb + " " + percent + "%,";
      });

      return style + getRGB(scheme[4]) + " 100%);";
    },
  },
  watch: {
    colorScheme(newCS) {
      this.changeColorScheme(newCS);
    },
    isDark() {
      this.toggleDarkMode();
    },
  },
};
</script>
