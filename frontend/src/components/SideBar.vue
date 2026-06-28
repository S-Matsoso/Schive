<template>
  <v-navigation-drawer permanent width="200" class="custom-sidebar">
    
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
      <div v-if="isLoading" class="px-3 py-1 text-caption text-grey">Loading...</div>
      <router-link
        v-for="mod in modules"
        :key="mod.id"
        :to="`/module/${mod.id}`"
        class="nav-item"
        active-class="active"
      >
        <span class="module-dot" :style="{ backgroundColor: mod.color }"></span>
        <span class="label">{{ mod.code }}</span>
      </router-link>
    </nav>

    
    <hr>
    
    <v-btn variant="text" prepend-icon="mdi-plus-circle-outline" class="new-module-btn" @click="dialog = true">
    new module
    </v-btn>
    <v-btn variant="text" prepend-icon="mdi-table-large" class="new-module-btn" @click="timetable">
      timetable
    </v-btn>
    
   

  </v-navigation-drawer>

   <v-dialog v-model="dialog" width="auto">
    <v-card
      max-width="400"
      prepend-icon="mdi-update"
      text="Your application will relaunch automatically after the update is complete."
      title="Update in progress"
    >
      <template v-slot:actions>
        <v-btn class="ms-auto" text="Ok" @click="dialog = false"></v-btn>
      </template>
    </v-card>
  </v-dialog>

  <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
    {{ alertMsg }}
  </v-snackbar>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { NavItems } from '../types';
  import { GetAllModules } from '../../wailsjs/go/controllers/ModuleController';
  import {models} from '../../wailsjs/go/models'


  const dialog = ref(false)
  const navItems = ref<NavItems[]|null>(null)
  const alertMsg = ref("")
  const alertType = ref("")
  const showAlert = ref(false);

  const modules = ref<models.Module[]|null>([])
  const isLoading = ref(true)
  
  async function fetchModules(){
    try{
      isLoading.value = true;
      modules.value = await GetAllModules();
      isLoading.value = false

    } catch (err) {
      alertMsg.value = "Module created successfully " + (err as string)
      alertType.value = "error"
      showAlert.value = true;
    }
  }

  const timetable = () => {
    console.log("12")
  }

  onMounted(()=>{
    fetchModules()
  })
  
  console.log("Modules ", modules)
</script>

<style scoped>
.custom-sidebar {
  background-color: var(--color-secondary) !important;
  border-right: 3px solid var(--color-bg-4) !important;
  color: var(--color-text-light);
}

.sidebar-title {
  margin: 20px 0;
  color: black;
}

.sidebar-subtitle {
  margin: 20px 5px;
  color: var(--color-text-light);
  font-size: small;
}



.module-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
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
  margin: 20px auto;
  border: none;
  border-top: 1px solid var(--color-bg-4);
}

.new-module-btn {
  display: flex;
  align-items: center;     
  cursor: pointer;
  font-family: inherit;  
  gap: 10px;
  margin: 0 10px; 
  padding: 8px 10px;
}

.new-module-btn:hover {
  background-color: var(--color-primary);
  border-radius: 7px; 
}
</style>