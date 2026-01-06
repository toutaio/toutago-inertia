import { describe, it, expect } from 'vitest'
import { createInertiaApp } from '../src/app'

describe('createInertiaApp', () => {
  it('should be defined', () => {
    expect(createInertiaApp).toBeDefined()
    expect(typeof createInertiaApp).toBe('function')
  })

  it('should be a function', () => {
    expect(typeof createInertiaApp).toBe('function')
  })

  it('should accept config parameter', () => {
    // Just check the function signature
    expect(createInertiaApp.length).toBeGreaterThan(0)
  })

  it('should return a promise', () => {
    // Mock DOM element for testing
    const mockEl = document.createElement('div')
    mockEl.id = 'app'
    mockEl.dataset.page = JSON.stringify({
      component: 'Test',
      props: {},
      url: '/test',
      version: '1'
    })
    document.body.appendChild(mockEl)

    const config = {
      resolve: () => ({ name: 'Test' }),
      setup: () => {}
    }
    
    const result = createInertiaApp(config)
    expect(result).toBeInstanceOf(Promise)
    
    document.body.removeChild(mockEl)
  })
})
