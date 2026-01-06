import { createApp, h } from 'vue'
import { createInertiaApp, withLayout, resolvePageLayout } from '@toutaio/inertia-vue'
import AppLayout from './layouts/App.vue'

createInertiaApp({
  resolve: async (name) => {
    const pages = import.meta.glob('./pages/**/*.vue')
    const page = await pages[`./pages/${name}.vue`]()
    
    // Get the page's custom layout (if any)
    const pageLayout = resolvePageLayout(page.default)
    
    // If page has a custom layout, nest it inside AppLayout
    if (pageLayout) {
      page.default = withLayout(page.default, pageLayout)
    }
    
    // Always wrap with AppLayout (outer layout)
    page.default = withLayout(page.default, AppLayout)
    
    return page
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el)
  },
})
