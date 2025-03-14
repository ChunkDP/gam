<template>
  <div class="ma-search-box">
    <!-- 搜索表单 -->

    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd">新增角色</el-button>
        </div>
      </template>

     
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="角色名称">
          <el-input v-model="searchInfo.name" placeholder="角色名称" />
        </el-form-item>
        <el-form-item label="状态" >
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable class="adaptive-select">
            <el-option label="启用" :value="1"  />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
  



      <el-table 
        :data="roleList" 
        :border="true" 
        style="width: 100%"
        row-key="id"
        @sort-change="handleSortChange"
      >
        <el-table-column 
          prop="sort" 
          label="排序" 
          width="120"
          sortable="custom"
        >
          <template #default="scope">
           
            <div 
              v-if="!scope.row.is_preset" 
              @click.self="handleStartEdit(scope.row)"
              class="sort-cell"
            >
              <span v-if="editingId !== scope.row.id">
                {{ scope.row.sort }} 
              </span>
              <el-input-number 
                v-else
                v-model="scope.row.sort" 
                :min="0"
                :max="9999"
                size="small"
                
                @blur="() => handleEndEdit(scope.row)"
                @mousedown.stop="(event) => handleMouseDown(event)"
                v-focus
              />
            </div>
            <span v-else>{{ scope.row.sort }}</span>
          </template>
        </el-table-column>
        <el-table-column 
          prop="name" 
          label="角色名称"
          sortable="custom"
        />
        <el-table-column 
          prop="code" 
          label="角色编码"
          sortable="custom"
        />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              :active-value="1"
              :inactive-value="0"
              :disabled="scope.row.is_preset"
              @change="(value) => handleStatusChange(scope.row.id, value)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button 
              size="small" 
              type="primary"
              @click="handlePermission(scope.row)"
            >权限设置</el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(scope.row)"
              :disabled="scope.row.is_preset"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 添加分页组件 -->
      <div class="pagination">
        <el-pagination background 
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          layout="total, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange">
         </el-pagination>
          </div>
</el-card>
    <!-- 角色表单对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="roleForm.name" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="roleForm.code" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input 
            v-model="roleForm.description" 
            type="textarea" 
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 权限设置对话框 -->
    <el-dialog
      title="权限设置"
      v-model="permissionVisible"
      width="600px"
    
    >
      <el-tree
        ref="menuTreeRef"
        :data="menuTree"
        show-checkbox
        node-key="id"
       
        :props="{
          label: 'title',
          children: 'children'
        }"
      />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="permissionVisible = false">取消</el-button>
          <el-button type="primary" @click="handlePermissionSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import {  getRoleMenus, updateRoleMenus,getRoleList, createRole, updateRole, deleteRole, updateRoleStatus, updateRoleSort } from '@/services/role'

import { checkRoleFieldUnique } from '@/services/role'

// 状态定义
const roleList = ref([])
// 处理表格排序变化
const handleSortChange = ({ prop, order }) => {
  // prop: 排序的字段名
  // order: ascending/descending/null
  
  let sortOrder = null
  if (order === 'ascending') {
    sortOrder = 'asc'
  } else if (order === 'descending') {
    sortOrder = 'desc'
  }

  // 更新搜索条件
  searchInfo.value.sortField = prop
  searchInfo.value.sortOrder = sortOrder
  
  // 重新获取列表
  fetchRoleList()
}

// 添加分页相关的状态
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

// 修改 searchInfo 的初始值，添加分页相关字段
const searchInfo = ref({
  name: '',
  status: '',
  sortField: '',
  sortOrder: '',
  page: 1,
  pageSize: 10
})
const dialogVisible = ref(false)
const dialogTitle = ref('')
const permissionVisible = ref(false)
const currentRole = ref(null)
const menuTree = ref([])
const roleFormRef = ref(null)
const roleForm = ref({
  name: '',
  code: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' }
  ]
}

const menuTreeRef = ref(null)

// 获取角色列表
const fetchRoleList = async () => {
  try {
    const params = { 
      ...searchInfo.value,
      page: pagination.value.page,
      pageSize: pagination.value.pageSize
    }
    const res = await getRoleList(params)
    roleList.value = res.roles
    pagination.value.total = res.total
    
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  }
}

// 搜索
const handleSearch = () => {
 
  fetchRoleList()
}

// 重置搜索
const resetSearch = () => {
  searchInfo.value = {
    name: '',
    status: '',
    sortField: '',
    sortOrder: ''
  }
  pagination.value.page = 1 // 重置到第一页
  fetchRoleList()
}

// 新增角色
const handleAdd = () => {
  dialogTitle.value = '新增角色'
  roleForm.value = {
    name: '',
    code: '',
    description: ''
  }
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row) => {
  dialogTitle.value = '编辑角色'
  roleForm.value = { ...row }
  
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  const formEl = roleFormRef.value
  if (!formEl) return

  await formEl.validate(async (valid) => {
    if (valid) {
      try {
        // 检查名称唯一性
        
        const nameCheck = await checkRoleFieldUnique('name', roleForm.value.name, roleForm.value.id)
        if (!nameCheck.unique) {
          ElMessage.error('角色名称已存在')
          return
        }

        // 检查编码唯一性
        const codeCheck = await checkRoleFieldUnique('code', roleForm.value.code, roleForm.value.id)
        if (!codeCheck.unique) {
          ElMessage.error('角色编码已存在')
          return
        }

        if (roleForm.value.id) {
          await updateRole(roleForm.value.id, roleForm.value)
          ElMessage.success('更新成功')
        } else {
          await createRole(roleForm.value)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchRoleList()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      }
    }
  })
}

// 删除角色
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该角色吗？', '提示', {
    type: 'warning',
    
      confirmButtonText: '确定',
      cancelButtonText: '取消',
    
  }).then(async () => {
    try {
      await deleteRole(row.id)
      ElMessage.success('删除成功')
      fetchRoleList()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 打开权限设置
const handlePermission = async (row) => {
  currentRole.value = row
  try {
   
    const res = await getRoleMenus(row.id)
   
    menuTree.value = res.menuTree
    permissionVisible.value = true
   
    // 等待 DOM 更新后设置选中状态
    await nextTick()
  
    if (menuTreeRef.value) {
      res.checkedMenus?.forEach(id => {
        menuTreeRef.value.setChecked(id, true, false)
      })
    }
    await nextTick()
  } catch (error) {
    ElMessage.error('获取权限数据失败')
  }
}

// 提交权限设置
const handlePermissionSubmit = async () => {
  const menuTree = menuTreeRef.value
  if (!menuTree || !currentRole.value) return
  
  const checkedKeys = menuTree.getCheckedKeys()
  const halfCheckedKeys = menuTree.getHalfCheckedKeys()
  const menuIds = [...checkedKeys, ...halfCheckedKeys]
  
  try {
    await updateRoleMenus(currentRole.value.id, menuIds)
    ElMessage.success('权限设置成功')
    permissionVisible.value = false
  } catch (error) {
    ElMessage.error(error.message || '权限设置失败')
  }
}

// 处理状态变更
const handleStatusChange = async (id, status) => {
  
  try {
    await updateRoleStatus(id, { status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    // 恢复原状态
    const role = roleList.value.find(r => r.id === id)
    if (role) {
      role.status = status === 1 ? 0 : 1
    }
  }
}

const originalSort = ref(null)


// 添加编辑状态控制
const editingId = ref(null)

// 开始编辑
const handleStartEdit = (row) => {


  editingId.value = row.id

  originalSort.value = row.sort
}


// 结束编辑并更新排序
const handleEndEdit = async (row) => {
  // 检查是否点击了输入框的按钮
 
  // 如果排序值没有变化，直接结束编辑

  if (!editingId.value) return

  // 如果排序值没有变化，直接结束编辑
  if (originalSort.value === row.sort) {
    editingId.value = null
    return
  }

  const loading = ElLoading.service({
    lock: true,
    text: '更新中...',
    background: 'rgba(0, 0, 0, 0.7)'
  })
  
  try {
    await updateRoleSort(row.id, { sort: row.sort })
    loading.close()
    ElMessage.success('排序更新成功')
    // 重新获取列表以确保排序正确
    await fetchRoleList()
  } catch (error) {
    loading.close()
    ElMessage.error(error.message || '排序更新失败')
    // 重新获取列表以恢复原始状态
    await fetchRoleList()
  }
  
  editingId.value = null
  originalSort.value = null  // 清除原始值
}

// 处理 mousedown 事件
const handleMouseDown = (event) => {
  // 阻止默认行为，防止 blur 事件立即触发
  event.preventDefault()
  event.stopPropagation()
}


// 自定义指令：自动聚焦
// 添加自定义指令
const vFocus = {
  mounted: (el) => {
    // 找到实际的输入框元素并聚焦
    const input = el.querySelector('input')
    if (input) {
      input.focus()

    }
  }
}

// 页面加载时获取角色列表
onMounted(() => {
  fetchRoleList()
})

// 处理每页条数变化
const handleSizeChange = (val) => {
  pagination.value.pageSize = val
  pagination.value.page = 1 // 重置到第一页
  fetchRoleList()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  pagination.value.page = val
  fetchRoleList()
}
</script>

<style scoped>

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sort-cell {
  cursor: pointer;
  padding: 5px;
}

.sort-cell:hover {
  background-color: #f5f7fa;
}

:deep(.el-input-number) {
  width: 80px;
}

:deep(.el-input-number.is-disabled) {
  opacity: 0.7;
}

:deep(.adaptive-select) {
  width: fit-content;
  min-width: 120px;
 
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>