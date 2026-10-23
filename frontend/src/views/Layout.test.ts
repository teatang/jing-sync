import { describe, it, expect, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import ElementPlus from "element-plus";
import { createRouter, createMemoryHistory } from "vue-router";
import Layout from "./Layout.vue";
import i18n from "@/i18n";

describe("Layout.vue", () => {
  let router: any;

  beforeEach(() => {
    setActivePinia(createPinia());
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: "/", component: { template: "<div>Home</div>" } },
        { path: "/cron-job", component: { template: "<div>Jobs</div>" } },
        { path: "/engine", component: { template: "<div>Engines</div>" } },
        { path: "/sync-logs", component: { template: "<div>Logs</div>" } },
        { path: "/user-setting", component: { template: "<div>Users</div>" } },
      ],
    });
  });

  it("renders navigation items including cron-job, engine, and sync-logs", async () => {
    const wrapper = mount(Layout, {
      global: {
        plugins: [ElementPlus, router, i18n],
      },
    });

    const text = wrapper.text();
    expect(text).toContain("定时");
    expect(text).toContain("引擎");
    expect(text).toContain("日志");
  });
});
