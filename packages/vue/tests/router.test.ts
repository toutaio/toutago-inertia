import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { router } from '../src/router';
import type { Page } from '../src/types';

describe('Router', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({
        component: 'Test',
        props: {},
        url: '/test',
        version: '1.0.0',
      }),
    });
    (global as any).fetch = fetchMock;
    (global as any).window = {
      history: {
        pushState: vi.fn(),
        replaceState: vi.fn(),
      },
      location: {
        href: 'http://localhost/current',
      },
      dispatchEvent: vi.fn(),
    };

    const mockPage: Page = {
      component: 'Home',
      props: {},
      url: '/current',
      version: '1.0.0',
    };
    router.init(mockPage);
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('visit', () => {
    it('makes GET request with correct headers', () => {
      router.visit('/users');

      expect(fetchMock).toHaveBeenCalledWith('/users', {
        method: 'GET',
        headers: expect.objectContaining({
          'X-Inertia': 'true',
          'X-Inertia-Version': '1.0.0',
        }),
        credentials: 'same-origin',
      });
    });

    it('makes POST request with data', () => {
      const data = { name: 'John', email: 'john@example.com' };
      router.visit('/users', { method: 'post', data });

      expect(fetchMock).toHaveBeenCalledWith('/users', {
        method: 'POST',
        headers: expect.objectContaining({
          'X-Inertia': 'true',
          'Content-Type': 'application/json',
        }),
        credentials: 'same-origin',
        body: JSON.stringify(data),
      });
    });

    it('includes query parameters for GET requests', () => {
      const data = { page: '2', search: 'test' };
      router.visit('/users', { method: 'get', data });

      expect(fetchMock).toHaveBeenCalledWith(
        '/users?page=2&search=test',
        expect.objectContaining({
          method: 'GET',
        })
      );
    });

    it('includes partial data headers when only is specified', () => {
      router.visit('/users', { only: ['users', 'meta'] });

      expect(fetchMock).toHaveBeenCalledWith('/users', {
        method: 'GET',
        headers: expect.objectContaining({
          'X-Inertia-Partial-Data': 'users,meta',
          'X-Inertia-Partial-Component': 'Home',
        }),
        credentials: 'same-origin',
      });
    });

    it('uses method override for PUT requests', () => {
      router.visit('/users/1', { method: 'put', data: { name: 'John' } });

      expect(fetchMock).toHaveBeenCalledWith('/users/1', {
        method: 'POST',
        headers: expect.objectContaining({
          'X-HTTP-Method-Override': 'PUT',
        }),
        credentials: 'same-origin',
        body: JSON.stringify({ name: 'John' }),
      });
    });

    it('uses method override for PATCH requests', () => {
      router.visit('/users/1', { method: 'patch', data: { name: 'John' } });

      expect(fetchMock).toHaveBeenCalledWith('/users/1', {
        method: 'POST',
        headers: expect.objectContaining({
          'X-HTTP-Method-Override': 'PATCH',
        }),
        credentials: 'same-origin',
        body: JSON.stringify({ name: 'John' }),
      });
    });

    it('uses method override for DELETE requests', () => {
      router.visit('/users/1', { method: 'delete' });

      expect(fetchMock).toHaveBeenCalledWith('/users/1', {
        method: 'POST',
        headers: expect.objectContaining({
          'X-HTTP-Method-Override': 'DELETE',
        }),
        credentials: 'same-origin',
        body: JSON.stringify({}),
      });
    });
  });

  describe('convenience methods', () => {
    it('get() calls visit with GET method', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      router.get('/users', { page: 2 });

      expect(visitSpy).toHaveBeenCalledWith('/users', {
        method: 'get',
        data: { page: 2 },
      });
    });

    it('post() calls visit with POST method', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      const data = { name: 'John' };
      router.post('/users', data);

      expect(visitSpy).toHaveBeenCalledWith('/users', {
        method: 'post',
        data,
      });
    });

    it('put() calls visit with PUT method', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      const data = { name: 'John' };
      router.put('/users/1', data);

      expect(visitSpy).toHaveBeenCalledWith('/users/1', {
        method: 'put',
        data,
      });
    });

    it('patch() calls visit with PATCH method', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      const data = { name: 'John' };
      router.patch('/users/1', data);

      expect(visitSpy).toHaveBeenCalledWith('/users/1', {
        method: 'patch',
        data,
      });
    });

    it('delete() calls visit with DELETE method', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      router.delete('/users/1');

      expect(visitSpy).toHaveBeenCalledWith('/users/1', {
        method: 'delete',
      });
    });

    it('reload() preserves scroll and state', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      router.reload();

      expect(visitSpy).toHaveBeenCalledWith('http://localhost/current', {
        preserveScroll: true,
        preserveState: true,
      });
    });

    it('replace() sets replace option', () => {
      const visitSpy = vi.spyOn(router, 'visit');
      router.replace('/users');

      expect(visitSpy).toHaveBeenCalledWith('/users', {
        replace: true,
      });
    });
  });

  describe('response handling', () => {
    it('pushes state to history for non-replace visits', async () => {
      const newPage: Page = {
        component: 'Users',
        props: { users: [] },
        url: '/users',
        version: '1.0.0',
      };

      fetchMock.mockResolvedValueOnce({
        json: () => Promise.resolve(newPage),
      });

      router.visit('/users');

      await new Promise((resolve) => setTimeout(resolve, 10));

      expect((global as any).window.history.pushState).toHaveBeenCalledWith(newPage, '', '/users');
    });

    it('replaces state in history for replace visits', async () => {
      const newPage: Page = {
        component: 'Users',
        props: { users: [] },
        url: '/users',
        version: '1.0.0',
      };

      fetchMock.mockResolvedValueOnce({
        json: () => Promise.resolve(newPage),
      });

      router.visit('/users', { replace: true });

      await new Promise((resolve) => setTimeout(resolve, 10));

      expect((global as any).window.history.replaceState).toHaveBeenCalledWith(newPage, '', '/users');
    });

    it('dispatches inertia:navigate event', async () => {
      const newPage: Page = {
        component: 'Users',
        props: { users: [] },
        url: '/users',
        version: '1.0.0',
      };

      fetchMock.mockResolvedValueOnce({
        json: () => Promise.resolve(newPage),
      });

      router.visit('/users');

      await new Promise((resolve) => setTimeout(resolve, 10));

      expect((global as any).window.dispatchEvent).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'inertia:navigate',
          detail: { page: newPage },
        })
      );
    });

    it('handles fetch errors', async () => {
      const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});
      fetchMock.mockRejectedValueOnce(new Error('Network error'));

      router.visit('/users');

      await new Promise((resolve) => setTimeout(resolve, 10));

      expect(consoleError).toHaveBeenCalledWith(
        'Inertia request failed:',
        expect.any(Error)
      );

      consoleError.mockRestore();
    });
  });
});
