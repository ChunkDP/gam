<template>
  <div class="config-container">
    <el-row :gutter="20">
      <!-- 左侧菜单 -->
      <el-col :span="6">
        <div class="config-sidebar">
          <el-menu
            :default-active="activeGroup"
            class="config-menu"
            @select="handleGroupSelect"
          >
            <el-menu-item
              v-for="group in configGroups"
              :key="group.id"
              :index="String(group.id)"
            >
              <el-icon v-if="group.icon"><component :is="group.icon" /></el-icon>
              <span>{{ group.config_name }}</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-col>

      <!-- 右侧表单 -->
      <el-col :span="18">
        <div class="config-content" v-loading="loading">
          <template v-if="currentGroup">
            <div class="group-header">
              <h2>{{ currentGroup.config_name }}</h2>
              <p class="description">{{ currentGroup.description }}</p>
            </div>

            <el-form
              ref="formRef"
              :model="formData"
              label-width="120px"
              class="config-form"
              :rules="rules"
            >
              <template v-for="item in configItems" :key="item.id">
                <el-form-item
                  v-if="checkVisible(item)"
                  :label="item.item_name"
                  :prop="item.item_key"
                  :required="item.required === 1"
                >
                  <!-- 文本输入 -->
                  <el-input
                    v-if="item.value_type === 'string'"
                    v-model="formData[item.item_key]"
                    :placeholder="'请输入' + item.item_name"
                  />
    <!-- 数字输入 -->
    <el-input-number
      v-if="item.value_type === 'number'"
      v-model="formData[item.item_key]"
      :min="0"
      :precision="0"
      :placeholder="'请输入' + item.item_name"
    />
                  <!-- 密码输入 -->
                  <el-input
                    v-else-if="item.value_type === 'password'"
                    v-model="formData[item.item_key]"
                    type="password"
                    show-password
                    :placeholder="'请输入' + item.item_name"
                  />

                  <!-- 开关 -->
                  <el-switch
                    v-else-if="item.value_type === 'switch'"
                    v-model="formData[item.item_key]"
                    :active-value="1"
                    :inactive-value="0"
                  />

                  <!-- 选择器 -->
                  <el-select
                    v-else-if="item.value_type === 'select'"
                    v-model="formData[item.item_key]"
                    :placeholder="'请选择' + item.item_name"
                    @change="handleSelectChange(item.item_key, $event)"
                  >
                    <el-option
                      v-for="opt in parseOptions(item.options)"
                      :key="opt.value"
                      :label="opt.label"
                      :value="opt.value"
                    />
                  </el-select>

                  <!-- 图片上传 -->
                  <ImageUploader
                    v-else-if="item.value_type === 'upload'"
                    v-model="formData[item.item_key]"
                  />

                  <!-- 富文本编辑器 -->
                  <div v-else-if="item.value_type === 'editor'" class="editor-wrapper">
                    <WangEditor
                      v-model="formData[item.item_key]"
                      :height="320"
                    />
                  </div>

                  <!-- 文本域 -->
                  <el-input
                    v-else-if="item.value_type === 'textarea'"
                    v-model="formData[item.item_key]"
                    type="textarea"
                    :rows="4"
                    :placeholder="'请输入' + item.item_name"
                  />

                  <!-- JSON编辑器 -->
                  <div v-else-if="item.value_type === 'json'" class="json-editor">
                    <el-table
                      :data="parseJsonValue(formData[item.item_key])"
                      border
                      style="width: 100%"
                    >
                      <el-table-column label="标签" width="320">
                        <template #default="{ row, $index }">
                          <el-input 
                            v-model="row.label" 
                            placeholder="请输入标签" 
                            @blur="handleJsonChange(item.item_key)"
                          />
                        </template>
                      </el-table-column>
                      <el-table-column label="值" width="120">
                        <template #default="{ row, $index }">
                          <el-input 
                            v-model="row.value" 
                            placeholder="请输入值" 
                            @blur="handleJsonChange(item.item_key)"
                          />
                        </template>
                      </el-table-column>
                      <el-table-column label="操作" width="120">
                        <template #default="{ row, $index }">
                          <el-button type="danger" link @click="removeJsonItem(item.item_key, $index)">删除</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                    <div class="json-toolbar">
                      <el-button type="primary" link @click="addJsonItem(item.item_key)">添加选项</el-button>
                    </div>
                  </div>

                  <div class="item-description" v-if="item.description">
                    {{ item.description }}
                  </div>
                </el-form-item>
              </template>

              <el-form-item>
                <el-button type="primary" @click="handleSave">保存配置</el-button>
              </el-form-item>
            </el-form>
          </template>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import WangEditor from '@/components/WangEditor/index.vue'
import ImageUploader from '@/components/ImageUploader/index.vue'
import {
  getConfigGroups,
  getConfigItems,
  batchUpdateConfigs
} from '@/services/system'

const configGroups = ref([])
const configItems = ref([])
const currentGroup = ref(null)
const activeGroup = ref('')
const formData = ref({})
const formRef = ref(null)
const loading = ref(false)

// 添加错误处理
const handleError = (error) => {
  console.error('Error:', error)
  if (error.response?.data?.message) {
    ElMessage.error(error.response.data.message)
  } else if (error.message) {
    ElMessage.error(error.message)
  } else {
    ElMessage.error('操作失败')
  }
}

// 加载配置组
const loadConfigGroups = async () => {
  try {
    loading.value = true
    const res = await getConfigGroups()
    
    configGroups.value = res.data.groups
   
    if (configGroups.value.length > 0) {
      handleGroupSelect(String(configGroups.value[0].id))
    }
  } catch (error) {
    handleError(error)
  } finally {
  
    loading.value = false
  }
}

// 格式化配置值
const formatConfigValue = (item) => {
  const value = formData.value[item.item_key] || item.item_value
  
  switch (item.value_type) {
    case 'select':
      try {
        const options = parseOptions(item.options)
        const validValues = options.map(opt => opt.value)
        return validValues.includes(value) ? value : options[0]?.value || ''
      } catch (e) {
        console.error('Invalid options format:', e)
        return ''
      }
    case 'switch':
      // 将布尔值转换为字符串 '0' 或 '1'
      
      return value === '1' || value === 1 || value === true || value === 'true' ? 1 : 0
    case 'number':
      const num = Number(value)
      return isNaN(num) ? 0 : num
    case 'password':
      return value || ''
    case 'json':
      try {
        if (typeof value === 'string') {
          return value.trim() ? JSON.parse(value) : []
        }
        return Array.isArray(value) ? value : []
      } catch (e) {
        console.error('Invalid JSON format:', e)
        return []
      }
    case 'editor':
      return value || ''
    case 'textarea':
      return value || ''
    case 'string':
    default:
      return value || ''
  }
}

// 加载配置项
const loadConfigItems = async (groupId) => {
 
  try {
    const res = await getConfigItems(groupId)
    configItems.value = res.data.items
    formData.value = {}
    
    // 先设置所有默认值
    res.data.items.forEach(item => {
      formData.value[item.item_key] = formatConfigValue(item)
    })
    
   
  } catch (error) {
    handleError(error)
  } finally {
   
    loading.value = false
  }
}

// 选择配置组
const handleGroupSelect = (groupId) => {
  activeGroup.value = groupId
  currentGroup.value = configGroups.value.find(g => String(g.id) === groupId)
  loadConfigItems(groupId)
}

// 处理选择器值变化
const handleSelectChange = (key, value) => {
  const clearKeys = {
    // 存储方式配置
    'upload_driver': {
      'aliyun': ['aliyun_oss_key', 'aliyun_oss_secret', 'aliyun_oss_bucket', 'aliyun_oss_endpoint'],
      'tencent': ['tencent_cos_key', 'tencent_cos_secret', 'tencent_cos_bucket', 'tencent_cos_region'],
      'qiniu': ['qiniu_access_key', 'qiniu_secret_key', 'qiniu_bucket', 'qiniu_domain']
    },
    // 支付方式配置
    'payment_driver': {
      'alipay': ['alipay_app_id', 'alipay_private_key', 'alipay_public_key'],
      'wxpay': ['wxpay_app_id', 'wxpay_mch_id', 'wxpay_key', 'wxpay_cert_path'],
      'paypal': ['paypal_client_id', 'paypal_secret', 'paypal_mode'],
      'unionpay': ['unionpay_mch_id', 'unionpay_key', 'unionpay_cert_path']
    },
    // 短信服务商配置
    'sms_driver': {
      'aliyun': ['aliyun_sms_key', 'aliyun_sms_secret', 'aliyun_sms_sign'],
      'tencent': ['tencent_sms_id', 'tencent_sms_key', 'tencent_sms_sign'],
      'huawei': ['huawei_sms_key', 'huawei_sms_secret', 'huawei_sms_sign', 'huawei_sms_channel'],
      'qiniu': ['qiniu_sms_key', 'qiniu_sms_secret', 'qiniu_sms_sign'],
      'yunpian': ['yunpian_api_key', 'yunpian_sms_sign']
    }
  }

  if (key in clearKeys) {
    // 清空所有相关配置
    Object.values(clearKeys[key]).flat().forEach(k => {
      if (k in formData.value) {
        formData.value[k] = ''
      }
    })
  }
}

// 优化检查显示条件的函数
const checkVisible = (item) => {
  if (!item.visible_condition) return true
  try {
    // 安全检查
    const condition = item.visible_condition
    if (typeof condition !== 'string') return true
    
    // 检查条件中是否只包含允许的操作符
    const allowedOperators = ['===', '!==', '&&', '||', '(', ')', '.']
    const isValid = allowedOperators.every(op => 
      condition.includes(op) ? condition.indexOf(op) !== -1 : true
    )
    
    if (!isValid) return false
    
    // 执行条件检查
    return new Function('formData', `return ${condition}`)(formData.value)
  } catch (e) {
    console.error('条件表达式错误:', e)
    return false
  }
}

// 添加表单验证规则
const rules = {
  system_name: [{ required: true, message: '请输入系统名称', trigger: 'blur' }],
  site_name: [{ required: true, message: '请输入站点名称', trigger: 'blur' }],
  upload_max_size: [
    { required: true, message: '请输入最大尺寸', trigger: 'blur' },
    { type: 'number', message: '请输入数字', trigger: 'blur' }
  ],
  upload_mime_types: [{ required: true, message: '请输入允许的文件类型', trigger: 'blur' }],
  upload_driver: [{ required: true, message: '请选择存储方式', trigger: 'change' }],
  aliyun_oss_key: [{ required: true, message: '请输入AccessKey', trigger: 'blur' }],
  aliyun_oss_secret: [{ required: true, message: '请输入AccessSecret', trigger: 'blur' }],
  aliyun_oss_bucket: [{ required: true, message: '请输入Bucket名称', trigger: 'blur' }],
  aliyun_oss_endpoint: [{ required: true, message: '请输入OSS域名', trigger: 'blur' }],
  qiniu_access_key: [{ required: true, message: '请输入AccessKey', trigger: 'blur' }],
  qiniu_secret_key: [{ required: true, message: '请输入SecretKey', trigger: 'blur' }],
  qiniu_bucket: [{ required: true, message: '请输入Bucket名称', trigger: 'blur' }],
  qiniu_domain: [{ required: true, message: '请输入访问域名', trigger: 'blur' }],
  payment_driver: [{ required: true, message: '请选择支付方式', trigger: 'change' }],
  alipay_app_id: [{ required: true, message: '请输入支付宝AppID', trigger: 'blur' }],
  wxpay_app_id: [{ required: true, message: '请输入微信AppID', trigger: 'blur' }],
  paypal_client_id: [{ required: true, message: '请输入PayPal Client ID', trigger: 'blur' }],
  unionpay_mch_id: [{ required: true, message: '请输入银联商户号', trigger: 'blur' }],
  sms_driver: [{ required: true, message: '请选择短信服务商', trigger: 'change' }],
  aliyun_sms_key: [{ required: true, message: '请输入阿里云AccessKey', trigger: 'blur' }],
  tencent_sms_id: [{ required: true, message: '请输入腾讯云AppID', trigger: 'blur' }],
  huawei_sms_key: [{ required: true, message: '请输入华为云AppKey', trigger: 'blur' }],
  qiniu_sms_key: [{ required: true, message: '请输入七牛云AccessKey', trigger: 'blur' }],
  yunpian_api_key: [{ required: true, message: '请输入云片网APIKey', trigger: 'blur' }],
  site_announcement: [
    { required: true, message: '请输入内容', trigger: 'blur' },
    { min: 1, max: 100000, message: '内容长度必须在1-100000之间', trigger: 'blur' }
  ]
}

// 修改保存方法，添加表单验证
const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    const processedData = {}
    for (const key in formData.value) {
      const item = configItems.value.find(i => i.item_key === key)
      if (item) {
        
        switch (item.value_type) {
          case 'number':
            processedData[key] = String(formData.value[key])
          
            break
          case 'switch':
           
            processedData[key] = String(formData.value[key])
            break
          case 'json':
            // 确保数据是数组格式
            const jsonValue = Array.isArray(formData.value[key]) ? formData.value[key] : []
            processedData[key] = JSON.stringify(jsonValue)
            break
          default:
            processedData[key] = formData.value[key]
        }
      }
    }


    loading.value = true
    await batchUpdateConfigs({
      group_id: currentGroup.value.id,
      configs: processedData
    })
    loading.value = false
    ElMessage.success('保存成功')
  } catch (error) {
    handleError(error)
  }
}

// 添加 options 解析函数
const parseOptions = (options) => {
  try {
    return typeof options === 'string' ? 
      JSON.parse(options || '[]') : 
      (Array.isArray(options) ? options : [])
  } catch (e) {
    console.error('Failed to parse options:', e)
    return []
  }
}

// JSON 值解析
const parseJsonValue = (value) => {
  try {
    if (Array.isArray(value)) {
      return value
    }
    if (typeof value === 'string' && value) {
      return JSON.parse(value)
    }
    return []
  } catch (e) {
    return []
  }
}

// 添加 JSON 选项
const addJsonItem = (key) => {
  const currentValue = parseJsonValue(formData.value[key])
  currentValue.push({ label: '', value: '' })
  formData.value[key] = currentValue
}

// 删除 JSON 选项
const removeJsonItem = (key, index) => {
  const currentValue = parseJsonValue(formData.value[key])
  currentValue.splice(index, 1)
  formData.value[key] = currentValue
}

// 处理 JSON 值变更
const handleJsonChange = (key) => {
  const currentValue = parseJsonValue(formData.value[key])
  formData.value[key] = currentValue
}

onMounted(() => {
  loadConfigGroups()
})
</script>

<style scoped>
.config-container {
  padding: 20px;
  height: 100%;
  background: #fff;
}

.config-sidebar {

  border-radius: 4px;
  height: 100%;
}

.config-menu {
  border-right: none;
}
.config-menu li.is-active{
  color: #fff;
}
.config-content {
  min-height: 500px;
  padding: 20px;
}

.group-header {
  margin-bottom: 30px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.group-header h2 {
  margin: 0;
  font-size: 18px;
}

.group-header .description {
  margin: 10px 0 0;
  color: #666;
}

.config-form {
  max-width: 800px;
}

.item-description {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}

.config-form .el-form-item {
  margin-bottom: 22px;
  transition: all 0.3s ease;
}

/* 添加配置项切换动画 */
.el-form-item-enter-active,
.el-form-item-leave-active {
  transition: all 0.3s ease;
}

.el-form-item-enter-from,
.el-form-item-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 编辑器样式 */
.el-textarea {
  font-family: Consolas, Monaco, 'Courier New', monospace;
}

.el-textarea__inner {
  padding: 10px;
  line-height: 1.6;
}

.editor-wrapper {
  margin-bottom: 20px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.el-textarea.is-textarea {
  font-family: inherit;
}

.el-textarea.is-textarea .el-textarea__inner {
  padding: 8px 12px;
  line-height: 1.5;
  min-height: 100px;
  resize: vertical;
}

.json-editor {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  margin-bottom: 10px;
  display: flex;
}

.json-toolbar {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
}
</style> 