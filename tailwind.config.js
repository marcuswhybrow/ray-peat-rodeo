/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    // Scanning source code is an order of magnitude faster than scanning build. 
    // But we cannot use runtime values in class names. For example,
    // 
    // ```
    // // Tailwind cannot understand this at compile time
    // fmt.Sprintf(`<span class="before:content-['%v']">`, count)
    //
    // // So far, CSS counters have sufficed in such situations:
    // `<span class="[counter-increment:count] before:content-[counter(count)]">`
    // ```
    './cmd/**/*.{templ,go}',
    './internal/**/*.{templ,go}',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('tailwind-scrollbar'),
  ],
}

