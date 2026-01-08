export { createInertiaApp, router } from './app'
export { Link } from './components/Link'
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
export type { LinkProps } from './components/Link'
