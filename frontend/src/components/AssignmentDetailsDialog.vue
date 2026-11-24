<template>
  <v-dialog
    :model-value="modelValue"
    width="600"
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card v-if="assignment" class="rounded-xl overflow-hidden" elevation="6">
      <v-card-title class="d-flex justify-space-between align-start pa-5 border-bottom" style="white-space: normal;">
        <div>
          <div class="text-h6 font-weight-bold mb-1">{{ assignment.title }}</div>
          <div class="d-flex align-center ga-2 flex-wrap">
            <v-chip v-if="assignment.assignmentNumber" size="x-small" variant="tonal" color="grey">
              #{{ assignment.assignmentNumber }}
            </v-chip>
            <v-chip size="x-small" :color="statusColor(assignment.status)" variant="tonal" class="text-uppercase font-weight-medium">
              {{ assignment.status }}
            </v-chip>
            <v-chip size="x-small" :color="priorityColor(assignment.priority)" variant="flat" class="text-uppercase font-weight-bold">
              {{ assignment.priority }}
            </v-chip>
            <v-chip
              v-if="assignment.dueDate"
              size="x-small"
              variant="tonal"
              :color="isOverdue(assignment.dueDate) ? 'error' : 'grey'"
              prepend-icon="mdi-calendar-clock"
            >
              Due {{ formatDueDate(assignment.dueDate) }}
            </v-chip>
          </div>
        </div>
        <v-btn icon="mdi-close" size="small" variant="text" @click="emit('update:modelValue', false)"></v-btn>
      </v-card-title>

      <v-card-text class="pa-5">
        <p v-if="assignment.description" class="text-body-2 text-grey-darken-1 mb-4 assignment-detail-text">
          {{ assignment.description }}
        </p>

        <div class="mb-2">
          <div class="d-flex justify-space-between text-caption mb-1">
            <span class="text-grey">Progress</span>
            <span class="font-weight-bold">{{ assignment.progress || 0 }}%</span>
          </div>
          <v-progress-linear
            :model-value="assignment.progress || 0"
            :color="progressColor(assignment.progress)"
            height="6"
            rounded
          ></v-progress-linear>
        </div>

        <div class="d-flex justify-space-between align-center mt-5 mb-2">
          <h4 class="text-subtitle-2 font-weight-bold">Instructions & files</h4>
          <v-btn
            size="small"
            variant="tonal"
            color="color-bg-2"
            prepend-icon="mdi-file-plus-outline"
            class="text-none"
            @click="openCreateFileDialog"
          >
            New File
          </v-btn>
        </div>

        <div v-if="assignment.files?.length" class="d-flex flex-column ga-2">
          <v-card
            v-for="file in assignment.files"
            :key="file.id"
            class="pa-3 rounded-lg d-flex align-center ga-3 cursor-pointer assignment-file-row"
            variant="outlined"
            @click="openFile(file)"
          >
            <v-icon :icon="fileInfo(file.fileType).icon" :color="fileInfo(file.fileType).color" size="30"></v-icon>
            <div class="flex-grow-1 overflow-hidden">
              <div class="text-body-2 font-weight-medium text-truncate">{{ file.name }}</div>
              <div class="text-caption text-grey">
                {{ formatFileSize(file.fileSize) }} • Modified {{ formatDateTime(file.updatedAt || file.createdAt) }}
              </div>
            </div>
            <v-btn
              icon="mdi-download-outline"
              size="small"
              variant="text"
              title="Export a copy"
              @click.stop="exportFile(file)"
            ></v-btn>
            <v-btn icon="mdi-open-in-new" size="small" variant="text" title="Open file"></v-btn>
          </v-card>
        </div>

        <v-card v-else variant="tonal" class="pa-6 text-center rounded-lg">
          <v-icon icon="mdi-paperclip" size="32" color="grey" class="mb-2"></v-icon>
          <div class="text-body-2 text-grey">No instructions attached yet.</div>
        </v-card>

        <v-alert density="compact" variant="tonal" color="info" class="text-caption mt-4">
          Files created here live inside Schive. Saving edits in Word or your default editor updates the Schive copy.
        </v-alert>
      </v-card-text>

      <v-card-actions class="pa-4 border-top">
        <v-btn v-if="showEdit" variant="text" class="text-none" @click="emit('edit', assignment)">
          Edit details
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
          v-if="assignment.status !== 'completed'"
          color="success"
          variant="flat"
          prepend-icon="mdi-check-circle"
          class="text-none font-weight-medium me-2"
          :loading="isUpdatingStatus"
          @click="toggleComplete"
        >
          Mark Complete
        </v-btn>
        <v-btn
          v-else
          color="grey"
          variant="tonal"
          prepend-icon="mdi-undo"
          class="text-none me-2"
          :loading="isUpdatingStatus"
          @click="toggleComplete"
        >
          Mark Incomplete
        </v-btn>
        <v-btn variant="text" class="text-none" @click="emit('update:modelValue', false)">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- New file dialog -->
  <v-dialog v-model="fileDialog" width="440">
    <v-card class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-1">New file</h3>
      <p class="text-body-2 text-grey-darken-1 mb-4">
        Creates the file inside Schive and opens it in your default application.
      </p>
      <v-form @submit.prevent="createFile">
        <v-select
          v-model="fileType"
          :items="fileTypes"
          item-title="title"
          item-value="value"
          label="File type"
          variant="outlined"
          density="comfortable"
          class="mb-2"
          @update:model-value="onFileTypeChange"
        ></v-select>
        <v-text-field
          v-model="fileName"
          label="File name"
          variant="outlined"
          density="comfortable"
          class="mb-2"
          :rules="[(v: string) => !!(v && v.trim()) || 'File name is required']"
        ></v-text-field>
        <p class="text-caption text-grey mb-4">
          The file is stored in this module's Assignments folder and attached to this assignment.
        </p>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" class="text-none" @click="fileDialog = false">Cancel</v-btn>
          <v-btn type="submit" color="color-bg-2" variant="flat" class="text-none" :loading="isCreating">
            Create & Open
          </v-btn>
        </div>
      </v-form>
    </v-card>
  </v-dialog>

  <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="4000" location="bottom right">
    {{ snackbarMessage }}
  </v-snackbar>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { models } from '../../wailsjs/go/models'
import { CreateAssignmentFile, OpenFile } from '../../wailsjs/go/controllers/FileController'
import { UpdateAssignmentStatus } from '../../wailsjs/go/controllers/AssignmentController'
import { exportSchiveFile } from '../utils/exportFile'
import { eventBus, SYNC_EVENTS } from '../utils/eventBus'

const props = withDefaults(defineProps<{
  modelValue: boolean
  assignment: models.Assignment | null
  showEdit?: boolean
}>(), {
  showEdit: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'edit', assignment: models.Assignment): void
  (e: 'files-changed'): void
  (e: 'status-changed', status: string): void
}>()

const isUpdatingStatus = ref(false)
const fileDialog = ref(false)
const fileType = ref('docx')
const fileName = ref('')
const isCreating = ref(false)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const fileTypes = [
  { title: 'Word document (.docx)', value: 'docx' },
  { title: 'Text file (.txt)', value: 'txt' },
  { title: 'Markdown (.md)', value: 'md' }
]

function showMessage(message: string, color: string) {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

function slug(title: string): string {
  const value = title.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
  return value || 'assignment'
}

function openCreateFileDialog() {
  if (!props.assignment) return
  fileType.value = 'docx'
  fileName.value = `${slug(props.assignment.title)}.docx`
  fileDialog.value = true
}

function onFileTypeChange(type: string) {
  const ext = type || 'docx'
  const current = fileName.value.trim()
  const base = current.replace(/\.[^.]+$/, '') || slug(props.assignment?.title || 'assignment')
  fileName.value = `${base}.${ext}`
}

async function createFile() {
  const assignment = props.assignment
  const name = fileName.value.trim()
  if (!assignment || !name) return

  try {
    isCreating.value = true
    const file = await CreateAssignmentFile(assignment.moduleId, assignment.id, name)
    fileDialog.value = false
    emit('files-changed')
    try {
      await OpenFile(file.id)
      showMessage(`${file.name} created — opening it in your default app`, 'success')
    } catch (openErr: any) {
      showMessage(`${file.name} created, but it could not be opened: ${openErr}`, 'error')
    }
  } catch (err: any) {
    showMessage('Failed to create file: ' + err, 'error')
  } finally {
    isCreating.value = false
  }
}

async function openFile(file: models.File) {
  try {
    await OpenFile(file.id)
  } catch (err: any) {
    showMessage(`Could not open ${file.name}: ${err}`, 'error')
  }
}

async function exportFile(file: models.File) {
  try {
    const destination = await exportSchiveFile(file)
    if (!destination) return
    showMessage(`${file.name} exported to ${destination}`, 'success')
  } catch (err: any) {
    showMessage(`Could not export ${file.name}: ${err}`, 'error')
  }
}

async function toggleComplete() {
  if (!props.assignment) return
  const newStatus = props.assignment.status === 'completed' ? 'pending' : 'completed'
  try {
    isUpdatingStatus.value = true
    await UpdateAssignmentStatus(props.assignment.id, newStatus)
    props.assignment.status = newStatus
    props.assignment.progress = newStatus === 'completed' ? 100 : 0
    showMessage(newStatus === 'completed' ? 'Assignment marked as complete!' : 'Assignment marked as incomplete', 'success')
    emit('files-changed')
    emit('status-changed', newStatus)
    eventBus.emit(SYNC_EVENTS.ASSIGNMENTS_CHANGED, { id: props.assignment.id, status: newStatus })
  } catch (err: any) {
    showMessage('Failed to update status: ' + err, 'error')
  } finally {
    isUpdatingStatus.value = false
  }
}

function formatFileSize(bytes: number): string {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatDateTime(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatDueDate(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  const diffTime = d.getTime() - new Date().getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  if (diffDays < 0) return `Overdue by ${Math.abs(diffDays)}d (${d.toLocaleDateString()})`
  if (diffDays === 0) return 'Due Today!'
  if (diffDays === 1) return 'Due Tomorrow'
  return `${diffDays} days left (${d.toLocaleDateString()})`
}

function isOverdue(dateStr: any): boolean {
  if (!dateStr) return false
  const d = new Date(dateStr)
  return d.getTime() < new Date().getTime()
}

function progressColor(progress?: number): string {
  const p = progress || 0
  if (p >= 100) return 'success'
  if (p >= 50) return 'color-bg-2'
  if (p > 0) return 'warning'
  return 'grey'
}

function priorityColor(priority?: string): string {
  switch ((priority || '').toLowerCase()) {
    case 'high': return 'error'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'grey'
  }
}

function statusColor(status?: string): string {
  switch ((status || '').toLowerCase()) {
    case 'completed': return 'success'
    case 'in_progress': return 'warning'
    default: return 'grey'
  }
}

function fileInfo(ext?: string): { icon: string; color: string; label: string } {
  const e = (ext || '').toLowerCase()
  switch (e) {
    case 'pdf': return { icon: 'mdi-file-pdf-box', color: '#D32F2F', label: 'PDF' }
    case 'doc':
    case 'docx': return { icon: 'mdi-file-word-box', color: '#1976D2', label: 'DOC' }
    case 'odt': return { icon: 'mdi-file-document-edit', color: '#0288D1', label: 'ODT' }
    case 'xls':
    case 'xlsx':
    case 'ods':
    case 'csv': return { icon: 'mdi-file-excel-box', color: '#388E3C', label: 'SHEET' }
    case 'ppt':
    case 'pptx':
    case 'odp': return { icon: 'mdi-file-powerpoint-box', color: '#F57C00', label: 'SLIDES' }
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'gif':
    case 'svg':
    case 'webp': return { icon: 'mdi-file-image', color: '#00897B', label: 'IMAGE' }
    case 'zip':
    case 'rar':
    case '7z': return { icon: 'mdi-folder-zip', color: '#FFA000', label: 'ARCHIVE' }
    case 'txt': return { icon: 'mdi-file-document-outline', color: '#607D8B', label: 'TEXT' }
    case 'md': return { icon: 'mdi-language-markdown', color: '#546E7A', label: 'MD' }
    default: return { icon: 'mdi-file-outline', color: '#757575', label: 'FILE' }
  }
}
</script>

<style scoped>
.assignment-detail-text {
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.assignment-file-row {
  transition: background-color 0.2s ease;
}

.assignment-file-row:hover {
  background-color: rgba(122, 139, 105, 0.08);
}
</style>
