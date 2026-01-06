<template>
  <!-- Nothing rendered -->
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, useSlots } from 'vue'

interface HeadProps {
  title?: string
}

const props = defineProps<HeadProps>()
const slots = useSlots()

let titleElement: HTMLTitleElement | null = null
const metaElements: HTMLMetaElement[] = []

onMounted(() => {
  // Set title
  if (props.title) {
    titleElement = document.createElement('title')
    titleElement.textContent = props.title
    document.head.appendChild(titleElement)
  }

  // Process slot content for meta tags
  if (slots.default) {
    const content = slots.default()
    // In a real implementation, we'd parse and insert meta tags
    // For now, this is a simplified version
  }
})

onBeforeUnmount(() => {
  // Clean up title
  if (titleElement && titleElement.parentNode) {
    titleElement.parentNode.removeChild(titleElement)
  }

  // Clean up meta tags
  metaElements.forEach(el => {
    if (el.parentNode) {
      el.parentNode.removeChild(el)
    }
  })
})
</script>
