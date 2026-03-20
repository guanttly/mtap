import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ruleApi } from '@/api/rule'
import type { ConflictRule, ConflictPackage, DependencyRule, PriorityTag } from '@/types/rule'

export const useRuleStore = defineStore('rule', () => {
  const conflictRules = ref<ConflictRule[]>([])
  const conflictPackages = ref<ConflictPackage[]>([])
  const dependencyRules = ref<DependencyRule[]>([])
  const priorityTags = ref<PriorityTag[]>([])
  const loading = ref(false)

  async function fetchConflictRules(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const res = await ruleApi.listConflictRules({ page, page_size: pageSize })
      conflictRules.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  async function fetchConflictPackages(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const res = await ruleApi.listConflictPackages({ page, page_size: pageSize })
      conflictPackages.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  async function fetchDependencyRules(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const res = await ruleApi.listDependencyRules({ page, page_size: pageSize })
      dependencyRules.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  async function fetchPriorityTags() {
    loading.value = true
    try {
      const res = await ruleApi.listPriorityTags()
      priorityTags.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  return { conflictRules, conflictPackages, dependencyRules, priorityTags, loading, fetchConflictRules, fetchConflictPackages, fetchDependencyRules, fetchPriorityTags }
})
