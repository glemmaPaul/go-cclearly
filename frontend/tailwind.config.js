/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        emerald: {
          300: '#6ee7b7',
          700: '#047857',
          900: '#064e3b',
        },
        amber: {
          300: '#fcd34d',
          700: '#b45309',
          900: '#78350f',
        },
        red: {
          300: '#fca5a5',
          700: '#b91c1c',
          900: '#7f1d1d',
        },
        gray: {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',
          800: '#1f2937',
          900: '#111827',
          950: '#030712',
        },
      },
    },
  },
  plugins: [],
} 