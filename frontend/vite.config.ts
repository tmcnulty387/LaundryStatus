import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"

export default defineConfig(({ mode }) => ({
  plugins: [react()],
  // Use vite proxy instead of nginx
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
    : undefined,
}))