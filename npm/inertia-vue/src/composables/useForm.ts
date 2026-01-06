import { reactive, watch } from 'vue'
import { router } from '../router'
import type { FormDataType, InertiaForm, VisitOptions, ValidationErrors, Progress } from '../types'

export function useForm<T extends FormDataType>(initialData: T): InertiaForm<T> {
  const defaults = JSON.parse(JSON.stringify(initialData))
  
  const form = reactive({
    data: { ...initialData },
    errors: {} as ValidationErrors,
    hasErrors: false,
    processing: false,
    progress: null as Progress | null,
    wasSuccessful: false,
    recentlySuccessful: false,
    isDirty: false,
    _transformCallback: null as ((data: T) => T) | null,

    get(url: string, options: VisitOptions = {}) {
      this.submit('get', url, options)
    },

    post(url: string, options: VisitOptions = {}) {
      this.submit('post', url, options)
    },

    put(url: string, options: VisitOptions = {}) {
      this.submit('put', url, options)
    },

    patch(url: string, options: VisitOptions = {}) {
      this.submit('patch', url, options)
    },

    delete(url: string, options: VisitOptions = {}) {
      this.submit('delete', url, options)
    },

    submit(method: string, url: string, options: VisitOptions = {}) {
      this.processing = true
      this.wasSuccessful = false
      this.recentlySuccessful = false

      let data = this.data
      if (this._transformCallback) {
        data = this._transformCallback(data)
      }

      const visitOptions: VisitOptions = {
        ...options,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        method: method as any,
        data,
        onProgress: (progress) => {
          this.progress = progress
          options.onProgress?.(progress)
        },
        onSuccess: (page) => {
          this.wasSuccessful = true
          this.recentlySuccessful = true
          this.errors = {}
          this.hasErrors = false
          this.isDirty = false
          
          setTimeout(() => {
            this.recentlySuccessful = false
          }, 2000)
          
          options.onSuccess?.(page)
        },
        onError: (errors) => {
          this.errors = errors
          this.hasErrors = Object.keys(errors).length > 0
          options.onError?.(errors)
        },
        onFinish: () => {
          this.processing = false
          this.progress = null
          options.onFinish?.()
        }
      }

      router.visit(url, visitOptions)
    },

    reset(...fields: (keyof T)[]) {
      if (fields.length === 0) {
        Object.assign(this.data, defaults)
      } else {
        fields.forEach(field => {
          this.data[field] = defaults[field]
        })
      }
      this.isDirty = false
    },

    clearErrors(...fields: (keyof T)[]) {
      if (fields.length === 0) {
        this.errors = {}
      } else {
        fields.forEach(field => {
          delete this.errors[field as string]
        })
      }
      this.hasErrors = Object.keys(this.errors).length > 0
    },

    setError(field: keyof T, message: string | string[]) {
      this.errors[field as string] = message
      this.hasErrors = true
    },

    transform(callback: (data: T) => T) {
      this._transformCallback = callback
      return this
    }
  })

  // Watch for data changes to track isDirty
  watch(
    () => form.data,
    () => {
      form.isDirty = JSON.stringify(form.data) !== JSON.stringify(defaults)
    },
    { deep: true }
  )

  return form as InertiaForm<T>
}
