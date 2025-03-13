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

<script>
import { defineComponent, onMounted, ref } from 'vue';
import * as echarts from 'echarts';
import { User, ShoppingCart, Ticket, Money } from '@element-plus/icons-vue';

export default defineComponent({
  name: 'Home',
  setup() {
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

    // 图表初始化
    onMounted(() => {
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
    });

    return {
      statisticsData,
      activities,
      todoList,
    };
  }
});
</script>

<style scoped>
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