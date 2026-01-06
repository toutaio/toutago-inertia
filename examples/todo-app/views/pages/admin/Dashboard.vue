<template>
  <div>
    <h1 class="title">Admin Dashboard</h1>
    <p class="subtitle">Demonstrating nested layouts</p>
    
    <div class="columns">
      <div class="column">
        <div class="notification is-info is-light">
          <h2 class="subtitle">Layout Hierarchy</h2>
          <p>This page uses:</p>
          <ol>
            <li>App Layout (outer)</li>
            <li>Admin Layout (inner)</li>
            <li>Dashboard Page (content)</li>
          </ol>
          <p class="mt-3">
            <small>Check <code>views/pages/admin/Dashboard.vue</code> to see the layout configuration.</small>
          </p>
        </div>
      </div>
      
      <div class="column">
        <div class="notification is-success is-light">
          <h2 class="subtitle">Benefits</h2>
          <ul>
            <li>✅ Reusable layout composition</li>
            <li>✅ Type-safe with TypeScript</li>
            <li>✅ Clean separation of concerns</li>
            <li>✅ Easy to understand and maintain</li>
          </ul>
        </div>
      </div>
    </div>

    <div class="box">
      <h3 class="title is-5">Stats</h3>
      <p v-if="stats">
        Users: <strong>{{ stats.users }}</strong>, 
        Todos: <strong>{{ stats.todos }}</strong>
      </p>
      <p v-else class="has-text-grey">No stats available</p>
    </div>

    <Link href="/" class="button">← Back to Home</Link>
  </div>
</template>

<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
import AdminLayout from '../../layouts/Admin.vue'

// Page props (can be typed)
defineProps<{
  stats?: {
    users: number
    todos: number
  }
}>()

// Use layout property - will be wrapped by app.ts resolver
// app.ts adds AppLayout by default, and this adds AdminLayout inside it
defineOptions({
  layout: AdminLayout
})
</script>

<style scoped>
ol, ul {
  margin-left: 1.5rem;
  margin-top: 0.5rem;
}

li {
  margin-bottom: 0.25rem;
}

.mt-3 {
  margin-top: 1rem;
}
</style>
