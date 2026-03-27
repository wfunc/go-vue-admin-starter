/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#0f172a',
        slate: '#475569',
        line: '#e2e8f0',
        dark: {
          400: '#94a3b8',
          500: '#64748b',
          700: '#334155',
          800: '#1e293b',
          900: '#0f172a'
        },
        primary: {
          50: '#f0fdfa',
          100: '#ccfbf1',
          200: '#99f6e4',
          300: '#5eead4',
          400: '#2dd4bf',
          500: '#14b8a6',
          600: '#0d9488',
          700: '#0f766e',
          800: '#115e59',
          900: '#134e4a'
        }
      },
      fontFamily: {
        sans: ['"Manrope"', '"PingFang SC"', '"Hiragino Sans GB"', 'sans-serif']
      },
      boxShadow: {
        panel: '0 10px 30px rgba(15, 23, 42, 0.08)',
        glass: '0 12px 40px rgba(15, 23, 42, 0.08)'
      }
    }
  },
  plugins: []
}
