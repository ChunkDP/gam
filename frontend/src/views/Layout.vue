<!-- Layout.vue -->
<template>
  <el-container class="layout">
    <!-- 顶部导航栏 -->
    <el-header height="60px" class="header mb-1" :class="{ 'header-light': theme === 'light' }">
      <div class="logo">
        <img src="@/assets/logo.png" alt="Logo" height="40" />
        <h1>Go-admin-manage</h1>
      </div>
      
      <!-- 顶部菜单模式 -->
      <el-menu
        v-if="menuMode === 'horizontal'"
        mode="horizontal"
        :default-active="activeMenu"
        :menu-trigger="theme=== 'light' ? 'click' : 'hover'"
        class="horizontal-menu"
        :background-color="theme === 'light' ? '#f5f5f5' : '#242424'"
        :text-color="theme === 'light' ? '#303133' : '#e5eaf3'"
        :active-text-color="primaryColor"
        @select="handleMenuSelect"
      >
        <template v-for="item in menuItems" :key="item.path">
          <el-menu-item v-if="item.children.length==0" :index="item.path">
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.title }}</span>
          </el-menu-item>
          <el-sub-menu v-else :index="item.path">
            <template #title>
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item 
              v-for="child in item.children" 
              :key="child.path" 
              :index="child.path"
            >
              <el-icon><component :is="child.icon" /></el-icon>
              <span>{{ child.title }}</span>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
      
      <div class="header-right">
        <!-- 设置按钮 -->
        <el-button 
          :type="theme === 'light' ? 'primary' : 'default'" 
          circle 
          @click="drawerVisible = true"
          class="setting-btn"
        >
          <el-icon><Setting /></el-icon>
        </el-button>
        <!-- 通知中心组件 -->
        <NotificationCenter />
        <el-dropdown>
          <span class="user-info">
            <el-avatar :size="32" :src="AvatarJpg"></el-avatar>
            
            <span>管理员</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人信息</el-dropdown-item>
              <el-dropdown-item>修改密码</el-dropdown-item>
              <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    
    <el-container>
      <!-- 侧边菜单 -->
      <el-aside 
        v-if="menuMode === 'vertical'" 
        :width="isCollapse ? '64px' : '220px'" 
        class="sidebar" 
        :class="{ 'sidebar-light': theme === 'light' }"
      >
        <el-menu
          :default-active="activeMenu"
         
          class="vertical-menu"
          :background-color="theme === 'light' ? '#ffffff' : '#252525'"
          :text-color="theme === 'light' ? '#303133' : '#e5eaf3'"
          :active-text-color="primaryColor"
          :collapse="isCollapse"
          @select="handleMenuSelect"
        >
          <template v-for="item in menuItems" :key="item.path">
            <el-menu-item v-if="item.children.length==0" :index="item.path">
              <el-icon><component :is="item.icon" /></el-icon>
              <template #title>{{ item.title }}</template>
            </el-menu-item>
            <el-sub-menu v-else :index="item.path">
              <template #title>
                <el-icon><component :is="item.icon" /></el-icon>
                <span>{{ item.title }}</span>
              </template>
              <el-menu-item 
                v-for="child in item.children" 
                :key="child.path" 
                :index="child.path"
              >
                <el-icon><component :is="child.icon" /></el-icon>
                <template #title>{{ child.title }}</template>
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
        <div class="collapse-btn" :class="{ 'collapse-btn-light': theme === 'light' }" @click="toggleCollapse">
          <el-icon v-if="isCollapse"><Expand /></el-icon>
          <el-icon v-else><Fold /></el-icon>
        </div>
      </el-aside>
      
      <!-- 主内容区域 -->
      <el-container class="main-container">
        <!-- 标签页导航 -->
        <div class="tabs-bar">
          <el-tabs 
            v-model="activeTab" 
            type="card" 
            closable 
            @tab-remove="removeTab"
            @tab-click="clickTab"
            @contextmenu.prevent="showContextMenu($event)"
          >
            <el-tab-pane
              v-for="tab in tabs"
              :key="tab.path"
              :label="tab.title"
              :name="tab.path"
              :closable="tab.path !== '/layout/home'"
            />
          </el-tabs>
        </div>
        
        <!-- 内容显示区域 -->
        <el-main class="content" :class="{ 'content-light': theme === 'light', 'content-dark': theme === 'dark' }">
          <router-view v-slot="{ Component }">
            <keep-alive :include="cachedViews">
              <component :is="Component" />
            </keep-alive>
          </router-view>
        </el-main>
      </el-container>
    </el-container>
    
    <!-- 右键菜单 -->
    <ul 
      class="contextmenu" 
      v-show="contextMenuVisible" 
      :style="{ left: contextMenuLeft + 'px', top: contextMenuTop + 'px' }"
      @click="hideContextMenu"
    >
      <li v-for="menu in contextMenus" :key="menu.text" @click="menu.handler">
        {{ menu.text }}
      </li>
    </ul>
    
    <!-- 设置抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="系统设置"
      direction="rtl"
      size="300px"
      :with-header="false"
      class="setting-drawer"
    >
      <div class="drawer-content">
        <div class="drawer-header">
          <h2 class="drawer-title">系统设置</h2>
          <el-button @click="drawerVisible = false" circle >
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
        
        <h3 class="drawer-section-title">界面显示</h3>
        
        <!-- 菜单模式设置 -->
        <div class="setting-item">
          <span class="setting-label">菜单模式</span>
          <el-radio-group v-model="menuMode" size="small">
            <el-radio-button value="vertical">侧边菜单</el-radio-button>
            <el-radio-button value="horizontal">顶部菜单</el-radio-button>
          </el-radio-group>
        </div>
        
        <!-- 主题设置 -->
        <div class="setting-item">
          <span class="setting-label">系统主题</span>
          <el-radio-group v-model="theme" size="small" @change="changeTheme">
            <el-radio-button value="light">浅色</el-radio-button>
            <el-radio-button value="dark">深色</el-radio-button>
          </el-radio-group>
        </div>
        
        <!-- 主题色设置 -->
        <div class="setting-item theme-color-setting">
          <span class="setting-label">主题色</span>
          <div class="theme-colors">
            <div 
              v-for="color in presetColors" 
              :key="color" 
              class="color-block" 
              :style="{ backgroundColor: color }"
              :class="{ active: primaryColor === color }"
              @click="selectThemeColor(color)"
            ></div>
            <el-color-picker v-model="primaryColor" @change="changePrimaryColor"    :predefine="presetColors"  class="custom-color-picker" />
          
          </div>
        </div>
      </div>
    </el-drawer>
  </el-container>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import '@/assets/style/dark-theme.css' // 引入深色主题样式
import '@/assets/style/light-theme.css' // 引入浅色主题样式
import AvatarJpg from '@/assets/avatar.png';

import NotificationCenter from '../components/NotificationCenter.vue';
import { useUserPermissionsStore } from '../stores/userPermissions';
const userPermissions = useUserPermissionsStore();

const route = useRoute()
const router = useRouter()

// 菜单模式：vertical(侧边栏) 或 horizontal(顶部)
const menuMode = ref('vertical')
const isCollapse = ref(false)

// 设置抽屉
const drawerVisible = ref(false)
const theme = ref('light')
const primaryColor = ref('#409EFF')
// 菜单数据
const menuItems = computed(() => userPermissions.getRoleMenu);

// 退出登录
const logout = () => {
  userPermissions.clearAuth();
  router.push('/login');
};


// 标签页相关
const activeTab = ref('/layout/home')
const tabs = ref([
  { title: '首页', path: '/layout/home',component: "Home"}
])
const cachedViews = ref(['Home'])

// 当前激活的菜单
const activeMenu = computed(() =>{
  console.log(route.path);
  return route.path
} )

// 监听路由变化，添加标签页
watch(
  () => route.path,
  (newPath) => {
    addTab(newPath)
  }
)

// 切换菜单折叠状态
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
  localStorage.setItem('isCollapse', isCollapse.value ? 'true' : 'false')
}

// 处理菜单选择
const handleMenuSelect = (path) => {
  router.push(path)
}

// 添加标签页
const addTab = (path) => {
  // 查找菜单项获取标题
  let title = '未知页面'
  let found = false
  let component = ''

  // 查找一级菜单
  for (const item of menuItems.value) {
   
    if (item.path === path) {
      title = item.title
      component = item.component
      found = true
      break
    }
    
    // 查找二级菜单
    if (item.children) {
      for (const child of item.children) {
        if (child.path === path) {
          title = child.title
          component = child.component
          found = true
          break
        }
      }
      if (found) break
    }
  }
  
  // 检查标签是否已存在
  const isExist = tabs.value.some(tab => tab.path === path)
  if (!isExist) {
    
    tabs.value.push({ title:title, path:path,component: component})
    
    // 添加到缓存视图
    const componentName = path.split('/').pop().charAt(0).toUpperCase() + 
                          path.split('/').pop().slice(1)
    if (!cachedViews.value.includes(componentName)) {
      cachedViews.value.push(componentName)
    }
     // 保存到localStorage
     saveTabs()
  }
  
  activeTab.value = path
  // 保存当前激活的标签
  localStorage.setItem('activeTab', path)
}

// 移除标签页
const removeTab = (targetPath) => {
  // 不允许关闭最后一个标签
  if (tabs.value.length === 1) {
    return
  }
  
  const tabIndex = tabs.value.findIndex(tab => tab.path === targetPath)
  
  // 从缓存视图中移除
  const componentName = targetPath.split('/').pop().charAt(0).toUpperCase() + 
                        targetPath.split('/').pop().slice(1)
  const cacheIndex = cachedViews.value.indexOf(componentName)
  if (cacheIndex > -1) {
    cachedViews.value.splice(cacheIndex, 1)
  }
  
  // 移除标签
  tabs.value.splice(tabIndex, 1)
  
  // 如果关闭的是当前激活的标签，则激活前一个标签
  if (activeTab.value === targetPath) {
    activeTab.value = tabs.value[tabIndex - 1 >= 0 ? tabIndex - 1 : 0].path
    router.push(activeTab.value)
  }
   // 保存到localStorage
   saveTabs()
  localStorage.setItem('activeTab', activeTab.value)
}
// 保存标签页到localStorage
const saveTabs = () => {
  localStorage.setItem('tabs', JSON.stringify(tabs.value))
  localStorage.setItem('cachedViews', JSON.stringify(cachedViews.value))
}
// 从localStorage恢复标签页
const restoreTabs = () => {
  const savedTabs = localStorage.getItem('tabs')
  const savedCachedViews = localStorage.getItem('cachedViews')
  const savedActiveTab = localStorage.getItem('activeTab')
  
  if (savedTabs) {
    try {
      const parsedTabs = JSON.parse(savedTabs)
      // 确保至少有首页标签
      if (parsedTabs.length > 0) {
        tabs.value = parsedTabs
      }
    } catch (e) {
      console.error('Failed to parse saved tabs:', e)
    }
  }
  
  if (savedCachedViews) {
    try {
      const parsedCachedViews = JSON.parse(savedCachedViews)
      if (parsedCachedViews.length > 0) {
        cachedViews.value = parsedCachedViews
      }
    } catch (e) {
      console.error('Failed to parse saved cached views:', e)
    }
  }
  
  if (savedActiveTab) {
    activeTab.value = savedActiveTab
    // 如果当前路由不是激活的标签，则导航到激活的标签
    if (route.path !== savedActiveTab) {
      router.push(savedActiveTab)
    }
  }
}
// 点击标签
const clickTab = (tab) => {
  router.push(tab.props.name)
}

// 右键菜单相关
const contextMenuVisible = ref(false)
const contextMenuLeft = ref(0)
const contextMenuTop = ref(0)
const selectedTabPath = ref(null)

// 右键菜单选项
const contextMenus = [
  { text: '关闭所有', handler: () => closeMultipleTabs('all') },
  { text: '关闭左侧', handler: () => closeMultipleTabs('left') },
  { text: '关闭右侧', handler: () => closeMultipleTabs('right') },
  { text: '关闭其他', handler: () => closeMultipleTabs('other') }
]

// 显示右键菜单
const showContextMenu = (event) => {
  event.preventDefault()
  
  // 获取目标元素ID
  let id = ''
  if (event.target.nodeName === 'SPAN') {
    id = event.target.offsetParent.id
  } else {
    id = event.target.id
  }
  
  // 如果找到ID，显示右键菜单
  if (id) {
    contextMenuLeft.value = event.clientX
    contextMenuTop.value = event.clientY + 10
    contextMenuVisible.value = true
    
    // 从ID中提取标签路径
    const match = id.match(/tab-(.+)/)
   
    if (match && match[1]) {
      selectedTabPath.value = match[1].replace(/-/g, '/')
      
    }
  }
}

// 隐藏右键菜单
const hideContextMenu = () => {
  contextMenuVisible.value = false
}

// 关闭多个标签的统一处理函数
const closeMultipleTabs = (type) => {
 
  const selectedIndex = tabs.value.findIndex(tab => tab.path === selectedTabPath.value)

 
  if (selectedIndex === -1) return
  
  const tabsToClose = tabs.value.reduce((acc, tab, index) => {
    // 不关闭首页标签
    if (tab.path === '/layout/home') return acc
   
    const shouldClose = {
      'all': true,
      'left': index < selectedIndex,
      'right': index > selectedIndex,
      'other': index !== selectedIndex
    }[type]
    
    return shouldClose ? [...acc, tab.path] : acc
  }, [])
  
  // 关闭所有需要关闭的标签
  tabsToClose.forEach(path => {
    removeTab(path)
  })
  
  // 如果关闭了当前激活的标签，需要激活一个新标签
  if (!tabs.value.some(tab => tab.path === activeTab.value)) {
    activeTab.value = '/layout/home'
    router.push('/layout/home')
  }
}

// 点击外部关闭右键菜单
const handleClickOutside = (event) => {
  if (!event.target.closest('.contextmenu')) {
    hideContextMenu()
  }
}

// 预设主题色
const presetColors = [
  '#409EFF', // 默认蓝色
  '#67C23A', // 绿色
  '#E6A23C', // 黄色
  '#F56C6C', // 红色
  '#909399', // 灰色
  '#9370DB', // 紫色
  '#13C2C2', // 青色
  '#1890FF'  // 亮蓝色
]

// 选择预设主题色
const selectThemeColor = (color) => {
  primaryColor.value = color
  changePrimaryColor(color)
}

// 修改主题色
const changePrimaryColor = (color) => {
 
  document.documentElement.style.setProperty('--el-color-primary', color)
 
  // 生成不同深浅的主题色
  const colorObj = generateThemeColors(color)
  for (const key in colorObj) {
    document.documentElement.style.setProperty(`--el-color-primary-${key}`, colorObj[key])
    
  }
   // 将十六进制颜色转换为RGB格式并设置--el-color-primary-rgb
   const rgbValues = hexToRgb(color)
  if (rgbValues) {
    document.documentElement.style.setProperty('--el-color-primary-rgb', `${rgbValues.r}, ${rgbValues.g}, ${rgbValues.b}`)
  }
  localStorage.setItem('primaryColor', color)
}

// 生成主题色的不同深浅变体
const generateThemeColors = (primary) => {
  const colors = {}
  
  // 转换为HSL
  const hsl = hexToHsl(primary)
  
  // 生成不同亮度的变体
  colors['light-3'] = hslToHex(hsl.h, hsl.s, Math.min(hsl.l + 30, 95))
  colors['light-5'] = hslToHex(hsl.h, hsl.s, Math.min(hsl.l + 50, 95))
  colors['light-7'] = hslToHex(hsl.h, hsl.s, Math.min(hsl.l + 70, 95))
  colors['light-8'] = hslToHex(hsl.h, hsl.s, Math.min(hsl.l + 80, 95))
  colors['light-9'] = hslToHex(hsl.h, hsl.s, Math.min(hsl.l + 90, 95))
  colors['dark-2'] = hslToHex(hsl.h, hsl.s, Math.max(hsl.l - 20, 10))
  
  return colors
}
// 添加十六进制颜色转RGB的函数
const hexToRgb = (hex) => {
  // 移除#号
  hex = hex.replace(/^#/, '')
  
  // 处理简写形式 (#RGB)
  if (hex.length === 3) {
    hex = hex.split('').map(char => char + char).join('')
  }
  
  const bigint = parseInt(hex, 16)
  const r = (bigint >> 16) & 255
  const g = (bigint >> 8) & 255
  const b = bigint & 255
  
  return { r, g, b }
}
// HEX转HSL
const hexToHsl = (hex) => {
  // 移除#号
  hex = hex.replace(/^#/, '')
  
  // 解析RGB值
  let r = parseInt(hex.substring(0, 2), 16) / 255
  let g = parseInt(hex.substring(2, 4), 16) / 255
  let b = parseInt(hex.substring(4, 6), 16) / 255
  
  // 找出最大和最小RGB值
  let max = Math.max(r, g, b)
  let min = Math.min(r, g, b)
  let h, s, l = (max + min) / 2
  
  if (max === min) {
    // 灰色
    h = s = 0
  } else {
    let d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break
      case g: h = (b - r) / d + 2; break
      case b: h = (r - g) / d + 4; break
    }
    
    h /= 6
  }
  
  return { h: h * 360, s: s * 100, l: l * 100 }
}

// HSL转HEX
const hslToHex = (h, s, l) => {
  h /= 360
  s /= 100
  l /= 100
  
  let r, g, b
  
  if (s === 0) {
    r = g = b = l
  } else {
    const hue2rgb = (p, q, t) => {
      if (t < 0) t += 1
      if (t > 1) t -= 1
      if (t < 1/6) return p + (q - p) * 6 * t
      if (t < 1/2) return q
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6
      return p
    }
    
    const q = l < 0.5 ? l * (1 + s) : l + s - l * s
    const p = 2 * l - q
    
    r = hue2rgb(p, q, h + 1/3)
    g = hue2rgb(p, q, h)
    b = hue2rgb(p, q, h - 1/3)
  }
  
  const toHex = (x) => {
    const hex = Math.round(x * 255).toString(16)
    return hex.length === 1 ? '0' + hex : hex
  }
  
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

// 切换主题
const changeTheme = (val) => {
  document.documentElement.classList.remove('light', 'dark')
  document.documentElement.classList.add(val)
  localStorage.setItem('theme', val)
}

// 初始化设置
const initSettings = () => {
  // 恢复菜单模式
  const savedMenuMode = localStorage.getItem('menuMode')
  if (savedMenuMode) {
    menuMode.value = savedMenuMode
  }
  
  // 恢复侧边栏折叠状态
  const savedIsCollapse = localStorage.getItem('isCollapse')
  if (savedIsCollapse !== null) {
   
    isCollapse.value = savedIsCollapse === 'true'
    
  }
  
  // 恢复主题
  const savedTheme = localStorage.getItem('theme') || 'light'
  theme.value = savedTheme
  changeTheme(savedTheme)
  
  // 恢复主题色
  const savedPrimaryColor = localStorage.getItem('primaryColor')
  if (savedPrimaryColor) {
    primaryColor.value = savedPrimaryColor
    changePrimaryColor(savedPrimaryColor)
  }

    // 恢复标签页
    restoreTabs()
}

// 监听菜单模式变化
watch(menuMode, (val) => {
  localStorage.setItem('menuMode', val)
})

onMounted(() => {
  // 初始化设置
  initSettings()
  
  // 添加点击外部关闭右键菜单的事件监听
  document.addEventListener('click', handleClickOutside)



})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)

})
</script>

<style scoped>
.layout {
  height: 100vh;
  width: 100%;
}

.header {
  background-color: #304156;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-light {
  background-color: #fff;
  color: #303133;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}


.logo h1 {
  margin-left: 10px;
  font-size: 18px;
  font-weight: 600;
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
  width: 80px;
  height: 80px;
  margin-right: 10px;
}
.header-right {
  display: flex;
  align-items: center;
}

.setting-btn {
  margin-right: 15px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}
.user-info:focus-visible{
  outline: none;
}
.user-info span {
  margin-left: 8px;
  color: inherit;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  overflow-x: hidden;
}

.sidebar-light {
  background-color: #fff;
  border-right: 1px solid #e6e6e6;
}

.vertical-menu {
  height: calc(100% - 50px);
  border-right: none;
}

.horizontal-menu {
  flex: 1;
  margin-left: 30px;
}

.collapse-btn {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #bfcbd9;
}

.collapse-btn-light {
  color: #606266;
}

.main-container {
  flex-direction: column;
  height: 100%;
}

.tabs-bar {
  background: #fff;
  padding: 6px 6px 0;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.12);
}

.content {
  padding: 10px;
  height: calc(100vh - 110px);
  overflow-y: auto;
  /* 美化滚动条 */
  scrollbar-width: thin; /* Firefox */
  scrollbar-color: transparent transparent; /* Firefox */
}

/* Webkit浏览器的滚动条样式 (Chrome, Safari, Edge等) */
.content::-webkit-scrollbar {
  width: 6px; /* 滚动条宽度 */
}

.content::-webkit-scrollbar-track {
  background: transparent; /* 滚动条轨道背景 */
}

.content::-webkit-scrollbar-thumb {
  background-color: rgba(144, 147, 153, 0.3); /* 滚动条颜色 */
  border-radius: 3px; /* 滚动条圆角 */
  border: none;
}

.content::-webkit-scrollbar-thumb:hover {
  background-color: rgba(144, 147, 153, 0.5); /* 悬停时的滚动条颜色 */
}

.content-light {
  background-color: #f5f5f5;
}

.content-dark {
  background-color: #1e2a39;
}

:deep(.el-tabs__header) {
  margin-bottom: 0;
}

:deep(.el-tabs__item) {
  height: 32px;
  line-height: 32px;
}

/* 右键菜单样式 */
.contextmenu {
  position: fixed;
  z-index: 50;
  margin: 0;
  width: 100px;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
  background-color: #fff;
  padding: 0;
  font-size: 14px;
  line-height: 30px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
}

.contextmenu li {
  cursor: pointer;
  list-style-type: none;
  padding: 3px 8px;
  font-size: 14px;
  line-height: 30px;
  color: #334155;
}

.contextmenu li:hover {
  background-color: #f3f4f6;
}

/* 深色模式支持 */
:root.dark .contextmenu {
  border-color: #1e293b;
  background-color: #0f172a;
}

:root.dark .contextmenu li {
  color: #e2e8f0;
}

:root.dark .contextmenu li:hover {
  background-color: #4b5563;
}

/* 设置抽屉样式 */
.setting-drawer :deep(.el-drawer__body) {
  padding: 0;
  overflow-y: auto;
}

.drawer-content {
  padding: 0;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.drawer-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
}

.drawer-section-title {
  margin: 20px 0 15px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  padding: 0 20px 10px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.setting-item {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.setting-label {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

/* 主题色设置 */
.theme-color-setting {
  flex-direction: column;
  align-items: flex-start;
}

.theme-colors {
  display: flex;
  flex-wrap: wrap;
  margin-top: 10px;
  width: 100%;
}

.color-block {
  width: 20px;
  height: 20px;
  margin-right: 8px;
  margin-bottom: 8px;
  border-radius: 2px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  border: 1px solid #e0e0e0;
}

.color-block.active {
  box-shadow: 0 0 0 2px #fff, 0 0 0 4px var(--el-color-primary);
}

.color-block:hover {
  transform: scale(1.1);
}

.custom-color-picker {
  margin-left: auto;
}
</style>
