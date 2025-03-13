<!-- 404.vue -->
<template>
  <el-container class="not-found-container">
    <el-main>
      <div class="not-found-content">
        <h1>404</h1>
        <p>页面未找到</p>
        <el-button type="primary" @click="goBack">返回上一页</el-button>
        <el-button @click="goHome">返回首页</el-button>
      </div>
    </el-main>
  </el-container>
</template>

<script>
import { defineComponent } from 'vue';
import { useRouter } from 'vue-router';

export default defineComponent({
  setup() {
    const router = useRouter();

    const goBack = () => {
      router.go(-1);
    };

    const goHome = () => {
      router.push('/');
    };

    return {
      goBack,
      goHome,
    };
  },
});
</script>

<style scoped>
.not-found-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.not-found-content {
  text-align: center;
}

.not-found-content h1 {
  font-size: 100px;
  margin-bottom: 20px;
}

.not-found-content p {
  font-size: 24px;
  margin-bottom: 40px;
}
</style>