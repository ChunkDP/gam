<!-- App.vue -->
<template>

    <ErrorBoundary>
    <router-view></router-view>
  </ErrorBoundary>

</template>

<script>
import { defineComponent, onMounted, onBeforeUnmount } from 'vue';
import ErrorBoundary from './components/ErrorBoundary.vue';
import websocketService from '@/services/websocket';
import { ElMessage } from 'element-plus';

export default defineComponent({
  name: 'App',
  components: {
    ErrorBoundary
  },
  setup() {
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

    let connectionChecker = null;

    onMounted(async () => {
      const token = sessionStorage.getItem('token');
      
      if (token) {
        try {
          await websocketService.connect(token);
          websocketService.on('notification', handleNewNotification);
          
          // 每30秒检查一次连接状态
          connectionChecker = setInterval(checkConnection, 30000);
        } catch (error) {
          console.error('Failed to connect WebSocket:', error);
        }
      }
    });

    onBeforeUnmount(() => {
      if (connectionChecker) {
        clearInterval(connectionChecker);
      }
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
