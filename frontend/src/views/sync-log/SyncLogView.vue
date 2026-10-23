<template>
  <div class="space-y-4">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between bg-white dark:bg-slate-900 p-4 rounded-xl border border-slate-200 dark:border-slate-800 shadow-xs">
      <div>
        <div class="flex items-center gap-2">
          <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">同步执行日志</h2>
          <el-tag v-if="filteredJobId" type="info" closable size="small" @close="clearJobFilter">
            过滤任务: #{{ filteredJobId }}
          </el-tag>
        </div>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">记录每次定时触发或手动执行的完整同步记录、传输文件数与耗时明细</p>
      </div>

      <div class="flex items-center gap-2">
        <el-button :icon="Refresh" :loading="loading" @click="fetchLogs">
          刷新
        </el-button>
      </div>
    </div>

    <!-- 日志列表卡片 -->
    <div class="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-4 shadow-xs">
      <el-table
        :data="logs"
        v-loading="loading"
        style="width: 100%"
        class="!bg-transparent"
      >
        <el-table-column prop="id" label="ID" width="80" align="center">
          <template #default="{ row }">
            <span class="font-mono text-xs text-slate-400">#{{ row.id }}</span>
          </template>
        </el-table-column>

        <el-table-column label="所属任务" min-width="160">
          <template #default="{ row }">
            <div class="font-medium text-slate-800 dark:text-slate-200">
              {{ parseJobRemark(row) }}
            </div>
            <div class="text-xs text-slate-400 font-mono">
              任务 ID: {{ row.job_id }}
            </div>
          </template>
        </el-table-column>

        <el-table-column label="执行状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="files_added" label="新增传输文件" width="120" align="center">
          <template #default="{ row }">
            <span class="font-semibold text-green-600 dark:text-green-400">
              {{ row.files_added }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="执行时间" min-width="160">
          <template #default="{ row }">
            <div class="text-xs text-slate-700 dark:text-slate-300 font-mono">
              {{ row.start_time ? new Date(row.start_time).toLocaleString() : '-' }}
            </div>
            <div class="text-xxs text-slate-400">
              耗时: {{ calculateDuration(row.start_time, row.end_time) }}
            </div>
          </template>
        </el-table-column>

        <el-table-column label="错误信息" min-width="180">
          <template #default="{ row }">
            <span v-if="row.error_msg" class="text-xs text-red-500 line-clamp-1" :title="row.error_msg">
              {{ row.error_msg }}
            </span>
            <span v-else class="text-xs text-slate-400">-</span>
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" plain @click="handleOpenDetail(row)">
              查看详情
            </el-button>
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
          @current-change="fetchLogs"
        />
      </div>
    </div>
  </div>

  <!-- 详情抽屉 -->
  <SyncLogDetailDrawer
    v-model="drawerVisible"
    :log="selectedLog"
  />
</template>

<script lang="ts" setup>
import { ref, onMounted, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import RestClient from "@/utils/rest-client";
import type { SyncLog, Pagination, InfoList } from "@/types";
import SyncLogDetailDrawer from "./SyncLogDetailDrawer.vue";
import { ElMessage } from "element-plus";
import { Refresh } from "@element-plus/icons-vue";

const route = useRoute();
const router = useRouter();

const logs = ref<SyncLog[]>([]);
const loading = ref(false);
const drawerVisible = ref(false);
const selectedLog = ref<SyncLog | null>(null);

const pagination = ref<Pagination>({
  page: 1,
  size: 10,
  total: 0,
});

const filteredJobId = computed(() => {
  return route.query.job_id ? Number(route.query.job_id) : null;
});

function clearJobFilter() {
  router.push({ path: "/sync-logs" });
}

function parseJobRemark(row: SyncLog) {
  if (row.details) {
    try {
      const data = JSON.parse(row.details);
      if (data?.job_remark) return data.job_remark;
    } catch {}
  }
  return `任务 #${row.job_id}`;
}

function calculateDuration(start?: string, end?: string) {
  if (!start || !end) return "-";
  const s = new Date(start).getTime();
  const e = new Date(end).getTime();
  const diffSec = Math.max(0, Math.round((e - s) / 1000));
  if (diffSec < 60) return `${diffSec} 秒`;
  const min = Math.floor(diffSec / 60);
  const sec = diffSec % 60;
  return `${min}分${sec}秒`;
}

const fetchLogs = async () => {
  loading.value = true;
  try {
    let url = "/sync-logs";
    if (filteredJobId.value) {
      url = `/job/${filteredJobId.value}/sync-logs`;
    }
    const infoList = await new RestClient().get<InfoList<SyncLog>>(url, {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    logs.value = (infoList?.list || []) as SyncLog[];
    pagination.value = (infoList?.pagination || { page: 1, size: 10, total: 0 }) as Pagination;
  } catch (error) {
    ElMessage.error(`${error}`);
  } finally {
    loading.value = false;
  }
};

function handleOpenDetail(row: SyncLog) {
  selectedLog.value = row;
  drawerVisible.value = true;
}

watch(filteredJobId, () => {
  pagination.value.page = 1;
  fetchLogs();
});

onMounted(() => {
  fetchLogs();
});
</script>
