<template>
  <el-drawer
    v-model="visible"
    title="同步任务执行日志详情"
    size="560px"
    class="rounded-l-2xl"
    @close="handleClose"
  >
    <div v-if="log" class="space-y-5">
      <!-- 状态摘要横幅 -->
      <div
        class="p-4 rounded-xl border flex items-center justify-between"
        :class="log.status === 1 ? 'bg-green-50/70 border-green-200 dark:bg-green-950/20 dark:border-green-900/50' : 'bg-red-50/70 border-red-200 dark:bg-red-950/20 dark:border-red-900/50'"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-full flex items-center justify-center text-lg"
            :class="log.status === 1 ? 'bg-green-100 text-green-600 dark:bg-green-900/40 dark:text-green-400' : 'bg-red-100 text-red-600 dark:bg-red-900/40 dark:text-red-400'"
          >
            <el-icon><Check v-if="log.status === 1" /><Close v-else /></el-icon>
          </div>
          <div>
            <div class="font-semibold text-sm text-slate-800 dark:text-slate-100">
              {{ parsedResult?.job_remark || `任务 #${log.job_id}` }}
            </div>
            <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              {{ log.status === 1 ? '任务全部执行完毕' : (log.error_msg || '任务执行过程中发生异常') }}
            </div>
          </div>
        </div>
        <el-tag :type="log.status === 1 ? 'success' : 'danger'" size="small">
          {{ log.status === 1 ? '执行成功' : '执行失败' }}
        </el-tag>
      </div>

      <!-- 统计指标网格 -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5">
        <div class="p-3 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800 text-center">
          <div class="text-xs text-slate-400">发现文件</div>
          <div class="text-lg font-bold text-slate-800 dark:text-slate-200 mt-0.5">
            {{ parsedResult?.files_found ?? '-' }}
          </div>
        </div>
        <div class="p-3 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800 text-center">
          <div class="text-xs text-green-600 dark:text-green-400">成功新增</div>
          <div class="text-lg font-bold text-green-600 dark:text-green-400 mt-0.5">
            {{ log.files_added }}
          </div>
        </div>
        <div class="p-3 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800 text-center">
          <div class="text-xs text-amber-600 dark:text-amber-400">过滤跳过</div>
          <div class="text-lg font-bold text-amber-600 dark:text-amber-400 mt-0.5">
            {{ parsedResult?.files_skipped ?? '-' }}
          </div>
        </div>
        <div class="p-3 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800 text-center">
          <div class="text-xs text-red-600 dark:text-red-400">传输失败</div>
          <div class="text-lg font-bold text-red-600 dark:text-red-400 mt-0.5">
            {{ parsedResult?.files_failed ?? (log.status === 0 ? 1 : 0) }}
          </div>
        </div>
      </div>

      <!-- 任务时间与路径信息 -->
      <div class="text-xs space-y-2 p-3 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800">
        <div class="flex justify-between">
          <span class="text-slate-400">开始时间:</span>
          <span class="font-mono text-slate-700 dark:text-slate-300">{{ formatTime(log.start_time) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">完成时间:</span>
          <span class="font-mono text-slate-700 dark:text-slate-300">{{ formatTime(log.end_time) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">执行耗时:</span>
          <span class="font-mono text-slate-700 dark:text-slate-300">{{ calculateDuration(log.start_time, log.end_time) }}</span>
        </div>
        <div v-if="parsedResult?.src_path" class="flex justify-between">
          <span class="text-slate-400">源目录路径:</span>
          <code class="text-slate-700 dark:text-slate-300 font-mono">{{ parsedResult.src_path }}</code>
        </div>
        <div v-if="parsedResult?.dst_path" class="flex justify-between">
          <span class="text-slate-400">目标目录路径:</span>
          <code class="text-slate-700 dark:text-slate-300 font-mono">{{ parsedResult.dst_path }}</code>
        </div>
      </div>

      <!-- 详细步骤日志 (Terminal 风格) -->
      <div class="space-y-1.5">
        <div class="text-xs font-semibold text-slate-700 dark:text-slate-300 flex items-center justify-between">
          <span>执行过程输出明细</span>
          <span class="text-slate-400 text-xxs font-normal">共 {{ detailLogs.length }} 条记录</span>
        </div>
        <div class="p-3 bg-slate-900 dark:bg-black rounded-lg border border-slate-800 font-mono text-xs text-slate-200 max-h-80 overflow-y-auto space-y-1 select-text">
          <div
            v-for="(item, idx) in detailLogs"
            :key="idx"
            class="leading-relaxed break-all"
            :class="{
              'text-red-400': item.includes('失败') || item.includes('错误') || item.includes('failed'),
              'text-green-400': item.includes('成功') || item.includes('新增'),
              'text-amber-300': item.includes('跳过') || item.includes('排除'),
              'text-blue-300': item.includes('扫描') || item.includes('开始'),
            }"
          >
            <span class="text-slate-500 mr-2 select-none">[{{ idx + 1 }}]</span>
            <span>{{ item }}</span>
          </div>
          <div v-if="detailLogs.length === 0" class="text-slate-500 py-4 text-center">
            暂无详细单条操作日志
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import type { SyncLog } from "@/types";
import { Check, Close } from "@element-plus/icons-vue";

const props = defineProps<{
  modelValue: boolean;
  log: SyncLog | null;
}>();

const emit = defineEmits(["update:modelValue"]);

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit("update:modelValue", val),
});

const parsedResult = computed(() => {
  if (!props.log?.details) return null;
  try {
    return JSON.parse(props.log.details);
  } catch {
    return null;
  }
});

const detailLogs = computed<string[]>(() => {
  if (parsedResult.value?.details && Array.isArray(parsedResult.value.details)) {
    return parsedResult.value.details;
  }
  return [];
});

function formatTime(isoStr?: string) {
  if (!isoStr) return "-";
  return new Date(isoStr).toLocaleString();
}

function calculateDuration(start?: string, end?: string) {
  if (!start || !end) return "-";
  const s = new Date(start).getTime();
  const e = new Date(end).getTime();
  const diffSec = Math.max(0, Math.round((e - s) / 1000));
  if (diffSec < 60) return `${diffSec} 秒`;
  const min = Math.floor(diffSec / 60);
  const sec = diffSec % 60;
  return `${min} 分 ${sec} 秒`;
}

function handleClose() {
  emit("update:modelValue", false);
}
</script>
