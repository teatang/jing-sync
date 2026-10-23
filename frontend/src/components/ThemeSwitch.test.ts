import { describe, it, expect, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import ElementPlus from "element-plus";
import ThemeSwitch from "./ThemeSwitch.vue";
import { useThemeStore } from "@/stores/themeStore";
import i18n from "@/i18n";

describe("ThemeSwitch.vue", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders the theme switch trigger button", () => {
    const wrapper = mount(ThemeSwitch, {
      global: {
        plugins: [ElementPlus, i18n],
      },
    });
    expect(wrapper.find("button").exists() || wrapper.find(".theme-switch").exists()).toBe(true);
  });

  it("changes theme when setTheme is called via store", async () => {
    const wrapper = mount(ThemeSwitch, {
      global: {
        plugins: [ElementPlus, i18n],
      },
    });
    const store = useThemeStore();
    store.setTheme("dark");
    expect(store.theme).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });
});
