<template>
  <div class="notification-center">
    <el-badge :value="unreadCount" :max="99" class="notification-badge">
      <el-popover
        placement="bottom"
        :width="350"
        trigger="click"
        popper-class="notification-popover"
        @show="loadNotifications"
      >
        <template #reference>
          <el-button class="notification-btn" circle>
            <el-icon><Bell /></el-icon>
          </el-button>
        </template>
        
        <template #default>
          <div class="notification-header">
            <h3>通知中心</h3>
            <div class="notification-actions">
              <el-button type="primary" link @click="markAllAsRead" :disabled="!hasUnread">
                全部已读
              </el-button>
              <el-button type="primary" link @click="goToNotificationPage">
                查看全部
              </el-button>
            </div>
          </div>
          
          <el-divider />
          
          <div v-if="loading" class="notification-loading">
            <el-skeleton :rows="3" animated />
          </div>
          
          <div v-else-if="notifications.length === 0" class="notification-empty">
            <el-empty description="暂无通知" />
          </div>
          
          <div v-else class="notification-list">
            <div
              v-for="notification in notifications"
              :key="notification.id"
              class="notification-item"
              :class="{ 'is-read': notification.is_read }"
              @click="viewNotification(notification)"
            >
              <div class="notification-icon">
                <el-icon :size="20" :color="getNotificationColor(notification)">
                  <component :is="getNotificationIcon(notification)" />
                </el-icon>
              </div>
              <div class="notification-content">
                <div class="notification-title">{{ notification.notification.title }}</div>
                <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
              </div>
              <div class="notification-actions">
                <el-button
                  v-if="!notification.is_read"
                  type="primary"
                  size="small"
                  circle
                  @click.stop="markAsRead(notification)"
                >
                  <el-icon><Check /></el-icon>
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  circle
                  @click.stop="deleteNotification(notification)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </template>
      </el-popover>
    </el-badge>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { Bell, Check, Delete, Warning, Message, InfoFilled } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { notificationApi } from '@/services/notification';
import websocketService from '@/services/websocket';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import eventBus from '@/utils/eventBus';
dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

const router = useRouter();

const notifications = ref([]);
const unreadCount = ref(0);
const loading = ref(false);
const hasUnread = computed(() => unreadCount.value > 0);

// 加载通知列表
const loadNotifications = async () => {
  loading.value = true;
  try {
    const params = {
      page: 1,
      page_size: 10
    };
    const data = await notificationApi.getUserNotifications(params);
   
    notifications.value = data.items;
  } catch (error) {
    console.error('Failed to load notifications:', error);
    ElMessage.error('加载通知失败');
  } finally {
    loading.value = false;
  }
};

// 获取未读通知数量
const getUnreadCount = async () => {

  try {
    const data = await notificationApi.getUnreadNotificationCount();
   
    unreadCount.value = data.count;
  } catch (error) {
    console.error('Failed to get unread count:', error);
  }
};

// 标记通知为已读
const markAsRead = async (notification) => {
  try {
    await notificationApi.markNotificationAsRead(notification.notification_id);
    notification.is_read = true;
    unreadCount.value = Math.max(0, unreadCount.value - 1);
    ElMessage.success('已标记为已读');
  } catch (error) {
    console.error('Failed to mark as read:', error);
    ElMessage.error('标记已读失败');
  }
};

// 标记所有通知为已读
const markAllAsRead = async () => {
  try {
    await notificationApi.markAllNotificationsAsRead();
    notifications.value.forEach(notification => {
      notification.is_read = true;
    });
    unreadCount.value = 0;
    ElMessage.success('已全部标记为已读');
  } catch (error) {
    console.error('Failed to mark all as read:', error);
    ElMessage.error('标记全部已读失败');
  }
};

// 删除通知
const deleteNotification = async (notification) => {
  try {
    await ElMessageBox.confirm('确定要删除这条通知吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await notificationApi.deleteUserNotification(notification.notification_id);
    const index = notifications.value.findIndex(n => n.id === notification.id);
    if (index !== -1) {
      notifications.value.splice(index, 1);
    }
    
    if (!notification.is_read) {
      unreadCount.value = Math.max(0, unreadCount.value - 1);
    }
    
    ElMessage.success('删除成功');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete notification:', error);
      ElMessage.error('删除通知失败');
    }
  }
};

// 查看通知详情
const viewNotification = async (notification) => {
  if (!notification.is_read) {
    await markAsRead(notification);
  }
  
  // 根据通知类型跳转到不同页面或显示详情
  ElMessageBox.alert(notification.notification.content, notification.notification.title, {
    confirmButtonText: '确定',
    dangerouslyUseHTMLString: true
  });
};

// 跳转到通知页面
const goToNotificationPage = () => {
  router.push('/layout/user/notifications');
};

// 获取通知图标
const getNotificationIcon = (notification) => {
  
  const level = notification.notification.level;
 
  if (level === 3) return Warning;
  if (level === 2) return Message;
  return InfoFilled;
};

// 获取通知颜色
const getNotificationColor = (notification) => {
  const level = notification.notification.level;
  if (level === 3) return '#E6A23C';
  if (level === 2) return '#409EFF';
  return '#909399';
};

// 格式化时间
const formatTime = (time) => {
  return dayjs(time).fromNow();
};

// 处理新通知
const handleNewNotification = (data) => {
  getUnreadCount();
  
  // 如果通知列表已加载，则添加新通知到列表顶部
  if (notifications.value.length > 0) {
    loadNotifications();
  }
};

onMounted(async () => {
  await getUnreadCount();
  // 只需要监听通知事件，不需要重复建立连接
  websocketService.on('notification', handleNewNotification);
  eventBus.on('notification-read', getUnreadCount);
});

</script>

<style scoped>
.notification-center {
  display: inline-block;
  position: relative;
}

.notification-badge {
  margin-right: 20px;
}

.notification-btn {
  background: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
  height: 40px;
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 修正徽标位置 */
:deep(.el-badge__content.is-fixed) {
  top: 8px;
  right: 8px;
  transform: translateY(-50%) translateX(50%);
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 10px;
}

.notification-header h3 {
  margin: 0;
  font-size: 16px;
}

.notification-list {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  padding: 10px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.is-read {
  opacity: 0.7;
}

.notification-icon {
  margin-right: 10px;
  display: flex;
  align-items: center;
}

.notification-content {
  flex: 1;
}

.notification-title {
  font-weight: bold;
  margin-bottom: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notification-time {
  font-size: 12px;
  color: #909399;
}

.notification-actions {
  display: flex;
  align-items: center;
}

.notification-empty,
.notification-loading {
  padding: 20px;
  text-align: center;
}
</style>

<style>
.notification-popover {
  padding: 0;
  overflow: hidden;
}
</style> 