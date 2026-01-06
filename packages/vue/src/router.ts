import { Method, Page, VisitOptions } from './types';

export interface Router {
  visit(url: string, options?: VisitOptions): void;
  get(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void;
  post(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void;
  put(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void;
  patch(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void;
  delete(url: string, options?: Omit<VisitOptions, 'method' | 'data'>): void;
  reload(options?: Omit<VisitOptions, 'preserveScroll' | 'preserveState'>): void;
  replace(url: string, options?: Omit<VisitOptions, 'replace'>): void;
}

class InertiaRouter implements Router {
  private page: Page | null = null;

  public init(page: Page): void {
    this.page = page;
  }

  public visit(url: string, options: VisitOptions = {}): void {
    const method = options.method || 'get';
    const data = options.data || {};
    const headers: Record<string, string> = {
      'X-Inertia': 'true',
      'X-Inertia-Version': this.page?.version || '',
      ...options.headers,
    };

    if (options.only) {
      headers['X-Inertia-Partial-Data'] = options.only.join(',');
      headers['X-Inertia-Partial-Component'] = this.page?.component || '';
    }

    this.makeRequest(url, method, data, headers, options);
  }

  public get(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void {
    this.visit(url, { ...options, method: 'get', data });
  }

  public post(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void {
    this.visit(url, { ...options, method: 'post', data });
  }

  public put(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void {
    this.visit(url, { ...options, method: 'put', data });
  }

  public patch(url: string, data?: any, options?: Omit<VisitOptions, 'method' | 'data'>): void {
    this.visit(url, { ...options, method: 'patch', data });
  }

  public delete(url: string, options?: Omit<VisitOptions, 'method' | 'data'>): void {
    this.visit(url, { ...options, method: 'delete' });
  }

  public reload(options?: Omit<VisitOptions, 'preserveScroll' | 'preserveState'>): void {
    this.visit(window.location.href, {
      ...options,
      preserveScroll: true,
      preserveState: true,
    });
  }

  public replace(url: string, options?: Omit<VisitOptions, 'replace'>): void {
    this.visit(url, { ...options, replace: true });
  }

  private makeRequest(
    url: string,
    method: Method,
    data: any,
    headers: Record<string, string>,
    options: VisitOptions
  ): void {
    const isGet = method.toLowerCase() === 'get';

    const fetchOptions: RequestInit = {
      method: isGet ? 'GET' : 'POST',
      headers: {
        Accept: 'text/html, application/xhtml+xml',
        'Content-Type': 'application/json',
        ...headers,
      },
      credentials: 'same-origin',
    };

    let fetchUrl = url;

    if (isGet && Object.keys(data).length > 0) {
      const params = new URLSearchParams(data);
      fetchUrl = `${url}?${params.toString()}`;
    } else if (!isGet) {
      fetchOptions.body = JSON.stringify(data);
      if (method.toLowerCase() !== 'post') {
        fetchOptions.headers = {
          ...fetchOptions.headers,
          'X-HTTP-Method-Override': method.toUpperCase(),
        };
      }
    }

    fetch(fetchUrl, fetchOptions)
      .then((response) => response.json())
      .then((newPage: Page) => {
        this.handleResponse(newPage, options);
      })
      .catch((error) => {
        console.error('Inertia request failed:', error);
      });
  }

  private handleResponse(newPage: Page, options: VisitOptions): void {
    if (options.replace) {
      window.history.replaceState(newPage, '', newPage.url);
    } else {
      window.history.pushState(newPage, '', newPage.url);
    }

    this.page = newPage;
    this.updatePage(newPage);
  }

  private updatePage(page: Page): void {
    (window as any).INERTIA_PAGE = page;
    window.dispatchEvent(new CustomEvent('inertia:navigate', { detail: { page } }));
  }
}

export const router = new InertiaRouter();
