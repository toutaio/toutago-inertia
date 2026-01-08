import { defineComponent, h } from 'vue'
import { router } from '../router'
import type { VisitOptions } from '../types'

export interface LinkProps {
  href: string
  method?: 'get' | 'post' | 'put' | 'patch' | 'delete'
  data?: Record<string, unknown>
  replace?: boolean
  preserveScroll?: boolean
  preserveState?: boolean
  only?: string[]
  headers?: Record<string, string>
}

export const Link = defineComponent({
  name: 'Link',
  props: {
    href: {
      type: String,
      required: true
    },
    method: {
      type: String as () => 'get' | 'post' | 'put' | 'patch' | 'delete',
      default: 'get'
    },
    data: {
      type: Object as () => Record<string, unknown>,
      default: () => ({})
    },
    replace: {
      type: Boolean,
      default: false
    },
    preserveScroll: {
      type: Boolean,
      default: false
    },
    preserveState: {
      type: Boolean,
      default: false
    },
    only: {
      type: Array as () => string[],
      default: () => []
    },
    headers: {
      type: Object as () => Record<string, string>,
      default: () => ({})
    }
  },
  emits: ['click', 'before', 'start', 'success', 'error'],
  setup(props, { slots, emit }) {
    function handleClick(event: MouseEvent) {
      event.preventDefault()
      emit('click', event)

      const options: VisitOptions = {
        method: props.method,
        data: props.data,
        replace: props.replace,
        preserveScroll: props.preserveScroll,
        preserveState: props.preserveState,
        only: props.only,
        headers: props.headers,
        onBefore: () => {
          emit('before')
        },
        onStart: () => {
          emit('start')
        },
        onSuccess: () => {
          emit('success')
        },
        onError: () => {
          emit('error')
        }
      }

      router.visit(props.href, options)
    }

    return () => h('a', {
      href: props.href,
      onClick: handleClick
    }, slots.default?.())
  }
})
