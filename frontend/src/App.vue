
<template>
  <div class="app-layout">
    <SideBar />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>


<script setup>
  import SideBar from './components/SideBar.vue';
</script>

<style scoped>
.app-layout {
  color: var(--color-text-dark);
  display: flex;
  height: 100vh;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
  padding: 24px;
}

</style>
