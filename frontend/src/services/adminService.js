import apiService from './apiService'

// 获取管理员列表
export const getAdminList = async (params) => {
  
  const response = await apiService.get('/admins', { params })
  return response.data.data
}

// 创建管理员
export const createAdmin = async (data) => {
  const response = await apiService.post('/admins', data)
  return response.data.data
}

// 更新管理员
export const updateAdmin = async (id, data) => {
  const response = await apiService.put(`/admins/${id}`, data)
  return response.data.data
}

// 删除管理员
export const deleteAdmin = async (id) => {
  const response = await apiService.delete(`/admins/${id}`)
  return response.data.data
}

// 更新管理员状态
export const updateAdminStatus = async (id, status) => {
  const response = await apiService.put(`/admins/${id}/status`, { status })
  return response.data.data
}

export default {
  getAdminList,
  createAdmin,
  updateAdmin,
  deleteAdmin,
  updateAdminStatus
} 