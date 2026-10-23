<template>
  <div class="theme-switch inline-flex items-center">
    <el-dropdown trigger="click" @command="handleThemeChange">
      <button
        type="button"
        class="inline-flex items-center justify-center p-2 rounded-lg text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors focus:outline-none"
        :title="themeTitle"
      >
        <el-icon :size="18">
          <Sunny v-if="themeStore.theme === 'light'" />
          <Moon v-else-if="themeStore.theme === 'dark'" />
          <Monitor v-else />
        </el-icon>
      </button>

      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="light" :class="{ 'font-semibold text-primary': themeStore.theme === 'light' }">
            <div class="flex items-center gap-2">
              <el-icon><Sunny /></el-icon>
              <span>{{ $t ? $t('theme.light') : '浅色模式' }}</span>
              <el-icon v-if="themeStore.theme === 'light'" class="ml-auto"><Check /></el-icon>
            </div>
          </el-dropdown-item>
          <el-dropdown-item command="dark" :class="{ 'font-semibold text-primary': themeStore.theme === 'dark' }">
            <div class="flex items-center gap-2">
              <el-icon><Moon /></el-icon>
              <span>{{ $t ? $t('theme.dark') : '深色模式' }}</span>
              <el-icon v-if="themeStore.theme === 'dark'" class="ml-auto"><Check /></el-icon>
            </div>
          </el-dropdown-item>
          <el-dropdown-item command="auto" :class="{ 'font-semibold text-primary': themeStore.theme === 'auto' }">
            <div class="flex items-center gap-2">
              <el-icon><Monitor /></el-icon>
              <span>{{ $t ? $t('theme.auto') : '跟随系统' }}</span>
              <el-icon v-if="themeStore.theme === 'auto'" class="ml-auto"><Check /></el-icon>
            </div>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { Sunny, Moon, Monitor, Check } from "@element-plus/icons-vue";
import { useThemeStore, type ThemeMode } from "@/stores/themeStore";

const themeStore = useThemeStore();

const themeTitle = computed(() => {
  if (themeStore.theme === "light") return "浅色模式";
  if (themeStore.theme === "dark") return "深色模式";
  return "跟随系统";
});

function handleThemeChange(command: string) {
  themeStore.setTheme(command as ThemeMode);
}
</script>

<style scoped>
.text-primary {
  color: var(--el-color-primary, #409eff);
}
</style>
