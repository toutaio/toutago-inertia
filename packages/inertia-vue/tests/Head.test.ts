import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Head from '../src/components/Head.vue'

describe('Head', () => {
  it('should render nothing in template', () => {
    const wrapper = mount(Head, {
      props: { title: 'Test Page' }
    })
    
    expect(wrapper.html()).toContain('<!-- Nothing rendered -->')
  })

  it('should accept title prop', () => {
    const wrapper = mount(Head, {
      props: { title: 'Dashboard' }
    })
    
    expect(wrapper.props('title')).toBe('Dashboard')
  })

  it('should accept slots for meta tags', () => {
    const wrapper = mount(Head, {
      slots: {
        default: () => '<meta name="description" content="Test">'
      }
    })
    
    expect(wrapper.html()).toContain('<!-- Nothing rendered -->')
  })
})
