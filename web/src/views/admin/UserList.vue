<!-- 核心目的：用户管理页面 -->
<!-- 模块功能：用户列表CRUD、状态管理、重置密码 -->
<script setup lang="ts">
import { message, Modal } from 'ant-design-vue'
import { h, onMounted, ref } from 'vue'
import { adminApi, type RoleInfo, type UserInfo } from '@/api/admin'

const loading = ref(false)
const users = ref<UserInfo[]>([])
const total = ref(0)
const roles = ref<RoleInfo[]>([])

// 弹窗
const showCreateModal = ref(false)
const showPasswordModal = ref(false)
const saving = ref(false)
const selectedUserId = ref('')

const createForm = ref({
  username: '',
  password: '',
  real_name: '',
  role_id: '',
  department_id: '',
})
const passwordForm = ref({ new_password: '' })

const STATUS_MAP: Record<string, { color: string, label: string }> = {
  active: { label: '正常', color: 'green' },
  inactive: { label: '停用', color: 'red' },
}

const columns = [
  { title: '用户名', dataIndex: 'username', key: 'username' },
  { title: '真实姓名', dataIndex: 'real_name', key: 'real_name' },
  { title: '角色', dataIndex: 'role_name', key: 'role_name' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }: { text: string }) => {
      const s = STATUS_MAP[text] ?? { label: text, color: 'default' }
      return h('a-tag', { color: s.color }, s.label)
    },
  },
  { title: '最后登录', dataIndex: 'last_login_at', key: 'last_login_at' },
  { title: '操作', key: 'actions' },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await adminApi.listUsers()
    users.value = res.items
    total.value = res.total
  }
  finally { loading.value = false }
}

async function fetchRoles() {
  const res = await adminApi.listRoles()
  roles.value = res.items
}

onMounted(() => {
  fetchData()
  fetchRoles()
})

function openCreate() {
  createForm.value = { username: '', password: '', real_name: '', role_id: '', department_id: '' }
  showCreateModal.value = true
}

async function handleCreate() {
  if (!createForm.value.username || !createForm.value.password || !createForm.value.role_id) {
    message.warning('用户名、密码、角色为必填项')
    return
  }
  saving.value = true
  try {
    await adminApi.createUser(createForm.value)
    message.success('创建成功')
    showCreateModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function openResetPassword(record: UserInfo) {
  selectedUserId.value = record.id
  passwordForm.value = { new_password: '' }
  showPasswordModal.value = true
}

async function handleResetPassword() {
  if (!passwordForm.value.new_password || passwordForm.value.new_password.length < 6) {
    message.warning('密码至少6位')
    return
  }
  saving.value = true
  try {
    await adminApi.resetPassword(selectedUserId.value, passwordForm.value.new_password)
    message.success('密码已重置')
    showPasswordModal.value = false
  }
  finally { saving.value = false }
}

function handleToggleStatus(record: UserInfo) {
  const next = record.status === 'active' ? 'inactive' : 'active'
  const label = next === 'inactive' ? '停用' : '启用'
  Modal.confirm({
    title: `确认${label}`,
    content: `确定要${label}用户「${record.username}」吗？`,
    okType: next === 'inactive' ? 'danger' : 'primary',
    onOk: async () => {
      await adminApi.updateUser(record.id, { status: next })
      message.success(`已${label}`)
      fetchData()
    },
  })
}
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>用户管理</template>
    <template #extra>
      <a-space>
        <span class="text-muted" style="font-size:13px;">共 {{ total }} 位用户</span>
        <a-button @click="fetchData">刷新</a-button>
        <a-button type="primary" @click="openCreate">新建用户</a-button>
      </a-space>
    </template>

    <a-table
      :columns="columns"
      :data-source="users"
      :loading="loading"
      row-key="id"
      size="middle"
      :pagination="{ pageSize: 20, total, showTotal: (t: number) => `共 ${t} 条` }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openResetPassword(record as UserInfo)">
              重置密码
            </a-button>
            <a-button
              type="link"
              size="small"
              :danger="(record as UserInfo).status === 'active'"
              @click="handleToggleStatus(record as UserInfo)"
            >
              {{ (record as UserInfo).status === 'active' ? '停用' : '启用' }}
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <!-- 新建用户弹窗 -->
    <a-modal
      v-model:open="showCreateModal"
      title="新建用户"
      :confirm-loading="saving"
      @ok="handleCreate"
    >
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="用户名" required>
          <a-input v-model:value="createForm.username" placeholder="3~50个字符" />
        </a-form-item>
        <a-form-item label="初始密码" required>
          <a-input-password v-model:value="createForm.password" placeholder="至少6位" />
        </a-form-item>
        <a-form-item label="真实姓名">
          <a-input v-model:value="createForm.real_name" />
        </a-form-item>
        <a-form-item label="角色" required>
          <a-select v-model:value="createForm.role_id" placeholder="请选择角色">
            <a-select-option v-for="r in roles" :key="r.id" :value="r.id">
              {{ r.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 重置密码弹窗 -->
    <a-modal
      v-model:open="showPasswordModal"
      title="重置密码"
      :confirm-loading="saving"
      @ok="handleResetPassword"
    >
      <a-form :model="passwordForm" layout="vertical">
        <a-form-item label="新密码" required>
          <a-input-password v-model:value="passwordForm.new_password" placeholder="至少6位" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>
