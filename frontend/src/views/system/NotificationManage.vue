<template>
  <div class="ma-search-box">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>通知管理</span>
          <el-button type="primary" @click="openCreateDialog">创建通知</el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <el-form :model="searchForm" inline>
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="请输入标题" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type_id" placeholder="请选择类型" clearable class="adaptive-select">
            <el-option
              v-for="type in notificationTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="重要程度">
          <el-select v-model="searchForm.level" placeholder="请选择重要程度" clearable class="adaptive-select">
            <el-option :value="1" label="普通" />
            <el-option :value="2" label="重要" />
            <el-option :value="3" label="紧急" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable class="adaptive-select">
            <el-option :value="0" label="草稿" />
            <el-option :value="1" label="已发布" />
            <el-option :value="2" label="已撤回" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchNotifications">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 通知列表 -->
      <el-table
        v-loading="loading"
        :data="notifications"
        border
        style="width: 100%"
      >
        <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
        <el-table-column label="类型" min-width="100">
          <template #default="{ row }">
            {{ row.type ? row.type.name : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="重要程度" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.level === 3 ? 'danger' : (row.level === 2 ? 'warning' : 'info')"
            >
              {{ row.level === 3 ? '紧急' : (row.level === 2 ? '重要' : '普通') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 1 ? 'success' : (row.status === 2 ? 'info' : 'primary')"
            >
              {{ row.status === 1 ? '已发布' : (row.status === 2 ? '已撤回' : '草稿') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="发送范围" width="120">
          <template #default="{ row }">
            {{ row.receiver_type === 'all' ? '全部用户' : (row.receiver_type === 'members' ? '仅会员' : '仅管理员') }}
          </template>
        </el-table-column>
        <el-table-column prop="read_count" label="已读数量" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              type="primary"
              size="small"
              @click="publishNotification(row)"
            >
              发布
            </el-button>
            <el-button
              v-if="row.status === 1"
              type="warning"
              size="small"
              @click="recallNotification(row)"
            >
              撤回
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="viewNotification(row)"
            >
              查看
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="primary"
              size="small"
              @click="editNotification(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="danger"
              size="small"
              @click="deleteNotification(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination">
        <el-pagination background
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>  

<!-- 通知表单对话框 -->
<el-dialog
  v-model="dialogVisible"
  :title="isEdit ? '编辑通知' : '创建通知'"
  width="700px"
>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
  >
    <el-form-item label="标题" prop="title">
      <el-input v-model="form.title" placeholder="请输入通知标题" />
    </el-form-item>

    <el-form-item label="通知类型" prop="type_id">
      <el-select v-model="form.type_id" placeholder="请选择通知类型" style="width: 100%">
        <el-option
          v-for="type in notificationTypes"
          :key="type.id"
          :label="type.name"
          :value="type.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="重要程度" prop="level">
      <el-select v-model="form.level" placeholder="请选择重要程度" style="width: 100%">
        <el-option :value="1" label="普通" />
        <el-option :value="2" label="重要" />
        <el-option :value="3" label="紧急" />
      </el-select>
    </el-form-item>

    <el-form-item label="接收范围" prop="receiver_type">
      <el-select v-model="form.receiver_type" placeholder="请选择接收范围" style="width: 100%">
        <el-option value="all" label="所有用户" />
        <el-option value="members" label="仅会员" />
        <el-option value="admins" label="仅管理员" />
      </el-select>
    </el-form-item>

    <el-form-item label="过期时间" prop="expiration_time">
      <el-date-picker
        v-model="form.expiration_time"
        type="datetime"
        placeholder="请选择过期时间"
        style="width: 100%"
        :min-date="new Date()"
        value-format="YYYY-MM-DD HH:mm:ss"
      />
    </el-form-item>

    <el-form-item label="通知内容" prop="content">
      <el-input
        v-model="form.content"
        type="textarea"
        :rows="6"
        placeholder="请输入通知内容"
      />
    </el-form-item>
  </el-form>
  <template #footer>
    <span class="dialog-footer">
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitForm">确认</el-button>
    </span>
  </template>
</el-dialog>
  
  <!-- 通知详情对话框 -->
  <el-dialog
    v-model="detailDialogVisible"
    title="通知详情"
    width="800px"
    class="notification-detail-dialog"
  >
    <div v-if="currentNotification" class="notification-detail">
      <!-- 基本信息卡片 -->
      <el-card class="mb-4">
        <template #header>
          <div class="card-header">
            <h3 class="notification-title">{{ currentNotification.title }}</h3>
            <el-tag
              :type="currentNotification.status === 1 ? 'success' : (currentNotification.status === 2 ? 'info' : 'primary')"
            >
              {{ currentNotification.status === 1 ? '已发布' : (currentNotification.status === 2 ? '已撤回' : '草稿') }}
            </el-tag>
          </div>
        </template>

        <div class="notification-meta">
          <el-tag
            :type="currentNotification.level === 3 ? 'danger' : (currentNotification.level === 2 ? 'warning' : 'info')"
            size="small"
          >
            {{ currentNotification.level === 3 ? '紧急' : (currentNotification.level === 2 ? '重要' : '普通') }}
          </el-tag>
          <span class="notification-type">{{ currentNotification.type?.name || '-' }}</span>
          <span class="notification-time">创建时间：{{ currentNotification.created_at }}</span>
        </div>

        <div class="notification-content">{{ currentNotification.content }}</div>

        <div class="notification-info">
          <div>发送范围：{{ currentNotification.receiver_type === 'all' ? '全部用户' : 
                        (currentNotification.receiver_type === 'members' ? '仅会员' : '仅管理员') }}</div>
          <div v-if="currentNotification.expiration_time">过期时间：{{ currentNotification.expiration_time }}</div>
        </div>
      </el-card>

      <!-- 阅读统计卡片 -->
      <el-card v-if="currentNotification.status !== 0" class="stats-card">
        <template #header>
          <div class="card-header">
            <span>阅读统计</span>
            <el-button type="primary" link @click="refreshStats(currentNotification.id)">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </template>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <el-statistic title="总接收人数" :value="notificationStats.total_receivers">
              <template #suffix>人</template>
            </el-statistic>
          </el-col>
          <el-col :span="8">
            <el-statistic title="已读人数" :value="notificationStats.read_count">
              <template #suffix>人</template>
            </el-statistic>
          </el-col>
          <el-col :span="8" v-if="currentNotification.status === 2">
            <el-statistic title="撤回前已读" :value="notificationStats.recalled_read">
              <template #suffix>人</template>
            </el-statistic>
          </el-col>
        </el-row>

        <div class="read-progress mt-4">
          <div class="progress-label">
            阅读率：{{ calculateReadRate(notificationStats.read_count, notificationStats.total_receivers) }}%
          </div>
          <el-progress 
            :percentage="calculateReadRate(notificationStats.read_count, notificationStats.total_receivers)"
            :status="getProgressStatus(notificationStats.read_count, notificationStats.total_receivers)"
          />
        </div>
      </el-card>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { notificationApi } from '@/services/notification';
import { getMemberList } from '@/services/member';
import { getRoleList } from '@/services/role';
import { Refresh } from '@element-plus/icons-vue';

// 数据定义
const notifications = ref([]);
const notificationTypes = ref([]);
const roles = ref([]);
const userOptions = ref([]);
const currentNotification = ref(null);
const loading = ref(false);
const userSearchLoading = ref(false);
const dialogVisible = ref(false);
const detailDialogVisible = ref(false);
const isEdit = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const dateRange = ref([]);

// 表单相关
const formRef = ref(null);
const form = reactive({
  title: '',
  content: '',
  type_id: '',
  level: 1,
  receiver_type: 'all', // 默认发送给所有用户
  expiration_time: null
});

const rules = {
  title: [{ required: true, message: '请输入通知标题', trigger: 'blur' }],
  type_id: [{ required: true, message: '请选择通知类型', trigger: 'change' }],
  level: [{ required: true, message: '请选择重要程度', trigger: 'change' }],
  receiver_type: [{ required: true, message: '请选择接收范围', trigger: 'change' }],
  content: [{ required: true, message: '请输入通知内容', trigger: 'blur' }]
};

// 搜索表单
const searchForm = reactive({
  title: '',
  type_id: '',
  level: '',
  status: '',
  start_time: '',
  end_time: '',
  page: 1,
  page_size: 10
});

// 添加统计数据
const notificationStats = ref({
  totalReceivers: 0,
  readCount: 0,
  recalledRead: 0
});

// 计算阅读率
const calculateReadRate = (readCount, total) => {
  if (!total) return 0;
  return Math.round((readCount / total) * 100);
};

// 获取进度条状态
const getProgressStatus = (readCount, total) => {
  const rate = calculateReadRate(readCount, total);
  if (rate >= 80) return 'success';
  if (rate >= 50) return 'warning';
  return 'exception';
};

// 刷新统计数据
const refreshStats = async (notificationId) => {
  try {
    const stats = await notificationApi.getNotificationStats(notificationId);
   
    notificationStats.value = stats;
  } catch (error) {
    console.error('Failed to refresh notification stats:', error);
    ElMessage.error('获取统计数据失败');
  }
};

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

// 加载角色列表
const loadRoles = async () => {
  try {
    const { data } = await getRoleList();
    roles.value = data;
  } catch (error) {
    console.error('Failed to load roles:', error);
    ElMessage.error('加载角色列表失败');
  }
};

// 远程搜索用户
const remoteSearchUsers = async (query) => {
  if (query.trim() === '') {
    userOptions.value = [];
    return;
  }
  
  userSearchLoading.value = true;
  try {
    const { data } = await getMemberList({ username: query });
    userOptions.value = data;
  } catch (error) {
    console.error('Failed to search users:', error);
  } finally {
    userSearchLoading.value = false;
  }
};

// 加载通知列表
const loadNotifications = async () => {
  loading.value = true;
  try {
    searchForm.page = currentPage.value;
    searchForm.page_size = pageSize.value;
    
    const data = await notificationApi.getNotifications(searchForm);
    
    notifications.value = data.notifications;
    total.value = data.total;
  } catch (error) {
    console.error('Failed to load notifications:', error);
    ElMessage.error('加载通知列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索通知
const searchNotifications = () => {
  currentPage.value = 1;
  loadNotifications();
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
  loadNotifications();
};

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size;
  loadNotifications();
};

const handleCurrentChange = (page) => {
  currentPage.value = page;
  loadNotifications();
};

// 打开创建对话框
const openCreateDialog = () => {
  isEdit.value = false;
  // 重置表单
  Object.keys(form).forEach(key => {
    form[key] = key === 'level' ? 1 : 
                (key === 'receiver_type' ? 'all' : 
                (key === 'expiration_time' ? null : ''));
  });
  dialogVisible.value = true;
};

// 编辑通知
const editNotification = (row) => {
  isEdit.value = true;
  Object.keys(form).forEach(key => {
    form[key] = row[key] !== undefined ? row[key] : 
                (key === 'level' ? 1 : 
                (key === 'receiver_type' ? 'all' : ''));
  });
  dialogVisible.value = true;
};

// 查看通知详情
const viewNotification = async (row) => {
  try {
    const [notification, stats] = await Promise.all([
      notificationApi.getNotification(row.id),
      notificationApi.getNotificationStats(row.id)
    ]);

    console.log(stats)
    currentNotification.value = notification;
    notificationStats.value = stats;
    detailDialogVisible.value = true;
  } catch (error) {
    console.error('Failed to get notification details:', error);
    ElMessage.error('获取通知详情失败');
  }
};

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const formData = { ...form };
        
        if (isEdit.value) {
          await notificationApi.updateNotification(currentNotification.value.id, formData);
          ElMessage.success('更新通知成功');
        } else {
          await notificationApi.createNotification(formData);
          ElMessage.success('创建通知成功');
        }
        
        dialogVisible.value = false;
        loadNotifications();
      } catch (error) {
        console.error('Failed to save notification:', error);
        ElMessage.error(isEdit.value ? '更新通知失败' : '创建通知失败');
      }
    }
  });
};

// 发布通知
const publishNotification = async (row) => {
  try {
    await ElMessageBox.confirm('确定要发布该通知吗？发布后将立即推送给目标用户。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await notificationApi.publishNotification(row.id);
    ElMessage.success('发布通知成功');
    loadNotifications();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to publish notification:', error);
      ElMessage.error('发布通知失败');
    }
  }
};

// 撤回通知
const recallNotification = async (row) => {
  try {
    await ElMessageBox.confirm('确定要撤回该通知吗？撤回后用户将不再看到该通知。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await notificationApi.recallNotification(row.id);
    ElMessage.success('撤回通知成功');
    loadNotifications();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to recall notification:', error);
      ElMessage.error('撤回通知失败');
    }
  }
};

// 删除通知
const deleteNotification = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该通知吗？删除后无法恢复。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'danger'
    });
    
    await notificationApi.deleteNotification(row.id);
    ElMessage.success('删除通知成功');
    loadNotifications();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete notification:', error);
      ElMessage.error('删除通知失败');
    }
  }
};

// 页面加载时初始化
onMounted(() => {
  loadNotificationTypes();
  loadRoles();
  loadNotifications();
});
</script>

<style scoped>




.notification-detail {
  padding: 0 20px;
}

.notification-title {
  margin-top: 0;
  margin-bottom: 10px;
  font-size: 20px;
}

.notification-meta {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  color: #909399;
  font-size: 14px;
}

.notification-type {
  margin-left: 10px;
  margin-right: 10px;
}

.notification-content {
  line-height: 1.6;
  white-space: pre-wrap;
  margin: 15px 0;
}

.notification-footer {
  color: #909399;
  font-size: 14px;
}

.notification-footer > div {
  margin-bottom: 5px;
}

.notification-detail-dialog {
  .stats-card {
    margin-top: 20px;
  }

  .read-progress {
    margin-top: 20px;
    
    .progress-label {
      margin-bottom: 10px;
      color: #606266;
    }
  }


}

.notification-info {
  margin-top: 15px;
  color: #606266;
  font-size: 14px;
  
  > div {
    margin-bottom: 5px;
  }
}
</style> 