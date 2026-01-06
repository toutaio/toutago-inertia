export { createInertiaApp, router } from './app'
export { default as Link } from './components/Link.vue'
export { default as Head } from './components/Head.vue'
export { useForm } from './composables/useForm'
export { usePage } from './composables/usePage'
export { useRemember } from './composables/useRemember'
export { useLiveUpdate } from './composables/useLiveUpdate'
export { withLayout, resolvePageLayout } from './utils/layout'

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

export type { PageWithLayout } from './utils/layout'
