<!-- Layout.vue -->
<template>
  <el-container class="layout">
    <!-- 顶部导航栏 -->
    <el-header class="el-header dark:bg-slate-900" style="height: 50px;">
      <div class="header">
        <div class="logo">
          <el-image :src="LogoJpg" />
          <span @click="goToHomePage">Go-admin-manage</span>
        </div>
        <div class="user-info">
          <!-- 布局切换按钮 -->
          <el-tooltip content="切换布局模式" placement="bottom">
            <el-button 
              @click="toggleLayoutMode" 
              :icon="layoutMode === 'sidebar' ? 'Menu' : 'Operation'" 
              circle 
              class="mr-3"
            />
          </el-tooltip>
          
          <!-- 通知中心组件 -->
          <NotificationCenter />
          
          <!-- 用户下拉菜单 -->
          <el-dropdown trigger="click" popper-class="popclass">
            <span class="cursor-pointer flex justify-center items-center text-black dark:text-gray-100 el-dropdown-link">
              <el-avatar :size="40" :src="AvatarJpg"></el-avatar>
              <span>{{ userInfo.username || '用户名' }} <el-icon><arrow-down /></el-icon></span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goToUserProfile">个人中心</el-dropdown-item>
                <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>

    <el-container>
      <!-- 顶部导航模式 -->
      <template v-if="layoutMode === 'top'">
        <el-menu
          mode="horizontal"
          :default-active="navigationStore.activeIndex"
          class="el-menu-horizontal-demo w-full"
          @select="onMenuItemClick"
        >
          <template v-for="item in menuData" :key="item.name">
            <!-- 没有子菜单的项目 -->
            <el-menu-item v-if="!(item.children?.length)" :index="item.name">
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </el-menu-item>
            
            <!-- 有子菜单的项目 -->
            <el-sub-menu v-else :index="item.name">
              <template #title>
                <el-icon><component :is="item.icon" /></el-icon>
                <span>{{ item.title }}</span>
              </template>
              <template v-for="child in item.children" :key="child.name">
                <el-menu-item :index="child.name" :data-parent="item.name">
                  <el-icon><component :is="child.icon" /></el-icon>
                  <span>{{ child.title }}</span>
                </el-menu-item>
              </template>
            </el-sub-menu>
          </template>
        </el-menu>
      </template>
      
      <!-- 侧边栏模式 -->
      <el-aside v-if="layoutMode === 'sidebar'" :class="isCollapsed ? 'aside64' : 'aside250'">
        <el-menu 
          :default-active="navigationStore.activeIndex" 
          class="el-menu-vertical-demo dark:text-slate-300 overflow-hidden"
          :unique-opened="true" 
          style="border: none;" 
          active-text-color="#fff" 
          :collapse="isCollapsed"
          @select="onMenuItemClick"
        >
          <template v-for="item in menuData" :key="item.name">
            <!-- 折叠状态下的单个菜单项 -->
            <el-tooltip v-if="isCollapsed && !(item.children?.length)" effect="dark" :content="item.title" placement="right">
              <el-menu-item 
                :index="item.name"
                :class="navigationStore.activeIndex === item.name ? 'singlemenu' : ''"
              >
                <el-icon><component :is="item.icon" /></el-icon>
              </el-menu-item>
            </el-tooltip>

            <!-- 展开状态下的单个菜单项 -->
            <el-menu-item 
              v-else-if="!isCollapsed && !(item.children?.length)" 
              :index="item.name"
            >
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </el-menu-item>

            <!-- 折叠状态下的子菜单 -->
            <el-popover 
              v-else-if="isCollapsed && item.children?.length" 
              placement="right-end" 
              :width="200"
              trigger="click" 
              :offset="5" 
              transition="el-slide-in-right" 
              :show-arrow="false" 
              :hide-after="0"
              popper-class="popperclass"
            >
              <template #reference>
                <el-menu-item 
                  :index="item.name" 
                  :class="navigationStore.parentActiveIndex === item.name ? 'parentmenu' : ''"
                >
                  <el-icon><component :is="item.icon" /></el-icon>
                </el-menu-item>
              </template>
              <el-menu 
                :default-active="navigationStore.activeIndex" 
                class="el-menu-vertical-demo dark:text-slate-300 overflow-hidden"
                :unique-opened="true" 
                style="border: none;" 
                active-text-color="#fff"
              >
                <template v-for="child in item.children" :key="child.name">
                  <el-menu-item 
                    :index="child.name" 
                    @click="onMenuItemClick(child.name, item.name)"
                  >
                    <el-icon><component :is="child.icon" /></el-icon>
                    <span>{{ child.title }}</span>
                  </el-menu-item>
                </template>
              </el-menu>
            </el-popover>

            <!-- 展开状态下的子菜单 -->
            <el-sub-menu 
              v-else 
              :index="item.name" 
              :class="navigationStore.parentActiveIndex === item.name ? 'parentmenu' : ''"
            >
              <template #title>
                <el-icon><component :is="item.icon" /></el-icon>
                <span>{{ item.title }}</span>
              </template>
              <template v-for="child in item.children" :key="child.name">
                <el-menu-item 
                  :index="child.name" 
                  :data-parent="item.name"
                >
                  <el-icon><component :is="child.icon" /></el-icon>
                  <span>{{ child.title }}</span>
                </el-menu-item>
              </template>
            </el-sub-menu>
          </template>
        </el-menu>

        <!-- 侧边栏折叠按钮 -->
        <span class="toggleside" @click="toggleAside">
          <el-icon v-if="!isCollapsed"><DArrowLeft /></el-icon>
          <el-icon v-else><DArrowRight /></el-icon>
        </span>
      </el-aside>

      <!-- 主内容区域 -->
      <el-main style="padding: 0;">
        <!-- 标签页导航 -->
        <Tabs 
          :tabs="navigationStore.tabs" 
          :active-tab="navigationStore.activeIndex" 
          @close-tab="handleCloseTab" 
          @click-tab="handleClickTab" 
        />
        
        <!-- 路由视图容器 -->
        <div class="router-view-container">
          <router-view v-slot="{ Component }">
            <keep-alive>
              <component :is="Component" />
            </keep-alive>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import Tabs from '../components/Tabs.vue';
import AvatarJpg from '@/assets/avatar.png';
import LogoJpg from '@/assets/logo.png';
import NotificationCenter from '../components/NotificationCenter.vue';
import { useUserPermissionsStore } from '../stores/userPermissions';
import { useNavigationStore } from '@/stores/navigation';

// 使用 store
const navigationStore = useNavigationStore();
const userPermissions = useUserPermissionsStore();
const router = useRouter();

// 菜单数据
const menuData = computed(() => userPermissions.getRoleMenu);

// 用户信息
const userInfo = computed(() => userPermissions.getUserInfo || {});

// 侧边栏折叠状态
const isCollapsed = ref(false);

// 布局模式状态 (sidebar 或 top)
const layoutMode = ref(localStorage.getItem('layoutMode') || 'sidebar');

// 导航到首页
const goToHomePage = () => {
  navigationStore.navigateToPage({
    title: '首页',
    componentName: 'Home'
  });
};

// 导航到个人中心
const goToUserProfile = () => {
  navigationStore.navigateToPage({
    title: '个人中心',
    componentName: 'UserProfile'
  });
};

// 退出登录
const logout = () => {
  userPermissions.clearAuth();
  router.push('/login');
};

// 菜单项点击处理
const onMenuItemClick = (index, parentName) => {
  // 如果是从事件中获取，需要检查是否有父菜单
  if (!parentName && event && event.target) {
    const menuItem = event.target.closest('.el-menu-item');
    if (menuItem) {
      parentName = menuItem.dataset.parent;
    }
  }
  
  // 查找菜单项标题
  let title = '';
  let menuItem;
  
  if (parentName) {
    // 查找子菜单项
    const parentItem = menuData.value.find(item => item.name === parentName);
    if (parentItem && parentItem.children) {
      menuItem = parentItem.children.find(child => child.name === index);
      if (menuItem) {
        title = menuItem.title;
      }
    }
  } else {
    // 查找顶级菜单项
    menuItem = menuData.value.find(item => item.name === index);
    if (menuItem) {
      title = menuItem.title;
    }
  }
  
  // 导航到对应页面
  if (title) {
    navigationStore.navigateToPage({
      title: title,
      componentName: index,
      parentName: parentName
    });
  }
};

// 标签页点击处理
const handleClickTab = (tabName) => {
  navigationStore.clickTab(tabName);
};

// 关闭标签页处理
const handleCloseTab = (tabName, selectIndex) => {
  navigationStore.closeTab(tabName, selectIndex);
};

// 切换侧边栏折叠状态
const toggleAside = () => {
  isCollapsed.value = !isCollapsed.value;
};

// 切换布局模式
const toggleLayoutMode = () => {
  layoutMode.value = layoutMode.value === 'sidebar' ? 'top' : 'sidebar';
  // 保存用户偏好到本地存储
  localStorage.setItem('layoutMode', layoutMode.value);
};

// 响应式调整侧边栏
const handleResize = () => {
  const width = window.innerWidth;
  if (width <= 768) { // 小屏幕自动折叠
    if (!isCollapsed.value) {
      isCollapsed.value = true;
    }
  } else if (width > 1200) { // 大屏幕自动展开
    if (isCollapsed.value) {
      isCollapsed.value = false;
    }
  }
};

// 生命周期钩子
onMounted(() => {
  window.addEventListener('resize', handleResize);
  // 初始化时检查一次
  handleResize();
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

// 监听布局模式变化，更新相关样式
watch(layoutMode, (newMode) => {
  // 可以在这里添加额外的布局调整逻辑
  document.body.className = newMode === 'sidebar' ? 'layout-sidebar-mode' : 'layout-top-mode';
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

.layout {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.el-container {
  height: 100%;
}

.el-header {
  background-color: #fff;
  color: #333;
  line-height: 50px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 0 20px;
  z-index: 10;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.logo {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: bold;
  color: #333;
  cursor: pointer;
}

.logo img, .logo .el-image {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

/* 侧边栏样式 */
.el-aside {
  background-color: #fff;
  color: #333;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-top: 2px;
  position: relative;
  overflow: hidden;
  transition: width 0.3s;
  z-index: 9;
}

.aside250 {
  width: 250px;
  padding: 5px 0;
}

.aside64 {
  width: 64px;
  padding: 5px 0;
}

.toggleside {
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  bottom: 20px;
  right: 20px;
  width: 30px;
  height: 30px;
  background-color: #f5f5f5;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 菜单样式 */
.el-menu-vertical-demo {
  border-right: none !important;
  height: calc(100% - 50px);
}

.el-menu-horizontal-demo {
  border-bottom: solid 1px var(--el-menu-border-color);
}

.el-menu-item.is-active {
  background-color: #337ecc !important;
  color: #fff !important;
  border-radius: 4px;
}

.el-menu-horizontal-demo .el-menu-item.is-active {
  border-bottom: 2px solid #337ecc;
  background-color: transparent !important;
  color: #337ecc !important;
}

.singlemenu.is-active.el-tooltip__trigger,
.singlemenu.el-tooltip__trigger {
  background-color: #337ecc;
  color: #fff;
}

.parentmenu,
.parentmenu .el-sub-menu__title {
  color: #337ecc;
}

/* 主内容区域样式 */
.el-main {
  padding: 0;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 50px);
  overflow: hidden;
}

.router-view-container {
  flex: 1;
  overflow-y: auto;
  background-color: #fff;
  margin: 0 10px 10px;
  border-radius: 4px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

/* 隐藏滚动条 */
.router-view-container::-webkit-scrollbar {
  display: none;
}

.router-view-container {
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

/* 响应式调整 */
@media (max-width: 768px) {
  .el-menu-horizontal-demo {
    overflow-x: auto;
    white-space: nowrap;
  }
  
  .el-menu-horizontal-demo::-webkit-scrollbar {
    display: none;
  }
  
  .logo span {
    display: none;
  }
  
  .user-info .el-dropdown-link span {
    display: none;
  }
}

/* 深色模式适配 */
.dark .el-header,
.dark .el-aside,
.dark .router-view-container {
  background-color: #1e293b;
  color: #e2e8f0;
}

.dark .el-menu {
  background-color: #1e293b;
  color: #e2e8f0;
}

.dark .el-menu-item:not(.is-active) {
  color: #e2e8f0;
}

.dark .el-sub-menu__title {
  color: #e2e8f0;
}

.dark .toggleside {
  background-color: #334155;
  color: #e2e8f0;
}
</style>