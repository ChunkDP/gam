<!-- Layout.vue -->
<template>
  <el-container class="layout">
    <el-header class="el-header dark:bg-slate-900" style="height: 50px;">

      <div class="header">
        <div class="logo">
          <el-image :src="LogoJpg" />
          <span @click="goToMonitorPage">Go-admin-manage</span>
         
          
        </div>
        <div class="user-info">
           <!-- 添加通知中心组件 -->
           <NotificationCenter />
          <el-dropdown trigger="click" popper-class="popclass">
            <span
              class="cursor-pointer flex justify-center items-center text-black dark:text-gray-100 el-dropdown-link">
              <el-avatar :size="40" :src="AvatarJpg"></el-avatar>
              <span>用户名 <el-icon>
                  <arrow-down />
                </el-icon></span>

            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>

    <el-container>
      <el-aside :class="isCollapsed ? 'aside64' : 'aside250'">
        <el-menu :default-active="navigationStore.activeIndex" class="el-menu-vertical-demo dark:text-slate-300 overflow-hidden"
          :unique-opened="true" style="border: none;" active-text-color="#fff" :collapse="isCollapsed">
          <template v-for="item in menuData" :key="item.name">

            <el-tooltip v-if="isCollapsed && !(item.children?.length)" effect="dark" :content="item.title"
              placement="right">

              <el-menu-item v-if="!(item.children?.length)" :index="item.name"
                @click="onMenuItemClick(item.title, item.name)" :class="navigationStore.activeIndex == item.name ? 'singlemenu' : ''">
                <el-icon>
                  <component :is="item.icon" />
                </el-icon>

              </el-menu-item>


            </el-tooltip>


            <el-menu-item v-else-if="!isCollapsed && !(item.children?.length)" :index="item.name"
              @click="onMenuItemClick(item.title, item.name)">
              <el-icon>
                <component :is="item.icon" />
              </el-icon>
              <span>{{ item.title }}</span>
            </el-menu-item>





            <el-popover v-else-if="isCollapsed && item.children?.length" placement="right-end" :width="200"
              trigger="click" :offset=5 transition="el-slide-in-right" :show-arrow="false" :hide-after=0
              popper-class="popperclass">
              <template #reference>
                <el-menu-item :index="item.name" :class="navigationStore.parentActiveIndex == item.name ? 'parentmenu' : ''">
                  <el-icon>
                    <component :is="item.icon" />
                  </el-icon>
                </el-menu-item>
              </template>
              <el-menu :default-active="navigationStore.activeIndex" class="el-menu-vertical-demo dark:text-slate-300 overflow-hidden"
                :unique-opened="true" style="border: none;" active-text-color="#fff">
                <template v-for="child in item.children" :key="child.name">
                  <el-menu-item :index="child.name"
                    @click="onMenuItemClick(child.title, child.name, child.parent_name)">
                    <el-icon>
                      <component :is="child.icon" />
                    </el-icon>
                    <span>{{ child.title }}</span>
                  </el-menu-item>
                </template>
              </el-menu>
            </el-popover>



            <el-sub-menu v-else :index="item.name" :class="navigationStore.parentActiveIndex == item.name ? 'parentmenu' : ''">
              <template #title>
                <el-icon>
                  <component :is="item.icon" />
                </el-icon>
                <span>{{ item.title }}</span>
              </template>
              <template v-for="child in item.children" :key="child.name">
                <el-menu-item :index="child.name" @click="onMenuItemClick(child.title, child.name, child.parent_name)">
                  <el-icon>
                    <component :is="child.icon" />
                  </el-icon>
                  <span>{{ child.title }}</span>
                </el-menu-item>
              </template>
            </el-sub-menu>
          </template>
        </el-menu>







        <span class="toggleside" @click="toggleAside">

          <el-icon v-if="!isCollapsed">
            <DArrowLeft />
          </el-icon>
          <el-icon v-else>
            <DArrowRight />
          </el-icon>
        </span>
      </el-aside>
      <el-main style="padding: 0;">
        <Tabs :tabs="navigationStore.tabs" :active-tab="navigationStore.activeIndex" @close-tab="HandleCloseTab" @click-tab="HandleClickTab" />
        <div class="router-view-container">
          <router-view></router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref,  onMounted, onBeforeUnmount, computed,watch,watchEffect } from 'vue';
import { useRouter } from 'vue-router';
import Tabs from '../components/Tabs.vue';
import AvatarJpg from '@/assets/avatar.png';
import LogoJpg from '@/assets/logo.png';
import NotificationCenter from '../components/NotificationCenter.vue'
import { useUserPermissionsStore } from '../stores/userPermissions';

import { useNavigationStore } from '@/stores/navigation';

// 使用 navigation store
const navigationStore = useNavigationStore();




    const userPermissions = useUserPermissionsStore();

    const menuData = computed(() => userPermissions.getRoleMenu);


    const router = useRouter();
    const isCollapsed = ref(false);


    const goToMonitorPage = () => {
  
      navigationStore.navigateToPage({
    title: '系统监控',
    componentName: 'SystemMonitor',
    parentName: 'System'
  });
  
};
    //console.log('tabs', tabs.value);
    const logout = () => {
    
      userPermissions.clearAuth();
      router.push('/');

    };



    /**
     * 处理菜单点击事件，用于切换标签页和路由。
     * @param {string} tabName - 要激活的标签页名称。
     * @param {string} componentName - 对应的组件名称，用于路由跳转。
     * 如果标签页不存在，则将其添加到标签页数组中，并激活该标签页及对应的路由。
     */
    const onMenuItemClick = (tabName, componentName, parentName) => {

      navigationStore.navigateToPage({
    title: tabName,
    componentName: componentName,
    parentName: parentName
  });
    };
    const HandleClickTab = (tabName) => {
      navigationStore.clickTab(tabName);

    };
    const HandleCloseTab = (tabName, selectIndex) => {
       
      navigationStore.closeTab(tabName, selectIndex);
        
    };
    const toggleAside = () => {
      isCollapsed.value = !isCollapsed.value;
    }
    const handleResize = () => {
      const width = window.innerWidth;
      if (width > 768) { // 假设768px是触发折叠的阈值
        if (isCollapsed.value) {
          toggleAside();
        }
      } else {
        if (!isCollapsed.value) {
          toggleAside();
        }
      }
    };
    onMounted(() => {
      window.addEventListener('resize', handleResize);
      // 初始化时也检查一次
      handleResize();






    });

    onBeforeUnmount(() => {
      window.removeEventListener('resize', handleResize);
    });


</script>

<style>
#app,
html,
body {
  height: 100%;
  margin: 0;
  padding: 0;
  background-color: #f5f5f5;
}

.el-container {
  height: 100%;
}

.el-header {
  background-color: #fff;
  color: #333;
  line-height: 50px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

}

.el-aside {
  background-color: #fff;
  color: #333;
  text-align: center;
  line-height: 200px;
  padding: 5px 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-top: 2px;
  position: relative;
  overflow: hidden;
  transition: width 0.3s;
}

.el-menu-vertical-demo {
  min-height: 400px;
}

.aside250 {

  width: 250px;
}

.aside64 {
  width: 64px;
  padding: 0;
}

.toggleside {
  display: block;

  position: absolute;
  bottom: 20px;
  right: 20px;
  width: 30px;
  height: 30px;
  line-height: 30px;
  background-color: #f5f5f5;
  cursor: pointer;
}

.el-menu-item {
  margin: 5px 0;
}

.el-menu-item.is-active {
  background-color: #337ecc;
  border-radius: 4px;
}

/* 添加以下样式以防止选中的菜单项在悬停时改变颜色 */

.singlemenu.is-active.el-tooltip__trigger,
.singlemenu.el-tooltip__trigger {
  background-color: #337ecc;
  color: #fff;
  border-radius: 0;
}

.el-popper .el-menu-item.is-active {
  border-radius: 0;
}

.parentmenu,
.parentmenu .el-sub-menu__title {
  color: #337ecc;
}

.is-active.el-tooltip__trigger {
  background: none;
  color: #000;
}




.popperclass.el-popover.el-popper {

  padding: 0;
}

.el-main {
  color: #333;
  text-align: center;
  margin: 5px 0.5rem 1px;
  height: calc(100vh - 60px);
  --el-main-padding: 0px;
  overflow-y: auto;
  /* 启用垂直滚动 */
}

.router-view-container {
  height: calc(100% - 110px);
  /* 减去标签栏的高度 */
  overflow-y: scroll;
  /* 启用垂直滚动 */
  background-color: #fff;
  padding: 1rem;
  border-radius: 4px;
}

/* 隐藏滚动条（WebKit 浏览器） */
.router-view-container::-webkit-scrollbar {
  display: none;
}

/* 隐藏滚动条（Firefox） */
.router-view-container {
  scrollbar-width: none;
  /* 对于 Firefox */
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  font-size: 20px;
  font-weight: bold;
  color: #333;
  /* 改变字体颜色 */
}

.logo img {
  width: 50px;
  height: 50px;
}

.logo span {
  height: 50px;
  line-height: 50px;
}

.user-info {
  display: flex;
  align-items: center;
  color: #333;
  cursor: pointer;
  /* 改变字体颜色 */
}

.el-dropdown-link {
  display: flex;
  align-items: center;
}

.el-dropdown-link .el-icon-arrow-down {
  margin-left: 5px;
  font-size: 12px;
}
</style>