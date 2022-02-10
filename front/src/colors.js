export const colorSchemes = {
  first: ["51 54 82", "250 208 44", "189 239 83", "144 173 198", "213 227 172"],
  second: [
    "51 54 82",
    "246 81 100",
    "90 101 197",
    "189 213 234",
    "229 229 243",
  ],
  third: [
    "44 47 62",
    "252 157 154",
    "145 109 196",
    "191 149 249",
    "207 184 239",
  ],
  fourth: ["37 41 52", "241 177 39", "235 158 5", "227 185 107", "227 185 147"],
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

export function setColorScheme() {
  let scheme = localStorage.getItem("colorScheme");
  switch (scheme) {
    case "second":
      writeColors(colorSchemes.second);
      break;
    case "third":
      writeColors(colorSchemes.third);
      break;
    case "fourth":
      writeColors(colorSchemes.fourth);
      break;
    default:
      writeColors(colorSchemes.first);
  }
}
