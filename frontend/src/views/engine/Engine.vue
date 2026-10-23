<template>
  <div class="space-y-4">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between bg-white dark:bg-slate-900 p-4 rounded-xl border border-slate-200 dark:border-slate-800 shadow-xs">
      <div>
        <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">存储引擎管理</h2>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">配置 OpenList / AList 实例地址与授权令牌，作为数据同步的数据源与目标端</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="handleAdd">
        新增引擎
      </el-button>
    </div>

    <!-- 列表数据卡片 -->
    <div class="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-4 shadow-xs">
      <el-table
        :data="engines"
        v-loading="loading"
        style="width: 100%"
        class="!bg-transparent"
      >
        <el-table-column prop="remark" label="备注说明" min-width="140">
          <template #default="{ row }">
            <span class="font-medium text-slate-800 dark:text-slate-200">{{ row.remark || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="url" label="服务地址" min-width="220">
          <template #default="{ row }">
            <code class="px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-xs text-slate-700 dark:text-slate-300 font-mono">
              {{ row.url }}
            </code>
          </template>
        </el-table-column>

        <el-table-column prop="token" label="访问令牌" min-width="160">
          <template #default="{ row }">
            <span class="text-xs text-slate-400 font-mono">
              {{ maskToken(row.token) }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="更新时间" min-width="160">
          <template #default="{ row }">
            <span class="text-xs text-slate-500">
              {{ row.update_time ? new Date(row.update_time).toLocaleString() : '-' }}
            </span>
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" min-width="200">
          <template #default="{ row }">
            <div class="flex items-center gap-1.5">
              <el-button
                size="small"
                type="success"
                plain
                :loading="testingId === row.id"
                @click="handleTestConnection(row)"
              >
                测试连接
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
          @current-change="fetchEngines"
        />
      </div>
    </div>
  </div>

  <engine-dialog
    :modelValue="dialogVisible"
    :engine="currentEngine"
    :mode="dialogMode"
    @success="fetchEngines"
    @update:modelValue="(val) => (dialogVisible = val)"
  />
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import RestClient from "@/utils/rest-client";
import type { Engine, Pagination, InfoList } from "@/types";
import EngineDialog from "@/views/engine/EngineDialog.vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus } from "@element-plus/icons-vue";

const engines = ref<Engine[]>([]);
const currentEngine = ref<Partial<Engine>>({});
const dialogVisible = ref(false);
const dialogMode = ref<"add" | "edit">("add");
const loading = ref(false);
const testingId = ref<number | null>(null);

const pagination = ref<Pagination>({
  page: 1,
  size: 10,
  total: 0,
});

function maskToken(token?: string) {
  if (!token) return "-";
  if (token.length <= 8) return "******";
  return token.substring(0, 4) + "..." + token.substring(token.length - 4);
}

const fetchEngines = async () => {
  loading.value = true;
  try {
    const infoList = await new RestClient().get<InfoList<Engine>>("/engine", {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    engines.value = (infoList?.list || []) as Engine[];
    pagination.value = (infoList?.pagination || { page: 1, size: 10, total: 0 }) as Pagination;
  } catch (error) {
    ElMessage.error(`${error}`);
  } finally {
    loading.value = false;
  }
};

const handleTestConnection = async (engine: Engine) => {
  testingId.value = engine.id ?? null;
  try {
    await new RestClient().post("/engine/test", { id: engine.id });
    ElMessage.success(`引擎 [${engine.remark || engine.url}] 连接测试成功！`);
  } catch (error: any) {
    ElMessage.error(`连接测试失败: ${error?.message || error}`);
  } finally {
    testingId.value = null;
  }
};

const handleAdd = () => {
  currentEngine.value = {};
  dialogMode.value = "add";
  dialogVisible.value = true;
};

const handleEdit = (engine: Engine) => {
  currentEngine.value = { ...engine };
  dialogMode.value = "edit";
  dialogVisible.value = true;
};

const handleDelete = async (engine: Engine) => {
  try {
    await ElMessageBox.confirm(`确认删除引擎 "${engine.remark || engine.url}" 吗？`, "删除警告", {
      type: "warning",
      confirmButtonText: "确认删除",
      cancelButtonText: "取消",
    });
    await new RestClient().delete<Engine[]>("/engine", engine);
    ElMessage.success("删除成功");
    fetchEngines();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error(`${error}`);
    }
  }
};

onMounted(() => {
  fetchEngines();
});
</script>
