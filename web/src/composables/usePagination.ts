// 核心目的：分页组合式函数
// 模块功能：通用分页状态管理、页码切换、每页数量配置
import { ref, reactive, computed } from 'vue'
import type { PageResult } from '@/api/request'

export interface PaginationOptions {
  pageSize?: number
}

export function usePagination<T>(
  fetcher: (params: Record<string, unknown>) => Promise<PageResult<T>>,
  opts: PaginationOptions = {},
) {
  const loading = ref(false)
  const items = ref<T[]>([]) as { value: T[] }
  const total = ref(0)
  const pagination = reactive({
    current: 1,
    pageSize: opts.pageSize ?? 20,
    total: 0,
    showSizeChanger: true,
    showTotal: (t: number) => `共 ${t} 条`,
  })

  const extraParams = ref<Record<string, unknown>>({})
  const hasData = computed(() => items.value.length > 0)

  async function fetchData(params?: Record<string, unknown>) {
    loading.value = true
    try {
      const res = await fetcher({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...extraParams.value,
        ...params,
      })
      items.value = res.items
      total.value = res.total
      pagination.total = res.total
    }
    finally {
      loading.value = false
    }
  }

  function onTableChange(pag: { current?: number, pageSize?: number }) {
    if (pag.current) pagination.current = pag.current
    if (pag.pageSize) pagination.pageSize = pag.pageSize
    fetchData()
  }

  function search(params: Record<string, unknown>) {
    pagination.current = 1
    extraParams.value = params
    fetchData()
  }

  function refresh() {
    fetchData()
  }

  return { loading, items, total, pagination, hasData, fetchData, onTableChange, search, refresh }
}
