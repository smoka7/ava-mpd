function withOpacityValue(variable) {
  return ({ opacityValue }) => {
    if (opacityValue === undefined) {
      return `rgb(var(${variable}))`;
    }
    return `rgb(var(${variable}) / ${opacityValue})`;
  };
}
module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,jsx}"],
  darkMode: "class", // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        primary: withOpacityValue("--primary-color"),
        secondary: withOpacityValue("--secondary-color"),
        lighter: withOpacityValue("--lighter-color"),
        lightest: withOpacityValue("--lightest-color"),
        accent: withOpacityValue("--accent-color"),
      },
    },
  },
};
