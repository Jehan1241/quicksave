/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],

  theme: {
    extend: {
      colors: {
        primary: '#0a0a0a',
        gameView: '#101010',
        secondary: '#8d6e63',
        gameBG: '#0a0a0a'
      },
      fontFamily: {
        sans: ['Inter', 'system-ui']
      }
    }
  },

  plugins: []
}
