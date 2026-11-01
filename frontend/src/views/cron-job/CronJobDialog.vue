<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'add' ? '新建同步任务' : '编辑同步任务'"
    width="680px"
    class="rounded-xl"
    @close="handleClose"
  >
    <el-form :model="form" label-position="top" ref="formRef" class="space-y-1">
      <!-- 基础描述 -->
      <el-form-item
        label="任务名称 / 备注"
        prop="remark"
        :rules="[{ required: true, message: '请输入任务名称' }]"
      >
        <el-input v-model="form.remark" placeholder="例如：照片备份到阿里云盘、每周工作文档归档" />
      </el-form-item>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- 存储引擎 -->
        <el-form-item
          label="存储引擎"
          prop="engine_id"
          :rules="[{ required: true, message: '请选择存储引擎' }]"
        >
          <el-select v-model="form.engine_id" placeholder="请选择关联的存储引擎" class="w-full">
            <el-option
              v-for="item in allEngineListStore.list"
              :key="item.id"
              :label="`${item.remark || '引擎'} (${item.url})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <!-- 同步模式 -->
        <el-form-item
          label="同步模式"
          prop="sync_mode"
          :rules="[{ required: true, message: '请选择同步模式' }]"
        >
          <el-select v-model="form.sync_mode" placeholder="选择同步模式" class="w-full">
            <el-option
              v-for="item in SYNC_MODE_CONFIGS"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
              <div class="flex items-center justify-between">
                <span>{{ item.name }}</span>
                <span class="text-xs text-slate-400">{{ item.remark }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
      </div>

      <!-- 路径配置与安全检查 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <el-form-item
          label="源路径 (Source Path)"
          prop="src_path"
          :rules="[{ required: true, message: '请输入或选择源目录路径' }]"
        >
          <el-input
            v-model="form.src_path"
            placeholder="直接输入或点击右侧浏览选择"
            :disabled="!form.engine_id"
            clearable
          >
            <template #append>
              <el-button :disabled="!form.engine_id" @click="openPathSelector('src_path')">
                浏览
              </el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item
          label="目标路径 (Destination Path)"
          prop="dst_path"
          :rules="[{ required: true, message: '请输入或选择目标目录路径' }]"
        >
          <el-input
            v-model="form.dst_path"
            placeholder="直接输入或点击右侧浏览选择"
            :disabled="!form.engine_id"
            clearable
          >
            <template #append>
              <el-button :disabled="!form.engine_id" @click="openPathSelector('dst_path')">
                浏览
              </el-button>
            </template>
          </el-input>
        </el-form-item>
      </div>

      <!-- 路径重叠警告 -->
      <el-alert
        v-if="pathsOverlapWarning"
        title="路径重叠安全警告：源路径与目标路径相互包含或相同，同步时可能导致死循环传输！"
        type="error"
        show-icon
        :closable="false"
        class="!my-2"
      />

      <!-- 文件过滤规则 -->
      <div class="p-3 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800 space-y-3">
        <div class="text-xs font-semibold text-slate-700 dark:text-slate-300">文件过滤与排他规则</div>

        <el-form-item label="排除通配规则 (Glob Pattern)" class="!mb-2">
          <el-input
            v-model="form.excludes"
            placeholder="多条规则用逗号分隔，如: *.tmp, *.bak, .git/**, @eaDir/**"
            clearable
          />
        </el-form-item>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <el-form-item label="最小文件大小限制" class="!mb-0">
            <el-input-number
              v-model="form.min_size"
              :min="0"
              placeholder="字节 (Byte)"
              class="!w-full"
            />
            <span class="text-xs text-slate-400 mt-1 block">
              {{ form.min_size ? `约 ${formatBytes(form.min_size)}` : '无下限' }}
            </span>
          </el-form-item>

          <el-form-item label="最大文件大小限制" class="!mb-0">
            <el-input-number
              v-model="form.max_size"
              :min="0"
              placeholder="字节 (Byte)"
              class="!w-full"
            />
            <span class="text-xs text-slate-400 mt-1 block">
              {{ form.max_size ? `约 ${formatBytes(form.max_size)}` : '无上限' }}
            </span>
          </el-form-item>
        </div>
      </div>

      <!-- 调度配置 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
        <el-form-item label="扫描速率">
          <el-select v-model="form.speed" class="w-full">
            <el-option
              v-for="item in JOB_SPEED_CONFIGS"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="调度模式">
          <el-select v-model="form.is_cron" class="w-full">
            <el-option
              v-for="item in JOB_IS_CRON_CONFIGS"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </div>

      <!-- 动态间隔或Cron表达式 -->
      <div v-if="form.is_cron === 0">
        <el-form-item label="同步间隔 (分钟)">
          <el-input-number v-model="form.interval" :min="1" :max="10080" class="!w-full" />
        </el-form-item>
      </div>

      <div v-else-if="form.is_cron === 1">
        <el-form-item label="Cron表达式 (分 时 日 月 周)">
          <el-input v-model="form.cron_expr" placeholder="例如: 0 2 * * * (每日凌晨2点执行)" />
        </el-form-item>
      </div>

      <!-- 启用状态开关 -->
      <div class="pt-2 flex items-center justify-between border-t border-slate-100 dark:border-slate-800">
        <span class="text-sm text-slate-700 dark:text-slate-300">保存后自动启用定时调度</span>
        <el-switch v-model="enabledBoolean" />
      </div>
    </el-form>

    <!-- 路径选择弹窗 -->
    <PathSelectorDialog
      v-model="pathSelectorVisible"
      :engine-id="form.engine_id || 0"
      :title="pathSelectorTitle"
      :initial-path="pathSelectorInitial"
      @select="handlePathSelected"
    />

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确认保存</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, watch, computed } from "vue";
import { ElMessage } from "element-plus";
import type { Job } from "@/types/index";
import RestClient from "@/utils/rest-client";
import {
  JOB_SPEED_CONFIGS,
  JOB_IS_CRON_CONFIGS,
  SYNC_MODE_CONFIGS,
} from "@/const/index";
import { useAllEngineListStore } from "@/stores";
import { formatBytes } from "@/utils/format";
import PathSelectorDialog from "./PathSelectorDialog.vue";

const props = defineProps({
  modelValue: Boolean,
  job: {
    type: Object as () => Partial<Job>,
    default: () => ({}),
  },
  mode: {
    type: String,
    default: "add",
  },
});

const emit = defineEmits(["update:modelValue", "success"]);

const visible = ref(false);
const form = ref<Partial<Job>>({});
const formRef = ref();
const submitting = ref(false);

const pathSelectorVisible = ref(false);
const pathSelectorTitle = ref("选择源路径");
const pathSelectorField = ref<"src_path" | "dst_path">("src_path");
const pathSelectorInitial = ref("/");

const allEngineListStore = useAllEngineListStore() as ReturnType<
  typeof useAllEngineListStore
> & {
  fetchAllEngines: () => Promise<void>;
};

allEngineListStore.fetchAllEngines();

const enabledBoolean = computed({
  get: () => (form.value.enabled ?? 1) === 1,
  set: (val) => {
    form.value.enabled = val ? 1 : 0;
  },
});

const pathsOverlapWarning = computed(() => {
  const src = form.value.src_path?.trim();
  const dst = form.value.dst_path?.trim();
  if (!src || !dst) return false;
  if (src === dst || src === "/" || dst === "/") return true;
  const cleanSrc = src.endsWith("/") ? src : src + "/";
  const cleanDst = dst.endsWith("/") ? dst : dst + "/";
  return cleanDst.startsWith(cleanSrc) || cleanSrc.startsWith(cleanDst);
});

const openPathSelector = (field: "src_path" | "dst_path") => {
  if (!form.value.engine_id) {
    ElMessage.warning("请先选择存储引擎");
    return;
  }
  pathSelectorField.value = field;
  pathSelectorTitle.value =
    field === "src_path" ? "选择源路径 (Source Path)" : "选择目标路径 (Destination Path)";
  pathSelectorInitial.value =
    (field === "src_path" ? form.value.src_path : form.value.dst_path) || "/";
  pathSelectorVisible.value = true;
};

const handlePathSelected = (path: string) => {
  if (pathSelectorField.value === "src_path") {
    form.value.src_path = path;
  } else {
    form.value.dst_path = path;
  }
};

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val && props.job) {
      form.value = {
        sync_mode: 1,
        speed: 0,
        interval: 60,
        is_cron: 0,
        enabled: 1,
        ...props.job,
      };
    }
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

const handleClose = () => {
  formRef.value?.resetFields();
};

const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    if (pathsOverlapWarning.value) {
      ElMessage.warning("源路径与目标路径存在重叠，请修正后再保存！");
      return;
    }
    submitting.value = true;
    if (props.mode === "add") {
      await new RestClient().post<number>("/job", form.value);
      ElMessage.success("新增成功");
    } else {
      await new RestClient().put<number>("/job", form.value);
      ElMessage.success("更新成功");
    }
    emit("success");
    visible.value = false;
  } catch (e: any) {
    if (e) {
      ElMessage.error(`${e?.message || e}`);
    }
  } finally {
    submitting.value = false;
  }
};
</script>
