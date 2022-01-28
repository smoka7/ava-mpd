<template>
  <div class="flex flex-col space-y-4">
    <SwitchGroup as="div" class="flex items-center p-2">
      <SwitchLabel class="mr-4 text-lg md:mr-2">Dark Mode</SwitchLabel>
      <Switch
        isDark
        v-model="isDark"
        class="bg-lightest relative inline-flex h-10 w-16 items-center rounded shadow-inner md:h-7 md:w-14"
      >
        <span class="sr-only">Enable Dark Mode</span>
        <span
          class="bg-primary inline-block h-8 w-8 translate-x-1 transform rounded duration-300 dark:translate-x-7 md:h-5 md:w-6"
          isDark
        />
      </Switch>
    </SwitchGroup>
    <h2 class="text-lg">Color Schemes</h2>
    <RadioGroup v-model="colorScheme" class="flex flex-col space-y-3">
      <template v-for="(scheme, index) in colorSchemes" :key="index">
        <RadioGroupOption v-slot="{ checked }" :value="index">
          <div
            :class="[
              { 'text-primary rounded bg-blue-200': checked },
              'flex cursor-pointer flex-col space-y-2 p-2 duration-200',
            ]"
          >
            <div class="flex items-center justify-between">
              <RadioGroupLabel>{{ index }}</RadioGroupLabel>
              <font-awesome-icon
                v-show="checked"
                class="text-primary ml-auto"
                icon="check-circle"
              ></font-awesome-icon>
            </div>
            <div
              :style="renderScheme(scheme)"
              class="h-10 w-full rounded shadow md:h-8"
            ></div>
          </div>
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
      return (
        "background-image: linear-gradient(90deg," +
        getRGB(scheme[0]) +
        " 0%," +
        getRGB(scheme[0]) +
        " 20%," +
        getRGB(scheme[1]) +
        " 20%," +
        getRGB(scheme[1]) +
        " 40%," +
        getRGB(scheme[2]) +
        " 40%," +
        getRGB(scheme[2]) +
        " 60%," +
        getRGB(scheme[3]) +
        " 60%," +
        getRGB(scheme[3]) +
        " 80%," +
        getRGB(scheme[4]) +
        " 80%," +
        getRGB(scheme[4]) +
        " 100%);"
      );
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
