import { defineConfig } from "vitest/config";
import path from "path";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  test: {
    root: __dirname,
    include: ["test/**/*.{test,spec}.?(c|m)[jt]s?(x)"],
    testTimeout: 1000 * 29,
  },
});
