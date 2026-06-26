<template>
  <aside class="sidebar">
     <nav class="nav-list">
      <div class="sidebar-title">
        COURSERWORK
      </div>
      <div class="sidebar-subtitle">
        MODULES
      </div>
      <router-link
        v-for="item in navItems"
        :key="item.name"
        :to="item.path"
        class="nav-item"
        active-class="active"
      >
        <span class="label">{{ item.label }}</span>
      </router-link>
    </nav>
      <hr>
    <div class="new-module-btn" @click="newModule">
      <BadgePlus/>
      <span> new module</span>
    </div>
    <div class="new-module-btn" @click="timetable">
      <Sheet />
      <span> timetable</span>
    </div>
    <div>
      <span class="label"> {{ val }}</span>
    </div>
  </aside>
</template>

<script lang="ts" setup>
  import { BadgePlus, Sheet } from 'lucide-vue-next';
import { ref } from 'vue';
  

  const navItems = [
    {name: 'home', label: 'Home', path: '/'},
    {name: 'files', label: 'Files', path: '/files'}
  ]

  const val = ref("")
  const num = ref(0)

  const newModule = () => {
    val.value = "new module " + num.value
    num.value += 1
  }

  const timetable = () => {
    val.value = "timetable " + num.value
    num.value += 1
  }

</script>

<style scoped>
.sidebar {
  width: 200px;
  background: var(--color-secondary);
  border-right: 3px solid var(--color-bg-4);
  transition: width 0.2s ease;
  display: flex;
  flex-direction: column;
  height: 100vh;
  flex-shrink: 0;
}

.sidebar-title {
  margin: 20px  0;
  color: black;
}

.sidebar{
  margin: 20px 5px;
  color: var(--color-text-light);
  font-size: small;
}



.nav-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 0 8px;
  margin-top: 10px;
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

hr {
  width: 90%;
  margin: 20px 10px;
  color: var(--color-bg-4);
}

.new-module-btn {
  display: flex;
  align-items: center;     
  cursor: pointer;
  font-family: inherit;  
  gap: 10px;
  margin-left: 0 10px;
  padding: 8px 10px;

}

.new-module-btn:hover {
  background-color: var(--color-primary);
  border-radius: 7px; 
}
</style>