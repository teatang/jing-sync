<template>
  <el-container v-loading="mainStore.loading" class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '200px'">
      <el-scrollbar>
        <el-menu router :collapse="isCollapse" default-active="/cron-job">
          <el-menu-item index="/cron-job">
            <el-icon><AlarmClock /></el-icon>
            <template #title>定时任务</template>
          </el-menu-item>
          <el-menu-item index="/engine">
            <el-icon><Connection /></el-icon>
            <template #title>引擎管理</template>
          </el-menu-item>
          <el-sub-menu index="setting">
            <template #title>
              <el-icon><setting /></el-icon>
              <span>设置</span>
            </template>
            <el-menu-item index="/user-setting">用户设置</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>
    </el-aside>
    <el-container>
      <el-header class="header-toolbar">
        <div class="header-toolbar-left">
          <el-button
            @click="isCollapse = !isCollapse"
            :icon="isCollapse ? 'Expand' : 'Fold'"
          >
          </el-button>
        </div>
        <div class="header-toolbar-right">
          <el-dropdown>
            <el-button>
              <el-icon>
                <setting />
              </el-icon>
              <span>Tea</span>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>View</el-dropdown-item>
                <el-dropdown-item>Add</el-dropdown-item>
                <el-dropdown-item>Delete</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import { useMainStore } from "@/stores";
import { ref } from "vue";

const mainStore = useMainStore();
const isCollapse = ref<boolean>(false);
</script>

<style scoped>
.el-aside {
  transition: width 0.3s ease-in-out;
}

.layout-container .el-header {
  position: relative;
  background-color: var(--el-color-primary-light-7);
  color: var(--el-text-color-primary);
}
.layout-container .el-aside {
  color: var(--el-text-color-primary);
  background: var(--el-color-primary-light-8);
}
.layout-container .el-menu {
  border-right: none;
}
.layout-container .el-main {
  padding: 0;
}
.header-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
