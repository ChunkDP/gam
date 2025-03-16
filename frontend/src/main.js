import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import zhLocale from 'element-plus/es/locale/lang/zh-cn'


import './assets/style/main.css';


import router from './router'
import { createPinia } from 'pinia';
const app = createApp(App)
app.use(ElementPlus,  { locale: zhLocale })
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)

// 初始化主题
const initTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'light'
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
  }
}

// 在创建应用之前初始化主题
initTheme()

app.mount('#app')





