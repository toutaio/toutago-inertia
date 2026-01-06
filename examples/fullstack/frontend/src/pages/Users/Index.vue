<template>
  <Layout>
    <div class="users">
      <h1>Users</h1>
      
      <div class="user-grid">
        <Link
          v-for="user in users"
          :key="user.id"
          :href="`/users/${user.id}`"
          class="user-card"
        >
          <div class="user-avatar">
            {{ user.name.charAt(0) }}
          </div>
          <h3>{{ user.name }}</h3>
          <p>{{ user.email }}</p>
          <small>Joined {{ formatDate(user.createdAt) }}</small>
        </Link>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
import Layout from '../../components/Layout.vue'

interface User {
  id: number
  name: string
  email: string
  createdAt: string
}

interface Props {
  users: User[]
}

defineProps<Props>()

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<style scoped>
.users {
  padding: 2rem 0;
}

h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 2rem;
}

.user-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.user-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
  text-align: center;
}

.user-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #3498db;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0 auto 1rem;
}

.user-card h3 {
  margin: 0 0 0.5rem;
  color: #2c3e50;
}

.user-card p {
  margin: 0 0 0.5rem;
  color: #7f8c8d;
  font-size: 0.875rem;
}

.user-card small {
  color: #95a5a6;
  font-size: 0.75rem;
}
</style>
