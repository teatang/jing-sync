<template>
  <div class="p-[10px]">
    <div class="mb-[10px]">
      <el-button type="primary" @click="handleAdd">新增任务</el-button>
    </div>

    <el-table :data="jobs" style="width: 100%">
      <el-table-column type="expand">
        <template #default="props">
          <el-row class="p-[12px]">
            <el-col :span="4"
              >method:
              {{
                getSelectOptionName(props.row.speed, JOB_METHOD_CONFIGS)
              }}</el-col
            >
            <el-col :span="4"
              >speed:
              {{
                getSelectOptionName(props.row.speed, JOB_SPEED_CONFIGS)
              }}</el-col
            >
            <el-col :span="4"
              >method:
              {{
                getSelectOptionName(props.row.is_cron, JOB_IS_CRON_CONFIGS)
              }}</el-col
            >
            <el-col :span="4"
              >year: {{ props.row.year ? props.row.year : "-" }}</el-col
            >
            <el-col :span="4"
              >month: {{ props.row.month ? props.row.month : "-" }}</el-col
            >
            <el-col :span="4"
              >day: {{ props.row.day ? props.row.day : "-" }}</el-col
            >
          </el-row>
          <el-row class="p-[12px]">
            <el-col :span="4"
              >hour: {{ props.row.hour ? props.row.hour : "-" }}</el-col
            >
            <el-col :span="4"
              >minute: {{ props.row.minute ? props.row.minute : "-" }}</el-col
            >
            <el-col :span="4"
              >second: {{ props.row.second ? props.row.second : "-" }}</el-col
            >
            <el-col :span="4"
              >week: {{ props.row.week ? props.row.week : "-" }}</el-col
            >
            <el-col :span="4"
              >day_of_week:
              {{ props.row.day_of_week ? props.row.day_of_week : "-" }}</el-col
            >
            <el-col :span="4"
              >start_date:
              {{
                props.row.start_date
                  ? new Date(props.row.start_date)
                      .toLocaleString()
                      .substring(0, 10)
                  : "-"
              }}</el-col
            >
          </el-row>
          <el-row class="p-[12px]">
            <el-col :span="4"
              >end_date:
              {{
                props.row.end_date
                  ? new Date(props.row.end_date)
                      .toLocaleString()
                      .substring(0, 10)
                  : "-"
              }}</el-col
            >
            <el-col :span="4"
              >create_time:
              {{
                props.row.create_time
                  ? new Date(props.row.create_time).toLocaleString()
                  : "-"
              }}</el-col
            >
            <el-col :span="4"
              >update_time:
              {{
                props.row.update_time
                  ? new Date(props.row.update_time).toLocaleString()
                  : "-"
              }}</el-col
            >
          </el-row>
        </template>
      </el-table-column>
      <el-table-column label="状态">
        <template #default="scope">
          <el-text :type="scope.row.status == 1 ? 'success' : 'danger'">
            {{ getSelectOptionName(scope.row.status, JOB_IS_ENABLE_CONFIG) }}
          </el-text>
        </template>
      </el-table-column>
      <el-table-column prop="src_path" label="来源目录" />
      <el-table-column prop="dst_path" label="目标目录" />
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
    <div class="mt-[10px]">
      <el-pagination
        background
        v-model:current-page="pagination.page"
        :page-size="pagination.size"
        :total="pagination.total"
        @current-change="fetchJobs"
      />
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
import RestClient from "@/utils/rest-client";
import type {
  Job,
  Pagination,
  InfoList,
  JobSelectOptionConfig,
} from "@/types/index";
import CronJobDialog from "@/views/cron-job/CronJobDialog.vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  JOB_IS_ENABLE_CONFIG,
  JOB_SPEED_CONFIGS,
  JOB_METHOD_CONFIGS,
  JOB_IS_CRON_CONFIGS,
} from "@/const/index";

const jobs = ref<Job[]>([]);
const currentJob = ref<Partial<Job>>({});
const dialogVisible = ref(false);
const dialogMode = ref<"add" | "edit">("add");
const pagination = ref<Pagination>({
  page: 1,
  size: 10,
  total: 0,
});

const getSelectOptionName = (
  speedId: number,
  configs: JobSelectOptionConfig[]
) => {
  return configs.find((item) => item.id === speedId)?.name || "";
};

const fetchJobs = async () => {
  try {
    const infoList = await new RestClient().get<InfoList<Job>>("/job", {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    jobs.value = infoList.list as Job[];
    pagination.value = infoList.pagination as Pagination;
  } catch (error) {
    ElMessage.error(`${error}`);
  }
};

const handleAdd = () => {
  currentJob.value = {};
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
    await ElMessageBox.confirm("确认删除该数据吗？", "提示", {
      type: "warning",
    });
    await new RestClient().delete<Job[]>("/job", job);
    ElMessage.success("删除成功");
    fetchJobs();
  } catch (error) {
    ElMessage.error(`${error}`);
  }
};

onMounted(async () => {
  fetchJobs();
});
</script>

<style scoped></style>
