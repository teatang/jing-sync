<template>
  <div class="space-y-4">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between bg-white dark:bg-slate-900 p-4 rounded-xl border border-slate-200 dark:border-slate-800 shadow-xs">
      <div>
        <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">用户账号管理</h2>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">管理系统的登录用户与访问凭证权限</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="handleAdd">
        新增用户
      </el-button>
    </div>

    <!-- 列表数据卡片 -->
    <div class="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-4 shadow-xs">
      <el-table
        :data="users"
        v-loading="loading"
        style="width: 100%"
        class="!bg-transparent"
      >
        <el-table-column prop="username" :label="t('page.user.username')" min-width="160">
          <template #default="{ row }">
            <span class="font-medium text-slate-800 dark:text-slate-200">{{ row.username }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="password" :label="t('page.user.password')" min-width="160">
          <template #default>
            <span class="text-xs text-slate-400 font-mono">••••••••</span>
          </template>
        </el-table-column>

        <el-table-column :label="t('page.create_time')" min-width="180">
          <template #default="{ row }">
            <span class="text-xs text-slate-500">
              {{ row.create_time ? new Date(row.create_time).toLocaleString() : "-" }}
            </span>
          </template>
        </el-table-column>

        <el-table-column :label="t('page.update_time')" min-width="180">
          <template #default="{ row }">
            <span class="text-xs text-slate-500">
              {{ row.update_time ? new Date(row.update_time).toLocaleString() : "-" }}
            </span>
          </template>
        </el-table-column>

        <el-table-column fixed="right" :label="t('page.handle')" min-width="160">
          <template #default="{ row }">
            <div class="flex items-center gap-1.5">
              <el-button size="small" @click="handleEdit(row)">
                {{ t('page.handle_edit') }}
              </el-button>
              <el-button
                size="small"
                type="danger"
                plain
                @click="handleDelete(row)"
              >
                {{ t('page.handle_delete') }}
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
          @current-change="fetchUsers"
        />
      </div>
    </div>
  </div>

  <user-setting-dialog
    v-model="dialogVisible"
    :user="currentUser"
    :mode="dialogMode"
    @success="fetchUsers"
  />
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import RestClient from "@/utils/rest-client";
import type { User, Pagination, InfoList } from "@/types/index";
import UserSettingDialog from "@/views/setting/UserSettingDialog.vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import { useI18n } from 'vue-i18n'

const { t } = useI18n();
const users = ref<User[]>([]);
const currentUser = ref<Partial<User>>({});
const dialogVisible = ref(false);
const dialogMode = ref<"add" | "edit">("add");
const loading = ref(false);

const fetchUsers = async () => {
  loading.value = true;
  try {
    const infoList = await new RestClient().get<InfoList<User>>("/user", {
      page: pagination.value.page,
      size: pagination.value.size,
    });
    users.value = infoList.list as User[];
    pagination.value = infoList.pagination as Pagination;
  } catch (error) {
    ElMessage.error(`${error}`);
  } finally {
    loading.value = false;
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
    if (error !== "cancel") {
      ElMessage.error(`${error}`);
    }
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
