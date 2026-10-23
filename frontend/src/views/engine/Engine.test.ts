import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import ElementPlus from "element-plus";
import Engine from "./Engine.vue";
import RestClient from "@/utils/rest-client";

describe("Engine.vue", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.spyOn(RestClient.prototype, "get").mockResolvedValue({
      list: [
        {
          id: 1,
          url: "http://127.0.0.1:5244",
          token: "alist-xyz",
          remark: "My OpenList",
          create_time: "2026-10-23T10:00:00Z",
          update_time: "2026-10-23T10:00:00Z",
        },
      ],
      pagination: { page: 1, size: 10, total: 1 },
    });
  });

  it("renders engine list and test connection button", async () => {
    const wrapper = mount(Engine, {
      global: {
        plugins: [ElementPlus],
      },
    });

    await flushPromises();
    const text = wrapper.text();
    expect(text).toContain("测试连接");
  });
});
