/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      // 自定义颜色
      colors: {
              // 使用嵌套对象形式定义颜色
              'custom': {
                blue: '#1fb6ff',
                // 可以添加不同深浅度
                'blue-light': '#3cc5ff',
                'blue-dark': '#0da2eb',
              }
      },
      // 自定义字体大小
      fontSize: {
        'xxs': '.625rem',
        'mega': '4rem',
      },
      // 自定义间距
      spacing: {
        '128': '32rem',
        '144': '36rem',
      },
      // 自定义断点
      screens: {
        'tablet': '640px',
        'laptop': '1024px',
        'desktop': '1280px',
      },
    },
  },
  plugins: [],
}