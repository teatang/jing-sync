import { createMemoryHistory, createRouter } from "vue-router";
import Engine from "@/views/engine/Engine.vue";
import CronJob from "@/views/cron-job/CronJob.vue";
import UserSetting from "@/views/setting/UserSetting.vue";

const routes = [
  { path: "/", component: CronJob },
  { path: "/cron-job", component: CronJob },
  { path: "/engine", component: Engine },
  { path: "/user-setting", component: UserSetting },
];

export const router = createRouter({
  history: createMemoryHistory(),
  routes,
});
