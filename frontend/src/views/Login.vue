<template>
  <div class="login-container">
    <div class="login-wrapper">
      <!-- 左侧背景 -->
      <div class="login-background">
        <div class="background-content">
          <h2 class="welcome-text">Hello World!</h2>

        </div>
      </div>

      <!-- 右侧登录表单 -->
      <div class="login-form-container">
        <div class="login-header">
          <div class="logo-wrapper">
            <img src="@/assets/logo.png" alt="Logo" class="logo">
          </div>
    
        </div>

        <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form">
          <el-form-item prop="username">
            <el-input v-model="loginForm.username" placeholder="请输入用户名" prefix-icon="User" />
          </el-form-item>

          <el-form-item prop="password">
            <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" prefix-icon="Lock"
              show-password />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" class="login-button" :loading="loading" @click="submitForm">
              登录
            </el-button>
          </el-form-item>

          <el-form-item>
            <el-button class="reset-button" @click="resetForm">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useUserPermissionsStore } from '../stores/userPermissions';
import { login } from '@/services/apiService';

export default defineComponent({
  setup() {
    const router = useRouter();
    const loginForm = ref({
      username: '',
      password: ''
    });

    const loginRules = ref({
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' }
      ]
    });

    const loginFormRef = ref(null);

    const userPermissionsStore = useUserPermissionsStore();
    const loading = ref(false);

    const submitForm = async () => {
      try {
        loading.value = true;
        const loginResponse = await login(loginForm.value);
       
        await userPermissionsStore.login(loginResponse);
        router.push('/layout');
      } catch (error) {
        ElMessage.error(error.message || '登录失败');
      } finally {
        loading.value = false;
      }
    };

    const resetForm = () => {
      loginFormRef.value.resetFields();
    };

    return {
      loginForm,
      loginRules,
      loginFormRef,
      submitForm,
      resetForm,
      loading
    };
  }
});
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;

}

.login-wrapper {
  display: flex;
  width: 1000px;
  height: 600px;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.login-form-container {
  width: 400px;
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: white;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  /* 确保所有子元素水平居中 */
}

.logo-wrapper {
  width: 200px;
  height: 200px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 16px;
}

.logo {
  width: 200px;
  height: 200px;
  margin-bottom: 16px;
}

.title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #909399;
}

.login-form {
  :deep(.el-input) {
    height: 40px;

    .el-input__wrapper {
      background-color: #f5f7fa;
    }
  }
}

.login-button,
.reset-button {
  width: 100%;
  height: 40px;
  border-radius: 4px;
}

.login-button {
  background: #409eff;

  &:hover {
    background: #66b1ff;
  }
}

/* 左侧背景样式 */
/* .login-background {
  flex: 1;
  background: linear-gradient(135deg, #1890ff 0%, #3f51b5 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  padding: 40px;
} */
.login-background {
  flex: 1;
  /* 多重背景：第一层是渐变色，第二层是背景图片 */
  background:
    linear-gradient(135deg, #1890ff 0%, #3f51b5 100%),
    url('@/assets/login-bg.svg');

  /* 控制每个背景层的大小和位置 */
  background-size: cover, contain;
  /* 渐变覆盖整个区域，图片保持比例 */
  background-position: center, center;
  /* 都居中显示 */
  background-repeat: no-repeat, no-repeat;
  /* 都不重复 */

  /* 可以调整图片的混合模式 */
  background-blend-mode: overlay;
  /* 可以尝试：multiply, overlay, soft-light 等 */

  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  padding: 40px;
}

.background-content {
  max-width: 800px;
}

.welcome-text {
  font-size: 42px;
  margin-bottom: 40px;
  font-weight: 300;
}

.code-section {
  background: rgba(255, 255, 255, 0.1);
  padding: 20px;
  border-radius: 8px;
  font-family: 'Courier New', Courier, monospace;

  pre {
    margin: 0;
    color: #e6e6e6;
    font-size: 16px;
    line-height: 1.6;
  }
}

@media (max-width: 1200px) {
  .login-wrapper {
    width: 100%;
    height: auto;
    flex-direction: column;
  }

  .login-background {
    display: none;
  }

  .login-form-container {
    width: 100%;
    max-width: 400px;
    margin: 0 auto;
  }
}
</style>