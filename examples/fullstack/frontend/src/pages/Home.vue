<template>
  <Layout>
    <div class="home">
      <h1>{{ message }}</h1>
      
      <div class="stats">
        <div class="stat-card">
          <h3>{{ stats.users }}</h3>
          <p>Users</p>
        </div>
        <div class="stat-card">
          <h3>{{ stats.posts }}</h3>
          <p>Posts</p>
        </div>
      </div>

      <div class="cta">
        <Link href="/users" class="btn btn-primary">View Users</Link>
        <Link href="/posts" class="btn btn-secondary">View Posts</Link>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
import Layout from '../components/Layout.vue'

interface Props {
  message: string
  stats: {
    users: number
    posts: number
  }
}

defineProps<Props>()
</script>

<style scoped>
.home {
  text-align: center;
  padding: 2rem 0;
}

h1 {
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 2rem;
}

.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 2rem;
  margin: 2rem 0;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

.stat-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-card h3 {
  font-size: 3rem;
  color: #3498db;
  margin: 0;
}

.stat-card p {
  margin: 0.5rem 0 0;
  color: #7f8c8d;
  text-transform: uppercase;
  font-size: 0.875rem;
}

.cta {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 3rem;
}

.btn {
  padding: 0.75rem 2rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-primary {
  background: #3498db;
  color: white;
}

.btn-primary:hover {
  background: #2980b9;
}

.btn-secondary {
  background: #95a5a6;
  color: white;
}

.btn-secondary:hover {
  background: #7f8c8d;
}
</style>
