import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import ElementPlus from "element-plus";
import LoginView from "./LoginView.vue";
import i18n from "@/i18n";

describe("LoginView.vue", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.restoreAllMocks();
  });

  it("renders login form with title and input fields", () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [ElementPlus, i18n],
      },
    });

    expect(wrapper.text()).toContain("系统登录");
    const inputs = wrapper.findAll("input");
    expect(inputs.length).toBeGreaterThanOrEqual(2);
    expect(wrapper.find("button.login-btn").exists()).toBe(true);
  });

  it("updates text reactively when switching language to English", async () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [ElementPlus, i18n],
      },
    });

    i18n.global.locale.value = "en";
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain("System Sign In");
    expect(wrapper.text()).toContain("Username");
    expect(wrapper.text()).toContain("Password");

    // 切回中文
    i18n.global.locale.value = "zh-CN";
    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain("系统登录");
  });
});
