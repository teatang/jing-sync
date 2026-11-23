<template>
  <el-config-provider :locale="elementLocale">
    <Layout v-if="token" />
    <LoginView v-if="!token" />
  </el-config-provider>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";
import Layout from "@/views/Layout.vue";
import LoginView from "@/views/LoginView.vue";
import { getToken } from "@/utils/token";
import { useLocaleStore } from "@/stores/localeStore";
import { useThemeStore } from "@/stores/themeStore";
import zhCn from "element-plus/dist/locale/zh-cn.mjs";
import en from "element-plus/dist/locale/en.mjs";

const token = getToken();
const localeStore = useLocaleStore();
const themeStore = useThemeStore();

onMounted(() => {
  themeStore.applyDomTheme();
});

const elementLocale = computed(() => {
  return localeStore.currentLang === "zh-CN" ? zhCn : en;
});
</script>
