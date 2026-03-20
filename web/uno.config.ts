import {
  defineConfig,
  presetAttributify,
  presetIcons,
  presetUno,
  transformerAttributifyJsx,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      scale: 1.2,
      warn: true,
    }),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
    transformerAttributifyJsx(),
  ],
  shortcuts: {
    'btn-primary': 'px-4 py-2 rounded bg-blue-500 text-white hover:bg-blue-600 cursor-pointer transition',
    'btn-danger': 'px-4 py-2 rounded bg-red-500 text-white hover:bg-red-600 cursor-pointer transition',
    'card': 'bg-white rounded-lg shadow p-4',
    'flex-center': 'flex items-center justify-center',
  },
})
