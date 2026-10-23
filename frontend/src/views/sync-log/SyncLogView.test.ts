import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import { createRouter, createMemoryHistory } from "vue-router";
import ElementPlus from "element-plus";
import SyncLogView from "./SyncLogView.vue";
import RestClient from "@/utils/rest-client";

describe("SyncLogView.vue", () => {
  let router: any;

  beforeEach(() => {
    setActivePinia(createPinia());
    router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: "/sync-logs", component: { template: "<div />" } }],
    });
    vi.spyOn(RestClient.prototype, "get").mockImplementation((url: string) => {
      if (url.includes("/sync-logs")) {
        return Promise.resolve({
          list: [
            {
              id: 1,
              job_id: 10,
              status: 1,
              files_added: 8,
              start_time: "2026-10-23T10:00:00Z",
              end_time: "2026-10-23T10:00:05Z",
              details: JSON.stringify({
                job_id: 10,
                job_remark: "每日自动归档",
                files_found: 10,
                files_added: 8,
                files_skipped: 2,
                files_failed: 0,
                details: ["开始扫描目录", "跳过文件 exist.txt", "同步成功 photo1.jpg"],
              }),
            },
          ],
          pagination: { page: 1, size: 10, total: 1 },
        });
      }
      return Promise.resolve({ list: [], pagination: { page: 1, size: 10, total: 0 } });
    });
  });

  it("renders sync logs table and details button", async () => {
    const wrapper = mount(SyncLogView, {
      global: {
        plugins: [ElementPlus, router],
      },
    });

    await flushPromises();
    const text = wrapper.text();
    expect(text).toContain("同步执行日志");
    expect(text).toContain("查看详情");
  });
});
