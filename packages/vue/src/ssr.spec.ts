import { describe, it, expect } from 'vitest';
import { h, defineComponent } from 'vue';
import { createInertiaSSRApp, createSSRPage } from './ssr';
import type { Page } from './types';

describe('SSR Support', () => {
  describe('createInertiaSSRApp', () => {
    it('should render a simple component to HTML', async () => {
      const TestComponent = defineComponent({
        props: ['title', 'initialPage', 'initialComponent'],
        setup(props) {
          return () => h('div', { class: 'test' }, props.title);
        },
      });

      const page: Page = {
        component: 'Test',
        props: { title: 'Hello SSR' },
        url: '/test',
        version: '1',
      };

      const result = await createInertiaSSRApp({
        page,
        resolveComponent: (name) => {
          if (name === 'Test') return TestComponent;
          throw new Error(`Component ${name} not found`);
        },
      });

      expect(result.body).toContain('<div class="test">Hello SSR</div>');
    });

    it('should include version in head meta tags', async () => {
      const TestComponent = defineComponent({
        setup() {
          return () => h('div', 'Test');
        },
      });

      const page: Page = {
        component: 'Test',
        props: {},
        url: '/test',
        version: 'abc123',
      };

      const result = await createInertiaSSRApp({
        page,
        resolveComponent: () => TestComponent,
      });

      expect(result.head).toContain('<meta name="inertia-version" content="abc123">');
    });

    it('should throw error when component not found', async () => {
      const page: Page = {
        component: 'NonExistent',
        props: {},
        url: '/test',
        version: '1',
      };

      await expect(
        createInertiaSSRApp({
          page,
          resolveComponent: () => null as any,
        })
      ).rejects.toThrow('Component "NonExistent" not found');
    });

    it('should pass page props to component', async () => {
      let receivedProps: any = null;

      const TestComponent = defineComponent({
        props: ['user', 'settings', 'initialPage', 'initialComponent'],
        setup(props) {
          receivedProps = props;
          return () => h('div', JSON.stringify({ user: props.user, settings: props.settings }));
        },
      });

      const page: Page = {
        component: 'Test',
        props: {
          user: { name: 'John' },
          settings: { theme: 'dark' },
        },
        url: '/test',
        version: '1',
      };

      await createInertiaSSRApp({
        page,
        resolveComponent: () => TestComponent,
      });

      expect(receivedProps.user).toEqual({ name: 'John' });
      expect(receivedProps.settings).toEqual({ theme: 'dark' });
    });
  });

  describe('createSSRPage', () => {
    it('should create complete HTML page with body', () => {
      const body = '<div id="content">Test Content</div>';
      const page: Page = {
        component: 'Test',
        props: { title: 'Test' },
        url: '/test',
        version: '1',
      };

      const html = createSSRPage(body, page);

      expect(html).toContain('<!DOCTYPE html>');
      expect(html).toContain('<html>');
      expect(html).toContain('<div id="app"');
      expect(html).toContain(body);
    });

    it('should embed page data for hydration', () => {
      const body = '<div>Test</div>';
      const page: Page = {
        component: 'TestPage',
        props: { user: { name: 'Alice' } },
        url: '/dashboard',
        version: '2',
      };

      const html = createSSRPage(body, page);

      expect(html).toContain('data-page=');
      expect(html).toContain('TestPage');
      expect(html).toContain('Alice');
      expect(html).toContain('/dashboard');
    });

    it('should include custom head tags', () => {
      const body = '<div>Test</div>';
      const page: Page = {
        component: 'Test',
        props: {},
        url: '/test',
        version: '1',
      };

      const head = [
        '<title>Test Page</title>',
        '<meta name="description" content="Test">',
      ];

      const html = createSSRPage(body, page, head);

      expect(html).toContain('<title>Test Page</title>');
      expect(html).toContain('<meta name="description" content="Test">');
    });

    it('should escape < characters in page data to prevent XSS', () => {
      const body = '<div>Test</div>';
      const page: Page = {
        component: 'Test',
        props: { html: '<script>alert("xss")</script>' },
        url: '/test',
        version: '1',
      };

      const html = createSSRPage(body, page);

      // Should not contain literal script tags in data-page attribute
      expect(html).not.toContain('<script>alert');
      // Should contain escaped version
      expect(html).toContain('\\u003c');
    });

    it('should include viewport meta tag', () => {
      const html = createSSRPage('<div>Test</div>', {
        component: 'Test',
        props: {},
        url: '/test',
        version: '1',
      });

      expect(html).toContain('<meta name="viewport" content="width=device-width, initial-scale=1">');
    });

    it('should include charset meta tag', () => {
      const html = createSSRPage('<div>Test</div>', {
        component: 'Test',
        props: {},
        url: '/test',
        version: '1',
      });

      expect(html).toContain('<meta charset="utf-8">');
    });
  });

  describe('SSR Integration', () => {
    it('should complete full SSR flow', async () => {
      const UserProfile = defineComponent({
        props: ['user', 'initialPage', 'initialComponent'],
        setup(props) {
          return () => h('div', { class: 'profile' }, [
            h('h1', props.user.name),
            h('p', props.user.email),
          ]);
        },
      });

      const page: Page = {
        component: 'UserProfile',
        props: {
          user: {
            name: 'Jane Doe',
            email: 'jane@example.com',
          },
        },
        url: '/profile',
        version: 'v1.0',
      };

      const { head, body } = await createInertiaSSRApp({
        page,
        resolveComponent: (name) => {
          if (name === 'UserProfile') return UserProfile;
          throw new Error(`Component ${name} not found`);
        },
      });

      const html = createSSRPage(body, page, head);

      // Should contain rendered content
      expect(html).toContain('Jane Doe');
      expect(html).toContain('jane@example.com');
      
      // Should contain hydration data
      expect(html).toContain('UserProfile');
      expect(html).toContain('/profile');
      
      // Should contain version in head
      expect(html).toContain('<meta name="inertia-version" content="v1.0">');
    });
  });
});
