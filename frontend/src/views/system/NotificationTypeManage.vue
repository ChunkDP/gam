<template>
  <div class="ma-search-box">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>通知类型管理</span>
          <el-button type="primary" @click="openCreateDialog">创建通知类型</el-button>
        </div>
      </template>
      
      <!-- 通知类型列表 -->
      <el-table
        v-loading="loading"
        :data="notificationTypes"
        border
        style="width: 100%"
      >
        <el-table-column prop="name" label="类型名称" min-width="150" />
        <el-table-column prop="code" label="类型编码" min-width="150" />
        <el-table-column prop="description" label="类型描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="editNotificationType(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="deleteNotificationType(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 通知类型表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑通知类型' : '创建通知类型'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="类型名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入类型名称" />
        </el-form-item>
        <el-form-item label="类型编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入类型编码" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="类型描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入类型描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { notificationApi } from '@/services/notification';

// 数据定义
const notificationTypes = ref([]);
const loading = ref(false);
const dialogVisible = ref(false);
const isEdit = ref(false);
const currentId = ref(null);

// 表单相关
const formRef = ref(null);
const form = reactive({
  name: '',
  code: '',
  description: ''
});

const rules = {
  name: [{ required: true, message: '请输入类型名称', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入类型编码', trigger: 'blur' },
    { pattern: /^[A-Z_]+$/, message: '类型编码只能包含大写字母和下划线', trigger: 'blur' }
  ]
};

// 加载通知类型列表
const loadNotificationTypes = async () => {
  loading.value = true;
  try {
    const data = await notificationApi.getNotificationTypes();
    notificationTypes.value = data;
  } catch (error) {
    console.error('Failed to load notification types:', error);
    ElMessage.error('加载通知类型失败');
  } finally {
    loading.value = false;
  }
};

// 打开创建对话框
const openCreateDialog = () => {
  isEdit.value = false;
  currentId.value = null;
  resetForm();
  dialogVisible.value = true;
};

// 编辑通知类型
const editNotificationType = async (row) => {
  isEdit.value = true;
  currentId.value = row.id;
  resetForm();
  
  Object.assign(form, {
    name: row.name,
    code: row.code,
    description: row.description
  });
  
  dialogVisible.value = true;
};

// 删除通知类型
const deleteNotificationType = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该通知类型吗？删除后可能影响相关通知的显示。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await notificationApi.deleteNotificationType(row.id);
    ElMessage.success('删除通知类型成功');
    loadNotificationTypes();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete notification type:', error);
      ElMessage.error('删除通知类型失败');
    }
  }
};

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  } else {
    form.name = '';
    form.code = '';
    form.description = '';
  }
};

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          await notificationApi.updateNotificationType(currentId.value, form);
          ElMessage.success('更新通知类型成功');
        } else {
          await notificationApi.createNotificationType(form);
          ElMessage.success('创建通知类型成功');
        }
        
        dialogVisible.value = false;
        loadNotificationTypes();
      } catch (error) {
        console.error('Failed to save notification type:', error);
        ElMessage.error('保存通知类型失败');
      }
    }
  });
};

// 页面加载时初始化
onMounted(() => {
  loadNotificationTypes();
});
</script>

<style scoped>

</style> 