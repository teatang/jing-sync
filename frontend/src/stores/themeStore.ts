import { defineStore } from "pinia";
import { ref, computed } from "vue";

export type ThemeMode = "light" | "dark" | "auto";

const THEME_STORAGE_KEY = "jing-sync-theme";

export const useThemeStore = defineStore("theme", () => {
  // 从本地存储恢复或默认为 auto
  const savedTheme = (typeof localStorage !== "undefined"
    ? localStorage.getItem(THEME_STORAGE_KEY)
    : null) as ThemeMode | null;

  const theme = ref<ThemeMode>(
    savedTheme && ["light", "dark", "auto"].includes(savedTheme)
      ? savedTheme
      : "auto"
  );

  const systemIsDark = ref<boolean>(false);

  // 初始化媒体查询监听系统偏好
  if (typeof window !== "undefined" && window.matchMedia) {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    systemIsDark.value = mediaQuery.matches;

    mediaQuery.addEventListener("change", (e) => {
      systemIsDark.value = e.matches;
      if (theme.value === "auto") {
        applyDomTheme();
      }
    });
  }

  // 当前是否实际为暗色模式
  const isDark = computed(() => {
    if (theme.value === "dark") return true;
    if (theme.value === "light") return false;
    return systemIsDark.value;
  });

  // 更新 DOM 根节点的 class
  function applyDomTheme() {
    if (typeof document === "undefined") return;
    const root = document.documentElement;
    if (isDark.value) {
      root.classList.add("dark");
    } else {
      root.classList.remove("dark");
    }
  }

  // 显式设置模式
  function setTheme(newTheme: ThemeMode) {
    theme.value = newTheme;
    if (typeof localStorage !== "undefined") {
      localStorage.setItem(THEME_STORAGE_KEY, newTheme);
    }
    applyDomTheme();
  }

  // 循环切换 light -> dark -> auto -> light
  function toggleTheme() {
    if (theme.value === "light") {
      setTheme("dark");
    } else if (theme.value === "dark") {
      setTheme("auto");
    } else {
      setTheme("light");
    }
  }

  // 初始执行一次
  applyDomTheme();

  return {
    theme,
    isDark,
    setTheme,
    toggleTheme,
    applyDomTheme,
  };
});
