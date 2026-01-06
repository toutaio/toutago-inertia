import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import { h, defineComponent } from 'vue';
import { withLayout } from '../src/utils/layout';

describe('Layout Utils', () => {
  it('should wrap component with layout', () => {
    const Layout = defineComponent({
      template: '<div class="layout"><slot /></div>'
    });

    const Page = defineComponent({
      template: '<div>Page Content</div>'
    });

    const WrappedPage = withLayout(Page, Layout);
    const wrapper = mount(WrappedPage);

    expect(wrapper.find('.layout').exists()).toBe(true);
    expect(wrapper.text()).toContain('Page Content');
  });

  it('should pass props to page component', () => {
    const Layout = defineComponent({
      template: '<div class="layout"><slot /></div>'
    });

    const Page = defineComponent({
      props: ['title'],
      template: '<div>{{ title }}</div>'
    });

    const WrappedPage = withLayout(Page, Layout);
    const wrapper = mount(WrappedPage, {
      props: { title: 'Test Title' }
    });

    expect(wrapper.text()).toContain('Test Title');
  });

  it('should support layout as function', () => {
    const Layout = defineComponent({
      template: '<div class="layout"><slot /></div>'
    });

    const Page = defineComponent({
      template: '<div>Page Content</div>',
      layout: () => Layout
    });

    // Access the layout property
    expect(typeof Page.layout).toBe('function');
    expect(Page.layout()).toBe(Layout);
  });

  it('should support nested layouts', () => {
    const OuterLayout = defineComponent({
      template: '<div class="outer"><slot /></div>'
    });

    const InnerLayout = defineComponent({
      template: '<div class="inner"><slot /></div>'
    });

    const Page = defineComponent({
      template: '<div>Page Content</div>'
    });

    const WrappedInner = withLayout(Page, InnerLayout);
    const WrappedOuter = withLayout(WrappedInner, OuterLayout);
    const wrapper = mount(WrappedOuter);

    expect(wrapper.find('.outer').exists()).toBe(true);
    expect(wrapper.find('.inner').exists()).toBe(true);
    expect(wrapper.text()).toContain('Page Content');
  });

  it('should handle component without layout', () => {
    const Page = defineComponent({
      template: '<div>Page Content</div>'
    });

    const wrapper = mount(Page);
    expect(wrapper.text()).toBe('Page Content');
  });
});
