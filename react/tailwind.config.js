module.exports = {
  content: [
    "./dist/**/*.html",
    "./src/**/*.{js,jsx,ts,tsx}",
    "./*.html",
  ],
  plugins: [require("@tailwindcss/forms"), require("tw-elements/plugin.cjs")],
  variants: {
    extend: {
      opacity: ["disabled"],
    },
  },
  theme: {
    extend: {
      colors: {
        mongoose: "#b29f7e",
      },
    },
  },

};
