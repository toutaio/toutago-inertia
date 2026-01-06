<template>
  <div class="layout">
    <nav class="navbar">
      <div class="container">
        <Link href="/" class="logo">Toutago</Link>
        <div class="nav-links">
          <Link href="/">Home</Link>
          <Link href="/users">Users</Link>
          <Link href="/posts">Posts</Link>
          <Link href="/about">About</Link>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <div class="container">
        <div v-if="$page.props.flash?.success" class="alert alert-success">
          {{ $page.props.flash.success }}
        </div>
        <div v-if="$page.props.flash?.error" class="alert alert-error">
          {{ $page.props.flash.error }}
        </div>
        
        <slot />
      </div>
    </main>

    <footer class="footer">
      <div class="container">
        <p>&copy; 2024 Toutago. Built with ❤️ and Go.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: #2c3e50;
  color: white;
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.navbar .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: white;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
}

.nav-links a {
  color: white;
  text-decoration: none;
  transition: opacity 0.2s;
}

.nav-links a:hover {
  opacity: 0.8;
}

.main-content {
  flex: 1;
  padding: 2rem 0;
}

.alert {
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.alert-success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.footer {
  background: #34495e;
  color: white;
  padding: 1.5rem 0;
  text-align: center;
}
</style>
