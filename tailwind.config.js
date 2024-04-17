/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './cmd/ray-peat-rodeo/**/*.{templ,go}',
    './internal/**/*.{templ,go}',
    './internal/blog/**/*.md',
    './build/index.html'
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

