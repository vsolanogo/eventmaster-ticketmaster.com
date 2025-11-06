import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  server: {
    port: 3001,
    proxy: {
      "/api": {
        target: "http://localhost:3000",
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          // console.log({ path });
          // console.log(path.replace(/^\/api/, ""));
          return path.replace(/^\/api/, "");
        },
      },
    },
  },
  plugins: [react()],
});
