<template>
  <div class="home-container">

    <!-- 顶部数据卡片 -->
    <el-row :gutter="20">
      <el-col :span="6" v-for="(item, index) in statisticsData" :key="index">
        <el-card shadow="hover" class="statistics-card">
          <div class="card-content">
            <el-icon class="card-icon" :size="40">
              <component :is="item.icon"></component>
            </el-icon>
            <div class="card-info">
              <div class="card-value">{{ item.value }}</div>
              <div class="card-label">{{ item.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统监控卡片 -->
    <el-row :gutter="20" class="monitor-row">
      <el-col :span="24">
        <el-card shadow="hover" v-loading="monitorLoading">
          <template #header>
            <div class="card-header">
              <span>系统资源监控</span>
              <el-button type="primary" size="small" @click="refreshMonitorData">
                刷新
              </el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="monitor-item">
                <div class="monitor-title">CPU使用率</div>
                <el-progress 
                  type="dashboard" 
                  :percentage="monitorData.cpu_usage ? parseFloat(monitorData.cpu_usage.toFixed(2)) : 0" 
                  :color="getProgressColor"
                />
                <div class="monitor-value">{{ monitorData.cpu_usage ? monitorData.cpu_usage.toFixed(2) : 0 }}%</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="monitor-item">
                <div class="monitor-title">内存使用率</div>
                <el-progress 
                  type="dashboard" 
                  :percentage="monitorData.memory_usage ? parseFloat(monitorData.memory_usage.toFixed(2)) : 0" 
                  :color="getProgressColor"
                />
                <div class="monitor-value">{{ monitorData.memory_usage ? monitorData.memory_usage.toFixed(2) : 0 }}%</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="monitor-item">
                <div class="monitor-title">磁盘使用率</div>
                <el-progress 
                  type="dashboard" 
                  :percentage="monitorData.disk_usage ? parseFloat(monitorData.disk_usage.toFixed(2)) : 0" 
                  :color="getProgressColor"
                />
                <div class="monitor-value">{{ monitorData.disk_usage ? monitorData.disk_usage.toFixed(2) : 0 }}%</div>
              </div>
            </el-col>
          </el-row>
          <el-row :gutter="20" class="monitor-info">
            <el-col :span="8">
              <div class="info-item">
                <div class="info-label">网络IO：</div>
                <div class="info-value">{{ monitorData.network_io || '暂无数据' }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="info-item">
                <div class="info-label">进程数：</div>
                <div class="info-value">{{ monitorData.process_count || 0 }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="info-item">
                <div class="info-label">负载均衡：</div>
                <div class="info-value">{{ monitorData.load_average || '暂无数据' }}</div>
              </div>
            </el-col>
          </el-row>
          <div class="monitor-footer">
            <el-link type="primary" @click="goToMonitorPage">查看详细监控 <el-icon><ArrowRight /></el-icon></el-link>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>近7天访问趋势</span>
            </div>
          </template>
          <div class="chart" ref="visitChart"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>用户分布</span>
            </div>
          </template>
          <div class="chart" ref="userChart"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统信息 -->
    <el-row :gutter="20" class="system-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>系统公告</span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(activity, index) in activities"
              :key="index"
              :timestamp="activity.timestamp"
              :type="activity.type"
            >
              {{ activity.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>待办事项</span>
            </div>
          </template>
          <el-table :data="todoList" style="width: 100%">
            <el-table-column prop="title" label="标题"></el-table-column>
            <el-table-column prop="priority" label="优先级" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.priority === '高' ? 'danger' : scope.row.priority === '中' ? 'warning' : 'info'">
                  {{ scope.row.priority }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="deadline" label="截止日期" width="120"></el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import * as echarts from 'echarts';
import { User, ShoppingCart, Ticket, Money, ArrowRight } from '@element-plus/icons-vue';
import { systemMonitorApi } from '@/services/systemMonitor';
import { useNavigationStore } from '@/stores/navigation';
const navigationStore = useNavigationStore();


const router = useRouter();

// 统计数据
const statisticsData = ref([
  { label: '总用户数', value: '1,234', icon: 'User' },
  { label: '总订单数', value: '856', icon: 'ShoppingCart' },
  { label: '本月收入', value: '￥45,678', icon: 'Money' },
  { label: '待处理工单', value: '12', icon: 'Ticket' },
]);

// 系统公告
const activities = ref([
  { content: '系统更新维护通知', timestamp: '2024-03-20', type: 'primary' },
  { content: '新功能上线公告', timestamp: '2024-03-18', type: 'success' },
  { content: '安全更新提醒', timestamp: '2024-03-15', type: 'warning' },
]);

// 待办事项
const todoList = ref([
  { title: '系统升级', priority: '高', deadline: '2024-03-25' },
  { title: '数据备份', priority: '中', deadline: '2024-03-26' },
  { title: '用户反馈处理', priority: '低', deadline: '2024-03-27' },
]);

// 系统监控数据
const monitorData = ref({});
const monitorLoading = ref(false);
const refreshTimer = ref(null);

// 获取系统监控数据
const getMonitorData = async () => {
  try {
    monitorLoading.value = true;
    const data = await systemMonitorApi.getLatestSystemMonitor();
    monitorData.value = data;
  } catch (error) {
    console.error('获取系统监控数据失败:', error);
  } finally {
    monitorLoading.value = false;
  }
};

// 刷新监控数据
const refreshMonitorData = () => {
  getMonitorData();
};

// 跳转到监控详情页
const goToMonitorPage = () => {
  
  navigationStore.navigateToPage({
    title: '系统监控',
    componentName: 'SystemMonitor',
    parentName: 'System'
  });
  
};

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage < 60) {
    return '#67C23A';
  } else if (percentage < 80) {
    return '#E6A23C';
  } else {
    return '#F56C6C';
  }
};

// 图表初始化
onMounted(() => {
  // 获取系统监控数据
  getMonitorData();
  
  // 设置定时刷新
  refreshTimer.value = setInterval(() => {
    getMonitorData();
  }, 60000); // 每分钟刷新一次
  
  // 访问趋势图表
  const visitChart = echarts.init(document.querySelector('.chart'));
  visitChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
    yAxis: { type: 'value' },
    series: [{
      data: [820, 932, 901, 934, 1290, 1330, 1320],
      type: 'line',
      smooth: true
    }]
  });

  // 用户分布图表
  const userChart = echarts.init(document.querySelectorAll('.chart')[1]);
  userChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: [
        { value: 435, name: '新用户' },
        { value: 679, name: '活跃用户' },
        { value: 120, name: '沉睡用户' },
      ]
    }]
  });
  
  // 监听窗口大小变化，重绘图表
  window.addEventListener('resize', () => {
    visitChart.resize();
    userChart.resize();
  });
});

onUnmounted(() => {
  // 清除定时器
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value);
  }
});
</script>

<style lang="scss">

.home-container {
  padding: 20px;
}

.statistics-card {
  margin-bottom: 20px;
}

.card-content {
  display: flex;
  align-items: center;
}

.card-icon {
  margin-right: 15px;
  color: #409EFF;
}

.card-info {
  flex-grow: 1;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 5px;
}

.card-label {
  font-size: 14px;
  color: #909399;
}

.monitor-row {
  margin-bottom: 20px;
}

.monitor-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 0;
}

.monitor-title {
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: bold;
}

.monitor-value {
  font-size: 18px;
  margin-top: 10px;
  font-weight: bold;
}

.monitor-info {
  margin-top: 20px;
}

.info-item {
  display: flex;
  padding: 5px 0;
}

.info-label {
  font-weight: bold;
  width: 100px;
}

.info-value {
  flex: 1;
}

.monitor-footer {
  margin-top: 15px;
  text-align: right;
}

.chart-row {
  margin-bottom: 20px;
}

.chart {
  height: 300px;
}

.system-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>