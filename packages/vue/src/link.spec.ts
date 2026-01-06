import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Link } from './link'
import { router } from './router'

vi.mock('./router', () => ({
  router: {
    visit: vi.fn(),
  },
}))

describe('Link Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders a link with href', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
      slots: {
        default: 'Users',
      },
    })

    expect(wrapper.find('a').exists()).toBe(true)
    expect(wrapper.find('a').attributes('href')).toBe('/users')
    expect(wrapper.text()).toBe('Users')
  })

  it('prevents default and uses Inertia router on click', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
    })

    await wrapper.find('a').trigger('click')

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      method: 'get',
      preserveScroll: false,
      preserveState: false,
    }))
  })

  it('uses custom method', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users/1',
        method: 'delete',
      },
    })

    await wrapper.find('a').trigger('click')

    expect(router.visit).toHaveBeenCalledWith('/users/1', expect.objectContaining({
      method: 'delete',
    }))
  })

  it('passes data payload', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        method: 'post',
        data: { name: 'John' },
      },
    })

    await wrapper.find('a').trigger('click')

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      method: 'post',
      data: { name: 'John' },
    }))
  })

  it('preserves scroll when specified', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        preserveScroll: true,
      },
    })

    await wrapper.find('a').trigger('click')

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      preserveScroll: true,
    }))
  })

  it('preserves state when specified', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        preserveState: true,
      },
    })

    await wrapper.find('a').trigger('click')

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      preserveState: true,
    }))
  })

  it('allows normal navigation with ctrl/cmd/shift click', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
    })

    await wrapper.find('a').trigger('click', { ctrlKey: true })
    expect(router.visit).not.toHaveBeenCalled()

    await wrapper.find('a').trigger('click', { metaKey: true })
    expect(router.visit).not.toHaveBeenCalled()

    await wrapper.find('a').trigger('click', { shiftKey: true })
    expect(router.visit).not.toHaveBeenCalled()
  })

  it('allows normal navigation with middle mouse button', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
    })

    await wrapper.find('a').trigger('click', { button: 1 })
    expect(router.visit).not.toHaveBeenCalled()
  })

  it('applies custom class', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        class: 'btn btn-primary',
      },
    })

    expect(wrapper.find('a').classes()).toContain('btn')
    expect(wrapper.find('a').classes()).toContain('btn-primary')
  })
})
