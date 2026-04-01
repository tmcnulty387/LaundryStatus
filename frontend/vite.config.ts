import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"

export default defineConfig(({ mode }) => ({
  base: "/",
  plugins: [react()],
  preview: {
    port: 8080,
    strictPort: true,
  },
  // Use vite proxy instead of nginx in dev
  server: mode === "development"
    ? {
        host: "127.0.0.1",
        port: 5173,
        proxy: {
          "/api": {
            target: "http://127.0.0.1:8080",
            changeOrigin: true,
          },
        },
      }
    : {
      host: true,
      port: 8080,
      origin: "http://127.0.0.1:8080",
      strictPort: true,
    },
}))
