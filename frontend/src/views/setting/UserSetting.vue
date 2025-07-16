<template>
  <div class="t-management">
    <div class="t-management-header">
      <el-button type="primary" @click="handleAdd">新增用户</el-button>
    </div>

    <el-table :data="users" style="width: 100%">
      <el-table-column prop="username" label="用户名" />1
      <el-table-column prop="password" label="密码" />
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
        @current-change="fetchUsers"
      />
    </div>
  </div>

  <user-setting-dialog
    :modelValue="dialogVisible"
    :user="currentUser"
    :mode="dialogMode"
    @success="fetchUsers"
    @update:modelValue="(val) => (dialogVisible = val)"
  />
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import RestClient from "@/utils/rest-client";
import type { User, Pagination, InfoList } from "@/types/index";
import UserSettingDialog from "@/views/setting/UserSettingDialog.vue";
import { ElMessage, ElMessageBox } from "element-plus";

const users = ref<User[]>([]);
const currentUser = ref<Partial<User>>({});
const dialogVisible = ref(false);
const dialogMode = ref<"add" | "edit">("add");

const fetchUsers = async () => {
  try {
    const infoList = await new RestClient().get<InfoList<User>>("/user", {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    users.value = infoList.list as User[];
    pagination.value = infoList.pagination as Pagination;
  } catch (error) {
    console.error("获取engines失败:", error);
  }
};

const handleAdd = () => {
  currentUser.value = {};
  dialogMode.value = "add";
  dialogVisible.value = true;
};

const handleEdit = (user: User) => {
  currentUser.value = { ...user };
  dialogMode.value = "edit";
  dialogVisible.value = true;
};

const handleDelete = async (user: User) => {
  try {
    await ElMessageBox.confirm("确认删除该数据吗？", "提示", {
      type: "warning",
    });
    await new RestClient().delete<User[]>("/user", user);
    ElMessage.success("删除成功");
    fetchUsers();
  } catch (error) {
    console.error("delete engine 失败:", error);
    ElMessage.error(`取消删除`);
  }
};

const pagination = ref<Pagination>({
  page: 1,
  size: 10,
  total: 0,
});

onMounted(async () => {
  fetchUsers();
});
</script>

<style scoped></style>
