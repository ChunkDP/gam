<!-- App.vue -->
<template>

    <ErrorBoundary>
    <router-view></router-view>
  </ErrorBoundary>

</template>

<script>
import { defineComponent, onMounted, onBeforeUnmount, watch, ref } from 'vue';
import { useRoute } from 'vue-router';
import ErrorBoundary from './components/ErrorBoundary.vue';
import websocketService from '@/services/websocket';
import { ElMessage } from 'element-plus';

export default defineComponent({
  name: 'App',
  components: {
    ErrorBoundary
  },
  setup() {
    const route = useRoute();
    const isConnected = ref(false);

    // 检查WebSocket连接状态
    const checkConnection = () => {
      const token = sessionStorage.getItem('token');
      if (token && !websocketService.isConnected()) {
        console.log('WebSocket disconnected, attempting to reconnect...');
        websocketService.reconnect();
      }
    };

    // 处理新通知
    const handleNewNotification = (data) => {
      ElMessage({
        message: data.title,
        type: data.level === 3 ? 'warning' : (data.level === 2 ? 'info' : 'success'),
        duration: 5000,
        showClose: true
      });
    };

    // 初始化WebSocket连接
    const initWebSocket = async () => {
      const token = sessionStorage.getItem('token');
      
      if (token && !isConnected.value) {
        try {
          await websocketService.connect(token);
          websocketService.on('notification', handleNewNotification);
          isConnected.value = true;
          console.log('WebSocket connected successfully');
        } catch (error) {
          console.error('Failed to connect WebSocket:', error);
        }
      }
    };

    let connectionChecker = null;

    // 监听token变化
    watch(() => sessionStorage.getItem('token'), (newToken) => {
      if (newToken) {
        initWebSocket();
      } else if (isConnected.value) {
        // 如果token被移除，断开WebSocket连接
        websocketService.disconnect();
        isConnected.value = false;
      }
    });

    onMounted(async () => {
      console.log('onMounted');
      await initWebSocket();
      
      // 每30秒检查一次连接状态
      connectionChecker = setInterval(checkConnection, 30000);
    });

    onBeforeUnmount(() => {
      console.log('onBeforeUnmount');
      if (connectionChecker) {
        clearInterval(connectionChecker);
      }
      // 清理WebSocket监听器
      websocketService.off('notification', handleNewNotification);
    });

    return {};
  }
});
</script>
<style lang="scss">
// 可以直接使用scss语法
$primary-color: #409EFF;

.container {
  color: $primary-color;
}
</style>
