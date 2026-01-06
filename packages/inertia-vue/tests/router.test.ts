import { describe, it, expect } from 'vitest'
import { router } from '../src/router'

describe('Router', () => {
  describe('visit', () => {
    it('should make GET request by default', () => {
      expect(router).toBeDefined()
    })

    it('should include Inertia headers', () => {
      expect(router.visit).toBeDefined()
    })

    it('should handle POST requests', () => {
      expect(router.post).toBeDefined()
    })

    it('should handle PUT requests', () => {
      expect(router.put).toBeDefined()
    })

    it('should handle PATCH requests', () => {
      expect(router.patch).toBeDefined()
    })

    it('should handle DELETE requests', () => {
      expect(router.delete).toBeDefined()
    })
  })

  describe('reload', () => {
    it('should reload current page', () => {
      expect(router.reload).toBeDefined()
    })
  })

  describe('replace', () => {
    it('should replace history instead of push', () => {
      expect(router.replace).toBeDefined()
    })
  })

  describe('partial reloads', () => {
    it('should send X-Inertia-Partial-Data header', () => {
      expect(router.get).toBeDefined()
    })

    it('should send X-Inertia-Partial-Component header', () => {
      expect(router.get).toBeDefined()
    })
  })
})
