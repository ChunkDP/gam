// stores/userPermissions.js
import { defineStore } from 'pinia';
import { fetchRoleMenuData } from '@/services/apiService';
export const useUserPermissionsStore = defineStore('userPermissions', {
  state: () => ({
    isLoggedIn: localStorage.getItem('isLoggedIn') === 'true',
    roleMenu: [],
    permissions: [], // 存储用户的权限列表
  }),
  getters: {

    getRoleMenu: (state) => state.roleMenu,

  },
  actions: {

    setLoggedIn(value) {

      localStorage.setItem('isLoggedIn', value);
    },

    setToken(accessToken, refreshToken) {

      sessionStorage.setItem('token', accessToken);
      if (refreshToken) {
        localStorage.setItem('refreshToken', refreshToken);
      }
    },

    clearAuth() {
      localStorage.removeItem('isLoggedIn');
      localStorage.removeItem('ParentactiveIndex');
      localStorage.removeItem('activeIntx');
      sessionStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('tabs');

      
    },

    getToken() {
      return sessionStorage.getItem('token');
    },


    async getRolePermissions() {
      try {


        const response = await fetchRoleMenuData();

          const tree = [];
          const menuMap = new Map();
          response.menus.forEach(item => {
            menuMap.set(item.id, { ...item, children: [] });
          });
          menuMap.forEach(item => {
            if (item.parent_id === 0 || !item.parent_id) {
              // 顶级菜单
              tree.push(item);
            } else {
              // 子菜单，找到父菜单并添加到其children中
              const parent = menuMap.get(item.parent_id);
              if (parent) {
                parent.children.push(item);
              }
            }
          });

          this.roleMenu = tree

          this.permissions = response.permissions;

          return response;
        
      } catch (error) {
        console.error('Failed to fetch role permissions:', error);
        return null;
      }
    },

    async login(loginResponse) {
      try {
        // 1. 设置登录状态和token

        this.setToken(loginResponse.token, loginResponse.refreshToken);
        // 3. 设置登录状态（仅在成功获取权限后）
        this.setLoggedIn(true);


      } catch (error) {

        throw error;
      }
    }
  }
});
