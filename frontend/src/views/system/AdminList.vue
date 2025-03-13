<template>
  <div class="admin-list">
    <!-- 搜索表单 -->
    <div class="ma-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="searchInfo.email" placeholder="邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">
            查询
          </el-button>
          <el-button icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="ma-table-box">
      <div class="ma-btn-list">
        <Permission permission="system:user:create">
          <el-button type="primary" @click="handleAdd">新增管理员</el-button>
        </Permission>
      </div>

      <el-table :data="tableData" border v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="real_name" label="真实姓名" />
        <el-table-column prop="role.name" label="角色" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <Permission permission="system:user:update">
              <el-button type="primary" link @click="handleEdit(row)">
                编辑
              </el-button>
            </Permission>
            <Permission permission="system:user:delete">
              <el-button type="danger" link @click="handleDelete(row)">
                删除
              </el-button>
            </Permission>
          </template>
        </el-table-column>
      </el-table>

      <div class="ma-pagination">
        <el-pagination 
          background
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!form.id">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="form.real_name" />
        </el-form-item>
        <el-form-item label="角色" prop="role_id">
          <el-select v-model="form.role_id" placeholder="请选择角色">
            <el-option
              v-for="role in roleOptions"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Permission from '@/components/Permission.vue'
import * as adminService from '@/services/adminService'
import { getRoleList } from '@/services/role'

// 搜索表单
const searchInfo = ref({
  username: '',
  email: '',
  phone: '',
})

// 表格数据
const tableData = ref([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const form = ref({
  username: '',
  password: '',
  email: '',
  phone: '',
  real_name: '',
  role_id: '',
})

// 表单校验规则
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  role_id: [{ required: true, message: '请选择角色', trigger: 'change' }],
}

// 角色选项
const roleOptions = ref([])

// 替换原来的分页变量，使用统一的分页对象
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取角色列表
const getRoleOptions = async () => {
  try {
    const res = await getRoleList({pageSize:-1})
    roleOptions.value = res.roles
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

// 获取管理员列表
const fetchAdminList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      ...searchInfo.value,
    }
    const data = await adminService.getAdminList(params)
   
    tableData.value = data.admins
    pagination.value.total = data.total
  } catch (error) {
    ElMessage.error(error.message || '获取管理员列表失败')
  } finally {
    loading.value = false
  }
}

// 处理每页条数变化
const handleSizeChange = (val) => {
  pagination.value.pageSize = val
  pagination.value.page = 1 // 重置到第一页
  fetchAdminList()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  pagination.value.page = val
  fetchAdminList()
}

// 搜索时重置分页
const handleSearch = () => {
  pagination.value.page = 1
  fetchAdminList()
}

// 重置搜索时也重置分页
const handleReset = () => {
  searchInfo.value = {
    username: '',
    email: '',
    phone: '',
  }
  pagination.value.page = 1
  fetchAdminList()
}

// 添加管理员
const handleAdd = () => {
  dialogTitle.value = '新增管理员'
  dialogVisible.value = true
  form.value = {
    username: '',
    password: '',
    email: '',
    phone: '',
    real_name: '',
    role_id: '',
  }
}

// 编辑管理员
const handleEdit = (row) => {
  dialogTitle.value = '编辑管理员'
  dialogVisible.value = true
  form.value = { ...row }
}

// 删除管理员
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该管理员吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await adminService.deleteAdmin(row.id)
      ElMessage.success('删除成功')
      fetchAdminList()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

// 更新状态
const handleStatusChange = async (row) => {
  try {
   
    await adminService.updateAdminStatus(row.id, row.status)
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error(error.message || '状态更新失败')
    row.status = row.status === 1 ? 0 : 1 // 恢复状态
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    if (form.value.id) {
      // 编辑
      await adminService.updateAdmin(form.value.id, form.value)
      ElMessage.success('更新成功')
    } else {
      // 新增
      await adminService.createAdmin(form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchAdminList()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 重置表单
const resetForm = () => {
  form.value = {
    username: '',
    password: '',
    email: '',
    phone: '',
    real_name: '',
    role_id: '',
  }
}

// 初始化
onMounted(() => {
  fetchAdminList()
  getRoleOptions()
})
</script>

<style scoped>
.admin-list {
  padding: 20px;
}
.ma-search-box {
  margin-bottom: 20px;
}
.ma-btn-list {
  margin-bottom: 20px;
}
.ma-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>