module.exports = {
  content: ['./index.html', './src/**/*.{vue,js,jsx}'],
  darkMode: 'class', // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        primary: 'var(--primary-color)',
        secondary: 'var(--secondary-color)',
        lighter:'var(--lighter-color)',
        lightest:'var(--lightest-color)',
        foreground:'var(--foreground-color)',
        accent:'var(--accent-color)',
      }
    },
  },
  variants: {
    extend: {
    },
  },
  plugins: [],
}
