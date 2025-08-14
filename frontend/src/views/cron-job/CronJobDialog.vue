<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'add' ? '新增任务' : '编辑任务'"
    width="800px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="80px" ref="formRef">
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="是否启用">
            <el-select v-model="form.status" placeholder="是否启用">
              <el-option
                v-for="item in JOB_IS_ENABLE_CONFIG"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="engine">
            <el-select v-model="form.engine_id" placeholder="选择engine">
              <el-option
                v-for="item in allEngineListStore.list"
                :key="item.id"
                :label="item.url"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="src_path">
            <el-cascader style="width: 100%;"
              v-model="srcPathList"
              :props="cascaderSrcPathProps"
              :options="srcPathCascaderNode"
              :placeholder="form.src_path ? form.src_path : '请选择src_path'"
              :disabled="!form.engine_id"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="dst_path">
            <el-cascader style="width: 100%;"
              v-model="dstPathList"
              :props="cascaderSrcPathProps"
              :options="dstPathCascaderNode"
              :placeholder="form.dst_path ? form.dst_path : '请选择dst_path'"
              :disabled="!form.engine_id"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="speed">
            <el-select v-model="form.speed" placeholder="选择speed">
              <el-option
                v-for="item in JOB_SPEED_CONFIGS"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <span style="float: left; margin-right: 20px">{{
                  item.name
                }}</span>
                <span
                  style="
                    float: right;
                    color: var(--el-text-color-secondary);
                    font-size: 13px;
                  "
                >
                  {{ item.remark }}
                </span>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="method">
            <el-select v-model="form.method" placeholder="选择method">
              <el-option
                v-for="item in JOB_METHOD_CONFIGS"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <span style="float: left; margin-right: 20px">{{
                  item.name
                }}</span>
                <span
                  style="
                    float: right;
                    color: var(--el-text-color-secondary);
                    font-size: 13px;
                  "
                >
                  {{ item.remark }}
                </span>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="is_cron">
            <el-select v-model="form.is_cron" placeholder="选择is_cron">
              <el-option
                v-for="item in JOB_IS_CRON_CONFIGS"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <span style="float: left; margin-right: 20px">{{
                  item.name
                }}</span>
                <span
                  style="
                    float: right;
                    color: var(--el-text-color-secondary);
                    font-size: 13px;
                  "
                >
                  {{ item.remark }}
                </span>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="interval" v-if="form.is_cron === 0">
            <el-input v-model="form.interval" placeholder="请输入同步间隔">
              <template #append>分钟</template>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="year" v-if="form.is_cron === 1">
            <el-input v-model="form.year" placeholder="XXXX"/>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="month" v-if="form.is_cron === 1">
            <el-input v-model="form.month" placeholder="1-12"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="day" v-if="form.is_cron === 1">
            <el-input v-model="form.day" placeholder="1-31"/>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="week" v-if="form.is_cron === 1">
            <el-input v-model="form.week" placeholder="1-53"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="day_of_week" v-if="form.is_cron === 1">
            <el-input v-model="form.day_of_week" placeholder="0-6 or mon,tue,wed,thu,fri,sat,sun"/>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="hour" v-if="form.is_cron === 1">
            <el-input v-model="form.hour" placeholder="0-23"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="minute" v-if="form.is_cron === 1">
            <el-input v-model="form.minute" placeholder="0-59"/>
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="second" v-if="form.is_cron === 1">
            <el-input v-model="form.second" placeholder="0-59"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12" class="t-col">
          <el-form-item label="start_date" v-if="form.is_cron === 1">
            <el-date-picker style="width: 100%;" v-model="form.start_date" type="date" placeholder="XXXX-XX-XX" />
          </el-form-item>
        </el-col>
        <el-col :span="12" class="t-col">
          <el-form-item label="end_date" v-if="form.is_cron === 1">
            <el-date-picker style="width: 100%;" v-model="form.end_date" type="date" placeholder="XXXX-XX-XX"/>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确认</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, watch} from "vue";
import { ElMessage, type CascaderProps } from "element-plus";
import type { Job, CascaderNode } from "@/types/index";
import RestClient from "@/utils/rest-client";
import {
  JOB_IS_ENABLE_CONFIG,
  JOB_SPEED_CONFIGS,
  JOB_METHOD_CONFIGS,
  JOB_IS_CRON_CONFIGS,
} from "@/const/index";
import { useAllEngineListStore } from "@/stores";
import type { InfoList } from "@/types";

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
const allEngineListStore = useAllEngineListStore() as ReturnType<
  typeof useAllEngineListStore
> & {
  fetchAllEngines: () => Promise<void>;
};

allEngineListStore.fetchAllEngines();

const srcPathList = ref<string[]>([]);
const srcPathCascaderNode = ref<CascaderNode[]>([]);
const dstPathList = ref<string[]>([]);
const dstPathCascaderNode = ref<CascaderNode[]>([]);
const cascaderSrcPathProps: CascaderProps = {
  lazy: true,
  checkStrictly: true,
  lazyLoad: async (node, resolve) => {
    const { pathLabels } = node;
    const pathStr = pathLabels ? "/" + pathLabels.join("/") : "/";
    const engineId = form.value.engine_id || 0;

    let paths: string[] = [];
    if (engineId) {
      paths = await fetchPaths(engineId, pathStr);
    }

    const nodes = paths.map((path) => ({
      value: path,
      label: path,
      leaf: false,
    }));

    resolve(nodes);
  },
};

const fetchPaths = async (
  engineId: number,
  pathStr: string
): Promise<string[]> => {
  try {
    const infoList = await new RestClient().get<InfoList<string>>(
      "/open-list",
      {
        engine_id: engineId,
        path: pathStr,
      }
    );
    return infoList.list as string[];
  } catch (error) {
    return [];
  }
};

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val && props.job) {
      form.value = { ...props.job };
    }
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

watch(srcPathList, async (val) => {
  if (val) {
    form.value.src_path = "/" + val.join("/");
  }
});

watch(
  () => form.value.engine_id,
  async (val) => {
    if (val) {
      const paths = await fetchPaths(val, "/");
      dstPathCascaderNode.value = srcPathCascaderNode.value = paths.map(
        (path) => ({
          value: path,
          label: path,
          children: [],
        })
      );
    }
  }
);

watch(dstPathList, async (val) => {
  if (val) {
    form.value.dst_path = "/" + val.join("/");
  }
});

const handleClose = () => {
  formRef.value?.resetFields();
};

const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    if (props.mode === "add") {
      await new RestClient().post<number>("/job", form.value);
      ElMessage.success("新增成功");
    } else {
      await new RestClient().put<number>("/job", form.value);
      ElMessage.success("更新成功");
    }
    emit("success");
    visible.value = false;
  } catch (e) {
    ElMessage.error(`${e}`);
  }
};
</script>
