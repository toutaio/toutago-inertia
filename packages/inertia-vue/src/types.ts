import type { App, Component } from 'vue'

export interface Page<T = Record<string, unknown>> {
  component: string
  props: T
  url: string
  version: string | null
}

export interface PageProps {
  [key: string]: unknown
  errors?: ValidationErrors
  flash?: FlashMessages
}

export interface ValidationErrors {
  [key: string]: string | string[]
}

export interface FlashMessages {
  [key: string]: unknown
}

export interface InertiaAppProps {
  initialPage: Page
  initialComponent?: Component
  resolveComponent?: (name: string) => Component | Promise<Component>
  resolve?: (name: string) => Component | Promise<Component>
  setup: (options: {
    el: Element
    App: Component
    props: Record<string, unknown>
    plugin: InertiaPlugin
  }) => void | App
  title?: (title: string) => string
  progress?:
    | false
    | {
        delay?: number
        color?: string
        includeCSS?: boolean
        showSpinner?: boolean
      }
}

export interface InertiaPlugin {
  install: (app: App) => void
}

export interface Router {
  get(url: string, options?: VisitOptions): void
  post(url: string, options?: VisitOptions): void
  put(url: string, options?: VisitOptions): void
  patch(url: string, options?: VisitOptions): void
  delete(url: string, options?: VisitOptions): void
  reload(options?: VisitOptions): void
  visit(url: string, options?: VisitOptions): void
  replace(url: string, options?: VisitOptions): void
}

export interface VisitOptions {
  method?: 'get' | 'post' | 'put' | 'patch' | 'delete'
  data?: Record<string, unknown>
  replace?: boolean
  preserveScroll?: boolean
  preserveState?: boolean
  only?: string[]
  headers?: Record<string, string>
  errorBag?: string
  forceFormData?: boolean
  onCancelToken?: (cancelToken: { cancel: () => void }) => void
  onBefore?: (visit: PendingVisit) => boolean | void
  onStart?: (visit: PendingVisit) => void
  onProgress?: (progress: Progress) => void
  onSuccess?: (page: Page) => void
  onError?: (errors: ValidationErrors) => void
  onCancel?: () => void
  onFinish?: () => void
}

export interface PendingVisit {
  url: string
  method: string
  data: Record<string, unknown>
  replace: boolean
  preserveScroll: boolean
  preserveState: boolean
  only: string[]
  headers: Record<string, string>
}

export interface Progress {
  percentage: number
  loaded: number
  total: number
}

export interface FormDataType {
  [key: string]: unknown
}

export interface InertiaForm<T extends FormDataType> {
  data: T
  errors: ValidationErrors
  hasErrors: boolean
  processing: boolean
  progress: Progress | null
  wasSuccessful: boolean
  recentlySuccessful: boolean
  isDirty: boolean
  
  get(url: string, options?: VisitOptions): void
  post(url: string, options?: VisitOptions): void
  put(url: string, options?: VisitOptions): void
  patch(url: string, options?: VisitOptions): void
  delete(url: string, options?: VisitOptions): void
  submit(method: string, url: string, options?: VisitOptions): void
  
  reset(...fields: (keyof T)[]): void
  clearErrors(...fields: (keyof T)[]): void
  setError(field: keyof T, message: string | string[]): void
  transform(callback: (data: T) => T): this
}
