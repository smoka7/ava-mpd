import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/vite";
import { HeadlessUiResolver } from "unplugin-vue-components/resolvers";
import AutoImport from "unplugin-auto-import/vite";
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      dts: true,
      imports: ["vue", "vuex"],
    }),
    Components({
      dts: true,
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
