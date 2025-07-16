<template>
  <div class="t-management">
    <div class="t-management-header">
      <el-button type="primary" @click="handleAdd">新增引擎</el-button>
    </div>

    <el-table :data="engines" style="width: 100%">
      <el-table-column prop="url" label="地址" />
      <el-table-column prop="token" label="令牌" />
      <el-table-column prop="remark" label="备注" />
      <el-table-column label="创建时间">
        <template #default="scope">
          {{
            scope.row.create_time
              ? new Date(scope.row.create_time).toLocaleString()
              : "-"
          }}
        </template>
      </el-table-column>
      <el-table-column label="更新时间">
        <template #default="scope">
          {{
            scope.row.update_time
              ? new Date(scope.row.update_time).toLocaleString()
              : "-"
          }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" min-width="100">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            Edit
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.row)"
          >
            Delete
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="t-pagination">
      <el-pagination
        background
        v-model:current-page="pagination.page"
        :page-size="pagination.size"
        :total="pagination.total"
        @current-change="fetchEngines"
      />
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

const engines = ref<Engine[]>([]);
const currentEngine = ref<Partial<Engine>>({});
const dialogVisible = ref(false);
const dialogMode = ref<"add" | "edit">("add");
const pagination = ref<Pagination>({
  page: 1,
  size: 10,
  total: 0,
});

const fetchEngines = async () => {
  try {
    const infoList = await new RestClient().get<InfoList<Engine>>("/engine", {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    engines.value = infoList.list as Engine[];
    pagination.value = infoList.pagination as Pagination;
  } catch (error) {
    console.error("获取engines失败:", error);
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
    await ElMessageBox.confirm("确认删除该数据吗？", "提示", {
      type: "warning",
    });
    await new RestClient().delete<Engine[]>("/engine", engine);
    ElMessage.success("删除成功");
    fetchEngines();
  } catch (error) {
    console.error("delete engine 失败:", error);
    ElMessage.error(`取消删除`);
  }
};

onMounted(async () => {
  fetchEngines();
});
</script>

<style scoped></style>
