import apiService from './apiService';

export const notificationApi = {
  // 通知类型管理
  getNotificationTypes: () => {
    return apiService.get('/notifications/types')
      .then(res => res.data.data);
  },
  
  getNotificationType: (id) => {
    return apiService.get(`/notifications/types/${id}`)
      .then(res => res.data.data);
  },
  
  createNotificationType: (data) => {
    return apiService.post('/notifications/types', data)
      .then(res => res.data.data);
  },
  
  updateNotificationType: (id, data) => {
    return apiService.put(`/notifications/types/${id}`, data)
      .then(res => res.data.data);
  },
  
  deleteNotificationType: (id) => {
    return apiService.delete(`/notifications/types/${id}`)
      .then(res => res.data.data);
  },
  
  // 通知管理
  getNotifications: (params) => {
    return apiService.get('/notifications', { params, user_type: params.userType  })
      .then(res => res.data.data);
  },
  
  getNotification: (id) => {
    return apiService.get(`/notifications/${id}`)
      .then(res => res.data.data);
  },
  
  createNotification: (data) => {
    return apiService.post('/notifications', data)
      .then(res => res.data.data);
  },
  
  updateNotification: (id, data) => {
    return apiService.put(`/notifications/${id}`, data)
      .then(res => res.data.data);
  },
  
  deleteNotification: (id) => {
    return apiService.delete(`/notifications/${id}`)
      .then(res => res.data.data);
  },
  
  publishNotification: (id) => {
    return apiService.post(`/notifications/${id}/publish`)
      .then(res => res.data.data);
  },
  
  recallNotification: (id) => {
    return apiService.post(`/notifications/${id}/recall`)
      .then(res => res.data.data);
  },
  
  // 用户通知
  getUserNotifications: (params) => {
    return apiService.get('/user/notifications', { params })
      .then(res => res.data.data);
  },
  
  markNotificationAsRead: (id) => {
    return apiService.post(`/user/notifications/${id}/read`)
      .then(res => res.data.data);
  },
  
  markAllNotificationsAsRead: () => {
    return apiService.post('/user/notifications/read-all')
      .then(res => res.data.data);
  },
  
  deleteUserNotification: (id) => {
    return apiService.delete(`/user/notifications/${id}`)
      .then(res => res.data.data);
  },
  
  getUnreadNotificationCount: () => {
    return apiService.get('/user/notifications/unread-count')
      .then(res => res.data.data);
  }
}; 