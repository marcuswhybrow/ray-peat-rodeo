/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './build/**/*.html',
  ],
  theme: {
    extend: {
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
      }
    }
  },
  plugins: [
    require('tailwind-scrollbar'),
  ],
}

