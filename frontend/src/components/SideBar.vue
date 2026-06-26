<template>
  <aside class="sidebar" :class="{collapsed}">
    <button class="toggle-btn" @click="collapsed = !collapsed">
      {{ collapsed ? '<-' : '->' }}
    </button>
     <nav class="nav-list">
      <router-link
        v-for="item in navItems"
        :key="item.name"
        :to="item.path"
        class="nav-item"
        active-class="active"
      >
        <span class="label" v-if="!collapsed">{{ item.label }}</span>
      </router-link>
    </nav>
  </aside>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  
  const collapsed = ref(false)

  const navItems = [
    {name: 'home', label: 'Home', path: '/'},
    {name: 'files', label: 'Files', path: '/files'}
  ]
</script>

<style>
.sidebar {
  width: 200px;
  background: var(--color-secondary);
  border-right: 1px solid var(--color-border-3);
  transition: width 0.2s ease;
  display: flex;
  flex-direction: column;
  height: 100vh;
  flex-shrink: 0;
}

.sidebar.collapsed {
  width: 60px;
}

.toggle-btn {
  align-self: flex-end;
  background: none;
  border: none;
  padding: 12px;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text-dark);
}

.nav-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 0 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 7px;
  text-decoration: none;
  color: var(--color-text-dark);
  white-space: nowrap;
  overflow: hidden;
}

.nav-item:hover {
  background: var(--color-bg-3);
}

.nav-item.active {
  background: var(--color-bg-3);
  color: var(--color-text-dark);
}


.label {
  font-size: 13px;
}
</style>