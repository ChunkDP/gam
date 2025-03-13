import apiClient from './apiService'

// 获取角色列表
export const getRoleList = async (params) => {
  
    const response = await apiClient.get('/roles', { params })
    return response.data.data

}

// 创建角色
export const createRole = async (data) => {
  
    const response = await apiClient.post('/roles', data)
    return response.data.data

}

// 更新角色
export const updateRole = async (id, data) => {
 
    const response = await apiClient.put(`/roles/${id}`, data)
    return response.data.data

}

// 删除角色
export const deleteRole = async (id) => {
 
    const response = await apiClient.delete(`/roles/${id}`)
    return response.data.data

}

// 获取单个角色详情
export const getRoleDetail = async (id) => {

    const response = await apiClient.get(`/roles/${id}`)
    return response.data.data

}

// 检查角色字段是否唯一
export const checkRoleFieldUnique = async (field, value, excludeId) => {
  const params = { field, value }
  if (excludeId) {
    params.excludeId = excludeId
  }
  const res = await apiClient.get('/roles/check-field', { params })
  return res.data.data
}

// 更新角色状态
export const updateRoleStatus = async (id, data) => {
  const res = await apiClient.put(`/roles/${id}/status`, data)
  return res.data.data
}

// 更新角色排序
export const updateRoleSort = async (id, data) => {
  const res = await apiClient.put(`/roles/${id}/sort`, data)
  return res.data.data
} 

// 获取角色的菜单权限
export const getRoleMenus = async (roleId) => {
 
  
  const response = await apiClient.get(`/roles/permissions/${roleId}/menus`)
  return response.data.data

}

// 更新角色的菜单权限
export const updateRoleMenus = async (roleId, menuIds) => {


  const response = await apiClient.put(`/roles/permissions/${roleId}/menus`, { menuIds })
  return response.data.data

} 