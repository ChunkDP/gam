<template>
  <div class="upload-container">
    <draggable v-model="fileList" itemKey="url" class="draggable-list" @end="handleDragEnd" >

      <template #item="{ element }">
        <div class="upload-item">
          <img class="file-image" :src="element.url" />
          <div class="file-actions">
            <el-button type="primary" link @click.stop="handlePreview(element)">
              <el-icon><ZoomIn /></el-icon>
            </el-button>
            <el-button type="danger" link @click.stop="handleRemove(element)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </template>
  
    </draggable>

    <el-upload
      v-model:file-list="fileList"
      class="uploader"
      :show-file-list="false"
      :http-request="customUpload"
      :on-success="handleSuccess"
      :on-remove="handleRemove"
      :before-upload="beforeUpload"
      multiple
      list-type="picture-card"
    >
      <el-icon><Plus /></el-icon>

    </el-upload>

    <!-- 图片预览 -->
    <el-dialog v-model="previewVisible" width="800px">
      <img :src="previewUrl" alt="Preview Image" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Plus, ZoomIn, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { batchUploadFiles } from '@/services/upload'
import draggable from 'vuedraggable';
const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  maxSize: {
    type: Number,
    default: 2 // 默认2MB
  },
  maxCount: {
    type: Number,
    default: 10 // 默认最多10张
  },
  acceptTypes: {
    type: Array,
    default: () => ['image/jpeg', 'image/png', 'image/gif']
  }
})

const emit = defineEmits(['update:modelValue'])

// 文件列表
const fileList = ref([])
console.log("fileList", typeof(fileList.value ))
// 预览相关
const previewVisible = ref(false)
const previewUrl = ref('')

// 监听modelValue变化
watch(() => props.modelValue, (newVal) => {
  console.log("newVal", newVal)
  fileList.value = newVal.map((url, index) => {


    return {
      name: `image-${index}`,
      url,
      status: 'success'
    }
  }

)

}, { immediate: true })
const handleDragEnd = () => {

  const urls = fileList.value
      .filter(item => item.status === 'success')
      .map(item => item.url)
    emit('update:modelValue', urls)
  //console.log('Product.images 已更新:')
}
// 自定义上传方法
const customUpload = async ({ file }) => {
  try {
    const formData = new FormData()
    formData.append('files[]', file)
    const res = await batchUploadFiles(formData)
    return res[0] // 返回第一个上传结果
  } catch (error) {
    ElMessage.error('上传失败')
    return false
  }
}

// 上传成功回调
const handleSuccess = (res, file) => {
  if (res && res.url) {
    file.url = res.url
    const urls = fileList.value
      .filter(item => item.status === 'success')
      .map(item => item.url)
    emit('update:modelValue', urls)
    ElMessage.success('上传成功')
  } else {
    ElMessage.error('上传失败')
  }
}

// 删除图片
const handleRemove = (file) => {
  const index = fileList.value.indexOf(file)
  if (index !== -1) {
    fileList.value.splice(index, 1)
    const urls = fileList.value
      .filter(item => item.status === 'success')
      .map(item => item.url)
    emit('update:modelValue', urls)
  }
}

// 预览图片
const handlePreview = (file) => {
  previewUrl.value = file.url
  previewVisible.value = true
}

// 上传前检查
const beforeUpload = (file) => {
  // 检查文件类型
  const isAcceptType = props.acceptTypes.includes(file.type)
  if (!isAcceptType) {
    ElMessage.error('只能上传图片文件!')
    return false
  }

  // 检查文件大小
  const isLtSize = file.size / 1024 / 1024 < props.maxSize
  if (!isLtSize) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB!`)
    return false
  }

  // 检查数量限制
  if (fileList.value.length >= props.maxCount) {
    ElMessage.error(`最多只能上传 ${props.maxCount} 张图片!`)
    return false
  }

  return true
}
</script>
<style scoped>
.upload-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.draggable-list {
  display: flex;
  gap: 8px;
}

.upload-item {
  position: relative;
  width: 148px;
  height: 148px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  cursor: move;
}

.file-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  color: #909399;
}

.file-actions {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  background-color: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s;
}

.upload-item:hover .file-actions {
  opacity: 1;
}

.file-actions .el-button {
  padding: 4px;
  color: #fff;
}

.upload-button {
  display: inline-block;
}

.upload-trigger {
  width: 148px;
  height: 148px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  background-color: #f5f7fa;
  cursor: pointer;
}

.upload-trigger:hover {
  border-color: #409eff;
}
</style>