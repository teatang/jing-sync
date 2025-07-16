<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'add' ? '新增用户' : '编辑用户'"
    width="500px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="80px" ref="formRef">
      <el-form-item
        label="username"
        prop="username"
        :rules="[{ required: true, message: '请输入username' }]"
      >
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item
        label="password"
        prop="password"
        :rules="[{ required: true, message: '请输入password' }]"
      >
        <el-input v-model="form.password" />
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
import { type User } from "@/types/index";
import RestClient from "@/utils/rest-client";

const props = defineProps({
  modelValue: Boolean,
  user: {
    type: Object as () => Partial<User>,
    default: () => ({}),
  },
  mode: {
    type: String,
    default: "add",
  },
});

const emit = defineEmits(["update:modelValue", "success"]);

const visible = ref(false);
const form = ref<Partial<User>>({});
const formRef = ref();

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val && props.user) {
      form.value = { ...props.user };
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
      await new RestClient().post<number>("/user", form.value);
      ElMessage.success("新增成功");
    } else {
      await new RestClient().put<number>("/user", form.value);
      ElMessage.success("更新成功");
    }
    emit("success");
    visible.value = false;
  } catch (e) {
    console.error(e);
    ElMessage.error(`操作失败 msg: ${e}`);
  }
};
</script>
