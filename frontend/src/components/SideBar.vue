<template>
  <v-navigation-drawer permanent width="200" class="custom-sidebar">
    
    <nav class="nav-list">
      <div class="sidebar-title">
        COURSERWORK
      </div>
      <div class="sidebar-subtitle">
        MODULES
      </div>
      <div v-if="isLoading" class="px-3 py-1 text-caption text-grey">Loading...</div>
      <router-link
        v-for="mod in modules"
        :key="mod.id"
        :to="`/module/${mod.code}`"
        class="nav-item"
        active-class="active"
      >
        <span class="module-dot" :style="{ backgroundColor: mod.color }"></span>
        <span v-if="mod.name.trim().length == 0"  class="label">{{ mod.code }}</span>
        <span v-else class="label">{{ mod.name }}</span>
      </router-link>
    </nav>

    
    <hr>
    
    <div class="sidebar-actions d-flex flex-column ga-1">
      <v-btn
        variant="text"
        block
        prepend-icon="mdi-plus-circle-outline"
        class="new-module-btn"
        @click="dialog = true"
      >
        new module
      </v-btn>

      <v-btn
        to="/timetable"
        variant="text"
        block
        prepend-icon="mdi-table-large"
        class="new-module-btn"
      >
        timetable
      </v-btn>
    </div>
  </v-navigation-drawer>

  <v-dialog persistent v-model="dialog" width="420">
  <v-card class="pa-6 rounded-lg" elevation="4">
    <h3 class="text-h6 font-weight-medium mb-4">New module</h3>

    <v-form @submit.prevent="handleCreateModule">
      <v-text-field
        v-model="moduleData.code"
        label="Module code"
        variant="outlined"
        :rules="codeRules"
        density="comfortable"
        class="mb-2"
      ></v-text-field>

      <v-text-field
        v-model="moduleData.name"
        label="Module name"
        variant="outlined"
        density="comfortable"
        class="mb-2"
      ></v-text-field>

      <v-text-field
        v-model="moduleData.color"
        type="color"
        label="Module color"
        variant="outlined"
        density="comfortable"
        class="mb-4"
      ></v-text-field>

      <div class="d-flex gap-2 justify-end">
        <v-btn variant="text" @click="handleCancelCreateModule">
          Cancel
        </v-btn>
        <v-btn type="submit" color="color-bg-2" class="hover:" variant="flat">
          Create module
        </v-btn>
      </div>
    </v-form>
  </v-card>
</v-dialog>

  <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
    {{ alertMsg }}
  </v-snackbar>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { NavItems } from '../types';
  import { GetAllModules, CreateModule } from '../../wailsjs/go/controllers/ModuleController';
  import {models} from '../../wailsjs/go/models'

  const moduleData = ref(models.Module.createFrom({
    code: '',
    name: '',
    color: '#00000'
  }))

  const codeRules = [
  (value:string) => {
    if (value && value.trim().length > 0) return true
    return 'Module Code is required.'
  }
]
  

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
      alertMsg.value = "Failed to get available modules " + (err as string)
      alertType.value = "error"
      showAlert.value = true;
    }
  }

  async function handleCreateModule(){
    try {
      const newModule = await CreateModule(moduleData.value)
      await fetchModules()
      handleCancelCreateModule()
      alertMsg.value = "Module created successfully "
      alertType.value = "success"
      showAlert.value = true;
    } catch (err){
      alertMsg.value = "Failed to create module " + (err as string)
      alertType.value = "error"
      showAlert.value = true;
    }
  }

  function handleCancelCreateModule(){
    dialog.value = false
    moduleData.value.code = ''
    moduleData.value.name = ''
    moduleData.value.color = "#00000"
  }

  const router = useRouter()

  const timetable = () => {
    router.push('/timetable')
  }

  onMounted(()=>{
    fetchModules()
  })
  
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
  justify-content: flex-start;
  text-transform: none;
  letter-spacing: normal;
  font-size: 15px;
  font-weight: 500;
  height: 40px;
  padding: 0 16px;
  opacity: 1;
}

.new-module-btn.v-btn--active,
.new-module-btn.router-link-active {
  opacity: 1;
}
.new-module-btn.v-btn--active .v-btn__overlay {
  opacity: 0.08;
}

.new-module-btn :deep(.v-btn__prepend) {
  margin-inline: 0 12px;
}

.new-module-btn:hover {
  background-color: var(--color-primary);
  border-radius: 7px; 
}

.dialog-btns{
  flex: auto;
  flex-direction: row;
  gap: 20px;
}

</style>