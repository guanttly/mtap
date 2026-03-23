# MTAP 前端代码规范

> 版本：1.0.0 | 适用范围：`web/src/` 目录下所有 Vue 3 + TypeScript 代码

---

## 1. 图标使用规范

### 1.1 图标系统架构

项目采用 **本地 SVG 图标注册表** 方案，所有图标数据存储在 `src/assets/icons/index.ts`，通过 `SvgIcon` 组件统一渲染。

```
src/
├── assets/
│   └── icons/
│       └── index.ts      ← SVG body 注册表（唯一真相来源）
└── components/
    └── common/
        └── SvgIcon.vue   ← 统一图标渲染组件
```

图标源数据来自本地安装的 `@iconify-json/ant-design` npm 包（Ant Design 图标集），**不依赖任何 CDN 或网络请求**。

---

### 1.2 使用 SvgIcon 组件

所有图标场景（菜单、按钮、标签、装饰）均使用 `<SvgIcon>` 组件：

```vue
<script setup lang="ts">
import SvgIcon from '@/components/common/SvgIcon.vue'
</script>

<template>
  <!-- 行内图标（默认 1em，继承颜色） -->
  <SvgIcon name="calendar-outlined" />

  <!-- 自定义尺寸（数字=px，字符串=CSS 值） -->
  <SvgIcon name="user-outlined" :size="20" />
  <SvgIcon name="user-outlined" size="1.5rem" />

  <!-- 自定义颜色 -->
  <SvgIcon name="check-circle-outlined" color="#52c41a" />

  <!-- 大型装饰图标 -->
  <SvgIcon name="mobile-outlined" :size="72" color="#1677ff" />
</template>
```

#### Props

| Prop    | 类型               | 默认值           | 说明                           |
|---------|--------------------|------------------|--------------------------------|
| `name`  | `IconName`         | 必填             | 图标名称，IDE 有完整类型提示   |
| `size`  | `number \| string` | `'1em'`          | 数字单位为 px，字符串原样输出  |
| `color` | `string`           | `'currentColor'` | 继承父级文字颜色               |
| `class` | `string`           | —                | 追加自定义 class               |

---

### 1.3 在渲染函数（h()）中使用图标

菜单项、Tree 等需要通过渲染函数传递 icon 的场景：

```ts
import type { IconName } from '@/assets/icons/index'
import { h } from 'vue'
import SvgIcon from '@/components/common/SvgIcon.vue'

// ✅ 推荐：定义辅助函数
function icon(name: IconName) {
  return () => h(SvgIcon, { name, size: '1em' })
}

const menuItems = [
  {
    key: '/dashboard',
    icon: icon('dashboard-outlined'),
    label: '数据看板',
  },
]
```

---

### 1.4 添加新图标

1. 确认图标在 Ant Design 图标集中存在（参考 [https://ant.design/components/icon](https://ant.design/components/icon)）
2. 从本地包提取 SVG body：

   ```bash
   node -e "
   const icons = require('./node_modules/@iconify-json/ant-design/icons.json');
   console.log(icons.icons['your-icon-name'].body);
   "
   ```

3. 在 `src/assets/icons/index.ts` 中添加：
   - 在 `IconName` 类型联合中追加 `| 'your-icon-name'`
   - 在 `iconBodies` 对象中追加对应条目

4. 提交 `src/assets/icons/index.ts` 的改动

> ⚠️ **禁止**直接从 CDN（如 `cdn.jsdelivr.net/npm/@iconify`）引用图标

---

### 1.5 新增图标集

若需要使用 Ant Design 以外的图标集：

```bash
# 安装图标集数据包
pnpm add -D @iconify-json/<集合名>

# 示例：Material Design 图标
pnpm add -D @iconify-json/mdi
```

在 `uno.config.ts` 中注册后，按同样流程提取 body 写入 `src/assets/icons/index.ts`，并在 `IconName` 中用前缀区分，如 `'mdi:home'`。

---

## 2. 禁止使用 Emoji 作为 UI 图标

### ❌ 禁止

```vue
<!-- 禁止：Emoji 不支持主题色、不可缩放矢量、可访问性差 -->
<a-radio value="morning">
☀ 上午
</a-radio>

<div style="font-size: 48px;">
📱
</div>
```

### ✅ 正确

```vue
<!-- 正确：使用 SvgIcon 组件 -->
<a-radio value="morning">
  <SvgIcon name="sun-outlined" style="margin-right: 4px" />上午
</a-radio>

<SvgIcon name="mobile-outlined" :size="72" color="#1677ff" />
```

**原因**：
- Emoji 渲染因操作系统/浏览器不同而显示差异巨大
- 无法跟随主题色变化
- 矢量缩放质量差
- 可访问性（a11y）支持不完整

---

## 3. 组件规范

### 3.1 文件结构

```vue
<!-- 核心目的：一句话说明该组件的核心目标 -->
<!-- 模块功能：具体功能要点列表 -->
<script setup lang="ts">
// 1. Vue 核心
// 2. 第三方库
// 3. 内部 API / Store
// 4. 内部组件
// 5. 类型导入
</script>

<template>
  <!-- 模板内容 -->
</template>

<style scoped>
/* 仅当组件有独立样式时添加 */
</style>
```

### 3.2 Props 定义

```ts
// ✅ 使用 withDefaults + defineProps 泛型
const props = withDefaults(defineProps<{
  title: string
  visible?: boolean
  size?: 'small' | 'large'
}>(), {
  visible: false,
  size: 'small',
})
```

---

## 4. 列表页面规范

所有数据列表页面必须使用 `list-card` 模式：

```vue
<template>
  <a-card class="list-card" :bordered="false">
    <template #title>
      页面标题
    </template>
    <template #extra>
      <!-- 操作按钮区 -->
      <a-button type="primary">
        新增
      </a-button>
    </template>

    <!-- 搜索栏 -->
    <div class="list-toolbar">
      <a-input-search ... />
    </div>

    <!-- 数据表格 -->
    <a-table :columns="columns" :data-source="data" />
  </a-card>
</template>
```

---

## 5. 样式规范

### 5.1 全局 CSS 类

| 类名            | 用途                   |
|-----------------|------------------------|
| `.list-card`    | 列表页卡片容器         |
| `.list-toolbar` | 搜索/过滤工具栏        |
| `.help-card`    | 表单页右侧帮助说明面板 |

### 5.2 禁止内联关键样式

```vue
<!-- ❌ 禁止用内联样式定义主色、间距系统 -->
<div style="color: #1677ff; margin: 24px;">
...
</div>

<!-- ✅ 使用 CSS 变量或全局类 -->
<div class="primary-text mb-6">
...
</div>
```

### 5.3 颜色使用

统一使用 Ant Design 设计变量：
- 主色：`#1677ff`（`--ant-color-primary`）
- 成功：`#52c41a`
- 警告：`#faad14`
- 错误：`#ff4d4f`
- 文字次要：`#8c8c8c`
- 背景：`#f0f2f5`

---

## 6. API 调用规范

```ts
// ✅ 标准 loading + try/catch/finally 模式
const loading = ref(false)

async function handleSubmit() {
  loading.value = true
  try {
    await someApi.create(form.value)
    message.success('操作成功')
  }
  catch {
    // 错误已由 request.ts 拦截器统一处理
  }
  finally {
    loading.value = false
  }
}
```

---

## 7. 路由与权限

- 所有路由 `meta` 必须包含 `requiresAuth: boolean`
- 受保护页面不得在 `beforeRouteEnter` 之前渲染敏感数据
- 角色权限由后端返回，前端仅做显示控制（不作安全边界）

---

## 8. TypeScript 规范

- 禁止使用 `any`，使用 `unknown` 替代不确定类型
- 所有 API 请求/响应类型定义在 `src/types/` 对应模块文件中
- 组件 Props 必须有完整类型定义，不使用 `defineProps` 数组语法

---

*最后更新：请在修改后更新此日期。*
