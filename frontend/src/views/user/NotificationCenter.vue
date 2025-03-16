<template>
  <div class="ma-search-box">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的通知</span>
          <el-button-group>
            <el-button type="primary" @click="markAllAsRead">全部已读</el-button>
            <el-button @click="handleRefresh">刷新</el-button>
          </el-button-group>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="queryParams" >
        <el-form-item label="类型">
          <el-select v-model="queryParams.type_id" placeholder="请选择通知类型" clearable class="adaptive-select">
            <el-option
              v-for="type in notificationTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.is_read" placeholder="请选择状态" clearable class="adaptive-select">
            <el-option label="未读" :value="false" />
            <el-option label="已读" :value="true" />
          </el-select>
        </el-form-item>
        <el-form-item label="撤回状态">
          <el-select v-model="queryParams.show_recalled" placeholder="请选择状态" clearable class="adaptive-select">
            <el-option label="未撤回" :value="false" />
            <el-option label="已撤回" :value="true" />
          </el-select>
        </el-form-item>

        
        <el-form-item label="级别">
          <el-select v-model="queryParams.level" placeholder="请选择级别" clearable class="adaptive-select">
            <el-option label="普通" :value="1" />
            <el-option label="重要" :value="2" />
            <el-option label="紧急" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 通知列表 -->
      <el-table v-loading="loading" :data="notifications" border>
        <el-table-column width="50">
          <template #default="{ row }" >
            <el-badge :is-dot="!row.is_read" type="danger" >
              <el-icon><Bell /></el-icon>
            </el-badge>
          </template>
        </el-table-column>
        <el-table-column label="标题" prop="notification.title" >
          <template #default="{ row }">
          <span :class="{ 'recalled': row.is_recalled }">
            {{ row.notification.title }}
            <el-tag v-if="row.is_recalled" type="info" size="small">已撤回</el-tag>
          </span>
        </template>
        </el-table-column>
        <el-table-column label="类型" prop="notification.type.name" width="100" />
        <el-table-column label="级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.notification.level)">
              {{ getLevelLabel(row.notification.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="发送时间" prop="created_at" width="220" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_read ? 'info' : 'danger'">
              {{ row.is_read ? '已读' : '未读' }}
            </el-tag>
            <el-tooltip v-if="row.is_recalled && row.read_time && row.recall_time" 
                     :content="getReadStatus(row)" placement="top">
            <el-icon class="ml-2"><InfoFilled /></el-icon>
          </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="viewNotification(row)">查看</el-button>
              <el-button
                v-if="!row.is_read"
                type="primary"
                link
                @click="markAsRead(row)"
              >标记已读</el-button>
              <el-button type="primary" link @click="deleteNotification(row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
      <el-pagination 
      background
        v-model:current-page="queryParams.page"
        v-model:page-size="queryParams.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      </div>
    </el-card>

    <!-- 查看通知对话框 -->
    <el-dialog title="通知详情" v-model="detailDialogVisible" width="60%">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="标题">{{ currentNotification?.notification?.title }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ currentNotification?.notification?.type?.name }}</el-descriptions-item>
        <el-descriptions-item label="级别">
          {{ getLevelLabel(currentNotification?.notification?.level) }}
        </el-descriptions-item>
        <el-descriptions-item label="发送时间">{{ currentNotification?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="阅读时间">{{ currentNotification?.read_time || '-' }}</el-descriptions-item>
        <el-descriptions-item label="内容">{{ currentNotification?.notification?.content }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button
          v-if="!currentNotification?.is_read"
          type="primary"
          @click="markAsReadAndClose"
        >标记已读并关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, onBeforeUnmount } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Bell, InfoFilled } from '@element-plus/icons-vue';
import { notificationApi } from '@/services/notification';
import websocketService from '@/services/websocket';

const getReadStatus = (notification) => {
 
  const readTime = new Date(notification.read_time);
  const recallTime = new Date(notification.recall_time);
  return readTime < recallTime ? 
    '撤回前已读' : 
    '撤回后已读';
};
// 查询参数
const queryParams = ref({
  page: 1,
  page_size: 10,
  type_id: '',
  is_read: '',
  show_recalled:'',
  level: '',
  user_type: 'admins' // 从用户store中获取用户类型
});

const notifications = ref([]);
const notificationTypes = ref([]);
const currentNotification = ref(null);
const loading = ref(false);
const detailDialogVisible = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const dateRange = ref([]);

// 计算属性
const hasUnread = computed(() => {
  return notifications.value.some(item => !item.is_read);
});

// 搜索表单
const searchForm = reactive({
  type_id: '',
  is_read: '',
  level: '',
  start_time: '',
  end_time: '',
  page: 1,
  page_size: 10
});

// 监听日期范围变化
watch(dateRange, (newVal) => {
  if (newVal && newVal.length === 2) {
    searchForm.start_time = newVal[0];
    searchForm.end_time = newVal[1];
  } else {
    searchForm.start_time = '';
    searchForm.end_time = '';
  }
});

// 加载通知类型
const loadNotificationTypes = async () => {
  try {
    notificationTypes.value = await notificationApi.getNotificationTypes();
  } catch (error) {
    console.error('Failed to load notification types:', error);
    ElMessage.error('加载通知类型失败');
  }
};

// 获取通知列表
const fetchNotifications = async () => {
  try {
   
      // 只有当明确选择了已读/未读状态时才传递参数
  if (queryParams.value.is_read === '') {
    delete queryParams.value.is_read;
  }
  
  if (queryParams.value.show_recalled === '') {
    delete queryParams.value.show_recalled;
  }

    loading.value = true;
    const data = await notificationApi.getUserNotifications(queryParams.value);
 
    notifications.value = data.items;
    total.value = data.total;
  } catch (error) {
    ElMessage.error('获取通知列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索通知
const searchNotifications = () => {
  currentPage.value = 1;
  fetchNotifications();
};

// 重置搜索
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    if (key !== 'page' && key !== 'page_size') {
      searchForm[key] = '';
    }
  });
  dateRange.value = [];
  currentPage.value = 1;
  fetchNotifications();
};

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchNotifications();
};

const handleCurrentChange = (page) => {
  currentPage.value = page;
  fetchNotifications();
};

// 查看通知详情
const viewNotification = (row) => {
  currentNotification.value = row;
  detailDialogVisible.value = true;
  
  // 如果是未读通知，标记为已读
  // if (!row.is_read) {
  //   markAsRead(row, false);
  // }
};
const markAsReadAndClose = () => {
  markAsRead(currentNotification.value, false);
  detailDialogVisible.value = false;
};
// 标记通知为已读
const markAsRead = async (row, showMessage = true) => {
  try {
    await notificationApi.markNotificationAsRead(row.notification.id);
    row.is_read = true;
    row.read_time = new Date().toISOString();
    
    if (showMessage) {
      ElMessage.success('标记为已读成功');
    }
  } catch (error) {
    console.error('Failed to mark notification as read:', error);
    if (showMessage) {
      ElMessage.error('标记为已读失败');
    }
  }
};

// 标记所有通知为已读
const markAllAsRead = async () => {
  try {
    await ElMessageBox.confirm('确定要将所有通知标记为已读吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await notificationApi.markAllNotificationsAsRead();
    ElMessage.success('全部标记为已读成功');
    fetchNotifications();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to mark all notifications as read:', error);
      ElMessage.error('全部标记为已读失败');
    }
  }
};

// 删除通知
const deleteNotification = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该通知吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await notificationApi.deleteUserNotification(row.notification.id);
    ElMessage.success('删除通知成功');
    fetchNotifications();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete notification:', error);
      ElMessage.error('删除通知失败');
    }
  }
};

// 处理新通知
const handleNewNotification = (data) => {
  ElMessage({
    message: `收到新通知: ${data.title}`,
    type: data.level === 3 ? 'warning' : (data.level === 2 ? 'info' : 'success'),
    duration: 5000,
    showClose: true
  });
  
  fetchNotifications();
};

onMounted(() => {
  loadNotificationTypes();
  fetchNotifications();
  
  // 只需要监听通知事件，不需要重复建立连接
  websocketService.on('notification', handleNewNotification);
  websocketService.on('notification-recall', (data) => {
    // 从列表中移除已撤回的通知
    notifications.value = notifications.value.filter(
      item => item.notification.id !== data.id
    );
    
    // 显示撤回提示
    ElMessage({
      type: 'info',
      message: data.message,
      duration: 3000
    });
  });
});

// 获取级别类型
const getLevelType = (level) => {
  switch (level) {
    case 3:
      return 'danger';
    case 2:
      return 'warning';
    default:
      return 'info';
  }
};

// 获取级别标签
const getLevelLabel = (level) => {
  switch (level) {
    case 3:
      return '紧急';
    case 2:
      return '重要';
    default:
      return '普通';
  }
};

// 处理刷新
const handleRefresh = () => {
  fetchNotifications();
};

// 处理查询
const handleQuery = () => {
  currentPage.value = 1;
  fetchNotifications();
};

// 处理重置查询
const resetQuery = () => {
  Object.keys(queryParams.value).forEach(key => {
    if (key !== 'page' && key !== 'page_size') {
      queryParams.value[key] = '';
    }
  });
  currentPage.value = 1;
  fetchNotifications();
};
</script>

<style scoped>


/* 强制表格单元格内容可见 */
:deep(.el-table .cell) {
  overflow: visible !important;
  white-space: nowrap;
}

/* 确保徽标正常显示 */
:deep(.el-badge) {
  display: inline-block;
  vertical-align: middle;
}

/* 调整徽标和图标的位置 */
:deep(.el-badge .el-badge__content) {
  z-index: 2;
}
</style> 