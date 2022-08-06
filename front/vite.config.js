import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/vite";
import { HeadlessUiResolver } from "unplugin-vue-components/resolvers";

export default defineConfig({
  plugins: [
    vue(),

    Components({
      resolvers: [
        HeadlessUiResolver(),
        (componentName) => {
          if (componentName.startsWith("FontAwesomeIcon"))
            return {
              name: "FontAwesomeIcon",
              from: "@fortawesome/vue-fontawesome",
            };
        },
      ],
    }),
  ],
  server: {
    proxy: {
      "/api": "http://localhost:3001",
      "/coverart": "http://localhost:3001",
    },
    host: "0.0.0.0",
  },
});
