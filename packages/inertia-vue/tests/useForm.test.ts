import { describe, it, expect } from 'vitest'
import { useForm } from '../src/composables/useForm'

describe('useForm', () => {
  it('should create form with initial data', () => {
    const form = useForm({ name: 'John', email: 'john@example.com' })
    
    expect(form.data.name).toBe('John')
    expect(form.data.email).toBe('john@example.com')
  })

  it('should track processing state', () => {
    const form = useForm({ name: '' })
    expect(form.processing).toBe(false)
  })

  it('should track errors', () => {
    const form = useForm({ name: '' })
    expect(form.errors).toEqual({})
    expect(form.hasErrors).toBe(false)
  })

  it('should track success state', () => {
    const form = useForm({ name: '' })
    expect(form.wasSuccessful).toBe(false)
  })

  it('should track dirty state', async () => {
    const form = useForm({ name: 'John' })
    expect(form.isDirty).toBe(false)
    
    form.data.name = 'Jane'
    
    // Wait for Vue reactivity
    await new Promise(resolve => setTimeout(resolve, 10))
    expect(form.isDirty).toBe(true)
  })

  it('should set errors', () => {
    const form = useForm({ name: '' })
    form.setError('name', 'Name is required')
    
    expect(form.errors.name).toBe('Name is required')
    expect(form.hasErrors).toBe(true)
  })

  it('should clear errors', () => {
    const form = useForm({ name: '' })
    form.setError('name', 'Error')
    form.clearErrors('name')
    
    expect(form.errors.name).toBeUndefined()
    expect(form.hasErrors).toBe(false)
  })

  it('should reset fields', () => {
    const form = useForm({ name: 'John', email: 'john@example.com' })
    form.data.name = 'Jane'
    form.data.email = 'jane@example.com'
    
    form.reset('name')
    expect(form.data.name).toBe('John')
    expect(form.data.email).toBe('jane@example.com')
  })

  it('should reset all fields when no args', () => {
    const form = useForm({ name: 'John', email: 'john@example.com' })
    form.data.name = 'Jane'
    form.data.email = 'jane@example.com'
    
    form.reset()
    expect(form.data.name).toBe('John')
    expect(form.data.email).toBe('john@example.com')
  })

  it('should have submit methods', () => {
    const form = useForm({ name: '' })
    
    expect(typeof form.get).toBe('function')
    expect(typeof form.post).toBe('function')
    expect(typeof form.put).toBe('function')
    expect(typeof form.patch).toBe('function')
    expect(typeof form.delete).toBe('function')
  })

  it('should support transform', () => {
    const form = useForm({ name: 'john' })
    
    form.transform(data => ({
      name: data.name.toUpperCase()
    }))
    
    expect(form).toBeDefined()
  })
})
