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
    './build/index.html',
    './build/2018-01-22-when-western-medicine-isn\'t-working/index.html',
    './build/2022-11-18-untitled-interview-2022-11-18/index.html'
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

