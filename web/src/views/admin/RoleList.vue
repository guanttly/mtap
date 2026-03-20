<!-- 核心目的：角色管理页面 -->
<!-- 模块功能：角色列表、权限配置、新增/删除 -->
<script setup lang="ts">
import { message, Modal } from 'ant-design-vue'
import { onMounted, ref } from 'vue'
import { adminApi, type RoleInfo } from '@/api/admin'

const loading = ref(false)
const roles = ref<RoleInfo[]>([])
const showModal = ref(false)
const saving = ref(false)
const editingRole = ref<RoleInfo | null>(null)

const createForm = ref({ name: '', permissions: [] as string[] })
const editPermissions = ref<string[]>([])

// 预定义权限项（与后端角色权限对应）
const PERMISSION_OPTIONS = [
  { label: '规则引擎 - 读', value: 'rule:read' },
  { label: '规则引擎 - 写', value: 'rule:write' },
  { label: '规则引擎 - 全部', value: 'rule:*' },
  { label: '资源管理 - 读', value: 'resource:read' },
  { label: '资源管理 - 写', value: 'resource:write' },
  { label: '资源管理 - 全部', value: 'resource:*' },
  { label: '预约服务 - 读', value: 'appt:read' },
  { label: '预约服务 - 写', value: 'appt:write' },
  { label: '预约服务 - 全部', value: 'appt:*' },
  { label: '分诊执行 - 全部', value: 'triage:*' },
  { label: '统计分析 - 读', value: 'analytics:read' },
  { label: '效能优化 - 全部', value: 'optimization:*' },
  { label: '系统管理 - 全部', value: 'admin:*' },
  { label: '全部权限', value: '*' },
]

const columns = [
  { title: '角色名', dataIndex: 'name', key: 'name' },
  { title: '权限', dataIndex: 'permissions', key: 'permissions' },
  { title: '类型', dataIndex: 'is_preset', key: 'is_preset' },
  { title: '操作', key: 'actions' },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await adminApi.listRoles()
    roles.value = res.items
  }
  finally { loading.value = false }
}

onMounted(fetchData)

function openCreate() {
  editingRole.value = null
  createForm.value = { name: '', permissions: [] }
  showModal.value = true
}

function openEdit(record: RoleInfo) {
  editingRole.value = record
  editPermissions.value = [...record.permissions]
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editingRole.value) {
      await adminApi.updateRole(editingRole.value.id, editPermissions.value)
      message.success('权限已更新')
    }
    else {
      if (!createForm.value.name) {
        message.warning('请输入角色名称')
        return
      }
      await adminApi.createRole(createForm.value)
      message.success('角色已创建')
    }
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleDelete(record: RoleInfo) {
  if (record.is_preset) {
    message.warning('预置角色不可删除')
    return
  }
  Modal.confirm({
    title: '确认删除',
    content: `确定删除角色「${record.name}」吗？请确保该角色下无用户。`,
    okType: 'danger',
    onOk: async () => {
      await adminApi.deleteRole(record.id)
      message.success('删除成功')
      fetchData()
    },
  })
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center gap-2">
      <a-button type="primary" @click="openCreate">
        新增角色
      </a-button>
      <a-button @click="fetchData">
        刷新
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="roles"
      :loading="loading"
      row-key="id"
      size="middle"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'permissions'">
          <a-space wrap>
            <a-tag v-for="p in (record as RoleInfo).permissions" :key="p" color="blue">
              {{ p }}
            </a-tag>
            <span v-if="(record as RoleInfo).permissions.length === 0" class="text-gray-400">无</span>
          </a-space>
        </template>
        <template v-else-if="column.key === 'is_preset'">
          <a-tag :color="(record as RoleInfo).is_preset ? 'purple' : 'default'">
            {{ (record as RoleInfo).is_preset ? '系统预置' : '自定义' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button
              type="link"
              size="small"
              :disabled="(record as RoleInfo).is_preset"
              @click="openEdit(record as RoleInfo)"
            >
              编辑权限
            </a-button>
            <a-button
              type="link"
              danger
              size="small"
              :disabled="(record as RoleInfo).is_preset"
              @click="handleDelete(record as RoleInfo)"
            >
              删除
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="editingRole ? `编辑权限 - ${editingRole.name}` : '新增角色'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form layout="vertical">
        <template v-if="!editingRole">
          <a-form-item label="角色名称" required>
            <a-input v-model:value="createForm.name" placeholder="2~30个字符" />
          </a-form-item>
          <a-form-item label="初始权限">
            <a-checkbox-group v-model:value="createForm.permissions" :options="PERMISSION_OPTIONS" style="display: flex; flex-wrap: wrap; gap: 4px 16px;" />
          </a-form-item>
        </template>
        <template v-else>
          <a-form-item label="权限配置">
            <a-checkbox-group v-model:value="editPermissions" :options="PERMISSION_OPTIONS" style="display: flex; flex-wrap: wrap; gap: 4px 16px;" />
          </a-form-item>
        </template>
      </a-form>
    </a-modal>
  </div>
</template>
