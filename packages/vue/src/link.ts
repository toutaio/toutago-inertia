import { defineComponent, h, PropType } from 'vue'
import { router } from './router'

export interface LinkProps {
  href: string
  method?: 'get' | 'post' | 'put' | 'patch' | 'delete'
  data?: Record<string, any>
  preserveScroll?: boolean
  preserveState?: boolean
  replace?: boolean
  only?: string[]
  headers?: Record<string, string>
  as?: string
}

export const Link = defineComponent({
  name: 'InertiaLink',
  props: {
    href: {
      type: String,
      required: true,
    },
    method: {
      type: String as PropType<'get' | 'post' | 'put' | 'patch' | 'delete'>,
      default: 'get',
    },
    data: {
      type: Object as PropType<Record<string, any>>,
      default: () => ({}),
    },
    preserveScroll: {
      type: Boolean,
      default: false,
    },
    preserveState: {
      type: Boolean,
      default: false,
    },
    replace: {
      type: Boolean,
      default: false,
    },
    only: {
      type: Array as PropType<string[]>,
      default: () => [],
    },
    headers: {
      type: Object as PropType<Record<string, string>>,
      default: () => ({}),
    },
    as: {
      type: String,
      default: 'a',
    },
  },
  setup(props, { slots, attrs }) {
    const isExternal = (href: string) => {
      try {
        const url = new URL(href, window.location.origin)
        return url.origin !== window.location.origin
      } catch {
        return false
      }
    }

    const onClick = (event: MouseEvent) => {
      // Don't intercept if default is already prevented
      if (event.defaultPrevented) {
        return
      }

      // Don't intercept external links
      if (isExternal(props.href)) {
        return
      }

      // Allow normal browser behavior for modified clicks
      if (
        event.ctrlKey ||
        event.shiftKey ||
        event.metaKey ||
        event.altKey ||
        event.button !== 0
      ) {
        return
      }

      event.preventDefault()

      router.visit(props.href, {
        method: props.method,
        data: props.data,
        preserveScroll: props.preserveScroll,
        preserveState: props.preserveState,
        replace: props.replace,
        only: props.only.length > 0 ? props.only : undefined,
        headers: props.headers,
      })
    }

    return () => {
      const tag = props.as
      const children = slots.default?.()

      if (tag === 'a') {
        return h(
          'a',
          {
            ...attrs,
            href: props.href,
            onClick,
          },
          children
        )
      }

      return h(
        tag,
        {
          ...attrs,
          onClick,
        },
        children
      )
    }
  },
})
