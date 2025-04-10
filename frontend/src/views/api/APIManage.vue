<template>
  <div class="api-management">
    <el-container>
      <!-- API列表区域，默认占满宽度 -->
      <el-main :class="{ 'with-test-panel': showTestPanel }">
        <div class="toolbar">
          <el-input
            v-model="searchQuery"
            placeholder="搜索API名称"
            style="width: 200px"
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增API
          </el-button>
        </div>

        <el-table :data="apiList" style="width: 100%">
          <el-table-column prop="name" label="API名称" />
          <el-table-column prop="path" label="路径" />
          <el-table-column prop="method" label="请求方法">
            <template #default="{ row }">
              <el-tag :type="getMethodTagType(row.method)">{{ row.method }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="group" label="分组" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                @change="handleStatusChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280">
            <template #default="{ row }">
              <el-button-group>
                <el-button type="primary" @click="handleEdit(row)">编辑</el-button>
                <el-button type="success" @click="handleTest(row)">测试</el-button>
                <el-button type="danger" @click="handleDelete(row)">删除</el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-main>

      <!-- 右侧测试面板，只在showTestPanel为true时显示 -->
      <el-aside v-if="showTestPanel" width="40%">
        <div class="test-panel-header">
          <span>API测试</span>
          <el-button circle  @click="closeTestPanel">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
        <div class="api-test-panel">
          <div v-if="selectedAPI" class="api-info">
            <h3>{{ selectedAPI.name }}</h3>
            <p>{{ selectedAPI.description }}</p>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="请求路径">{{ selectedAPI.path }}</el-descriptions-item>
              <el-descriptions-item label="请求方法">
                <el-tag :type="getMethodTagType(selectedAPI.method)">{{ selectedAPI.method }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="API分组">{{ selectedAPI.group }}</el-descriptions-item>
            </el-descriptions>

            <div class="test-params">
              <h4>请求参数</h4>
              <el-input
                v-model="testParams"
                type="textarea"
                :rows="6"
                placeholder="请输入JSON格式的请求参数"
              />
            </div>

            <el-button 
              type="primary" 
              :loading="testing" 
              @click="handleTestSubmit"
            >
              发送请求
            </el-button>

            <div v-if="testResult" class="test-result">
              <h4>响应结果</h4>
              <el-alert
                :type="testSuccess ? 'success' : 'error'"
                :title="testSuccess ? '请求成功' : '请求失败'"
                :closable="false"
                show-icon
              />
              <pre class="response-data">{{ formatJSON(testResult) }}</pre>
            </div>
          </div>
        </div>
      </el-aside>
    </el-container>

    <!-- API表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑API' : '新增API'"
      width="50%"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="API名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="请求路径" prop="path">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item label="请求方法" prop="method">
          <el-select v-model="form.method">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="API分组" prop="group">
          <el-input v-model="form.group" />
        </el-form-item>
        <el-form-item label="API描述" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="请求参数" prop="parameters">
          <el-input v-model="form.parameters" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="响应示例" prop="response">
          <el-input v-model="form.response" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import { getAPIList, getAPI, createAPI, updateAPI, deleteAPI, testAPI } from '@/services/api'

// 状态变量
const apiList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const selectedAPI = ref(null)
const testing = ref(false)
const testParams = ref('')
const testResult = ref(null)
const testSuccess = ref(false)
const showTestPanel = ref(false)

// 表单相关
const formRef = ref(null)
const form = reactive({
  name: '',
  path: '',
  method: 'GET',
  group: '',
  description: '',
  parameters: '',
  response: '',
  status: 1
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入API名称', trigger: 'blur' }],
  path: [{ required: true, message: '请输入请求路径', trigger: 'blur' }],
  method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
  group: [{ required: true, message: '请输入API分组', trigger: 'blur' }]
}

// 加载API列表
const loadAPIList = async () => {
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      name: searchQuery.value
    }
    const { list, total: totalCount } = await getAPIList(params)
    apiList.value = list
    total.value = totalCount
  } catch (error) {
    ElMessage.error('获取API列表失败')
  }
}

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1
  loadAPIList()
}

// 处理分页
const handleSizeChange = (val) => {
  pageSize.value = val
  loadAPIList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  loadAPIList()
}

// 处理新增
const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, {
    name: '',
    path: '',
    method: 'GET',
    group: '',
    description: '',
    parameters: '',
    response: '',
    status: 1
  })
  dialogVisible.value = true
}

// 处理编辑
const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

// 处理删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该API吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteAPI(row.id)
      ElMessage.success('删除成功')
      loadAPIList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

// 处理状态变更
const handleStatusChange = async (row) => {
  try {
    await updateAPI(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1 // 恢复状态
    ElMessage.error('状态更新失败')
  }
}

// 处理表单提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          await updateAPI(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createAPI(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadAPIList()
      } catch (error) {
        ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
      }
    }
  })
}

// 处理API测试
const handleTest = (row) => {
  selectedAPI.value = row
  testParams.value = row.parameters || '{}'
  testResult.value = null
  showTestPanel.value = true
}

// 关闭测试面板
const closeTestPanel = () => {
  showTestPanel.value = false
  selectedAPI.value = null
  testResult.value = null
}

// 提交测试请求
const handleTestSubmit = async () => {
  if (!selectedAPI.value) return

  testing.value = true
  try {
    let params = {}
    try {
      params = JSON.parse(testParams.value)
    } catch (e) {
      ElMessage.error('请输入有效的JSON格式参数')
      return
    }

    const result = await testAPI(selectedAPI.value.id, params)
    testResult.value = result
    testSuccess.value = true
  } catch (error) {
    testResult.value = error.response?.data || error.message
    testSuccess.value = false
  } finally {
    testing.value = false
  }
}

// 格式化JSON显示
const formatJSON = (data) => {
  try {
    return JSON.stringify(data, null, 2)
  } catch (e) {
    return data
  }
}

// 获取请求方法对应的标签类型
const getMethodTagType = (method) => {
  const types = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger'
  }
  return types[method] || 'info'
}

// 生命周期
onMounted(() => {
  loadAPIList()
})
</script>

<style scoped>
.api-management {
  padding: 20px;
  height: 100%;
}

.el-container {
  height: 100%;
}

/* 主内容区域宽度控制 */
.el-main {
  transition: all 0.3s;
  padding: 0;
}

.el-main.with-test-panel {
  width: 60%;
}

/* 测试面板样式 */
.el-aside {
  background: var(--el-bg-color);
  border-left: 1px solid #dcdfe6;
  transition: all 0.3s;
}

.test-panel-header {
  padding: 15px;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.api-test-panel {
  padding: 20px;
  height: calc(100% - 52px);
  overflow-y: auto;
}

.toolbar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.test-params,
.test-result {
  margin: 20px 0;
}

.response-data {
  margin-top: 10px;
  padding: 10px;
  background: var(--el-text-color-primary);
  color: var(--el-bg-color);
  border-radius: 4px;
  overflow-x: auto;
  white-space: pre-wrap;
  font-family: monospace;
}

.api-info h3 {
  margin-top: 0;
}

.api-info p {
  color: #666;
  margin-bottom: 20px;
}
</style> 