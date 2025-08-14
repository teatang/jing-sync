<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'add' ? '新增引擎' : '编辑引擎'"
    width="500px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="80px" ref="formRef">
      <el-form-item
        label="url"
        prop="url"
        :rules="[{ required: true, message: '请输入url' }]"
      >
        <el-input v-model="form.url" />
      </el-form-item>
      <el-form-item
        label="token"
        prop="token"
        :rules="[{ required: true, message: '请输入token' }]"
      >
        <el-input v-model="form.token" />
      </el-form-item>
      <el-form-item
        label="备注"
        prop="remark"
        :rules="[{ required: true, message: '请输入备注' }]"
      >
        <el-input v-model="form.remark" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确认</el-button>
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

const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    if (props.mode === "add") {
      await new RestClient().post<number>("/engine", form.value);
      ElMessage.success("新增成功");
    } else {
      await new RestClient().put<number>("/engine", form.value);
      ElMessage.success("更新成功");
    }
    emit("success");
    visible.value = false;
  } catch (e) {
    ElMessage.error(`${e}`);
  }
};
</script>
