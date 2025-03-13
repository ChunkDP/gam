import apiService from './apiService'


// 获取配置组列表
export const getConfigGroups = async () => {
  const res = await apiService.get('/configs/groups')
  return res.data
}

// 获取配置项列表
export const getConfigItems = async (groupId) => {
  const res = await apiService.get(`/configs/items/${groupId}`)
  return res.data
}

// 更新单个配置
export const updateConfigValue = async (data) => {
  const res = await apiService.put('/configs/value', data)
  return res.data
}

// 批量更新配置
export const batchUpdateConfigs = async (data) => {
  console.log(data)
  const res = await apiService.put('/configs/batch', data)
  return res.data
} 