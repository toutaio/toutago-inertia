<template>
  <Layout>
    <div class="posts">
      <div class="posts-header">
        <h1>Blog Posts</h1>
        <Link href="/posts/create" class="btn btn-primary">Create Post</Link>
      </div>
      
      <div class="posts-list">
        <div v-for="post in posts" :key="post.id" class="post-card">
          <div class="post-author">
            <div class="author-avatar">
              {{ post.author.name.charAt(0) }}
            </div>
            <div>
              <strong>{{ post.author.name }}</strong>
              <small>{{ formatDate(post.created) }}</small>
            </div>
          </div>
          
          <h2>{{ post.title }}</h2>
          <p>{{ post.content }}</p>
          
          <div class="post-actions">
            <Link :href="`/posts/${post.id}`" class="btn-link">Read more →</Link>
          </div>
        </div>
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
}

interface Post {
  id: number
  title: string
  content: string
  author: User
  created: string
}

interface Props {
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
.posts {
  padding: 2rem 0;
}

.posts-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin: 0;
}

.btn-primary {
  padding: 0.75rem 1.5rem;
  background: #3498db;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  font-weight: 600;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #2980b9;
}

.posts-list {
  display: grid;
  gap: 1.5rem;
}

.post-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.post-author {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #3498db;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  flex-shrink: 0;
}

.post-author div {
  display: flex;
  flex-direction: column;
}

.post-author strong {
  color: #2c3e50;
}

.post-author small {
  color: #95a5a6;
  font-size: 0.875rem;
}

.post-card h2 {
  margin: 0 0 1rem;
  color: #2c3e50;
  font-size: 1.5rem;
}

.post-card p {
  margin: 0 0 1rem;
  color: #7f8c8d;
  line-height: 1.6;
}

.post-actions {
  border-top: 1px solid #ecf0f1;
  padding-top: 1rem;
}

.btn-link {
  color: #3498db;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s;
}

.btn-link:hover {
  color: #2980b9;
}
</style>
