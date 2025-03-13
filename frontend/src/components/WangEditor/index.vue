<template>
  <div class="editor-container">
    <Toolbar
      :editor="editorRef"
      :defaultConfig="toolbarConfig"
      :mode="mode"
      style="border-bottom: 1px solid #ccc"
    />
    <Editor
      v-model="valueHtml"
      :defaultConfig="editorConfig"
      :mode="mode"
      @onCreated="handleCreated"
      style="height: 300px"
    />
  </div>
</template>

<script setup>
import { ref, shallowRef, onBeforeUnmount, watch, onMounted } from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { ElMessage } from 'element-plus'
import {uploadFile,getUploadConfig}  from '@/services/upload'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: [String, Number],
    default: '300px'
  }
})

const emit = defineEmits(['update:modelValue'])

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()

// 内容 HTML
const valueHtml = ref('')

// 模式
const mode = ref('default')

// 工具栏配置
const toolbarConfig = {
  excludeKeys: []
}

// 编辑器配置
const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      // 单个文件的最大体积限制，默认为 2M
      maxFileSize: 10 * 1024 * 1024,
      // 最多可上传几个文件，默认为 100
      maxNumberOfFiles: 10,
      // 选择文件时的类型限制，默认为 ['image/*']
      allowedFileTypes: ['image/*'],
      // 自定义上传
      customUpload: async (file, insertFn) => {
        try {
          const formData = new FormData()
          formData.append('file', file)
          
          const data = await uploadFile(formData)
          if (data) {
            insertFn(data.url, data.name || file.name, data.url)
          } else {
            ElMessage.error('上传失败')
          }
        } catch (error) {
          console.error('Upload error:', error)
          ElMessage.error('上传失败：' + (error.message || '未知错误'))
        }
      },
      // 上传错误，或者触发 timeout 时会调用该方法
      onError(file, err, res) {
        console.error('Upload error:', file, err, res)
        ElMessage.error(`${file.name} 上传失败：${err.message || '未知错误'}`)
      },
      // 上传超时时间，默认为 10 秒
      timeout: 5 * 1000,
      // 单个文件上传失败
      onFailed(file, res) {
        console.error('Upload failed:', file, res)
        ElMessage.error(`${file.name} 上传失败`)
      },
      // 上传进度条
      onProgress(progress) {
        console.log('Upload progress:', progress)
      }
    }
  }
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

// 初始化创建编辑器
const handleCreated = (editor) => {
  editorRef.value = editor
  valueHtml.value = props.modelValue
}

// 监听内容变化
watch(valueHtml, (newValue) => {
  emit('update:modelValue', newValue)
})

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  if (newValue !== valueHtml.value) {
    valueHtml.value = newValue
  }
})

// 设置编辑器高度
watch(() => props.height, (newHeight) => {
  if (editorRef.value) {
    editorRef.value.setHeight(
      typeof newHeight === 'number' ? `${newHeight}px` : newHeight
    )
  }
}, { immediate: true })

// 在组件创建时获取上传配置
const initUploadConfig = async () => {
  try {
    const config = await getUploadConfig()
    console.log("config:",config)
    // 根据配置更新编辑器的上传设置
    editorConfig.MENU_CONF.uploadImage.maxFileSize = config.max_size * 1024 * 1024
   editorConfig.MENU_CONF.uploadImage.allowedFileTypes = config.allowed_types//.split(',').map(type => `image/${type}`)


   console.log("editorConfig.MENU_CONF.uploadImage:",editorConfig.MENU_CONF.uploadImage)
  } catch (error) {
    console.error('Failed to load upload config:', error)
  }
}

onMounted(() => {
  initUploadConfig()
})
</script>

<style>
.editor-container {
  border: 1px solid #ccc;
  z-index: 100;
}

.editor-container .w-e-toolbar {
  border-bottom: 1px solid #eee;
  background-color: #fafafa;
}

.editor-container .w-e-text-container {
  background-color: #fff;
}

/* 编辑器全屏样式 */
.editor-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10000;
  background: #fff;
}

.editor-container.fullscreen .w-e-text-container {
  height: calc(100vh - 40px) !important;
}
</style> 