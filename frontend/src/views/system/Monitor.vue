<template>
  <div class="ma-search-box">
    <el-card >
      <template #header>
        <div class="card-header">
          <span>系统实时监控</span>
          <div>
            <el-button type="primary" size="small" @click="refreshData">
              刷新数据
            </el-button>
            <el-button type="success" size="small" @click="collectSystemInfo">
              手动采集
            </el-button>
          </div>
        </div>
      </template>
      
      <el-row :gutter="20" v-loading="loading">
        <el-col :span="8">
          <div class="monitor-item">
            <div class="monitor-title">CPU使用率</div>
            <el-progress 
              type="dashboard" 
              :percentage="latestData.cpu_usage ? parseFloat(latestData.cpu_usage.toFixed(2)) : 0" 
              :color="getProgressColor"
            />
            <div class="monitor-value">{{ latestData.cpu_usage ? latestData.cpu_usage.toFixed(2) : 0 }}%</div>
          </div>
        </el-col>
        
        <el-col :span="8">
          <div class="monitor-item">
            <div class="monitor-title">内存使用率</div>
            <el-progress 
              type="dashboard" 
              :percentage="latestData.memory_usage ? parseFloat(latestData.memory_usage.toFixed(2)) : 0" 
              :color="getProgressColor"
            />
            <div class="monitor-value">{{ latestData.memory_usage ? latestData.memory_usage.toFixed(2) : 0 }}%</div>
          </div>
        </el-col>
        
        <el-col :span="8">
          <div class="monitor-item">
            <div class="monitor-title">磁盘使用率</div>
            <el-progress 
              type="dashboard" 
              :percentage="latestData.disk_usage ? parseFloat(latestData.disk_usage.toFixed(2)) : 0" 
              :color="getProgressColor"
            />
            <div class="monitor-value">{{ latestData.disk_usage ? latestData.disk_usage.toFixed(2) : 0 }}%</div>
          </div>
        </el-col>
      </el-row>
      
      <el-divider />
      
      <el-row :gutter="20">
        <el-col :span="8">
          <div class="info-item">
            <div class="info-label">网络IO：</div>
            <div class="info-value">{{ latestData.network_io || '暂无数据' }}</div>
          </div>
        </el-col>
        
        <el-col :span="8">
          <div class="info-item">
            <div class="info-label">进程数：</div>
            <div class="info-value">{{ latestData.process_count || 0 }}</div>
          </div>
        </el-col>
        
        <el-col :span="8">
          <div class="info-item">
            <div class="info-label">负载均衡：</div>
            <div class="info-value">{{ latestData.load_average || '暂无数据' }}</div>
          </div>
        </el-col>
      </el-row>
      
      <el-divider />
      
      <div class="chart-container">
        <div class="chart-title">系统资源使用趋势</div>
        <div ref="chartRef" style="width: 100%; height: 400px;"></div>
      </div>
    </el-card>
    
    <el-card class="history-card">
      <template #header>
        <div class="card-header">
          <span>历史监控数据</span>
          <el-form :inline="true" :model="queryParams" class="query-form">
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="dateRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="getHistoryData">查询</el-button>
              <el-button @click="resetQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </template>
      
      <el-table :data="historyData" style="width: 100%" v-loading="tableLoading" border>
        <el-table-column prop="createdAt" label="采集时间" width="180" sortable>
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="cpuUsage" label="CPU使用率" width="120">
          <template #default="scope">
            {{ scope.row.cpu_usage.toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column prop="memoryUsage" label="内存使用率" width="120">
          <template #default="scope">
            {{ scope.row.memory_usage.toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column prop="diskUsage" label="磁盘使用率" width="120">
          <template #default="scope">
            {{ scope.row.disk_usage.toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column prop="network_io" label="网络IO" width="200" />
        <el-table-column prop="process_count" label="进程数" width="100" />
        <el-table-column prop="load_average" label="负载均衡" />
      </el-table>
      
      <div class="pagination">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-size="queryParams.limit"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue';
import { ElMessage } from 'element-plus';
import { systemMonitorApi } from '@/services/systemMonitor';
import * as echarts from 'echarts';
import dayjs from 'dayjs';

// 数据定义
const loading = ref(false);
const tableLoading = ref(false);
const latestData = ref({});
const historyData = ref([]);
const total = ref(0);
const chartRef = ref(null);
const chart = ref(null);
const dateRange = ref([]);
const refreshTimer = ref(null);

const queryParams = reactive({
  startTime: '',
  endTime: '',
  limit: 10,
  page: 1
});

// 获取最新监控数据
const getLatestData = async () => {
  try {
    loading.value = true;
    const data = await systemMonitorApi.getLatestSystemMonitor();
    
    latestData.value = data;
   
  } catch (error) {
    ElMessage.error('获取最新监控数据失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 获取历史监控数据
const getHistoryData = async () => {
  try {
    tableLoading.value = true;
    
    if (dateRange.value && dateRange.value.length === 2) {
      queryParams.startTime = dateRange.value[0];
      queryParams.endTime = dateRange.value[1];
    }
    
    const data = await systemMonitorApi.getSystemMonitors(queryParams);
    historyData.value = data;
    total.value = data.length; // 实际项目中可能需要从后端获取总数
    
    // 更新图表
    updateChart();
  } catch (error) {
    ElMessage.error('获取历史监控数据失败');
    console.error(error);
  } finally {
    tableLoading.value = false;
  }
};

// 手动收集系统信息
const collectSystemInfo = async () => {
  try {
    loading.value = true;
    await systemMonitorApi.collectSystemInfo();
    ElMessage.success('系统信息采集成功');
    // 刷新数据
    await getLatestData();
    await getHistoryData();
  } catch (error) {
    ElMessage.error('系统信息采集失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 刷新数据
const refreshData = async () => {
  await getLatestData();
  await getHistoryData();
};

// 重置查询条件
const resetQuery = () => {
  dateRange.value = [];
  queryParams.startTime = '';
  queryParams.endTime = '';
  queryParams.page = 1;
  getHistoryData();
};

// 分页处理
const handleSizeChange = (size) => {
  queryParams.limit = size;
  getHistoryData();
};

const handleCurrentChange = (page) => {
  queryParams.page = page;
  getHistoryData();
};

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
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

// 初始化图表
const initChart = () => {
  if (chart.value) {
    chart.value.dispose();
  }
  
  nextTick(() => {
    if (chartRef.value) {
      chart.value = echarts.init(chartRef.value);
      updateChart();
    }
  });
};

// 更新图表数据
const updateChart = () => {
  if (!chart.value || historyData.value.length === 0) return;
  
  // 反转数据以便按时间顺序显示
  const chartData = [...historyData.value].reverse();
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    legend: {
      data: ['CPU使用率', '内存使用率', '磁盘使用率']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: chartData.map(item => dayjs(item.createdAt).format('HH:mm:ss'))
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '{value}%'
      },
      max: 100
    },
    series: [
      {
        name: 'CPU使用率',
        type: 'line',
        data: chartData.map(item => parseFloat(item.cpu_usage.toFixed(2)))
      },
      {
        name: '内存使用率',
        type: 'line',
        data: chartData.map(item => parseFloat(item.memory_usage.toFixed(2)))
      },
      {
        name: '磁盘使用率',
        type: 'line',
        data: chartData.map(item => parseFloat(item.disk_usage.toFixed(2)))
      }
    ]
  };
  
  chart.value.setOption(option);
};

// 自动刷新数据
const startAutoRefresh = () => {
  refreshTimer.value = setInterval(() => {
    getLatestData();
  }, 60000); // 每分钟刷新一次
};

// 生命周期钩子
onMounted(async () => {
  await getLatestData();
  await getHistoryData();
  initChart();
  startAutoRefresh();
  
  // 监听窗口大小变化，重绘图表
  window.addEventListener('resize', () => {
    if (chart.value) {
      chart.value.resize();
    }
  });
});

onUnmounted(() => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value);
  }
  
  if (chart.value) {
    chart.value.dispose();
    chart.value = null;
  }
  
  window.removeEventListener('resize', () => {
    if (chart.value) {
      chart.value.resize();
    }
  });
});
</script>

<style scoped>



.monitor-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
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

.info-item {
  display: flex;
  padding: 10px 0;
}

.info-label {
  font-weight: bold;
  width: 100px;
}

.info-value {
  flex: 1;
}

.chart-container {
  margin-top: 20px;
}

.chart-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.query-form {
  margin-bottom: 0;
}
</style> 