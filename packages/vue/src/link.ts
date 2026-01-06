import { DefineComponent, PropType, computed, defineComponent, h } from 'vue';
import { router } from './router';
import { Method } from './types';

export interface LinkProps {
  href: string;
  method?: Method;
  data?: any;
  replace?: boolean;
  preserveScroll?: boolean;
  preserveState?: boolean;
  only?: string[];
  headers?: Record<string, string>;
  as?: string;
}

export const Link: DefineComponent<LinkProps> = defineComponent({
  name: 'InertiaLink',
  props: {
    href: {
      type: String,
      required: true,
    },
    method: {
      type: String as PropType<Method>,
      default: 'get',
    },
    data: {
      type: Object,
      default: () => ({}),
    },
    replace: {
      type: Boolean,
      default: false,
    },
    preserveScroll: {
      type: Boolean,
      default: false,
    },
    preserveState: {
      type: Boolean,
      default: false,
    },
    only: {
      type: Array as PropType<string[]>,
      default: () => [],
    },
    headers: {
      type: Object,
      default: () => ({}),
    },
    as: {
      type: String,
      default: 'a',
    },
  },
  setup(props: LinkProps, { slots, attrs }) {
    const isExternal = computed(() => {
      return /^https?:\/\//.test(props.href);
    });

    const onClick = (event: MouseEvent) => {
      if (shouldIntercept(event)) {
        event.preventDefault();

        router.visit(props.href, {
          method: props.method,
          data: props.data,
          replace: props.replace,
          preserveScroll: props.preserveScroll,
          preserveState: props.preserveState,
          only: props.only && props.only.length > 0 ? props.only : undefined,
          headers: props.headers,
        });
      }
    };

    const shouldIntercept = (event: MouseEvent): boolean => {
      if (isExternal.value) {
        return false;
      }

      if (event.defaultPrevented) {
        return false;
      }

      if (event.button !== 0) {
        return false;
      }

      if (event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) {
        return false;
      }

      return true;
    };

    return () => {
      const tag = props.as || 'a';
      const elementAttrs: Record<string, any> = { ...attrs };

      if (tag === 'a') {
        elementAttrs.href = props.href;
      }

      elementAttrs.onClick = onClick;

      return h(tag, elementAttrs, slots.default?.());
    };
  },
});
