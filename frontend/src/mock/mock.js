// mock.js
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';

// 创建一个 axios 实例
const instance = axios.create();

// 创建一个 mock 实例
const mock = new MockAdapter(instance);

// 模拟获取菜单数据的请求
mock.onPost('/api/login').reply(200, {token:'token=faewefwfasdfasdfewfwqefwefwfewef',refreshToken :"dsdwewe"});
// 模拟获取菜单数据的请求
mock.onGet('/api/roleMenu').reply(200, [
 
  {id: 1,sort:1,parentName: '',title: '首页',icon: 'HomeFilled',path:'/home',name: 'Home',component:'Home'},
  {id: 2,sort:4,parentName: '',title: '超级管理员',icon: 'Avatar',path:'',name: 'Admin',component:''},
  {id: 3,sort:1,parentName: 'Admin',title: '用户管理',path: '/user-management', name: 'UserManagement', component: 'UserManagement' },
  {id: 4,sort:3,parentName: '',title: '功能区',icon: 'Grid',path:'',name: 'Function',component:''},
  {id: 5,sort:1,parentName: 'Function',title: '设置',path: '/settings', name: 'Settings', component: 'Settings' },

]);


export default instance;