<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'add' ? '新增存储引擎' : '编辑存储引擎'"
    width="520px"
    class="rounded-xl"
    @close="handleClose"
  >
    <el-form :model="form" label-position="top" ref="formRef" class="mt-2">
      <el-form-item
        label="服务地址 (URL)"
        prop="url"
        :rules="[{ required: true, message: '请输入 OpenList / AList 服务地址 (如 http://192.168.1.100:5244)' }]"
      >
        <el-input
          v-model="form.url"
          placeholder="http://127.0.0.1:5244"
          clearable
        />
      </el-form-item>

      <el-form-item
        label="访问令牌 (Token)"
        prop="token"
        :rules="[{ required: true, message: '请输入 OpenList / AList API 访问令牌' }]"
      >
        <el-input
          v-model="form.token"
          placeholder="alist-xxxxxxxxxxxx"
          type="password"
          show-password
          clearable
        />
      </el-form-item>

      <el-form-item
        label="备注名称"
        prop="remark"
        :rules="[{ required: true, message: '请输入备注标识此引擎用途' }]"
      >
        <el-input
          v-model="form.remark"
          placeholder="例如：家庭群晖NAS、百度网盘挂载端"
          clearable
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="flex items-center justify-between">
        <el-button
          type="success"
          plain
          :loading="testing"
          @click="handleTestConnection"
        >
          测试连接
        </el-button>
        <div class="flex gap-2">
          <el-button @click="visible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">确认保存</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { type Engine } from "@/types/index";
import RestClient from "@/utils/rest-client";

const props = defineProps({
  modelValue: Boolean,
  engine: {
    type: Object as () => Partial<Engine>,
    default: () => ({}),
  },
  mode: {
    type: String,
    default: "add",
  },
});

const emit = defineEmits(["update:modelValue", "success"]);

const visible = ref(false);
const form = ref<Partial<Engine>>({});
const formRef = ref();
const testing = ref(false);
const submitting = ref(false);

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val && props.engine) {
      form.value = { ...props.engine };
    }
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

const handleClose = () => {
  formRef.value?.resetFields();
};

const handleTestConnection = async () => {
  if (!form.value.url) {
    ElMessage.warning("请先填写服务地址");
    return;
  }
  testing.value = true;
  try {
    await new RestClient().post("/engine/test", {
      id: form.value.id || 0,
      url: form.value.url,
      token: form.value.token,
    });
    ElMessage.success("连通性测试通过！服务响应正常");
  } catch (error: any) {
    ElMessage.error(`连通性测试失败: ${error?.message || error}`);
  } finally {
    testing.value = false;
  }
};

const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    submitting.value = true;
    if (props.mode === "add") {
      await new RestClient().post<number>("/engine", form.value);
      ElMessage.success("新增成功");
    } else {
      await new RestClient().put<number>("/engine", form.value);
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
