import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import ElementPlus from "element-plus";
import LoginView from "./LoginView.vue";

describe("LoginView.vue", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.restoreAllMocks();
  });

  it("renders login form with title and input fields", () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [ElementPlus],
      },
    });

    expect(wrapper.text()).toContain("系统登录");
    const inputs = wrapper.findAll("input");
    expect(inputs.length).toBeGreaterThanOrEqual(2);
    expect(wrapper.find("button.login-btn").exists()).toBe(true);
  });
});
