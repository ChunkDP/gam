// router/index.js
import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import Layout from '../views/Layout.vue';
import Page404 from '../views/404.vue';
const routes = [
  {

    path: "/",
    redirect: "/login"
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: Page404
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/layout',
    name: 'Layout',
    component: Layout,
    meta: { requiresAuth: true},
  },
  // 其他路由配置
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});
import { useUserPermissionsStore } from '../stores/userPermissions';

// 动态加载路由
export async function loadDynamicRoutes() {
  //加载动态路由时，先恢复登录状态，再获取菜单，如果没有登录，跳转登录页
  const userPermissionsStore = useUserPermissionsStore();

  
  const menuData = await userPermissionsStore.getRolePermissions();

  if (menuData?.menus) {

    const routes = [];
   
    menuData.menus.filter(menu=>menu.component!='').forEach(menu => {
   
        routes.push({
          path: '/layout'+menu.path,
          name: menu.name,
          component: () => import(/* @vite-ignore */`../views/${menu.component}.vue`),
          meta: { 
            title: menu.title,
            icon: menu.icon,
            requiresAuth: true
          }
        })
  

       

    });






    router.addRoute({
      path: '/layout',
      name: 'Layout',
      component: Layout,
      children: routes,
      meta: { requiresAuth: true},
    })

   
  }


}

router.beforeEach(async (to, from, next) => {
  const userPermissionsStore = useUserPermissionsStore();
    // 未登录跳转登录页
    const token = sessionStorage.getItem('token');

      // 没有token且不是登录页，跳转到登录页
  if (!token && to.path !== '/login') {
   
    return next('/login');
  }

  if(token && userPermissionsStore.roleMenu.length === 0){ 
     
    await loadDynamicRoutes();
  
  }
    // 如果是登录页且有token，跳转到首页
    if (to.path === '/login' && token) {

    
      return next('/layout');
    }

  try {
    // 有token但路由不存在（可能是刷新导致的动态路由丢失）
    if (token && to.name === 'NotFound') {


      
    
      // 保存原始请求的路径
      const originalPath = to.path;
         return next({
          path: originalPath,
          replace: true
        });
   
      // const matchedRoute = router.resolve(originalPath);

      // if (matchedRoute && matchedRoute.name !== 'NotFound') {
      //   return next({
      //     path: originalPath,
      //     replace: true
      //   });
      // } else {
       
      //   return next('/layout'); // 如果路由加载后仍然不存在，跳转到首页
      // }
    }
    next();
  } catch (error) {
    console.error('Failed to load routes:', error);
    next('/login');
  }
});

export default router;