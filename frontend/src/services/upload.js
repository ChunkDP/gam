import apiService from './apiService'

/**
 * 上传服务
 */
export const uploadFile = async (data) => {
  const response = await apiService.post('/upload', data, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response.data.data
}

/**
 * 批量上传文件
 * @param {FormData} data 包含多个文件的 FormData 对象
 * @returns {Promise<Array>} 上传结果数组
 */
export const batchUploadFiles = async (data) => {
  const response = await apiService.post('/upload/batch', data, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response.data.data
}

export const uploadImage = async (data) => {
  const response = await apiService.post('/upload/image', data, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response.data.data
}

export const getUploadConfig = async () => {
  const response = await apiService.get('/upload/config')
  return response.data.data
}

export const deleteFile = async (data) => {
  const response = await apiService.delete('/upload', { data })
  return response.data.data
}

// 默认导出对象
export default {
  uploadFile,
  uploadImage,
  getUploadConfig,
  deleteFile
} 