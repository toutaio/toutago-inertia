export { createInertiaApp, router } from './app'
export { default as Link } from './components/Link.vue'
export { default as Head } from './components/Head.vue'
export { useForm } from './composables/useForm'
export { usePage } from './composables/usePage'

export type {
  Page,
  PageProps,
  InertiaAppProps,
  Router,
  VisitOptions,
  InertiaForm,
  FormDataType,
  ValidationErrors,
  FlashMessages
} from './types'
