<template>
  <Layout>
    <div class="user-detail">
      <div class="user-header">
        <div class="user-avatar-large">
          {{ user.name.charAt(0) }}
        </div>
        <div>
          <h1>{{ user.name }}</h1>
          <p class="email">{{ user.email }}</p>
          <small>Member since {{ formatDate(user.createdAt) }}</small>
        </div>
      </div>

      <div class="user-posts">
        <h2>Posts by {{ user.name }}</h2>
        
        <div v-if="posts.length === 0" class="no-posts">
          No posts yet.
        </div>

        <div v-else class="posts-list">
          <Link
            v-for="post in posts"
            :key="post.id"
            :href="`/posts/${post.id}`"
            class="post-card"
          >
            <h3>{{ post.title }}</h3>
            <p>{{ post.content.substring(0, 150) }}...</p>
            <small>{{ formatDate(post.created) }}</small>
          </Link>
        </div>
      </div>

      <div class="actions">
        <Link href="/users" class="btn">← Back to Users</Link>
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

interface Post {
  id: number
  title: string
  content: string
  created: string
}

interface Props {
  user: User
  posts: Post[]
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
.user-detail {
  padding: 2rem 0;
}

.user-header {
  display: flex;
  gap: 2rem;
  align-items: center;
  margin-bottom: 3rem;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.user-avatar-large {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: #3498db;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  font-weight: bold;
  flex-shrink: 0;
}

.user-header h1 {
  margin: 0 0 0.5rem;
  color: #2c3e50;
}

.email {
  margin: 0 0 0.5rem;
  color: #7f8c8d;
}

small {
  color: #95a5a6;
}

.user-posts {
  margin-bottom: 2rem;
}

.user-posts h2 {
  font-size: 1.5rem;
  color: #2c3e50;
  margin-bottom: 1.5rem;
}

.no-posts {
  text-align: center;
  padding: 3rem;
  color: #95a5a6;
  background: white;
  border-radius: 8px;
}

.posts-list {
  display: grid;
  gap: 1rem;
}

.post-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
}

.post-card:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.post-card h3 {
  margin: 0 0 0.5rem;
  color: #2c3e50;
}

.post-card p {
  margin: 0 0 0.5rem;
  color: #7f8c8d;
  line-height: 1.6;
}

.post-card small {
  color: #95a5a6;
  font-size: 0.875rem;
}

.actions {
  margin-top: 2rem;
}

.btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background: #95a5a6;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn:hover {
  background: #7f8c8d;
}
</style>
