<template>
  <div class="pa-6">
    <v-card v-if="currentModule" class="pa-6 rounded-lg" variant="tonal" elevation="2">
      <h1 class="text-h4 font-weight-bold mb-2" :style="{ color: currentModule.color }">
        {{ currentModule.code }}
      </h1>
      
      <h2 class="text-h6 text-grey-darken-1 mb-6">{{ currentModule.name }}</h2>

      <v-divider class="mb-4"></v-divider>
      
      <p>This is where assignments and files will go.</p>
    </v-card>
    
    <v-card v-else class="pa-6" variant="tonal">
      <p>Loading module data...</p>
    </v-card>
  </div>

  <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
    {{ alertMsg }}
  </v-snackbar>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { models } from '../../wailsjs/go/models' // Adjust namespace if needed
import { GetModule } from '../../wailsjs/go/controllers/ModuleController'

const route = useRoute()

const moduleCode = computed(() => route.params.code as string)

const currentModule = ref<models.Module | null>(null)

const alertMsg = ref("")
const alertType = ref("")
const showAlert = ref(false)

async function fetchModuleDetails() {
  try {
    currentModule.value = await GetModule(moduleCode.value)
    console.log("Loaded:", currentModule.value)
  } catch (err) {
    alertMsg.value = "Failed to get module details: " + err
    alertType.value = "error"
    showAlert.value = true
  }
}

onMounted(() => {
  fetchModuleDetails()
})


watch(moduleCode, (newCode, oldCode) => {
  if (newCode && newCode !== oldCode) {
    currentModule.value = null
    fetchModuleDetails()
  }
})
</script>