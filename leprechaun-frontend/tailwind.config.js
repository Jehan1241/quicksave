/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],

  theme: {
    extend: {
      colors: {
        primary: "#0a0a0a",
        secondary: "#8d6e63",
      },
    },
  },

  plugins: [],
};
