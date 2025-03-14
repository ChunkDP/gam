<!-- Tabs.vue -->
<template>
  <div class="tags">


 

    <div class="el-tabs__nav-scroll">
      <el-tabs v-model="activeTab" type="card"
        class="top-tabs bg-white text-slate-700 dark:text-slate-500 dark:bg-slate-900" @tab-click="handleTagClick"
        @tab-remove="handleClose" @contextmenu.prevent="showContextMenu($event)">
        <el-tab-pane v-for="(tab, index) in tabs" :key="tab.component" :label="tab.name" :name="tab.component"
          :closable="tab.component !== 'Home'"
           />

      </el-tabs>
      
    </div>


  </div>


  <ul 
  class="contextmenu" 
  v-show="contextMenuVisible" 
  :style="{ left: left + 'px', top: top + 'px' }"
  @click="hideContextMenu"
>

<li v-for="menu in contextMenus" 
        :key="menu.text" 
        @click="menu.handler">
      {{ menu.text }}
    </li>
   
</ul>
</template>

<script setup>
import { defineComponent, ref, watch, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';

 const props= defineProps({
    tabs: {
      type: Array,
      required: true
    },
    activeTab: {
      type: String,
      required: true
    }
  
  }) 
// 定义事件
const emit = defineEmits([
  'close-tab',  // 关闭标签页
  'click-tab',  // 点击标签页
  
])
 
    const activeTab = ref(props.activeTab);
  
    const tabs = ref(props.tabs);
  
    const router = useRouter();
    const left = ref(0)
    const top = ref(0)
    const contextMenuVisible = ref(false);

const selectedTab = ref(null);
// 右键菜单选项
const contextMenus = [
  { text: '关闭所有', handler: () => closeMultipleTabs('all') },
  { text: '关闭左侧', handler: () => closeMultipleTabs('left') },
  { text: '关闭右侧', handler: () => closeMultipleTabs('right') },
  { text: '关闭其他', handler: () => closeMultipleTabs('other') }
];

// 监听 props 变化
watch(
      [() => props.activeTab,() => props.tabs],
      ([newActiveTab,newTabs]) => {
     
        activeTab.value = newActiveTab
        tabs.value =newTabs
      },
     
      { deep:true,immediate: true }
    );
    const handleTagClick = (tabName) => {
  
  activeTab.value = tabName.props.name;
 router.push({ name: activeTab.value });

  emit('click-tab',activeTab.value);
};

const handleClose = (tabName) => {

emit('close-tab', tabName);

};
const showContextMenu = (event) => {
  event.preventDefault();



  let id = ''
    if (event.srcElement.nodeName === 'SPAN') {
      id = event.srcElement.offsetParent.id
    } else {
      id = event.srcElement.id
    }
  

    if (id) {
     left.value = event.clientX
     top.value = event.clientY + 10
      contextMenuVisible.value = true;
      selectedTab.value = id.substring(4)
    }

};

const hideContextMenu = () => {
  contextMenuVisible.value = false;
};





    // 关闭多个标签的统一处理函数
const closeMultipleTabs = (type) => {
  let selectedIndex = props.tabs.findIndex(tab => tab.component === selectedTab.value);
  
  const tabsToClose = props.tabs.reduce((acc, tab, index) => {
    if (tab.component === 'Home') return acc;
    
    const shouldClose = {
      'all': true,
      'left': index < selectedIndex,
      'right': index > selectedIndex,
      'other': index !== selectedIndex
    }[type];

    return shouldClose ? [...acc, tab.component] : acc;
  }, []);

 
 if(type != 'right'){
  selectedIndex=1
 }
  if(type === 'all'){
    selectedIndex=0
  }


  tabsToClose.forEach(component => {
    
    emit('close-tab', component, selectedIndex);
  });
};


// 点击外部关闭右键菜单
const handleClickOutside = (event) => {
  if (!event.target.closest('.contextmenu')) {
    hideContextMenu();
  }
};

    onMounted(() => {
      document.addEventListener('click', handleClickOutside);
    });

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside);
    });

 

</script>

<style scoped>
.contextmenu {
    position: absolute;
    z-index: 50;
    margin: 0;
    width: 100px;
    border-radius: 4px;
    border-width: 1px;
    border-style: solid;
    --tw-border-opacity: 1;
    border-color: rgb(229 231 235 / var(--tw-border-opacity));
    --tw-bg-opacity: 1;
    background-color: rgb(255 255 255 / var(--tw-bg-opacity));

    padding-left: 0;
    padding-right: 0;
    font-size: 14px;
    line-height: 30px;
    --tw-shadow: 0 4px 6px -1px rgb(0 0 0 / .1), 0 2px 4px -2px rgb(0 0 0 / .1);
    --tw-shadow-colored: 0 4px 6px -1px var(--tw-shadow-color), 0 2px 4px -2px var(--tw-shadow-color);
    box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000),var(--tw-ring-shadow, 0 0 #0000),var(--tw-shadow)
}

.contextmenu:is(.dark *) {
    --tw-border-opacity: 1;
    border-color: rgb(30 41 59 / var(--tw-border-opacity));
    --tw-bg-opacity: 1;
    background-color: rgb(15 23 42 / var(--tw-bg-opacity))
}

.contextmenu li {
    cursor: pointer;
    list-style-type: none;
    padding: 3px 8px;
    font-size: 14px;
    line-height: 30px;
    --tw-text-opacity: 1;
    color: rgb(51 65 85 / var(--tw-text-opacity))
}

.contextmenu li:hover {
    --tw-bg-opacity: 1;
    background-color: rgb(243 244 246 / var(--tw-bg-opacity))
}

.contextmenu li:is(.dark *) {
    --tw-text-opacity: 1;
    color: rgb(226 232 240 / var(--tw-text-opacity))
}

.contextmenu li:hover:is(.dark *) {
    --tw-bg-opacity: 1;
    background-color: rgb(75 85 99 / var(--tw-bg-opacity))
}


.tags {
  margin-bottom: 10px;
  text-align: left;
  background-color: #fff;
  padding: 5px;
  width: 100%;
  position: relative;
  box-sizing: border-box;

  flex: 1 auto;
  display: flex;
  overflow: hidden;
}

.el-tabs__nav-next,
.el-tabs__nav-prev {
  color: var(--el-text-color-secondary);
  cursor: pointer;
  font-size: 12px;
  line-height: 44px;
  position: absolute;
  text-align: center;
  width: 20px;
}

.tags {

  ::v-deep(.el-tabs__nav-next, .el-tabs__nav-prev) {
    color: var(--el-text-color-secondary);
    cursor: pointer;
    font-size: 12px;
    line-height: 44px;
    position: absolute;
    text-align: center;
    width: 20px;
  }

  ::v-deep(.el-tabs--card > .el-tabs__header) {
    border: none;
  }

  ::v-deep(.el-tabs__nav-scroll) {
    padding: 4px 4px;
  }

  ::v-deep(.el-tabs__nav) {
    border: 0;
  }

  ::v-deep(.el-tabs__header) {
    border-bottom: 0;
    margin-bottom: 0;
  }

  ::v-deep(.el-tabs__item) {
    box-sizing: border-box;
    border: 1px solid var(--el-border-color-darker);
    border-radius: 2px;
    margin-right: 5px;
    margin-left: 2px;
    transition: padding 0.3s cubic-bezier(0.645, 0.045, 0.355, 1) !important;
    height: 34px;

    &.is-active {
      border: 1px solid var(--el-color-primary);
    }
  }

  ::v-deep(.el-tabs__item):first-child {
    border: 1px solid var(--el-border-color-darker);

    &.is-active {
      border: 1px solid var(--el-color-primary);
    }
  }
}
</style>