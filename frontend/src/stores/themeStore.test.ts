import { describe, it, expect, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useThemeStore } from "./themeStore";

describe("useThemeStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    document.documentElement.className = "";
  });

  it("initializes with default auto/system theme", () => {
    const store = useThemeStore();
    expect(["light", "dark", "auto"]).toContain(store.theme);
  });

  it("updates theme to dark and sets dark class on html", () => {
    const store = useThemeStore();
    store.setTheme("dark");
    expect(store.theme).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });

  it("updates theme to light and removes dark class from html", () => {
    const store = useThemeStore();
    store.setTheme("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
    store.setTheme("light");
    expect(store.theme).toBe("light");
    expect(document.documentElement.classList.contains("dark")).toBe(false);
  });
});
