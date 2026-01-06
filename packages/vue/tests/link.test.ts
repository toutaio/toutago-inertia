import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { Link } from '../src/link';
import { router } from '../src/router';

vi.mock('../src/router', () => ({
  router: {
    visit: vi.fn(),
  },
}));

describe('Link Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders as anchor tag by default', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
      slots: {
        default: 'View Users',
      },
    });

    expect(wrapper.element.tagName).toBe('A');
    expect(wrapper.attributes('href')).toBe('/users');
    expect(wrapper.text()).toBe('View Users');
  });

  it('renders as custom element when "as" prop is provided', () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        as: 'button',
      },
      slots: {
        default: 'View Users',
      },
    });

    expect(wrapper.element.tagName).toBe('BUTTON');
    expect(wrapper.text()).toBe('View Users');
  });

  it('intercepts click and calls router.visit', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
      slots: {
        default: 'View Users',
      },
    });

    await wrapper.trigger('click');

    expect(router.visit).toHaveBeenCalledWith('/users', {
      method: 'get',
      data: {},
      replace: false,
      preserveScroll: false,
      preserveState: false,
      only: undefined,
      headers: {},
    });
  });

  it('passes method prop to router', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users/1',
        method: 'delete',
      },
      slots: {
        default: 'Delete',
      },
    });

    await wrapper.trigger('click');

    expect(router.visit).toHaveBeenCalledWith('/users/1', expect.objectContaining({
      method: 'delete',
    }));
  });

  it('passes data prop to router', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        method: 'post',
        data: { name: 'John', email: 'john@example.com' },
      },
      slots: {
        default: 'Create User',
      },
    });

    await wrapper.trigger('click');

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      data: { name: 'John', email: 'john@example.com' },
    }));
  });

  it('passes preserve options to router', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        preserveScroll: true,
        preserveState: true,
      },
      slots: {
        default: 'View Users',
      },
    });

    await wrapper.trigger('click');

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      preserveScroll: true,
      preserveState: true,
    }));
  });

  it('passes only prop to router when provided', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
        only: ['users', 'meta'],
      },
      slots: {
        default: 'View Users',
      },
    });

    await wrapper.trigger('click');

    expect(router.visit).toHaveBeenCalledWith('/users', expect.objectContaining({
      only: ['users', 'meta'],
    }));
  });

  it('does not intercept external links', async () => {
    const wrapper = mount(Link, {
      props: {
        href: 'https://example.com',
      },
      slots: {
        default: 'External',
      },
    });

    await wrapper.trigger('click');

    expect(router.visit).not.toHaveBeenCalled();
  });

  it('does not intercept when default is prevented', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
      slots: {
        default: 'View Users',
      },
    });

    const event = new MouseEvent('click', { cancelable: true });
    event.preventDefault();
    await wrapper.element.dispatchEvent(event);

    expect(router.visit).not.toHaveBeenCalled();
  });

  it('does not intercept when modifier keys are pressed', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
      slots: {
        default: 'View Users',
      },
    });

    // Test ctrl key
    await wrapper.trigger('click', { ctrlKey: true });
    expect(router.visit).not.toHaveBeenCalled();

    // Test meta key
    await wrapper.trigger('click', { metaKey: true });
    expect(router.visit).not.toHaveBeenCalled();

    // Test shift key
    await wrapper.trigger('click', { shiftKey: true });
    expect(router.visit).not.toHaveBeenCalled();

    // Test alt key
    await wrapper.trigger('click', { altKey: true });
    expect(router.visit).not.toHaveBeenCalled();
  });

  it('does not intercept non-left-button clicks', async () => {
    const wrapper = mount(Link, {
      props: {
        href: '/users',
      },
      slots: {
        default: 'View Users',
      },
    });

    // Right click (button 2)
    await wrapper.trigger('click', { button: 2 });
    expect(router.visit).not.toHaveBeenCalled();

    // Middle click (button 1)
    await wrapper.trigger('click', { button: 1 });
    expect(router.visit).not.toHaveBeenCalled();
  });
});
