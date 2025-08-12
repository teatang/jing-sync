
<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="logo">
        <img src="https://picsum.photos/150/150?blue" alt="蓝色主题系统Logo">
      </div>
      <h2 class="title">系统登录</h2>
      
      <el-form 
        :model="loginForm" 
        :rules="loginRules" 
        ref="loginFormRef"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/index'
import { ElMessage } from 'element-plus'

const loginForm = ref({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const loginFormRef = ref()
const loading = ref(false)
const router = useRouter()
const userStore = useUserStore()

const handleLogin = async () => {
  try {
    await loginFormRef.value.validate()
    loading.value = true
    await userStore.login(loginForm.value)
    router.push('/')
  } catch (error) {
    ElMessage.error(`${error}`)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1e88e5, #0d47a1);
  
  .login-card {
    width: 400px;
    padding: 30px;
    border-radius: 10px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    
    .logo {
      text-align: center;
      margin-bottom: 20px;
      
      img {
        border-radius: 50%;
        border: 3px solid #2196f3;
      }
    }
    
    .title {
      text-align: center;
      color: #2196f3;
      margin-bottom: 30px;
    }
    
    .login-btn {
      width: 100%;
      background: linear-gradient(to right, #2196f3, #1976d2);
      border: none;
      height: 45px;
      font-size: 16px;
      letter-spacing: 2px;
      transition: all 0.3s;
      
      &:hover {
        opacity: 0.9;
        transform: translateY(-2px);
      }
    }
  }
}
</style>
