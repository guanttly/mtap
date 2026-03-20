<!-- 核心目的：登录页面 -->
<!-- 模块功能：用户名/密码登录、JWT Token获取 -->
<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { message } from 'ant-design-vue'

const { login } = useAuth()
const loading = ref(false)

const form = reactive({ username: '', password: '' })

async function handleLogin() {
  if (!form.username || !form.password) {
    message.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    await login(form)
  }
  catch (e: unknown) {
    const err = e as { message?: string }
    message.error(err?.message ?? '登录失败，请检查用户名和密码')
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);">
    <div style="width: 380px;">
      <div style="text-align: center; margin-bottom: 32px; color: #fff;">
        <div style="font-size: 32px; font-weight: 800; letter-spacing: 2px; margin-bottom: 8px;">MTAP</div>
        <div style="font-size: 14px; opacity: .7;">医技预约管理平台</div>
      </div>
      <a-card :bodyStyle="{ padding: '32px' }" style="border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,.3);">
        <a-form :model="form" layout="vertical" @finish="handleLogin">
          <a-form-item label="用户名" name="username" :rules="[{ required: true, message: '请输入用户名' }]">
            <a-input v-model:value="form.username" placeholder="请输入用户名" size="large" allow-clear>
              <template #prefix><span class="i-ant-design:user-outlined" style="color: #bfbfbf;" /></template>
            </a-input>
          </a-form-item>
          <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
            <a-input-password v-model:value="form.password" placeholder="请输入密码" size="large">
              <template #prefix><span class="i-ant-design:lock-outlined" style="color: #bfbfbf;" /></template>
            </a-input-password>
          </a-form-item>
          <a-button type="primary" html-type="submit" size="large" block :loading="loading" style="margin-top: 8px; border-radius: 8px; height: 44px; font-size: 16px;">
            登 录
          </a-button>
        </a-form>
        <div style="text-align: center; margin-top: 16px; color: #8c8c8c; font-size: 12px;">默认管理员：admin / Admin@1234</div>
      </a-card>
    </div>
  </div>
</template>
