export const colorSchemes = {
  first: [
    "51 54 82",
    "144 173 198",
    "250 208 44",
    "189 213 234",
    "213 227 172",
  ],
  second: [
    "51 54 82",
    "90 101 197",
    "246 81 100",
    "189 213 234",
    "229 229 243",
  ],
  third: [
    "44 47 62",
    "145 109 196",
    "252 157 154",
    "191 149 249",
    "207 184 239",
  ],
  fourth: [
    "65 60 88",
    "116 124 163",
    "242 195 160",
    "145 186 172",
    "239 232 206",
  ],
};

export function getRGB(color) {
  return "rgb(" + color + ")";
}
export function writeColors(colors = colorSchemes.first) {
  document.getElementById("style").innerText =
    ":root{--primary-color:" +
    colors[0] +
    ";--secondary-color:" +
    colors[1] +
    ";--accent-color:" +
    colors[2] +
    ";--lighter-color:" +
    colors[3] +
    ";--lightest-color:" +
    colors[4] +
    ";}";
}

export function loadTheme() {
  if (
    localStorage.theme === "dark" ||
    (!("theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    document.documentElement.classList.add("dark");
  }
}

export function setColorScheme() {
  const scheme = localStorage.getItem("colorScheme");
  if (colorSchemes[scheme]) {
    writeColors(colorSchemes[scheme]);
    return;
  }
  writeColors(colorSchemes["first"]);
}
