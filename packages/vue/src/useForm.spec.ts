import { describe, it, expect, vi } from 'vitest'
import { useForm } from './useForm'
import { router } from './router'

vi.mock('./router', () => ({
  router: {
    visit: vi.fn(),
  },
}))

describe('useForm', () => {
  it('initializes form with data', () => {
    const form = useForm({
      name: 'John',
      email: 'john@example.com',
    })

    expect(form.data.name).toBe('John')
    expect(form.data.email).toBe('john@example.com')
  })

  it('tracks dirty state', () => {
    const form = useForm({
      name: 'John',
    })

    expect(form.isDirty).toBe(false)

    form.data.name = 'Jane'
    expect(form.isDirty).toBe(true)
  })

  it('resets form data', () => {
    const form = useForm({
      name: 'John',
      email: 'john@example.com',
    })

    form.data.name = 'Jane'
    form.data.email = 'jane@example.com'

    form.reset()

    expect(form.data.name).toBe('John')
    expect(form.data.email).toBe('john@example.com')
    expect(form.isDirty).toBe(false)
  })

  it('resets specific fields', () => {
    const form = useForm({
      name: 'John',
      email: 'john@example.com',
    })

    form.data.name = 'Jane'
    form.data.email = 'jane@example.com'

    form.reset('name')

    expect(form.data.name).toBe('John')
    expect(form.data.email).toBe('jane@example.com')
  })

  it('clears errors', () => {
    const form = useForm({
      name: 'John',
    })

    form.errors.name = 'Name is required'
    form.errors.email = 'Email is invalid'

    form.clearErrors()

    expect(form.errors.name).toBeUndefined()
    expect(form.errors.email).toBeUndefined()
  })

  it('clears specific error', () => {
    const form = useForm({
      name: 'John',
    })

    form.errors.name = 'Name is required'
    form.errors.email = 'Email is invalid'

    form.clearErrors('name')

    expect(form.errors.name).toBeUndefined()
    expect(form.errors.email).toBe('Email is invalid')
  })

  it('submits with get method', () => {
    const form = useForm({
      search: 'test',
    })

    form.get('/search')

    expect(router.visit).toHaveBeenCalledWith('/search', expect.objectContaining({
      method: 'get',
      data: { search: 'test' },
    }))
  })

  it('submits with post method', () => {
    const form = useForm({
      name: 'John',
    })

    form.post('/users')

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      method: 'post',
      data: { name: 'John' },
    }))
  })

  it('submits with put method', () => {
    const form = useForm({
      name: 'Jane',
    })

    form.put('/users/1')

    expect(router.visit).toHaveBeenCalledWith('/users/1', expect.objectContaining({
      method: 'put',
      data: { name: 'Jane' },
    }))
  })

  it('submits with patch method', () => {
    const form = useForm({
      name: 'Jane',
    })

    form.patch('/users/1')

    expect(router.visit).toHaveBeenCalledWith('/users/1', expect.objectContaining({
      method: 'patch',
      data: { name: 'Jane' },
    }))
  })

  it('submits with delete method', () => {
    const form = useForm({
      id: 1,
    })

    form.delete('/users/1')

    expect(router.visit).toHaveBeenCalledWith('/users/1', expect.objectContaining({
      method: 'delete',
    }))
  })

  it('tracks processing state', () => {
    const form = useForm({
      name: 'John',
    })

    expect(form.processing).toBe(false)

    form.post('/users')

    expect(form.processing).toBe(true)
  })

  it('sets recently successful state', () => {
    const form = useForm({
      name: 'John',
    })

    expect(form.recentlySuccessful).toBe(false)
  })

  it('preserves scroll by default', () => {
    const form = useForm({
      name: 'John',
    })

    form.post('/users')

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      preserveScroll: true,
    }))
  })

  it('accepts custom options', () => {
    const form = useForm({
      name: 'John',
    })

    const onSuccess = vi.fn()
    const onError = vi.fn()

    form.post('/users', {
      preserveScroll: false,
      onSuccess,
      onError,
    })

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      preserveScroll: false,
      onSuccess: expect.any(Function),
      onError: expect.any(Function),
    }))
  })
})
