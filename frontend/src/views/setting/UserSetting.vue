<template>
  <div class="p-[10px]">
    <div class="mb-[10px]">
      <el-button type="primary" @click="handleAdd">新增用户</el-button>
    </div>

    <el-table :data="users" style="width: 100%">
      <el-table-column prop="username" :label="t('page.user.username')" />1
      <el-table-column prop="password" :label="t('page.user.password')" />
      <el-table-column :label="t('page.create_time')">
        <template #default="scope">
          {{
            scope.row.create_time
              ? new Date(scope.row.create_time).toLocaleString()
              : "-"
          }}
        </template>
      </el-table-column>
      <el-table-column :label="t('page.update_time')">
        <template #default="scope">
          {{
            scope.row.update_time
              ? new Date(scope.row.update_time).toLocaleString()
              : "-"
          }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" :label="t('page.handle')" min-width="100">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            {{ t('page.handle_edit') }}
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.row)"
          >
            {{ t('page.handle_delete') }}
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
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
    ElMessage.error(`${error}`);
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
    await ElMessageBox.confirm(t("msg.delete_confirm_msg"), t("msg.delete_confirm_msg_title"), {
      type: "warning",
    });
    await new RestClient().delete<User[]>("/user", user);
    ElMessage.success(t("msg.delete_confirm_msg_success"));
    fetchUsers();
  } catch (error) {
    ElMessage.error(`${error}`);
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
