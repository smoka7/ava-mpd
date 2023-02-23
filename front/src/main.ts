import { createApp } from "vue";
import App from "./App.vue";
import "./index.css";
import "./fontAwesome";
const pinia = createPinia();

createApp(App).use(pinia).mount("#app");
