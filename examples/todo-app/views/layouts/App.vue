<template>
  <div id="app">
    <nav class="navbar">
      <div class="container">
        <div class="navbar-brand">
          <Link href="/" class="navbar-item">
            <strong>Toutago Todo</strong>
          </Link>
        </div>
        <div class="navbar-menu">
          <div class="navbar-end">
            <Link href="/todos" class="navbar-item">Todos</Link>
            <div v-if="$page.props.user" class="navbar-item">
              <span class="mr-2">{{ $page.props.user.name }}</span>
              <form @submit.prevent="logout" style="display: inline">
                <button type="submit" class="button is-small">Logout</button>
              </form>
            </div>
            <Link v-else href="/login" class="navbar-item">Login</Link>
          </div>
        </div>
      </div>
    </nav>

    <main class="section">
      <div class="container">
        <!-- Flash messages -->
        <div v-if="$page.props.flash?.success" class="notification is-success">
          <button class="delete" @click="clearFlash('success')"></button>
          {{ $page.props.flash.success }}
        </div>
        <div v-if="$page.props.flash?.error" class="notification is-danger">
          <button class="delete" @click="clearFlash('error')"></button>
          {{ $page.props.flash.error }}
        </div>

        <!-- Page content -->
        <slot />
      </div>
    </main>

    <footer class="footer">
      <div class="content has-text-centered">
        <p>
          Built with <strong>Toutago</strong> and <strong>Inertia.js</strong>
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@toutaio/inertia-vue'

const logout = () => {
  router.post('/logout')
}

const clearFlash = (key: string) => {
  // Flash messages are automatically cleared on next request
}
</script>

<style>
.navbar {
  background-color: #3273dc;
  color: white;
}

.navbar-brand .navbar-item,
.navbar-menu .navbar-item {
  color: white;
}

.navbar-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
}

.mr-2 {
  margin-right: 0.5rem;
}
</style>
