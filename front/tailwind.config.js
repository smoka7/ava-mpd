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
    screens: {
      sm: "768px",
      md: "900px",
    },
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
