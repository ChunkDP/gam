<template>
  <div class="image-uploader">
    <el-upload
      class="uploader"
      :show-file-list="false"
      :http-request="customUpload"
      :on-success="handleSuccess"
      :before-upload="beforeUpload"
    >
      <img v-if="modelValue" :src="modelValue" class="uploaded-image" />
      <el-icon v-else class="uploader-icon"><Plus /></el-icon>
    </el-upload>
  </div>
</template>

<script setup>
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { uploadFile } from '@/services/upload'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  maxSize: {
    type: Number,
    default: 2 // 默认2MB
  }
})

const emit = defineEmits(['update:modelValue'])

// 自定义上传方法
const customUpload = async ({ file }) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res = await uploadFile(formData)
    return res
  } catch (error) {
    ElMessage.error('上传失败')
    return false
  }
}

// 上传成功回调
const handleSuccess = (res) => {
  if (res && res.url) {
    emit('update:modelValue', res.url)
    ElMessage.success('上传成功')
  } else {
    ElMessage.error('上传失败')
  }
}

// 上传前检查
const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < props.maxSize

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB!`)
    return false
  }
  return true
}
</script>

<style scoped>
.image-uploader {
  display: inline-block;
}

.uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
}

.uploader:hover {
  border-color: var(--el-color-primary);
}

.uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.uploaded-image {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: contain;
}
</style> 