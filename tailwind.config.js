/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["frontend/*.{html,js}","frontend/frontend.go"],
  theme: {
    extend: {
      fontFamily: {
        popoins: ['"Poppins"', 'mono']
      },
      // that is animation class
      animation: {
        fade: 'fadeIn 0.2s ease-in',
      },
      keyframes: theme => ({
        fadeIn: {
          '100%': { opacity: '100%' },
          '0%': { opacity: '0%' },
        },
      }),
    },
  },
  plugins: [],
}
// npx tailwindcss -i ./frontend/app.css -o ./frontend/style.css --watch
