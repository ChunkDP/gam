// src/services/apiService.js
import axios from 'axios';
import { useUserPermissionsStore } from '../stores/userPermissions';
import ElMessage from 'element-plus';
const API_BASE_URL = '/gam';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    
    const token = sessionStorage.getItem('token');
   // console.log("apiClient.interceptors.request--token:",token);
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
// 在现有的 apiClient.interceptors.response.use 中添加错误处理
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // HTTP 错误处理
    if (error.response) {
      const { status, data } = error.response;
      const userPermissionsStore = useUserPermissionsStore();

      switch (status) {
        case 400:
          ElMessage.error('请求参数错误');
          break;
        case 401:
          ElMessage.error('未授权，请重新登录');
          // 处理token过期   // console.log("apiClient.interceptors.response", response)
   // console.log("apiClient.interceptors.error", error)
    

     
        // token 过期
        const refreshToken = userPermissionsStore.refreshToken;
        if (refreshToken) {
          try {
            const newTokens = await refreshToken(refreshToken);
            userPermissionsStore.setToken(newTokens.token, newTokens.refreshToken);
            // 重试失败的请求
            error.config.headers['Authorization'] = `Bearer ${newTokens.token}`;
            return apiClient(error.config);
          } catch (refreshError) {
           // userPermissionsStore.clearAuth();
          // window.location.href = '/';
            return Promise.reject(refreshError);
          }
        } else {
          // userPermissionsStore.clearAuth();
         // window.location.href = '/';
        }
      
    
          break;
        case 403:
          ElMessage.error('拒绝访问');
          break;
        case 404:
          ElMessage.error('请求的资源不存在');
          break;
        case 500:
          ElMessage.error('服务器内部错误');
          break;
        default:
          ElMessage.error('未知错误');
      }
    } else if (error.request) {
      ElMessage.error('网络错误，请检查网络连接');
    } else {
      ElMessage.error('请求配置错误');
    }
    
    return Promise.reject(error);
  }
)

export const login = async (credentials) => {
  const response = await apiClient.post('/login', credentials);
  return response.data.data;
};

export const refreshToken = async (refreshToken) => {
  const response = await apiClient.post('/refresh-token', { refreshToken });
  return response.data.data;
};

export const fetchRoleMenuData = async () => {
 
    const response = await apiClient.get('/authmenus');
    return response.data.data;  // 返回完整的响应数据，包含 code/msg/data

};


// 导出 apiClient 以便其他地方可以直接使用
export default apiClient;