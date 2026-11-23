<template>
  <div class="login-wrapper relative min-h-screen w-full flex items-center justify-center overflow-hidden bg-slate-50 dark:bg-slate-950 text-slate-800 dark:text-slate-100 transition-colors duration-500 select-none">
    <!-- 动态背景光晕动画层 (支持浅色/深色模式自适应) -->
    <div class="ambient-glow glow-1"></div>
    <div class="ambient-glow glow-2"></div>
    <div class="ambient-glow glow-3"></div>

    <!-- 赛博科技网格背景纹理 (深浅色模式自适应透明度) -->
    <div class="cyber-grid absolute inset-0 pointer-events-none"></div>

    <!-- 顶部状态栏：主题模式与多语言切换 -->
    <header class="absolute top-4 right-4 sm:top-6 sm:right-6 flex items-center gap-2.5 z-20">
      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200/80 dark:border-white/10 rounded-xl p-1 shadow-sm dark:shadow-lg flex items-center gap-1 transition-colors">
        <!-- 主题切换组件 -->
        <ThemeSwitch />

        <!-- 语言切换下拉菜单 -->
        <el-dropdown trigger="click" @command="handleLanguageChange">
          <button
            type="button"
            class="inline-flex items-center justify-center p-2 rounded-lg text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-white/10 transition-colors focus:outline-none"
            :title="localeStore.currentLang === 'zh-CN' ? '切换语言 / Language' : 'Switch Language'"
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
      </div>
    </header>

    <!-- 核心登录主卡片 (玻璃拟态 + 霓虹边框) -->
    <div class="login-card-container relative z-10 w-full max-w-[440px] px-4 py-8 sm:px-0">
      <div class="glass-card relative rounded-3xl p-8 sm:p-10 border border-slate-200/80 dark:border-white/10 shadow-[0_20px_50px_rgba(15,23,42,0.08),0_0_40px_rgba(59,130,246,0.06)] dark:shadow-[0_25px_60px_-15px_rgba(0,0,0,0.7),0_0_50px_rgba(59,130,246,0.15)] backdrop-blur-2xl bg-white/90 dark:bg-slate-900/80 transition-all duration-300">
        
        <!-- Logo 与系统标题 -->
        <div class="text-center mb-8">
          <div class="logo-box relative inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-tr from-blue-600 via-indigo-600 to-cyan-400 p-0.5 shadow-lg shadow-blue-500/25 mb-4 group">
            <div class="w-full h-full rounded-[14px] bg-white dark:bg-slate-950/80 backdrop-blur-sm flex items-center justify-center transition-colors">
              <el-icon class="text-blue-600 dark:text-white text-3xl icon-spin-pulse"><Refresh /></el-icon>
            </div>
            <!-- 外圈环绕呼吸光效 -->
            <div class="absolute -inset-1 rounded-2xl bg-gradient-to-r from-cyan-400 via-blue-500 to-purple-600 opacity-30 dark:opacity-40 blur-md -z-10 group-hover:opacity-75 transition-opacity duration-500"></div>
          </div>

          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-blue-600 via-indigo-600 to-cyan-500 dark:from-blue-300 dark:via-indigo-200 dark:to-cyan-300">
            JingSync
          </h1>
          
          <div class="mt-2.5">
            <span class="title inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-blue-50 dark:bg-blue-500/15 text-blue-600 dark:text-blue-300 border border-blue-200 dark:border-blue-400/25 tracking-wide shadow-xs transition-colors">
              <span class="w-1.5 h-1.5 rounded-full bg-cyan-500 animate-ping"></span>
              {{ t('login.badge') }}
            </span>
          </div>
        </div>

        <!-- 登录表单 -->
        <el-form
          :model="loginForm"
          :rules="loginRules"
          ref="loginFormRef"
          @keyup.enter="handleLogin"
          class="space-y-4"
        >
          <!-- 用户名输入框 -->
          <el-form-item prop="username" class="!mb-4">
            <div class="w-full">
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1.5 ml-1 transition-colors">
                {{ t('login.username') }}
              </label>
              <el-input
                v-model="loginForm.username"
                :placeholder="t('login.username_placeholder')"
                :prefix-icon="User"
                size="large"
                class="custom-tech-input"
              />
            </div>
          </el-form-item>

          <!-- 密码输入框 -->
          <el-form-item prop="password" class="!mb-6">
            <div class="w-full">
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1.5 ml-1 transition-colors">
                {{ t('login.password') }}
              </label>
              <el-input
                v-model="loginForm.password"
                type="password"
                :placeholder="t('login.password_placeholder')"
                :prefix-icon="Lock"
                show-password
                size="large"
                class="custom-tech-input"
              />
            </div>
          </el-form-item>

          <!-- 登录按钮 -->
          <el-form-item class="!mb-0">
            <button
              type="button"
              class="login-btn relative w-full h-12 rounded-xl text-white font-semibold text-base tracking-wider overflow-hidden transition-all duration-300 flex items-center justify-center gap-2 group cursor-pointer shadow-lg shadow-blue-600/30 active:scale-[0.98]"
              :disabled="loading"
              @click="handleLogin"
            >
              <span v-if="!loading" class="relative z-10 flex items-center gap-2">
                <span>{{ t('login.submit') }}</span>
                <el-icon class="group-hover:translate-x-1 transition-transform"><ArrowRight /></el-icon>
              </span>
              <span v-else class="relative z-10 flex items-center gap-2">
                <el-icon class="is-loading"><Loading /></el-icon>
                <span>{{ t('login.submitting') }}</span>
              </span>
              <!-- 按钮高光掠影扫光动画效果 -->
              <span class="btn-shine"></span>
            </button>
          </el-form-item>
        </el-form>

        <!-- 底部科技特性标签 -->
        <div class="mt-8 pt-6 border-t border-slate-100 dark:border-slate-800/80 flex items-center justify-between text-[11px] text-slate-500 dark:text-slate-400 transition-colors">
          <span class="flex items-center gap-1 hover:text-blue-600 dark:hover:text-cyan-400 transition-colors">
            <span class="text-blue-500">⚡</span> {{ t('login.feature_sync') }}
          </span>
          <span class="flex items-center gap-1 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">
            <span class="text-indigo-500">🔒</span> {{ t('login.feature_jwt') }}
          </span>
          <span class="flex items-center gap-1 hover:text-blue-600 dark:hover:text-cyan-400 transition-colors">
            <span class="text-cyan-500">⏰</span> {{ t('login.feature_cron') }}
          </span>
        </div>
      </div>

      <!-- 底部版权与版本提示 -->
      <div class="text-center mt-6 text-xs text-slate-400 dark:text-slate-500 font-mono tracking-wider transition-colors">
        {{ t('login.copyright') }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/index'
import { useLocaleStore } from '@/stores/localeStore'
import { langs, type Language } from '@/types/index'
import { ElMessage } from 'element-plus'
import ThemeSwitch from '@/components/ThemeSwitch.vue'
import {
  User,
  Lock,
  Refresh,
  ArrowRight,
  Loading,
  ChatDotRound
} from '@element-plus/icons-vue'

const { t, locale } = useI18n()
const userStore = useUserStore()
const localeStore = useLocaleStore()

const loginForm = ref({
  username: '',
  password: ''
})

// 响应式表单校验规则，跟随语言即时更新
const loginRules = computed(() => ({
  username: [
    { required: true, message: t('login.username_required'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.password_required'), trigger: 'blur' }
  ]
}))

const loginFormRef = ref()
const loading = ref(false)

function handleLanguageChange(command: Language) {
  localeStore.setLanguage(command)
  locale.value = command
}

const handleLogin = async () => {
  try {
    await loginFormRef.value.validate()
    loading.value = true
    await userStore.login(loginForm.value)
    ElMessage.success(t('login.success'))
    setTimeout(() => {
      window.location.reload()
    }, 300)
  } catch (error) {
    ElMessage.error(`${error}`)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 赛博网格背景 (深浅色模式自适应) */
.cyber-grid {
  background-image: 
    linear-gradient(to right, rgba(15, 23, 42, 0.05) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(15, 23, 42, 0.05) 1px, transparent 1px);
  background-size: 40px 40px;
  mask-image: radial-gradient(circle at center, black 40%, transparent 85%);
  -webkit-mask-image: radial-gradient(circle at center, black 40%, transparent 85%);
}

:global(html.dark) .cyber-grid {
  background-image: 
    linear-gradient(to right, rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
}

/* 动态环境光晕 */
.ambient-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
  opacity: 0.55;
  transition: opacity 0.5s ease;
}

:global(html.dark) .ambient-glow {
  opacity: 0.7;
}

.glow-1 {
  width: 480px;
  height: 480px;
  top: -120px;
  left: -100px;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.18) 0%, rgba(37, 99, 235, 0.08) 50%, transparent 70%);
  animation: floatOrb1 16s ease-in-out infinite alternate;
}

:global(html.dark) .glow-1 {
  background: radial-gradient(circle, rgba(56, 189, 248, 0.35) 0%, rgba(37, 99, 235, 0.15) 50%, transparent 70%);
}

.glow-2 {
  width: 520px;
  height: 520px;
  bottom: -150px;
  right: -120px;
  background: radial-gradient(circle, rgba(147, 51, 234, 0.15) 0%, rgba(99, 102, 241, 0.06) 50%, transparent 70%);
  animation: floatOrb2 20s ease-in-out infinite alternate;
}

:global(html.dark) .glow-2 {
  background: radial-gradient(circle, rgba(168, 85, 247, 0.3) 0%, rgba(99, 102, 241, 0.15) 50%, transparent 70%);
}

.glow-3 {
  width: 380px;
  height: 380px;
  top: 40%;
  left: 55%;
  background: radial-gradient(circle, rgba(20, 184, 166, 0.12) 0%, transparent 70%);
  animation: pulseGlow 12s ease-in-out infinite alternate;
}

:global(html.dark) .glow-3 {
  background: radial-gradient(circle, rgba(45, 212, 191, 0.22) 0%, transparent 70%);
}

/* 光晕浮动关键帧 */
@keyframes floatOrb1 {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(80px, 60px) scale(1.1); }
  100% { transform: translate(30px, 120px) scale(0.95); }
}

@keyframes floatOrb2 {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-90px, -70px) scale(1.15); }
  100% { transform: translate(-40px, -130px) scale(0.9); }
}

@keyframes pulseGlow {
  0% { transform: translate(-50%, -50%) scale(0.85); opacity: 0.3; }
  50% { transform: translate(-40%, -45%) scale(1.15); opacity: 0.65; }
  100% { transform: translate(-60%, -55%) scale(0.95); opacity: 0.4; }
}

/* 主卡片入场动画 */
.login-card-container {
  animation: cardEntrance 0.7s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes cardEntrance {
  from {
    opacity: 0;
    transform: translateY(24px) scale(0.96);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 旋转同步脉冲 */
.icon-spin-pulse {
  animation: spinSlow 12s linear infinite;
}

@keyframes spinSlow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 输入框定制样式：浅色模式 */
:deep(.custom-tech-input .el-input__wrapper) {
  background-color: rgba(248, 250, 252, 0.95) !important;
  border: 1px solid rgba(203, 213, 225, 0.8) !important;
  box-shadow: none !important;
  border-radius: 12px !important;
  padding: 4px 12px !important;
  transition: all 0.25s ease;
}

:deep(.custom-tech-input .el-input__wrapper:hover) {
  border-color: rgba(59, 130, 246, 0.5) !important;
  background-color: #ffffff !important;
}

:deep(.custom-tech-input .el-input__wrapper.is-focus) {
  border-color: #3b82f6 !important;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15) !important;
  background-color: #ffffff !important;
}

:deep(.custom-tech-input .el-input__inner) {
  color: #0f172a !important;
}

:deep(.custom-tech-input .el-input__inner::placeholder) {
  color: #94a3b8 !important;
}

/* 输入框定制样式：深色模式 */
:global(html.dark) :deep(.custom-tech-input .el-input__wrapper) {
  background-color: rgba(15, 23, 42, 0.65) !important;
  border: 1px solid rgba(255, 255, 255, 0.12) !important;
}

:global(html.dark) :deep(.custom-tech-input .el-input__wrapper:hover) {
  border-color: rgba(96, 165, 250, 0.45) !important;
  background-color: rgba(15, 23, 42, 0.85) !important;
}

:global(html.dark) :deep(.custom-tech-input .el-input__wrapper.is-focus) {
  border-color: #38bdf8 !important;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2) !important;
  background-color: rgba(15, 23, 42, 0.95) !important;
}

:global(html.dark) :deep(.custom-tech-input .el-input__inner) {
  color: #f8fafc !important;
}

:global(html.dark) :deep(.custom-tech-input .el-input__inner::placeholder) {
  color: #64748b !important;
}

/* 渐变流光登录按钮 */
.login-btn {
  background: linear-gradient(135deg, #2563eb 0%, #4f46e5 50%, #06b6d4 100%);
  border: none;
}

.login-btn:hover {
  background: linear-gradient(135deg, #1d4ed8 0%, #4338ca 50%, #0891b2 100%);
  box-shadow: 0 10px 25px -5px rgba(37, 99, 235, 0.4), 0 0 20px rgba(6, 182, 212, 0.3);
}

/* 按钮扫光光斑 */
.btn-shine {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(
    60deg,
    transparent 20%,
    rgba(255, 255, 255, 0.25) 50%,
    transparent 80%
  );
  transform: rotate(25deg) translateY(-100%);
  transition: transform 0.8s ease;
}

.login-btn:hover .btn-shine {
  transform: rotate(25deg) translateY(100%);
}
</style>
