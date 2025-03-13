import apiService from './apiService';

export const systemMonitorApi = {
  // 获取系统监控列表
  getSystemMonitors: (params) => {
    return apiService.get('/system/monitor/list', { params })
      .then(res => res.data.data);
  },
  
  // 获取最新的系统监控信息
  getLatestSystemMonitor: () => {
    return apiService.get('/system/monitor/latest')
      .then(res => res.data.data);
  },
  
  // 手动收集系统信息
  collectSystemInfo: () => {
    return apiService.post('/system/monitor/collect')
      .then(res => res.data.data);
  }
}; 