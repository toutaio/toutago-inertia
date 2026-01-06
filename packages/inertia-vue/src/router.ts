import axios, { type AxiosInstance } from 'axios'
import type { Page, VisitOptions, Router } from './types'

class InertiaRouter implements Router {
  private axios: AxiosInstance
  private currentPage: Page | null = null

  constructor() {
    this.axios = axios.create({
      headers: {
        'X-Requested-With': 'XMLHttpRequest',
        'X-Inertia': 'true',
        'X-Inertia-Version': '',
        'Accept': 'text/html, application/xhtml+xml'
      },
      withCredentials: true
    })
  }

  visit(url: string, options: VisitOptions = {}): void {
    const {
      method = 'get',
      data = {},
      replace = false,
      preserveScroll = false,
      preserveState = false,
      only = [],
      headers = {},
      onBefore,
      onStart,
      onSuccess,
      onError,
      onFinish
    } = options

    const visit = {
      url,
      method,
      data,
      replace,
      preserveScroll,
      preserveState,
      only,
      headers
    }

    if (onBefore && onBefore(visit) === false) {
      return
    }

    onStart?.(visit)

    const requestHeaders: Record<string, string> = {
      ...headers
    }

    if (only.length > 0) {
      requestHeaders['X-Inertia-Partial-Data'] = only.join(',')
      requestHeaders['X-Inertia-Partial-Component'] = this.currentPage?.component || ''
    }

    const requestMethod = method.toLowerCase()
    const promise = requestMethod === 'get'
      ? this.axios.get(url, { headers: requestHeaders })
      : this.axios[requestMethod as 'post' | 'put' | 'patch' | 'delete'](url, data, { headers: requestHeaders })

    promise
      .then(response => {
        const page = response.data as Page
        this.currentPage = page

        if (replace) {
          window.history.replaceState(page, '', page.url)
        } else {
          window.history.pushState(page, '', page.url)
        }

        onSuccess?.(page)
      })
      .catch(error => {
        if (error.response?.data?.errors) {
          onError?.(error.response.data.errors)
        }
      })
      .finally(() => {
        onFinish?.()
      })
  }

  get(url: string, options: VisitOptions = {}): void {
    this.visit(url, { ...options, method: 'get' })
  }

  post(url: string, options: VisitOptions = {}): void {
    this.visit(url, { ...options, method: 'post' })
  }

  put(url: string, options: VisitOptions = {}): void {
    this.visit(url, { ...options, method: 'put' })
  }

  patch(url: string, options: VisitOptions = {}): void {
    this.visit(url, { ...options, method: 'patch' })
  }

  delete(url: string, options: VisitOptions = {}): void {
    this.visit(url, { ...options, method: 'delete' })
  }

  reload(options: VisitOptions = {}): void {
    if (this.currentPage) {
      this.visit(this.currentPage.url, { ...options, preserveState: true })
    }
  }

  replace(url: string, options: VisitOptions = {}): void {
    this.visit(url, { ...options, replace: true })
  }

  setPage(page: Page): void {
    this.currentPage = page
  }

  getPage(): Page | null {
    return this.currentPage
  }
}

export const router = new InertiaRouter()
