import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import { createRouter, createMemoryHistory } from "vue-router";
import ElementPlus from "element-plus";
import CronJob from "./CronJob.vue";
import RestClient from "@/utils/rest-client";

describe("CronJob.vue", () => {
  let router: any;

  beforeEach(() => {
    setActivePinia(createPinia());
    router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: "/", component: { template: "<div />" } }],
    });
    vi.spyOn(RestClient.prototype, "get").mockImplementation((url: string) => {
      if (url === "/job") {
        return Promise.resolve({
          list: [
            {
              id: 1,
              remark: "照片增量备份",
              src_path: "/photos",
              dst_path: "/backup/photos",
              engine_id: 1,
              sync_mode: 1,
              speed: 0,
              enabled: 1,
              create_time: "2026-10-23T10:00:00Z",
              update_time: "2026-10-23T10:00:00Z",
            },
          ],
          pagination: { page: 1, size: 10, total: 1 },
        });
      }
      return Promise.resolve({ list: [], pagination: { page: 1, size: 10, total: 0 } });
    });
  });

  it("renders job list and displays immediate sync button", async () => {
    const wrapper = mount(CronJob, {
      global: {
        plugins: [ElementPlus, router],
      },
    });

    await flushPromises();
    const text = wrapper.text();
    expect(text).toContain("照片增量备份");
    expect(text).toContain("立即同步");
  });
});
