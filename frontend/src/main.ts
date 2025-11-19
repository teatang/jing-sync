import { createApp } from "vue";
import { createPinia } from "pinia";

import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import "element-plus/theme-chalk/dark/css-vars.css";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

import App from "@/App.vue";
import { router } from "@/router/index";
import "@/assets/tailwind.css";
import i18n from "@/i18n";

const app = createApp(App);

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

app
  .use(ElementPlus) //
  .use(router) //
  .use(createPinia()) //
  .use(i18n) //
  .mount("#app");
