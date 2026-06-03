/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        leonidas: {
          50:  '#fef9ec',
          100: '#fdf0c8',
          200: '#fae08d',
          300: '#f6c94a',
          400: '#f2b020',
          500: '#ec9009',
          600: '#d06b05',
          700: '#ad4c08',
          800: '#8d3b0d',
          900: '#74310f',
          950: '#421703',
        },
        surface: {
          900: '#0d0d0d',
          800: '#161616',
          700: '#1e1e1e',
          600: '#252525',
          500: '#2e2e2e',
        },
      },
    },
  },
  plugins: [],
}
