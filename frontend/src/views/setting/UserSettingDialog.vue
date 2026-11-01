<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'add' ? '新增用户' : '编辑用户'"
    width="480px"
    class="rounded-xl"
    @close="handleClose"
  >
    <el-form :model="form" label-position="top" ref="formRef" class="space-y-3">
      <el-form-item
        label="用户名 (Username)"
        prop="username"
        :rules="[{ required: true, message: '请输入用户名' }]"
      >
        <el-input v-model="form.username" placeholder="请输入用户名" clearable />
      </el-form-item>
      <el-form-item
        label="密码 (Password)"
        prop="password"
        :rules="[{ required: true, message: '请输入密码' }]"
      >
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入登录密码"
          show-password
          clearable
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确认保存</el-button>
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
const submitting = ref(false);
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
    submitting.value = true;
    if (props.mode === "add") {
      await new RestClient().post<number>("/user", form.value);
      ElMessage.success("新增用户成功");
    } else {
      await new RestClient().put<number>("/user", form.value);
      ElMessage.success("更新用户成功");
    }
    emit("success");
    visible.value = false;
  } catch (e) {
    if (e) {
      ElMessage.error(`${e}`);
    }
  } finally {
    submitting.value = false;
  }
};
</script>
