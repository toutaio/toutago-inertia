<template>
  <a
    :href="href"
    @click.prevent="handleClick"
  >
    <slot />
  </a>
</template>

<script setup lang="ts">
import { router } from '../router'
import type { VisitOptions } from '../types'

interface LinkProps {
  href: string
  method?: 'get' | 'post' | 'put' | 'patch' | 'delete'
  data?: Record<string, unknown>
  replace?: boolean
  preserveScroll?: boolean
  preserveState?: boolean
  only?: string[]
  headers?: Record<string, string>
}

const props = withDefaults(defineProps<LinkProps>(), {
  method: 'get',
  data: () => ({}),
  replace: false,
  preserveScroll: false,
  preserveState: false,
  only: () => [],
  headers: () => ({})
})

const emit = defineEmits<{
  click: [event: MouseEvent]
  before: []
  start: []
  success: []
  error: []
}>()

function handleClick(event: MouseEvent) {
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
</script>
