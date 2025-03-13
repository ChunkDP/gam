import apiClient from './apiService'

// 获取菜单树
export const getMenuTree = async () => {
  const response = await apiClient.get('/menus/tree')
  return response.data.data
}

// 获取菜单列表（不分层级）
export const getMenuList = async() => {
  const response = await apiClient.get('/menus')
  return response.data.data
}

// 创建菜单
export const createMenu = async (data) => {
  const response = await apiClient.post('/menus', data)
  return response.data.data
}

// 更新菜单
export const updateMenu = async (id, data) => {
  const response = await apiClient.put(`/menus/${id}`, data)
  return response.data.data
}

// 删除菜单
export const deleteMenu = async (id) => {
  const response = await apiClient.delete(`/menus/${id}`)
  return response.data.data
}

// 更新菜单状态
export const updateMenuStatus = async (id, data) => {
  const response = await apiClient.patch(`/menus/${id}/status`, data)
  return response.data.data
  }

// 更新菜单排序
export const updateMenuSort = async (id, data) => {
  const response = await apiClient.patch(`/menus/${id}/sort`, data)
  return response.data.data
}

// 更新菜单显示状态
export const updateMenuHidden = async (id, data) => {
  const response = await apiClient.patch(`/menus/${id}/hidden`, data)
  return response.data.data
} 