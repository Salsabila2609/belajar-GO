import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Semua request ke /api/... diteruskan ke server Go.
    // Dengan ini browser menganggapnya satu asal (origin), jadi tidak ada masalah CORS.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
