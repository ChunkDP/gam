<!-- Tabs.vue -->
<template>
  <div class="tags">
    <div class="el-tabs__nav-scroll">
      <el-tabs 
        v-model="activeTabValue" 
        type="card"
        class="top-tabs bg-white text-slate-700 dark:text-slate-500 dark:bg-slate-900" 
        @tab-click="handleTagClick"
        @tab-remove="handleClose" 
        @contextmenu.prevent="showContextMenu($event)"
      >
        <el-tab-pane 
          v-for="(tab, index) in tabsValue" 
          :key="tab.component" 
          :label="tab.name" 
          :name="tab.component"
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
    <li 
      v-for="menu in contextMenus" 
      :key="menu.text" 
      @click="menu.handler"
    >
      {{ menu.text }}
    </li>
  </ul>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, computed } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps({
  tabs: {
    type: Array,
    required: true
  },
  activeTab: {
    type: String,
    required: true
  }
});

// 定义事件
const emit = defineEmits([
  'close-tab',  // 关闭标签页
  'click-tab',  // 点击标签页
]);

// 使用计算属性来处理props，避免直接修改props
const activeTabValue = computed({
  get: () => props.activeTab,
  set: (val) => emit('click-tab', val)
});

const tabsValue = computed(() => props.tabs);

const router = useRouter();
const left = ref(0);
const top = ref(0);
const contextMenuVisible = ref(false);
const selectedTab = ref(null);

// 右键菜单选项
const contextMenus = [
  { text: '关闭所有', handler: () => closeMultipleTabs('all') },
  { text: '关闭左侧', handler: () => closeMultipleTabs('left') },
  { text: '关闭右侧', handler: () => closeMultipleTabs('right') },
  { text: '关闭其他', handler: () => closeMultipleTabs('other') }
];

const handleTagClick = (tabName) => {
  activeTabValue.value = tabName.props.name;
  router.push({ name: activeTabValue.value });
  emit('click-tab', activeTabValue.value);
};

const handleClose = (tabName) => {
  emit('close-tab', tabName);
};

const showContextMenu = (event) => {
  event.preventDefault();

  // 获取目标元素ID
  let id = '';
  if (event.target.nodeName === 'SPAN') {
    id = event.target.offsetParent.id;
  } else {
    id = event.target.id;
  }

  // 如果找到ID，显示右键菜单
  if (id) {
    left.value = event.clientX;
    top.value = event.clientY + 10;
    contextMenuVisible.value = true;
    // 从ID中提取标签名
    selectedTab.value = id.substring(id.indexOf('-') + 1);
  }
};

const hideContextMenu = () => {
  contextMenuVisible.value = false;
};

// 关闭多个标签的统一处理函数
const closeMultipleTabs = (type) => {
  let selectedIndex = props.tabs.findIndex(tab => tab.component === selectedTab.value);
  
  const tabsToClose = props.tabs.reduce((acc, tab, index) => {
    // 不关闭Home标签
    if (tab.component === 'Home') return acc;
    
    const shouldClose = {
      'all': true,
      'left': index < selectedIndex,
      'right': index > selectedIndex,
      'other': index !== selectedIndex
    }[type];

    return shouldClose ? [...acc, tab.component] : acc;
  }, []);

  // 设置关闭后选中的标签索引
  let targetIndex = selectedIndex;
  if (type !== 'right') {
    targetIndex = 1; // 默认选中第二个标签（通常是Home后的第一个）
  }
  if (type === 'all') {
    targetIndex = 0; // 如果关闭所有，选中Home标签
  }

  // 关闭所有需要关闭的标签
  tabsToClose.forEach(component => {
    emit('close-tab', component, targetIndex);
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
  box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000),var(--tw-ring-shadow, 0 0 #0000),var(--tw-shadow);
}

.contextmenu:is(.dark *) {
  --tw-border-opacity: 1;
  border-color: rgb(30 41 59 / var(--tw-border-opacity));
  --tw-bg-opacity: 1;
  background-color: rgb(15 23 42 / var(--tw-bg-opacity));
}

.contextmenu li {
  cursor: pointer;
  list-style-type: none;
  padding: 3px 8px;
  font-size: 14px;
  line-height: 30px;
  --tw-text-opacity: 1;
  color: rgb(51 65 85 / var(--tw-text-opacity));
}

.contextmenu li:hover {
  --tw-bg-opacity: 1;
  background-color: rgb(243 244 246 / var(--tw-bg-opacity));
}

.contextmenu li:is(.dark *) {
  --tw-text-opacity: 1;
  color: rgb(226 232 240 / var(--tw-text-opacity));
}

.contextmenu li:hover:is(.dark *) {
  --tw-bg-opacity: 1;
  background-color: rgb(75 85 99 / var(--tw-bg-opacity));
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
  :deep(.el-tabs__nav-next), :deep(.el-tabs__nav-prev) {
    color: var(--el-text-color-secondary);
    cursor: pointer;
    font-size: 12px;
    line-height: 44px;
    position: absolute;
    text-align: center;
    width: 20px;
  }

  :deep(.el-tabs--card > .el-tabs__header) {
    border: none;
  }

  :deep(.el-tabs__nav-scroll) {
    padding: 4px 4px;
  }

  :deep(.el-tabs__nav) {
    border: 0;
  }

  :deep(.el-tabs__header) {
    border-bottom: 0;
    margin-bottom: 0;
  }

  :deep(.el-tabs__item) {
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

  :deep(.el-tabs__item):first-child {
    border: 1px solid var(--el-border-color-darker);

    &.is-active {
      border: 1px solid var(--el-color-primary);
    }
  }
}
</style>