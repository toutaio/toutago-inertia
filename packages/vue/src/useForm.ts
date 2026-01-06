import { reactive, computed, watch, nextTick } from 'vue'
import { router, VisitOptions } from './router'

export interface FormData {
  [key: string]: any
}

export interface FormErrors {
  [key: string]: string | undefined
}

export interface Form<T extends FormData> {
  data: T
  errors: FormErrors
  isDirty: boolean
  processing: boolean
  progress: number | null
  recentlySuccessful: boolean
  wasSuccessful: boolean
  
  reset(...fields: (keyof T)[]): void
  clearErrors(...fields: (keyof T)[]): void
  submit(method: string, url: string, options?: Partial<VisitOptions>): void
  get(url: string, options?: Partial<VisitOptions>): void
  post(url: string, options?: Partial<VisitOptions>): void
  put(url: string, options?: Partial<VisitOptions>): void
  patch(url: string, options?: Partial<VisitOptions>): void
  delete(url: string, options?: Partial<VisitOptions>): void
}

export function useForm<T extends FormData>(initialData: T): Form<T> {
  const defaults = JSON.parse(JSON.stringify(initialData))
  
  const state = reactive({
    data: JSON.parse(JSON.stringify(initialData)) as T,
    errors: {} as FormErrors,
    processing: false,
    progress: null as number | null,
    recentlySuccessful: false,
    wasSuccessful: false,
  })

  const isDirty = computed(() => {
    return Object.keys(state.data).some(
      (key) => JSON.stringify(state.data[key]) !== JSON.stringify(defaults[key])
    )
  })

  const form = {
    get data() {
      return state.data
    },
    set data(value: T) {
      state.data = value
    },
    get errors() {
      return state.errors
    },
    set errors(value: FormErrors) {
      state.errors = value
    },
    get isDirty() {
      return isDirty.value
    },
    get processing() {
      return state.processing
    },
    set processing(value: boolean) {
      state.processing = value
    },
    get progress() {
      return state.progress
    },
    set progress(value: number | null) {
      state.progress = value
    },
    get recentlySuccessful() {
      return state.recentlySuccessful
    },
    set recentlySuccessful(value: boolean) {
      state.recentlySuccessful = value
    },
    get wasSuccessful() {
      return state.wasSuccessful
    },
    set wasSuccessful(value: boolean) {
      state.wasSuccessful = value
    },

    reset(...fields: (keyof T)[]) {
      if (fields.length === 0) {
        state.data = JSON.parse(JSON.stringify(defaults)) as T
      } else {
        fields.forEach((field) => {
          state.data[field] = JSON.parse(JSON.stringify(defaults[field]))
        })
      }
    },

    clearErrors(...fields: (keyof T)[]) {
      if (fields.length === 0) {
        state.errors = {}
      } else {
        fields.forEach((field) => {
          delete state.errors[field as string]
        })
      }
    },

    submit(method: string, url: string, options: Partial<VisitOptions> = {}) {
      state.processing = true
      state.recentlySuccessful = false
      state.wasSuccessful = false

      const visitOptions: VisitOptions = {
        method: method as any,
        data: state.data,
        preserveScroll: options.preserveScroll ?? true,
        onStart: () => {
          state.processing = true
          options.onStart?.()
        },
        onProgress: (progress) => {
          state.progress = progress
          options.onProgress?.(progress)
        },
        onSuccess: (response) => {
          state.processing = false
          state.recentlySuccessful = true
          state.wasSuccessful = true
          state.progress = null
          
          setTimeout(() => {
            state.recentlySuccessful = false
          }, 2000)
          
          options.onSuccess?.(response)
        },
        onError: (errors) => {
          state.processing = false
          state.progress = null
          state.errors = errors
          options.onError?.(errors)
        },
        onFinish: () => {
          state.processing = false
          state.progress = null
          options.onFinish?.()
        },
        ...options,
      }

      router.visit(url, visitOptions)
    },

    get(url: string, options?: Partial<VisitOptions>) {
      this.submit('get', url, options)
    },

    post(url: string, options?: Partial<VisitOptions>) {
      this.submit('post', url, options)
    },

    put(url: string, options?: Partial<VisitOptions>) {
      this.submit('put', url, options)
    },

    patch(url: string, options?: Partial<VisitOptions>) {
      this.submit('patch', url, options)
    },

    delete(url: string, options?: Partial<VisitOptions>) {
      this.submit('delete', url, options)
    },
  } as Form<T>

  return form
}
