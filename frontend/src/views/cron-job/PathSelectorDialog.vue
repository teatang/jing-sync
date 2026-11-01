<template>
  <el-dialog
    v-model="visible"
    :title="title || '选择路径'"
    width="600px"
    class="rounded-xl"
    @open="handleOpen"
  >
    <div class="space-y-4">
      <!-- 当前路径与导航 -->
      <div class="flex items-center space-x-2">
        <el-button
          type="default"
          size="default"
          :disabled="currentBrowsingPath === '/'"
          @click="goParent"
        >
          <el-icon class="mr-1"><ArrowUp /></el-icon>
          上一级
        </el-button>
        <el-input
          v-model="currentBrowsingPath"
          placeholder="/ 输入或浏览路径"
          clearable
          @keyup.enter="loadDirs(currentBrowsingPath)"
        >
          <template #append>
            <el-button @click="loadDirs(currentBrowsingPath)">跳转</el-button>
          </template>
        </el-input>
      </div>

      <!-- 目录列表区域 -->
      <div
        v-loading="loading"
        class="border border-slate-200 dark:border-slate-800 rounded-lg h-72 overflow-y-auto p-2 bg-slate-50 dark:bg-slate-900"
      >
        <div v-if="dirs.length === 0 && !loading" class="text-center text-slate-400 py-20">
          暂无子目录
        </div>
        <div
          v-for="dir in dirs"
          :key="dir"
          class="flex items-center justify-between p-2 hover:bg-slate-200 dark:hover:bg-slate-800 rounded cursor-pointer transition-colors"
          @dblclick="selectSubDir(dir)"
        >
          <div class="flex items-center space-x-2 truncate" @click="selectSubDir(dir)">
            <el-icon class="text-amber-500"><Folder /></el-icon>
            <span class="text-sm text-slate-700 dark:text-slate-200">{{ dir }}</span>
          </div>
          <el-button size="small" type="primary" link @click.stop="selectSubDir(dir)">
            进入
          </el-button>
        </div>
      </div>

      <div class="text-xs text-slate-500">
        当前选中路径: <span class="font-mono text-primary font-semibold">{{ currentBrowsingPath }}</span>
      </div>
    </div>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleConfirm">确认选择此路径</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { Folder, ArrowUp } from "@element-plus/icons-vue";
import RestClient from "@/utils/rest-client";
import type { InfoList } from "@/types/index";

const props = defineProps({
  modelValue: Boolean,
  engineId: {
    type: Number,
    required: true,
  },
  title: {
    type: String,
    default: "选择路径",
  },
  initialPath: {
    type: String,
    default: "/",
  },
});

const emit = defineEmits(["update:modelValue", "select"]);

const visible = ref(false);
const loading = ref(false);
const currentBrowsingPath = ref("/");
const dirs = ref<string[]>([]);

const joinPath = (base: string, sub: string) => {
  if (!base || base === "/") return "/" + sub;
  return base.endsWith("/") ? base + sub : base + "/" + sub;
};

const parentPath = (curr: string) => {
  if (!curr || curr === "/") return "/";
  const clean = curr.endsWith("/") && curr !== "/" ? curr.slice(0, -1) : curr;
  const idx = clean.lastIndexOf("/");
  if (idx <= 0) return "/";
  return clean.slice(0, idx);
};

const loadDirs = async (pathStr: string) => {
  if (!props.engineId) {
    ElMessage.warning("请先选择存储引擎");
    return;
  }
  loading.value = true;
  try {
    const res = await new RestClient().get<InfoList<string>>("/open-list", {
      engine_id: props.engineId,
      path: pathStr || "/",
    });
    dirs.value = res?.list || [];
    currentBrowsingPath.value = pathStr || "/";
  } catch (err: any) {
    ElMessage.error(`获取目录失败: ${err?.message || err}`);
    dirs.value = [];
  } finally {
    loading.value = false;
  }
};

const selectSubDir = (dirName: string) => {
  const next = joinPath(currentBrowsingPath.value, dirName);
  loadDirs(next);
};

const goParent = () => {
  const parent = parentPath(currentBrowsingPath.value);
  loadDirs(parent);
};

const handleOpen = () => {
  currentBrowsingPath.value = props.initialPath || "/";
  loadDirs(currentBrowsingPath.value);
};

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

const handleConfirm = () => {
  emit("select", currentBrowsingPath.value);
  visible.value = false;
};
</script>
