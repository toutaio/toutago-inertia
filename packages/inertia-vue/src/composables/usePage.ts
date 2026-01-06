import { computed, inject, type ComputedRef } from 'vue'
import type { Page, PageProps } from '../types'

const pageSymbol = Symbol('inertia-page')

export function usePage<T extends PageProps = PageProps>(): ComputedRef<Page<T>> {
  const page = inject<ComputedRef<Page<T>>>(pageSymbol, undefined, true) // silent in production
  
  if (!page) {
    // Return default page for testing
    return computed(() => ({
      component: '',
      props: {} as T,
      url: '',
      version: null
    }))
  }
  
  return page
}

export function providePage(): void {
  // This will be used by the plugin
}

export { pageSymbol }
