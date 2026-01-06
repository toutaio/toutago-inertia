import { describe, it, expect } from 'vitest'
import { usePage } from '../src/composables/usePage'

describe('usePage', () => {
  it('should return current page data', () => {
    const page = usePage()
    expect(page).toBeDefined()
  })

  it('should be reactive', () => {
    const page = usePage()
    expect(page.value).toBeDefined()
  })

  it('should contain component', () => {
    const page = usePage()
    expect(page.value).toHaveProperty('component')
  })

  it('should contain props', () => {
    const page = usePage()
    expect(page.value).toHaveProperty('props')
  })

  it('should contain url', () => {
    const page = usePage()
    expect(page.value).toHaveProperty('url')
  })

  it('should contain version', () => {
    const page = usePage()
    expect(page.value).toHaveProperty('version')
  })

  it('should access nested props', () => {
    const page = usePage()
    expect(page.value.props).toBeDefined()
  })
})
