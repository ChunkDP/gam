<template>
   <div class="ma-search-box">
    <!-- 搜索表单 -->

    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <el-button type="primary" @click="handleAdd">新增菜单</el-button>
        </div>
      </template>


    <el-table
      :data="menuList"
      row-key="id"
      border
      :tree-props="{ children: 'children' }"
    >

      <el-table-column prop="title" label="菜单名称" width="220" />
      <el-table-column prop="name" label="路由名称" />
      <el-table-column prop="path" label="路由路径" />
      <el-table-column prop="component" label="组件路径" />
      <el-table-column prop="type" label="类型" >
        <template #default="scope">
          <el-tag type="success" v-if="scope.row.type === 'menu'">菜单</el-tag>
          <el-tag type="danger" v-else>按钮</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="icon" label="图标" width="70">
       
        <template #default="scope">
           <el-icon>
          
               <component :is="getIconComponent(scope.row.icon)" />
           </el-icon>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="250" fixed="right">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="primary" @click="handleAdd(scope.row)">
            添加子菜单
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.row)"
            :disabled="hasChildren(scope.row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>
  </div>
    <!-- 菜单表单对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="menuFormRef"
        :model="menuForm"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="上级菜单" >
          <el-tree-select
            v-model="menuForm.parent_id"
            :data="menuOptions"
            :props="{ label: 'title', value: 'id' }"
            check-strictly
            :render-after-expand="false"
            placeholder="请选择上级菜单"
            clearable
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="type">
          <el-radio-group v-model="menuForm.type">
            <el-radio value="menu">菜单</el-radio>
            <el-radio value="button">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称" prop="title">
          <el-input v-model="menuForm.title" placeholder="请输入菜单名称" />
        </el-form-item>
        
        <!-- 仅当类型为 menu 时显示的字段 -->
        <template v-if="menuForm.type === 'menu'">
          <el-form-item label="路由名称" prop="name">
            <el-input v-model="menuForm.name" placeholder="请输入路由名称" />
          </el-form-item>
          <el-form-item label="路由路径" prop="path">
            <el-input v-model="menuForm.path" placeholder="请输入路由路径" />
          </el-form-item>
          <el-form-item label="组件路径" prop="component">
            <el-input v-model="menuForm.component" placeholder="请输入组件路径" />
          </el-form-item>
          <el-form-item label="图标" prop="icon">
            <el-popover
              placement="bottom"
              :width="400"
              trigger="click"
              popper-class="icon-popover"
            >
              <template #reference>
                <el-input v-model="menuForm.icon" placeholder="点击选择图标">
                  <template #prefix>
                    <el-icon v-if="menuForm.icon">
                      <component :is="menuForm.icon" />
                    </el-icon>
                  </template>
                </el-input>
              </template>
              <div class="icon-list">
                <div
                  v-for="icon in elementIcons"
                  :key="icon"
                  class="icon-item"
                  @click="selectIcon(icon)"
                >
                  <el-icon>
                    <component :is="icon" />
                  </el-icon>
                  <span>{{ icon }}</span>
                </div>
              </div>
            </el-popover>
          </el-form-item>
        </template>

        <!-- 仅当类型为 button 时显示的字段 -->
        <template v-if="menuForm.type === 'button'">
          <el-form-item label="权限标识" prop="permission">
            <el-input 
              v-model="menuForm.permission" 
              placeholder="请输入权限标识，如：system:user:add"
            />
          </el-form-item>
          <el-form-item label="接口方法" prop="api_method">
            <el-select v-model="menuForm.api_method" placeholder="请选择接口方法">
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="接口路径" prop="api_path">
            <el-input 
              v-model="menuForm.api_path" 
              placeholder="请输入接口路径，如：/api/users"
            />
          </el-form-item>
        </template>

        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="menuForm.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="menuForm.status"
            :active-value="1"
            :inactive-value="0"
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
  
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import {
  getMenuTree,
  createMenu,
  updateMenu,
  deleteMenu,
  updateMenuStatus,
  updateMenuSort
} from '@/services/menu'

const menuList = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const menuFormRef = ref(null)
const menuOptions = ref([])

const menuForm = ref({
  parent_id: null,
  type: 'menu', // 默认为菜单类型
  title: '',
  name: '',
  path: '',
  component: '',
  icon: '',
  permission: '',
  api_method: '',
  api_path: '',
  sort: 0,
  status: 1
})

const rules = {
  title: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }],
  name: [{ 
    required: true, 
    message: '请输入路由名称', 
    trigger: 'blur',
    validator: (rule, value, callback) => {
      if (menuForm.value.type === 'menu' && !value) {
        callback(new Error('请输入路由名称'))
      } else {
        callback()
      }
    }
  }],
  permission: [{
    required: true,
    message: '请输入权限标识',
    trigger: 'blur',
    validator: (rule, value, callback) => {
      if (menuForm.value.type === 'button' && !value) {
        callback(new Error('请输入权限标识'))
      } else {
        callback()
      }
    }
  }],
  api_method: [{
    required: true,
    message: '请选择接口方法',
    trigger: 'change',
    validator: (rule, value, callback) => {
      if (menuForm.value.type === 'button' && !value) {
        callback(new Error('请选择接口方法'))
      } else {
        callback()
      }
    }
  }],
  api_path: [{
    required: true,
    message: '请输入接口路径',
    trigger: 'blur',
    validator: (rule, value, callback) => {
      if (menuForm.value.type === 'button' && !value) {
        callback(new Error('请输入接口路径'))
      } else {
        callback()
      }
    }
  }]
}

// 获取菜单列表
const fetchMenuList = async () => {
  try {
    const res = await getMenuTree()
   
    menuList.value = res.menuTree
   
    menuOptions.value = [{ id: 0, title: '顶级菜单' }, ...res.menuTree]
  } catch (error) {
    console.log("error", error)
    ElMessage.error('获取菜单列表失败')
  }
}

// 检查是否有子菜单
const hasChildren = (row) => {
  return row.children && row.children.length > 0
}
const getIconComponent = (iconName) => {
  if (!iconName) {
        return 'DefaultIcon'; // 默认组件
      }
      return iconName;
}
// 新增菜单
const handleAdd = (row = null) => {
 
  dialogTitle.value = row.id ? '新增子菜单' : '新增菜单'
  menuForm.value = {
    parent_id: row.id ? row.id : 0,
    title: '',
    name: '',
    path: '',
    component: '',
    icon: '',
    sort: 0,
    status: 1
  }
  dialogVisible.value = true
}

// 编辑菜单
const handleEdit = (row) => {
  dialogTitle.value = '编辑菜单'
  menuForm.value = { ...row }
  dialogVisible.value = true
}

// 删除菜单
const handleDelete = (row) => {
  if (hasChildren(row)) {
    ElMessage.warning('该菜单包含子菜单，无法删除')
    return
  }

  ElMessageBox.confirm('确定要删除该菜单吗？', '提示', {
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteMenu(row.id)
        ElMessage.success('删除成功')
        fetchMenuList()
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '删除失败')
      }
    })
    .catch(() => {})
}

// 提交表单
const handleSubmit = async () => {
  const formEl = menuFormRef.value
  if (!formEl) return

  await formEl.validate(async (valid) => {
    if (valid) {
      try {
        if (menuForm.value.id) {
          await updateMenu(menuForm.value.id, menuForm.value)
          ElMessage.success('更新成功')
        } else {
          await createMenu(menuForm.value)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchMenuList()
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '操作失败')
      }
    }
  })
}




// 重置表单
const resetForm = () => {
  if (menuFormRef.value) {
    menuFormRef.value.resetFields()
  }
}

// 获取 Element Plus 的所有图标
const elementIcons = ref(Object.keys(ElementPlusIconsVue))

// 选择图标
const selectIcon = (icon) => {
  menuForm.value.icon = icon
}

onMounted(() => {
  fetchMenuList()
})
</script>

<style scoped>



.icon-list {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  cursor: pointer;
  border-radius: 4px;
}

.icon-item:hover {
  background-color: #f5f7fa;
}

.icon-item span {
  font-size: 12px;
  margin-top: 4px;
  color: #666;
}

:deep(.icon-popover) {
  max-width: 460px;
  padding: 12px;
}
</style> 