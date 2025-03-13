import apiService from './apiService'

// 获取日志列表
export const getLogList = async (params) => {
  const response = await apiService.get('/system/logs', { params })
  return response.data.data
}

// 删除指定时间之前的日志
export const deleteLogs = async (params) => {
  const response = await apiService.delete('/system/logs', { params })
  return response.data.data
}

// 导出日志
export const exportLogs = async (params) => {
  const response = await apiService.get('/system/logs/export', {
    params,
    responseType: 'blob'
  })
  return response.data
}

export default {
  getLogList,
  deleteLogs,
  exportLogs
}