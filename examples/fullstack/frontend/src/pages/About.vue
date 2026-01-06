<template>
  <Layout>
    <div class="about">
      <h1>About Toutago</h1>
      
      <div class="about-content">
        <div class="info-card">
          <h2>Framework Information</h2>
          <dl>
            <dt>Name:</dt>
            <dd>{{ framework.name }}</dd>
            <dt>Version:</dt>
            <dd>{{ framework.version }}</dd>
          </dl>
        </div>

        <div class="features">
          <h2>Features</h2>
          <ul>
            <li>🚀 Fast and modern Go framework</li>
            <li>⚡ Server-side rendering with Inertia.js</li>
            <li>🎨 Vue 3 integration</li>
            <li>📝 TypeScript type generation</li>
            <li>🔒 Type-safe data transfer</li>
            <li>🎯 Simple and intuitive API</li>
          </ul>
        </div>

        <div class="tech-stack">
          <h2>Technology Stack</h2>
          <div class="tech-grid">
            <div class="tech-item">
              <h3>Backend</h3>
              <p>Go 1.22+</p>
              <p>Toutago Framework</p>
            </div>
            <div class="tech-item">
              <h3>Frontend</h3>
              <p>Vue 3</p>
              <p>TypeScript</p>
              <p>Vite</p>
            </div>
            <div class="tech-item">
              <h3>Bridge</h3>
              <p>Inertia.js</p>
              <p>SSR Support</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import Layout from '../components/Layout.vue'

interface Props {
  framework: {
    name: string
    version: string
  }
}

defineProps<Props>()
</script>

<style scoped>
.about {
  padding: 2rem 0;
}

h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 2rem;
}

.about-content {
  display: grid;
  gap: 2rem;
}

.info-card,
.features,
.tech-stack {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

h2 {
  color: #2c3e50;
  margin: 0 0 1.5rem;
  font-size: 1.5rem;
}

dl {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 1rem;
  margin: 0;
}

dt {
  font-weight: 600;
  color: #2c3e50;
}

dd {
  margin: 0;
  color: #7f8c8d;
}

.features ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.features li {
  padding: 0.75rem 0;
  color: #2c3e50;
  font-size: 1.1rem;
}

.features li:not(:last-child) {
  border-bottom: 1px solid #ecf0f1;
}

.tech-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.tech-item {
  text-align: center;
  padding: 1.5rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.tech-item h3 {
  margin: 0 0 1rem;
  color: #3498db;
  font-size: 1.25rem;
}

.tech-item p {
  margin: 0.5rem 0;
  color: #7f8c8d;
}
</style>
