<template>
  <v-dialog
    :model-value="modelValue"
    width="420"
    persistent
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card :class="['pa-6', 'rounded-lg', 'confirm-card', `confirm-card--${tone}`]" elevation="8">
      <div class="d-flex align-start ga-3">
        <div class="confirm-icon" :class="`confirm-icon--${tone}`">
          <v-icon :icon="resolvedIcon" size="22"></v-icon>
        </div>
        <div class="flex-grow-1">
          <h3 class="text-h6 font-weight-medium mb-1">{{ title }}</h3>
          <p class="text-body-2 text-grey-darken-1 mb-0">{{ message }}</p>
        </div>
      </div>

      <div class="d-flex justify-end ga-2 mt-6">
        <v-btn variant="text" class="text-none" :disabled="loading" @click="cancel">Cancel</v-btn>
        <v-btn
          :color="tone === 'error' ? 'error' : 'warning'"
          variant="flat"
          class="text-none font-weight-medium"
          :loading="loading"
          @click="$emit('confirm')"
        >
          {{ confirmText }}
        </v-btn>
      </div>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title: string
  message: string
  confirmText?: string
  tone?: 'error' | 'warning'
  icon?: string
  loading?: boolean
}>(), {
  confirmText: 'Delete',
  tone: 'error',
  icon: '',
  loading: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
}>()

const resolvedIcon = computed(() => {
  if (props.icon) return props.icon
  return props.tone === 'warning' ? 'mdi-alert-outline' : 'mdi-alert-circle-outline'
})

function cancel() {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.confirm-card--error {
  border-top: 4px solid rgb(var(--v-theme-error));
}

.confirm-card--warning {
  border-top: 4px solid rgb(var(--v-theme-warning));
}

.confirm-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  flex-shrink: 0;
}

.confirm-icon--error {
  background: rgba(var(--v-theme-error), 0.12);
  color: rgb(var(--v-theme-error));
}

.confirm-icon--warning {
  background: rgba(var(--v-theme-warning), 0.14);
  color: rgb(var(--v-theme-warning));
}

</style>
