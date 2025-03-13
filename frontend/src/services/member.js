import apiService from './apiService'

// 会员管理相关接口
export const getMemberList = async (params) => {
  const res = await apiService.get('/members', { params })
  
  return res.data.data
}

export const createMember = async (data) => {
  const res = await apiService.post('/members', data)
  return res.data.data
}

export const updateMember = async (id, data) => {
  const res = await apiService.put(`/members/${id}`, data)
  return res.data.data
}

export const deleteMember = async (id) => {
  const res = await apiService.delete(`/members/${id}`)
  return res.data.data
}

// 检查角色字段是否唯一
export const checkMemberFieldUnique = async (field, value, excludeId) => {
  const params = { field, value }
  if (excludeId) {
    params.excludeId = excludeId
  }
  const res = await apiService.get('/members/check-field', { params })
  return res.data.data
}





// 会员标签相关接口
export const getMemberTags = async () => {
  const res = await apiService.get('/members/tags')
  return res.data.data
}

export const createMemberTag = async (data) => {
  
  const res = await apiService.post('/members/tags', data)
  return res.data.data
}

export const updateMemberTag = async (id, data) => {
  const res = await apiService.put(`/members/tags/${id}`, data)
  return res.data.data
}

export const deleteMemberTag = async (id) => {
  const res = await apiService.delete(`/members/tags/${id}`)
  return res.data.data
}

// 会员等级相关接口
export const getLevelList = async () => {
  const res = await apiService.get('/member/levels')
  return res.data.data
}

export const createLevel = async (data) => {
  const res = await apiService.post('/member/levels', data)
  return res.data.data
}

export const updateLevel = async (id, data) => {
  const res = await apiService.put(`/member/levels/${id}`, data)
  return res.data.data
}

export const deleteLevel = async (id) => {
  const res = await apiService.delete(`/member/levels/${id}`)
  return res.data.data
}

// 会员积分相关接口
export const getPointsList = async (params) => {
  const res = await apiService.get('/member/points', { params })
  return res.data.data
}

export const adjustPoints = async (data) => {
  const res = await apiService.post('/member/points/adjust', data)
  return res.data.data
}
