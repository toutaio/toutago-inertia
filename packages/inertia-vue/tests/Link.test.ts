import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Link from '../src/components/Link.vue'

describe('Link', () => {
  it('should render anchor tag', () => {
    const wrapper = mount(Link, {
      props: { href: '/test' },
      slots: { default: 'Test Link' }
    })
    
    expect(wrapper.find('a').exists()).toBe(true)
    expect(wrapper.text()).toBe('Test Link')
  })

  it('should set href attribute', () => {
    const wrapper = mount(Link, {
      props: { href: '/users' }
    })
    
    expect(wrapper.find('a').attributes('href')).toBe('/users')
  })

  it('should handle click events', async () => {
    const wrapper = mount(Link, {
      props: { href: '/test' }
    })
    
    await wrapper.find('a').trigger('click')
    expect(wrapper.emitted()).toHaveProperty('click')
  })

  it('should support method prop', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/posts',
        method: 'post'
      }
    })
    
    expect(wrapper.props('method')).toBe('post')
  })

  it('should support data prop', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/posts',
        data: { title: 'Test' }
      }
    })
    
    expect(wrapper.props('data')).toEqual({ title: 'Test' })
  })

  it('should support preserveScroll', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/test',
        preserveScroll: true
      }
    })
    
    expect(wrapper.props('preserveScroll')).toBe(true)
  })

  it('should support preserveState', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/test',
        preserveState: true
      }
    })
    
    expect(wrapper.props('preserveState')).toBe(true)
  })

  it('should support only prop for partial reloads', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/test',
        only: ['users', 'posts']
      }
    })
    
    expect(wrapper.props('only')).toEqual(['users', 'posts'])
  })

  it('should support custom classes', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/test',
        class: 'btn btn-primary'
      }
    })
    
    expect(wrapper.classes()).toContain('btn')
    expect(wrapper.classes()).toContain('btn-primary')
  })
})
