<!-- 核心目的：ItemAliasManager页面 -->
<!-- 模块功能：资源管理-ItemAlias管理 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { resourceApi } from '@/api/resource'
import { usePagination } from '@/composables/usePagination'
import type { ExamItem } from '@/types/resource'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<ExamItem>(
  params => resourceApi.listExamItems(params),
)
onMounted(() => fetchData())

// 展开行管理
const expandedRowKeys = ref<string[]>([])
const newAlias = ref<Record<string, string>>({})
const addingAlias = ref<Record<string, boolean>>({})

function toggleExpand(record: ExamItem) {
  const idx = expandedRowKeys.value.indexOf(record.id)
  if (idx >= 0) {
    expandedRowKeys.value.splice(idx, 1)
  }
  else {
    expandedRowKeys.value.push(record.id)
  }
}

async function addAlias(item: ExamItem) {
  const alias = newAlias.value[item.id]?.trim()
  if (!alias) {
    message.warning('别名不能为空')
    return
  }
  addingAlias.value[item.id] = true
  try {
    await resourceApi.addAlias(item.id, alias)
    message.success('别名添加成功')
    newAlias.value[item.id] = ''
    fetchData()
  }
  finally {
    addingAlias.value[item.id] = false
  }
}

function removeAlias(item: ExamItem, aliasName: string) {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除别名「${aliasName}」吗？`,
    okType: 'danger',
    onOk: async () => {
      await resourceApi.deleteAlias(item.id, aliasName)
      message.success('删除成功')
      fetchData()
    },
  })
}

const columns = [
  { title: '检查项目名称', dataIndex: 'name', key: 'name' },
  { title: '标准时长(分钟)', dataIndex: 'duration_min', key: 'duration_min' },
  { title: '别名数量', key: 'alias_count', customRender: ({ record }: { record: ExamItem }) => record.aliases?.length ?? 0 },
  { title: '操作', key: 'actions', width: 100 },
]
</script>

<template>
  <div>
    <div class="mb-4 flex items-center gap-2">
      <a-button :loading="loading" @click="fetchData">
        刷新
      </a-button>
      <span class="text-gray-400 text-sm">点击「管理别名」展开行内编辑</span>
    </div>

    <a-table
      :columns="columns"
      :data-source="items"
      :loading="loading"
      :pagination="pagination"
      :expanded-row-keys="expandedRowKeys"
      row-key="id"
      size="middle"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'alias_count'">
          <a-tag v-if="(record as ExamItem).aliases?.length" color="blue">
            {{ (record as ExamItem).aliases!.length }} 个别名
          </a-tag>
          <span v-else class="text-gray-400">—</span>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-button
            type="link"
            size="small"
            @click="toggleExpand(record as ExamItem)"
          >
            {{ expandedRowKeys.includes((record as ExamItem).id) ? '收起' : '管理别名' }}
          </a-button>
        </template>
      </template>

      <template #expandedRowRender="{ record }">
        <div style="padding: 12px 16px; background: #fafafa; border-radius: 6px;">
          <div style="margin-bottom: 12px; font-size: 13px; color: #595959; font-weight: 500;">
            「{{ (record as ExamItem).name }}」的别名列表
            <span style="color: #8c8c8c; font-weight: 400;">（支持HIS系统不同命名识别）</span>
          </div>
          <div style="display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 12px;">
            <a-tag
              v-for="alias in (record as ExamItem).aliases"
              :key="alias.alias"
              closable
              color="blue"
              style="margin-bottom: 4px;"
              @close="removeAlias(record as ExamItem, alias.alias)"
            >
              {{ alias.alias }}
            </a-tag>
            <span v-if="!(record as ExamItem).aliases?.length" style="color: #8c8c8c; font-size: 12px;">
              暂无别名，请在下方添加
            </span>
          </div>
          <a-input-search
            v-model:value="newAlias[(record as ExamItem).id]"
            placeholder="输入新别名，回车或点击添加"
            enter-button="添加"
            style="width: 300px;"
            :loading="addingAlias[(record as ExamItem).id]"
            @search="addAlias(record as ExamItem)"
          />
        </div>
      </template>
    </a-table>
  </div>
</template>
