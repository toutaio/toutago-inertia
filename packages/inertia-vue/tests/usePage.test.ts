import { describe, it, expect, vi } from 'vitest'
import { createApp, h } from 'vue'
import { usePage } from '../src/composables/usePage'

describe('usePage', () => {
  it('should return current page data', () => {
    // usePage needs to be called within a Vue component context
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page).toBeDefined()
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })

  it('should be reactive', () => {
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page.value).toBeDefined()
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })

  it('should contain component', () => {
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page.value).toHaveProperty('component')
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })

  it('should contain props', () => {
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page.value).toHaveProperty('props')
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })

  it('should contain url', () => {
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page.value).toHaveProperty('url')
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })

  it('should contain version', () => {
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page.value).toHaveProperty('version')
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })

  it('should access nested props', () => {
    const app = createApp({
      setup() {
        const page = usePage()
        expect(page.value.props).toBeDefined()
        return () => h('div')
      }
    })
    const el = document.createElement('div')
    app.mount(el)
    app.unmount()
  })
})
