import { describe, it, expect } from 'vitest';
import { defineComponent } from 'vue';
import { mount } from '@vue/test-utils';
import { usePageProps, usePage, PagePropsKey, PageKey } from '../src/app';
import type { Page } from '../src/types';

describe('usePageProps', () => {
  it('returns injected page props', () => {
    const props = { user: { name: 'John' }, flash: { success: 'Saved!' } };

    const TestComponent = defineComponent({
      setup() {
        const pageProps = usePageProps<typeof props>();
        return { pageProps };
      },
      template: '<div>{{ pageProps.user.name }}</div>',
    });

    const wrapper = mount(TestComponent, {
      global: {
        provide: {
          [PagePropsKey as symbol]: props,
        },
      },
    });

    expect(wrapper.vm.pageProps).toEqual(props);
    expect(wrapper.text()).toBe('John');
  });

  it('throws error when props are not available', () => {
    const TestComponent = defineComponent({
      setup() {
        usePageProps();
        return {};
      },
      template: '<div>Test</div>',
    });

    expect(() => {
      mount(TestComponent);
    }).toThrow('Page props not found. Make sure you are using Inertia app.');
  });
});

describe('usePage', () => {
  it('returns injected page object', () => {
    const page: Page = {
      component: 'Users/Index',
      props: { users: [] },
      url: '/users',
      version: '1.0.0',
    };

    const TestComponent = defineComponent({
      setup() {
        const currentPage = usePage();
        return { currentPage };
      },
      template: '<div>{{ currentPage.component }}</div>',
    });

    const wrapper = mount(TestComponent, {
      global: {
        provide: {
          [PageKey as symbol]: page,
        },
      },
    });

    expect(wrapper.vm.currentPage).toEqual(page);
    expect(wrapper.text()).toBe('Users/Index');
  });

  it('throws error when page is not available', () => {
    const TestComponent = defineComponent({
      setup() {
        usePage();
        return {};
      },
      template: '<div>Test</div>',
    });

    expect(() => {
      mount(TestComponent);
    }).toThrow('Page not found. Make sure you are using Inertia app.');
  });

  it('allows accessing typed page props', () => {
    interface UserPageProps {
      users: Array<{ id: number; name: string }>;
      meta: { total: number };
    }

    const page: Page<UserPageProps> = {
      component: 'Users/Index',
      props: {
        users: [
          { id: 1, name: 'John' },
          { id: 2, name: 'Jane' },
        ],
        meta: { total: 2 },
      },
      url: '/users',
      version: '1.0.0',
    };

    const TestComponent = defineComponent({
      setup() {
        const currentPage = usePage<UserPageProps>();
        return { currentPage };
      },
      template: '<div>{{ currentPage.props.users.length }}</div>',
    });

    const wrapper = mount(TestComponent, {
      global: {
        provide: {
          [PageKey as symbol]: page,
        },
      },
    });

    expect(wrapper.vm.currentPage.props.users).toHaveLength(2);
    expect(wrapper.text()).toBe('2');
  });
});
