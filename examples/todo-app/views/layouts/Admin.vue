<template>
  <div class="admin-layout">
    <div class="admin-sidebar">
      <aside class="menu p-4">
        <p class="menu-label">Admin</p>
        <ul class="menu-list">
          <li><Link href="/admin/dashboard">Dashboard</Link></li>
          <li><Link href="/admin/users">Users</Link></li>
          <li><Link href="/admin/settings">Settings</Link></li>
        </ul>
      </aside>
    </div>
    <div class="admin-content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: calc(100vh - 52px);
}

.admin-sidebar {
  width: 250px;
  background-color: #f5f5f5;
  border-right: 1px solid #dbdbdb;
}

.admin-content {
  flex: 1;
  padding: 2rem;
}

.p-4 {
  padding: 1.5rem;
}
</style>
