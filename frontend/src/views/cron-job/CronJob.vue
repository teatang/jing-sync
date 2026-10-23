<template>
  <div class="space-y-4">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between bg-white dark:bg-slate-900 p-4 rounded-xl border border-slate-200 dark:border-slate-800 shadow-xs">
      <div>
        <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">定时同步任务</h2>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">配置源与目标目录映射、差异同步规则、过滤排除项与自动调度计划</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="handleAdd">
        新建同步任务
      </el-button>
    </div>

    <!-- 任务列表卡片 -->
    <div class="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-4 shadow-xs">
      <el-table
        :data="jobs"
        v-loading="loading"
        style="width: 100%"
        class="!bg-transparent"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="p-4 bg-slate-50 dark:bg-slate-950/60 rounded-lg text-xs space-y-2 border border-slate-100 dark:border-slate-850">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <span class="text-slate-400">排除通配规则: </span>
                  <code class="text-slate-700 dark:text-slate-300">{{ row.excludes || '无 (全部同步)' }}</code>
                </div>
                <div>
                  <span class="text-slate-400">文件大小限制: </span>
                  <span class="text-slate-700 dark:text-slate-300">
                    {{ formatSizeLimit(row.min_size, row.max_size) }}
                  </span>
                </div>
                <div>
                  <span class="text-slate-400">扫描速率模式: </span>
                  <span class="text-slate-700 dark:text-slate-300">{{ getSpeedName(row.speed) }}</span>
                </div>
              </div>
              <div class="text-slate-400 pt-1 border-t border-slate-200/50 dark:border-slate-800/50">
                <span>更新时间: {{ row.update_time ? new Date(row.update_time).toLocaleString() : '-' }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="remark" label="任务名称" min-width="150">
          <template #default="{ row }">
            <div class="font-medium text-slate-800 dark:text-slate-200">
              {{ row.remark || '未命名任务' }}
            </div>
          </template>
        </el-table-column>

        <el-table-column label="同步模式" min-width="130">
          <template #default="{ row }">
            <el-tag :type="getSyncModeTagType(row.sync_mode)" size="small">
              {{ getSyncModeName(row.sync_mode) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="路径映射" min-width="260">
          <template #default="{ row }">
            <div class="flex items-center gap-1.5 text-xs font-mono">
              <span class="px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 truncate max-w-[120px]" :title="row.src_path">
                {{ row.src_path }}
              </span>
              <el-icon class="text-slate-400"><Right /></el-icon>
              <span class="px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 truncate max-w-[120px]" :title="row.dst_path">
                {{ row.dst_path }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="调度周期" min-width="130">
          <template #default="{ row }">
            <span class="text-xs text-slate-600 dark:text-slate-400">
              {{ getScheduleDesc(row) }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              :model-value="row.enabled === 1"
              :loading="togglingId === row.id"
              inline-prompt
              active-text="开"
              inactive-text="关"
              @change="() => handleToggle(row)"
            />
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" min-width="260">
          <template #default="{ row }">
            <div class="flex items-center gap-1.5">
              <el-button
                size="small"
                type="primary"
                plain
                :icon="RefreshRight"
                :loading="syncingId === row.id"
                @click="handleTriggerSync(row)"
              >
                立即同步
              </el-button>
              <el-button
                size="small"
                :icon="Document"
                @click="handleViewLogs(row)"
              >
                日志
              </el-button>
              <el-button
                size="small"
                @click="handleEdit(row)"
              >
                编辑
              </el-button>
              <el-button
                size="small"
                type="danger"
                plain
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页栏 -->
      <div class="mt-4 flex justify-end">
        <el-pagination
          background
          v-model:current-page="pagination.page"
          :page-size="pagination.size"
          :total="pagination.total"
          layout="total, prev, pager, next"
          @current-change="fetchJobs"
        />
      </div>
    </div>
  </div>

  <cron-job-dialog
    :modelValue="dialogVisible"
    :job="currentJob"
    :mode="dialogMode"
    @success="fetchJobs"
    @update:modelValue="(val) => (dialogVisible = val)"
  />
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import RestClient from "@/utils/rest-client";
import type { Job, Pagination, InfoList } from "@/types/index";
import CronJobDialog from "@/views/cron-job/CronJobDialog.vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, Right, RefreshRight, Document } from "@element-plus/icons-vue";
import { formatBytes } from "@/utils/format";
import { SYNC_MODE_CONFIGS, JOB_SPEED_CONFIGS } from "@/const/index";

const router = useRouter();
const jobs = ref<Job[]>([]);
const currentJob = ref<Partial<Job>>({});
const dialogVisible = ref(false);
const dialogMode = ref<"add" | "edit">("add");
const loading = ref(false);
const syncingId = ref<number | null>(null);
const togglingId = ref<number | null>(null);

const pagination = ref<Pagination>({
  page: 1,
  size: 10,
  total: 0,
});

function getSyncModeName(mode?: number) {
  const item = SYNC_MODE_CONFIGS.find((m) => m.id === mode);
  return item ? item.name : "仅新增";
}

function getSyncModeTagType(mode?: number) {
  if (mode === 0) return "warning";
  if (mode === 2) return "danger";
  return "success";
}

function getSpeedName(speed?: number) {
  const item = JOB_SPEED_CONFIGS.find((s) => s.id === speed);
  return item ? item.name : "标准速率";
}

function formatSizeLimit(min?: number | null, max?: number | null) {
  if (!min && !max) return "无限制";
  if (min && !max) return `>= ${formatBytes(min)}`;
  if (!min && max) return `<= ${formatBytes(max)}`;
  return `${formatBytes(min!)} ~ ${formatBytes(max!)}`;
}

function getScheduleDesc(job: Job) {
  if (job.is_cron === 2) return "仅手动触发";
  if (job.is_cron === 1) return job.cron_expr ? `Cron: ${job.cron_expr}` : "Cron定时";
  return `每 ${job.interval || 60} 分钟`;
}

const fetchJobs = async () => {
  loading.value = true;
  try {
    const infoList = await new RestClient().get<InfoList<Job>>("/job", {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    jobs.value = (infoList?.list || []) as Job[];
    pagination.value = (infoList?.pagination || { page: 1, size: 10, total: 0 }) as Pagination;
  } catch (error) {
    ElMessage.error(`${error}`);
  } finally {
    loading.value = false;
  }
};

const handleToggle = async (job: Job) => {
  togglingId.value = job.id!;
  try {
    await new RestClient().put(`/job/${job.id}/toggle`, {});
    job.enabled = job.enabled === 1 ? 0 : 1;
    ElMessage.success(job.enabled === 1 ? "已启用任务" : "已暂停任务");
  } catch (error: any) {
    ElMessage.error(`切换状态失败: ${error?.message || error}`);
  } finally {
    togglingId.value = null;
  }
};

const handleTriggerSync = async (job: Job) => {
  syncingId.value = job.id!;
  try {
    const res: any = await new RestClient().post(`/job/${job.id}/sync`, {});
    const added = res?.files_added ?? 0;
    ElMessage.success(`同步成功！新增/传输文件数: ${added}`);
  } catch (error: any) {
    ElMessage.error(`同步执行失败: ${error?.message || error}`);
  } finally {
    syncingId.value = null;
  }
};

const handleViewLogs = (job: Job) => {
  router.push({ path: "/sync-logs", query: { job_id: job.id } });
};

const handleAdd = () => {
  currentJob.value = {
    enabled: 1,
    sync_mode: 1,
    speed: 0,
    interval: 60,
    is_cron: 0,
  };
  dialogMode.value = "add";
  dialogVisible.value = true;
};

const handleEdit = (job: Job) => {
  currentJob.value = { ...job };
  dialogMode.value = "edit";
  dialogVisible.value = true;
};

const handleDelete = async (job: Job) => {
  try {
    await ElMessageBox.confirm(`确认删除同步任务 "${job.remark || job.src_path}" 吗？`, "删除警告", {
      type: "warning",
      confirmButtonText: "确认删除",
      cancelButtonText: "取消",
    });
    await new RestClient().delete<Job[]>("/job", job);
    ElMessage.success("删除成功");
    fetchJobs();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error(`${error}`);
    }
  }
};

onMounted(() => {
  fetchJobs();
});
</script>
