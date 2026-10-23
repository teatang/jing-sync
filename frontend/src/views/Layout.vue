<template>
  <el-container v-loading="mainStore.loading" class="h-screen w-screen overflow-hidden bg-slate-50 dark:bg-slate-950 text-slate-800 dark:text-slate-100 transition-colors">
    <!-- 左侧极简侧边栏 -->
    <el-aside
      :width="isCollapse ? '64px' : '220px'"
      class="border-r border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 flex flex-col transition-all duration-300 select-none"
    >
      <!-- Logo 区域 -->
      <div class="h-16 flex items-center px-4 border-b border-slate-100 dark:border-slate-800/80 gap-3 overflow-hidden">
        <div class="w-8 h-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center font-bold text-lg shrink-0">
          <el-icon><Refresh /></el-icon>
        </div>
        <div v-show="!isCollapse" class="font-bold text-base tracking-wide whitespace-nowrap overflow-hidden text-ellipsis">
          JingSync
        </div>
      </div>

      <!-- 导航菜单 -->
      <el-scrollbar class="flex-1 py-2">
        <el-menu
          router
          :collapse="isCollapse"
          :default-active="activeRoute"
          class="!border-none !bg-transparent"
        >
          <el-menu-item index="/cron-job">
            <el-icon><AlarmClock /></el-icon>
            <template #title>{{ $t('nav.jobs') || '定时同步' }}</template>
          </el-menu-item>

          <el-menu-item index="/engine">
            <el-icon><Connection /></el-icon>
            <template #title>{{ $t('nav.engines') || '存储引擎' }}</template>
          </el-menu-item>

          <el-menu-item index="/sync-logs">
            <el-icon><Document /></el-icon>
            <template #title>{{ $t('nav.logs') || '同步日志' }}</template>
          </el-menu-item>

          <el-menu-item index="/user-setting">
            <el-icon><User /></el-icon>
            <template #title>{{ $t('nav.users') || '用户管理' }}</template>
          </el-menu-item>
        </el-menu>
      </el-scrollbar>

      <!-- 侧边栏底部折叠按钮 -->
      <div class="p-3 border-t border-slate-100 dark:border-slate-800 flex justify-center">
        <button
          type="button"
          @click="isCollapse = !isCollapse"
          class="p-2 rounded-lg text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors w-full flex items-center justify-center"
          :title="isCollapse ? '展开侧边栏' : '折叠侧边栏'"
        >
          <el-icon :size="16">
            <Expand v-if="isCollapse" />
            <Fold v-else />
          </el-icon>
        </button>
      </div>
    </el-aside>

    <!-- 右侧主体内容区 -->
    <el-container class="flex flex-col flex-1 h-screen overflow-hidden">
      <!-- 极简头部工具栏 -->
      <el-header class="h-16 border-b border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur px-6 flex items-center justify-between z-10">
        <div class="text-sm font-medium text-slate-500 dark:text-slate-400">
          {{ pageTitle }}
        </div>

        <div class="flex items-center gap-3">
          <!-- 日夜模式切换 -->
          <ThemeSwitch />

          <!-- 多语言切换 -->
          <el-dropdown trigger="click" @command="handleLanguageChange">
            <button
              type="button"
              class="inline-flex items-center justify-center p-2 rounded-lg text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
              title="切换语言 / Language"
            >
              <el-icon :size="18"><ChatDotRound /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="lang in langs"
                  :key="lang"
                  :command="lang"
                  :class="{ 'font-semibold text-primary': localeStore.currentLang === lang }"
                >
                  {{ lang === 'zh-CN' ? '简体中文 (zh-CN)' : 'English (en)' }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 用户菜单 -->
          <el-dropdown trigger="click" @command="handleUserCommand">
            <button
              type="button"
              class="flex items-center gap-2 px-2 py-1 rounded-lg text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
            >
              <div class="w-7 h-7 rounded-full bg-primary/20 text-primary flex items-center justify-center text-xs font-semibold">
                {{ usernameInitial }}
              </div>
              <span class="text-sm font-medium hidden sm:inline">{{ currentUsername }}</span>
              <el-icon :size="12"><ArrowDown /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <div class="flex items-center gap-2">
                    <el-icon><User /></el-icon>
                    <span>{{ $t('nav.users') || '用户设置' }}</span>
                  </div>
                </el-dropdown-item>
                <el-dropdown-item divided command="logout" class="text-red-500">
                  <div class="flex items-center gap-2">
                    <el-icon><SwitchButton /></el-icon>
                    <span>{{ $t('nav.logout') || '退出登录' }}</span>
                  </div>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主视图 -->
      <el-main class="flex-1 overflow-auto p-6 bg-slate-50 dark:bg-slate-950/60">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useMainStore } from "@/stores";
import { useLocaleStore } from "@/stores/localeStore";
import { langs, type Language } from "@/types/index";
import { removeToken } from "@/utils/token";
import ThemeSwitch from "@/components/ThemeSwitch.vue";
import {
  AlarmClock,
  Connection,
  Document,
  User,
  Expand,
  Fold,
  ChatDotRound,
  Refresh,
  ArrowDown,
  SwitchButton,
} from "@element-plus/icons-vue";

const route = useRoute();
const router = useRouter();
const mainStore = useMainStore();
const localeStore = useLocaleStore();

const isCollapse = ref<boolean>(false);

const activeRoute = computed(() => {
  return route.path === "/" ? "/cron-job" : route.path;
});

const currentUsername = ref("Admin");
const usernameInitial = computed(() => currentUsername.value.charAt(0).toUpperCase());

const pageTitle = computed(() => {
  if (route.path.includes("cron-job") || route.path === "/") return "定时任务调度";
  if (route.path.includes("engine")) return "存储引擎管理";
  if (route.path.includes("sync-logs")) return "同步执行日志";
  if (route.path.includes("user-setting")) return "用户偏好设置";
  return "JingSync 数据同步中心";
});

function handleLanguageChange(command: Language) {
  localeStore.setLanguage(command);
}

function handleUserCommand(command: string) {
  if (command === "logout") {
    removeToken();
    window.location.reload();
  } else if (command === "profile") {
    router.push("/user-setting");
  }
}
</script>

<style scoped>
.text-primary {
  color: var(--el-color-primary, #3b82f6);
}
.bg-primary\/10 {
  background-color: rgba(59, 130, 246, 0.1);
}
.bg-primary\/20 {
  background-color: rgba(59, 130, 246, 0.2);
}
</style>
