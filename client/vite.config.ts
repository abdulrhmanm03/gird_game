import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Proxying requests that start with /api to your backend
      "/api": {
        target: "http://localhost:3000", // your backend server
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""), // optional: rewrite the path if needed
      },
    },
  },
});
