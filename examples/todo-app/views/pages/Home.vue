<template>
  <div class="home">
    <div class="hero is-primary">
      <div class="hero-body">
        <h1 class="title">{{ greeting }}</h1>
        <p class="subtitle">A full-stack example using Toutago and Vue 3</p>
      </div>
    </div>

    <div class="section">
      <div class="columns">
        <div class="column">
          <div class="card">
            <div class="card-content">
              <h2 class="title is-4">🚀 Features</h2>
              <ul>
                <li>✅ Server-side rendering (SSR)</li>
                <li>✅ Type-safe props with TypeScript</li>
                <li>✅ Form handling with validation</li>
                <li>✅ Flash messages</li>
                <li>✅ Authentication example</li>
              </ul>
            </div>
          </div>
        </div>

        <div class="column">
          <div class="card">
            <div class="card-content">
              <h2 class="title is-4">📚 Quick Start</h2>
              <p class="mb-3">
                This example demonstrates how to build a full-stack application
                with Toutago and Inertia.js.
              </p>
              <Link href="/todos" class="button is-primary">
                View Todos
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
import type { HomePageProps } from '../types'

defineProps<HomePageProps>()
</script>

<style scoped>
.home {
  min-height: 60vh;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  padding: 0.5rem 0;
}

.mb-3 {
  margin-bottom: 1rem;
}
</style>
