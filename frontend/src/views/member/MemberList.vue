<template>
  <!-- 查询表单 -->
  <div class="ma-search-box">
  <el-form ref="searchForm"  :model="searchInfo" :inline="true" >
    <el-form-item label="用户名">
      <el-input v-model="searchInfo.username" placeholder="请输入用户名" />
    </el-form-item>
    <el-form-item label="手机号">
      <el-input v-model="searchInfo.mobile" placeholder="请输入手机号" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <!-- 添加重置按钮 -->
      <el-button @click="handleReset">重置</el-button>
    </el-form-item>
  </el-form>
</div>
  <!-- 操作按钮 -->
  <div class="ma-btn-list">
    <el-button type="primary" @click="handleAdd">添加会员</el-button>
  </div>

  <!-- 数据表格 -->
  <el-table
    v-loading="loading"
    :data="tableData"
    border
     style="width: 100%"
  >
    <el-table-column prop="username" label="用户名"  />
    <el-table-column prop="mobile" label="手机号" />
    <el-table-column prop="email" label="邮箱"  />
    <el-table-column prop="status" label="状态" width="80">
      <template #default="scope">
        <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
          {{ scope.row.status === 1 ? '启用' : '禁用' }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="180" fixed="right">
      <template #default="scope">
        <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
        <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>

  <!-- 分页组件 -->
  <div class="ma-pagination">
    <el-pagination
      background
      layout="total, prev, pager, next"
      :total="pagination.total"
      :page-size="pagination.pageSize"
      v-model:current-page="pagination.page"
      @current-change="handleCurrentChange"
    />
  </div>

  <!-- 对话框 -->
  <el-dialog
    :title="dialogTitle"
    v-model="dialogVisible"
    width="500px"
    @close="resetForm"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="手机号" prop="mobile">
        <el-input v-model="form.mobile" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择状态">
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit">确 定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMemberList, createMember, updateMember, deleteMember, checkMemberFieldUnique } from '@/services/member'

// 搜索表单
const searchInfo = ref({
  username: '',
  mobile: '',
})

// 表格数据
const tableData = ref([])
const loading = ref(false)
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const form = ref({
  id: null,
  username: '',
  mobile: '',
  email: '',
  status: 1,
})

// 表单校验规则
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  mobile: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

// 获取会员列表
const fetchMemberList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      ...searchInfo.value,
    }
    const data = await getMemberList(params)
   
    tableData.value = data.members
    pagination.value.total = data.total
  } catch (error) {
    ElMessage.error(error.message || '获取会员列表失败')
  } finally {
    loading.value = false
  }
}


// 处理页码变化
const handleCurrentChange = (val) => {
  pagination.value.page = val
  fetchMemberList()
}

// 搜索时重置分页
const handleSearch = () => {
  pagination.value.page = 1
  fetchMemberList()
}

// 重置搜索时也重置分页
const handleReset = () => {
  searchInfo.value = {
    username: '',
    phone: '',
  }
  pagination.value.page = 1
  fetchMemberList()
}

// 添加会员
const handleAdd = () => {
  dialogTitle.value = '添加会员'
  form.value = {
    id: null,
    username: '',
    mobile: '',
    email: '',
    status: 1,
  }
  dialogVisible.value = true
}

// 编辑会员
const handleEdit = (row) => {
  dialogTitle.value = '编辑会员'
  form.value = {
    id: row.id,
    username: row.username,
    mobile: row.mobile,
    email: row.email,
    status: row.status,
  }
  dialogVisible.value = true
}

// 删除会员
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该会员吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await deleteMember(row.id)
      ElMessage.success('删除成功')
      fetchMemberList()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 提交表单
const handleSubmit = async () => {
  formRef.value.validate(async (valid) => {
    if (valid) {
      try {
            // 检查用户名唯一性
            const usernameCheck = await checkMemberFieldUnique('username', form.value.username, form.value.id)
        if (!usernameCheck.unique) {
          ElMessage.error('用户名已存在')
          return
        }

        // 检查手机号唯一性
        const mobileCheck = await checkMemberFieldUnique('mobile', form.value.mobile, form.value.id)
        if (!mobileCheck.unique) {
          ElMessage.error('手机号已存在')
          return
        }

        // 检查邮箱唯一性
        const emailCheck = await checkMemberFieldUnique('email', form.value.email, form.value.id)
        if (!emailCheck.unique) {
          ElMessage.error('邮箱已存在')
          return
        }

        if (form.value.id) {
          await updateMember(form.value.id, form.value)
          ElMessage.success('更新成功')
        } else {
          await createMember(form.value)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchMemberList()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      }
    } else {
      return false
    }
  })
}

// 重置表单
const resetForm = () => {
  formRef.value.resetFields()
}

onMounted(() => {
  fetchMemberList()
})
</script>

<style scoped>
.member-list {
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