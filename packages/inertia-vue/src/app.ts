import { computed, defineComponent, h, markRaw, reactive, type App, type Component } from 'vue'
import { router } from './router'
import { pageSymbol } from './composables/usePage'
import type { InertiaAppProps, Page } from './types'

export async function createInertiaApp(config: InertiaAppProps): Promise<void> {
  const { resolve, setup, title } = config
  
  // Get initial page from data attribute
  const el = document.getElementById('app')
  if (!el) {
    throw new Error('Element #app not found')
  }

  const initialPageData = el.dataset.page
  if (!initialPageData) {
    throw new Error('Inertia page data not found')
  }

  const initialPage: Page = JSON.parse(initialPageData)
  router.setPage(initialPage)

  // Create reactive page
  const page = reactive({ ...initialPage })
  const pageRef = computed(() => page)

  // Resolve component
  let component: Component
  if (config.initialComponent) {
    component = config.initialComponent
  } else if (resolve) {
    component = await Promise.resolve(resolve(initialPage.component))
  } else {
    throw new Error('Either initialComponent or resolve must be provided')
  }

  // Create app component  
  const AppComponent = defineComponent({
    name: 'InertiaApp',
    setup() {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      return () => h(markRaw(component as any), page.props)
    }
  })

  // Plugin
  const plugin = {
    install(app: App) {
      app.provide(pageSymbol, pageRef)
      
      // Global properties
      app.config.globalProperties.$page = pageRef
      app.config.globalProperties.$inertia = router
    }
  }

  // Call setup
  setup({
    el,
    App: AppComponent,
    props: {
      initialPage: page,
      resolveComponent: resolve
    },
    plugin
  })

  // Listen for popstate (browser back/forward)
  window.addEventListener('popstate', (event) => {
    if (event.state) {
      Object.assign(page, event.state)
      router.setPage(event.state)
    }
  })

  // Update title if callback provided
  if (title && page.props.title) {
    document.title = title(page.props.title as string)
  }
}

export { router }
