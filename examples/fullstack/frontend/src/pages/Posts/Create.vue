<template>
  <Layout>
    <div class="create-post">
      <h1>Create New Post</h1>
      
      <form @submit.prevent="submit" class="post-form">
        <div class="form-group">
          <label for="title">Title</label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            class="form-control"
            :class="{ 'is-invalid': form.errors.title }"
          />
          <div v-if="form.errors.title" class="error">
            {{ form.errors.title }}
          </div>
        </div>

        <div class="form-group">
          <label for="content">Content</label>
          <textarea
            id="content"
            v-model="form.content"
            rows="10"
            class="form-control"
            :class="{ 'is-invalid': form.errors.content }"
          ></textarea>
          <div v-if="form.errors.content" class="error">
            {{ form.errors.content }}
          </div>
        </div>

        <div class="form-actions">
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="form.processing"
          >
            {{ form.processing ? 'Creating...' : 'Create Post' }}
          </button>
          <Link href="/posts" class="btn btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { Link, useForm } from '@toutaio/inertia-vue'
import Layout from '../../components/Layout.vue'

interface Props {
  users: any[]
  errors?: Record<string, string>
  old?: {
    title: string
    content: string
  }
}

const props = defineProps<Props>()

const form = useForm({
  title: props.old?.title || '',
  content: props.old?.content || '',
})

function submit() {
  form.post('/posts')
}
</script>

<style scoped>
.create-post {
  padding: 2rem 0;
  max-width: 800px;
  margin: 0 auto;
}

h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 2rem;
}

.post-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2c3e50;
  font-weight: 600;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: #3498db;
}

.form-control.is-invalid {
  border-color: #e74c3c;
}

.error {
  color: #e74c3c;
  font-size: 0.875rem;
  margin-top: 0.5rem;
}

textarea.form-control {
  resize: vertical;
  min-height: 200px;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: all 0.2s;
}

.btn-primary {
  background: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2980b9;
}

.btn-primary:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #95a5a6;
  color: white;
}

.btn-secondary:hover {
  background: #7f8c8d;
}
</style>
