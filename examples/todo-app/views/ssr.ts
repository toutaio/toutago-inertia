import { createSSRApp, h } from 'vue'
import { renderToString } from 'vue/server-renderer'
import { createInertiaApp } from '@toutaio/inertia-vue'
import AppLayout from './layouts/App.vue'

export async function render(page) {
  return await createInertiaApp({
    page,
    resolve: async (name) => {
      const pages = import.meta.glob('./pages/**/*.vue', { eager: true })
      const pageModule = pages[`./pages/${name}.vue`]
      pageModule.default.layout = pageModule.default.layout || AppLayout
      return pageModule
    },
    setup({ App, props, plugin }) {
      return createSSRApp({ render: () => h(App, props) }).use(plugin)
    },
    render: renderToString,
  })
}
