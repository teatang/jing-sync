import { describe, it, expect, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useLocaleStore } from "./localeStore";

describe("useLocaleStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
  });

  it("defaults to zh-CN if not stored", () => {
    const store = useLocaleStore();
    expect(store.currentLang).toBe("zh-CN");
    expect(store.getLanguage()).toBe("zh-CN");
  });

  it("updates language and saves to localStorage", () => {
    const store = useLocaleStore();
    store.setLanguage("en");
    expect(store.currentLang).toBe("en");
    expect(localStorage.getItem("lang")).toBe("en");
  });
});
