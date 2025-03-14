import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';

export const useNavigationStore = defineStore('navigation', () => {
  const router = useRouter();
  
  // 基础状态声明
  const tabs = ref(JSON.parse(localStorage.getItem('tabs') || '[]'));
  const activeIndex = ref(localStorage.getItem('activeIntx') || '');
  const parentActiveIndex = ref(localStorage.getItem('ParentactiveIndex') || '');

  // 初始化首页标签
  if (tabs.value.length === 0) {
    tabs.value.push({ 
      name: '首页', 
      component: 'Home',
      parent: '' 
    });
    localStorage.setItem('tabs', JSON.stringify(tabs.value));
  }

  // 设置父级菜单激活状态
  const setParentactiveIndex = (parentName) => {
    parentActiveIndex.value = parentName || '';
    localStorage.setItem('ParentactiveIndex', parentActiveIndex.value);
  };

  // 保存当前激活的标签页
  const saveActiveIndex = (componentName) => {
    activeIndex.value = componentName;
    localStorage.setItem('activeIntx', componentName);
  };

  /**
   * 导航到指定页面并处理标签页
   */
  const navigateToPage = ({ title, componentName, parentName }) => {
   
    // 先检查标签页是否存在
    const existingTabIndex = tabs.value.findIndex(tab => tab.component === componentName);
    
    // 如果标签页不存在，则添加
    if (existingTabIndex === -1) {
      const newTab = {
        name: title,
        component: componentName,
        parent: parentName || ''
      };
      
      // 添加新标签
      tabs.value.push(newTab);
      
     
      // 保存到 localStorage
      localStorage.setItem('tabs', JSON.stringify(tabs.value));
    }
    
    // 设置激活状态
    saveActiveIndex(componentName);
    setParentactiveIndex(parentName || '');
    
    // 路由跳转
   
    router.push({ name: componentName }).catch(err => {
      console.error('路由跳转失败:', err);
    });
  };

  // 点击标签页
  const clickTab = (componentName) => {
    const tab = tabs.value.find(tab => tab.component === componentName);
    if (tab) {
      setParentactiveIndex(tab.parent || '');
      saveActiveIndex(componentName);
      router.push({ name: componentName });
    }
  };

  // 关闭标签页
  const closeTab = (componentName, selectIndex) => {
    const index = tabs.value.findIndex(tab => tab.component === componentName);
    if (index === -1) return;

    // 移除标签页
    tabs.value = tabs.value.filter(tab => tab.component !== componentName);
    localStorage.setItem('tabs', JSON.stringify(tabs.value));

    if (componentName === activeIndex.value) {
      if (selectIndex === undefined) {
        selectIndex = Math.min(index, tabs.value.length - 1);
      }
      
      if (tabs.value.length > 0) {
        const newTab = tabs.value[selectIndex];
        saveActiveIndex(newTab.component);
        setParentactiveIndex(newTab.parent || '');
        router.push({ name: newTab.component });
      } else {
        // 如果没有标签页了，返回首页
        const homeTab = { 
          name: '首页', 
          component: 'Home',
          parent: '' 
        };
        tabs.value = [homeTab];
        localStorage.setItem('tabs', JSON.stringify(tabs.value));
        saveActiveIndex('Home');
        setParentactiveIndex('');
        router.push({ name: 'Home' });
      }
    }
  };

  return {
    tabs,
    activeIndex,
    parentActiveIndex,
    navigateToPage,
    closeTab,
    clickTab
  };
}); 