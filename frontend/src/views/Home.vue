<template>
    <div>THIS IS HOME</div>
  <div class="new-module-form">
    <input v-model="code" placeholder="Module code (e.g. CSIS3734)" />
    <input v-model="name" placeholder="Module name" />
    <input v-model="color" type="color" />
    <button @click="handleCreate">Create Module</button>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { CreateModule } from '../../wailsjs/go/controllers/ModuleController'

const code = ref('')
const name = ref('')
const color = ref('#3D5A80')
const error = ref('')

async function handleCreate() {
  error.value = ''
  try {
    const newModule = await CreateModule(code.value, name.value, color.value)
    console.log('Created:', newModule)
    // e.g. push into a modules list, clear form, navigate, etc.
    code.value = ''
    name.value = ''
  } catch (err) {
    error.value = err // Go errors come through as strings here
  }
}
</script>

