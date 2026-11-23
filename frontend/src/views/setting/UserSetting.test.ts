import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import ElementPlus from "element-plus";
import UserSetting from "./UserSetting.vue";
import RestClient from "@/utils/rest-client";
import i18n from "@/i18n";

describe("UserSetting.vue", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.spyOn(RestClient.prototype, "get").mockResolvedValue({
      list: [
        {
          id: 1,
          username: "admin_user",
          password: "",
          create_time: "2026-10-23T10:00:00Z",
          update_time: "2026-10-23T10:00:00Z",
        },
      ],
      pagination: { page: 1, size: 10, total: 1 },
    });
  });

  it("renders user table with action buttons", async () => {
    const wrapper = mount(UserSetting, {
      global: {
        plugins: [ElementPlus, i18n],
      },
    });

    await flushPromises();
    const text = wrapper.text();
    expect(text).toContain("用户账号管理");
    expect(text).toContain("admin_user");
    expect(text).toContain("新增用户");
  });
});
