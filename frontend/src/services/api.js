import apiService from './apiService'

// 获取API列表
export const getAPIList = async (params) => {
  const response = await apiService.get('/apis', { params })
  return response.data.data
}

// 获取单个API详情
export const getAPI = async (id) => {
  const response = await apiService.get(`/apis/${id}`)
  return response.data.data
}

// 创建API
export const createAPI = async (data) => {
  const response = await apiService.post('/apis', data)
  return response.data.data
}

// 更新API
export const updateAPI = async (id, data) => {
  const response = await apiService.put(`/apis/${id}`, data)
  return response.data.data
}

// 删除API
export const deleteAPI = async (id) => {
  const response = await apiService.delete(`/apis/${id}`)
  return response.data.data
}

// 测试API
export const testAPI = async (id, data) => {
  const response = await apiService.post(`/apis/${id}/test`, data)
  return response.data.data
} 