import { describe, it, expect, beforeEach } from 'vitest'
import { useRemember } from './useRemember'
import { nextTick } from 'vue'

describe('useRemember', () => {
  beforeEach(() => {
    sessionStorage.clear()
    localStorage.clear()
  })

  it('should initialize with default value when no stored value exists', () => {
    const state = useRemember('default value', 'test-key')
    expect(state.value).toBe('default value')
  })

  it('should restore value from sessionStorage', () => {
    sessionStorage.setItem('inertia:remember:test-key', JSON.stringify('stored value'))
    const state = useRemember('default', 'test-key')
    expect(state.value).toBe('stored value')
  })

  it('should update sessionStorage when value changes', async () => {
    const state = useRemember('initial', 'test-key')
    state.value = 'updated'
    await nextTick()
    
    const stored = sessionStorage.getItem('inertia:remember:test-key')
    expect(JSON.parse(stored!)).toBe('updated')
  })

  it('should work with objects', async () => {
    const initial = { name: 'John', age: 30 }
    const state = useRemember(initial, 'user-data')
    
    expect(state.value).toEqual(initial)
    
    state.value = { name: 'Jane', age: 25 }
    await nextTick()
    
    const stored = sessionStorage.getItem('inertia:remember:user-data')
    expect(JSON.parse(stored!)).toEqual({ name: 'Jane', age: 25 })
  })

  it('should work with arrays', async () => {
    const initial = [1, 2, 3]
    const state = useRemember(initial, 'numbers')
    
    state.value = [4, 5, 6]
    await nextTick()
    
    const stored = sessionStorage.getItem('inertia:remember:numbers')
    expect(JSON.parse(stored!)).toEqual([4, 5, 6])
  })

  it('should handle null values', () => {
    const state = useRemember(null, 'nullable')
    expect(state.value).toBe(null)
    
    state.value = 'not null'
    expect(state.value).toBe('not null')
  })

  it('should use localStorage when storage type is local', async () => {
    const state = useRemember('value', 'local-key', 'local')
    expect(state.value).toBe('value')
    await nextTick()
    
    const stored = localStorage.getItem('inertia:remember:local-key')
    expect(JSON.parse(stored!)).toBe('value')
  })

  it('should restore from localStorage', () => {
    localStorage.setItem('inertia:remember:local-key', JSON.stringify('stored'))
    const state = useRemember('default', 'local-key', 'local')
    expect(state.value).toBe('stored')
  })

  it('should clear storage when set to undefined', async () => {
    const state = useRemember('value', 'clear-key')
    await nextTick()
    sessionStorage.setItem('inertia:remember:clear-key', JSON.stringify('value'))
    
    state.value = undefined
    await nextTick()
    
    expect(sessionStorage.getItem('inertia:remember:clear-key')).toBe(null)
  })

  it('should handle JSON parse errors gracefully', () => {
    sessionStorage.setItem('inertia:remember:bad-json', 'not valid json')
    const state = useRemember('default', 'bad-json')
    expect(state.value).toBe('default')
  })

  it('should work with reactive updates', async () => {
    const state = useRemember({ count: 0 }, 'counter')
    
    state.value.count++
    state.value = { ...state.value } // Trigger reactivity
    await nextTick()
    
    const stored = sessionStorage.getItem('inertia:remember:counter')
    expect(JSON.parse(stored!).count).toBe(1)
  })

  it('should maintain separate keys', async () => {
    const state1 = useRemember('value1', 'key1')
    const state2 = useRemember('value2', 'key2')
    
    expect(state1.value).toBe('value1')
    expect(state2.value).toBe('value2')
    
    state1.value = 'updated1'
    await nextTick()
    
    expect(state2.value).toBe('value2') // Should not be affected
  })
})
