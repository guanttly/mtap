<!-- 核心目的：PriorityTagList页面 -->
<!-- 模块功能：规则引擎-PriorityTag列表管理CRUD与配置 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ruleApi } from '@/api/rule'
import type { PriorityTag } from '@/types/rule'

const loading = ref(false)
const items = ref<PriorityTag[]>([])

async function fetchData() {
  loading.value = true
  try {
    const res = await ruleApi.listPriorityTags()
    items.value = res.items
  }
  finally { loading.value = false }
}
onMounted(fetchData)

const showModal = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref('')
const editForm = ref<Partial<PriorityTag>>({ color: '#1890ff', weight: 10 })

const columns = [
  { title: '标签名', dataIndex: 'name', key: 'name' },
  { title: '权重', dataIndex: 'weight', key: 'weight' },
  { title: '颜色', key: 'color' },
  { title: '类型', key: 'is_preset' },
  { title: '操作', key: 'actions', width: 120 },
]

function openCreate() {
  isEdit.value = false
  editId.value = ''
  editForm.value = { color: '#1890ff', weight: 10 }
  showModal.value = true
}

function openEdit(record: PriorityTag) {
  if (record.is_preset) {
    message.warning('预置标签仅允许调整权重和颜色')
  }
  isEdit.value = true
  editId.value = record.id
  editForm.value = { name: record.name, weight: record.weight, color: record.color }
  showModal.value = true
}

async function handleSave() {
  if (!editForm.value.name?.trim() && !isEdit.value) {
    message.warning('请填写标签名称')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await ruleApi.updatePriorityTag(editId.value, {
        name: editForm.value.name,
        weight: editForm.value.weight,
        color: editForm.value.color,
      })
      message.success('更新成功')
    }
    else {
      await ruleApi.createPriorityTag(editForm.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleDelete(record: PriorityTag) {
  if (record.is_preset) {
    message.warning('预置标签不可删除')
    return
  }
  Modal.confirm({
    title: '确认删除',
    content: `确定删除优先级标签「${record.name}」吗？`,
    okType: 'danger',
    onOk: async () => {
      await ruleApi.deletePriorityTag(record.id)
      message.success('删除成功')
      fetchData()
    },
  })
}
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>优先级标签</template>
    <template #extra>
      <a-space>
        <a-button @click="fetchData">刷新</a-button>
        <a-button type="primary" @click="openCreate">新增优先级标签</a-button>
      </a-space>
    </template>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="false" row-key="id" size="middle">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'color'">
          <span :style="{ display:'inline-block', width:'16px', height:'16px', borderRadius:'4px', background: (record as PriorityTag).color, verticalAlign:'middle', marginRight:'6px' }" />
          <span style="font-family: monospace;">{{ (record as PriorityTag).color }}</span>
        </template>
        <template v-else-if="column.key === 'is_preset'">
          <a-tag :color="(record as PriorityTag).is_preset ? 'purple' : 'default'">
            {{ (record as PriorityTag).is_preset ? '系统预置' : '自定义' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as PriorityTag)">编辑</a-button>
            <a-button
              type="link" danger size="small"
              :disabled="(record as PriorityTag).is_preset"
              @click="handleDelete(record as PriorityTag)"
            >删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="isEdit ? '编辑优先级标签' : '新增优先级标签'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="标签名称" :required="!isEdit">
          <a-input v-model:value="editForm.name" placeholder="例如：急诊优先、老年患者" :maxlength="20" show-count />
        </a-form-item>
        <a-form-item label="权重（1~100，越高越优先）">
          <a-input-number v-model:value="editForm.weight" :min="1" :max="100" style="width:100%;" />
        </a-form-item>
        <a-form-item label="标签颜色">
          <div style="display:flex;align-items:center;gap:12px;">
            <input v-model="editForm.color" type="color" style="width:48px;height:32px;border:none;cursor:pointer;" />
            <a-input v-model:value="editForm.color" style="width:120px;font-family:monospace;" placeholder="#1890ff" />
            <span :style="{ display:'inline-block', padding:'2px 12px', borderRadius:'12px', background: editForm.color, color:'#fff', fontSize:'12px' }">预览</span>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>
