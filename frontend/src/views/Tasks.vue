<template>
  <div class="tasks-page pa-6">
    <!-- Header -->
    <v-card class="pa-6 rounded-lg mb-6 tasks-hero" elevation="2">
      <div class="d-flex justify-space-between align-center flex-wrap ga-4">
        <div>
          <h1 class="text-h4 font-weight-bold mb-1">Assignments</h1>
          <p class="text-body-2 text-grey-darken-1">
            Every assignment across your modules, with live progress tracking. Module to-dos live in each module's Tasks tab.
          </p>
        </div>

        <div class="d-flex ga-2 flex-wrap">
          <v-btn
            variant="outlined"
            prepend-icon="mdi-refresh"
            size="small"
            class="text-none rounded-lg"
            :loading="isLoading"
            @click="fetchData"
          >
            Refresh
          </v-btn>
          <v-btn
            color="color-bg-2"
            variant="flat"
            prepend-icon="mdi-plus"
            class="text-none font-weight-medium rounded-lg"
            @click="openCreateDialog"
          >
            New Assignment
          </v-btn>
        </div>
      </div>
    </v-card>

    <!-- Stats -->
    <v-row class="mb-4">
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">All Assignments</div>
              <div class="text-h4 font-weight-bold">{{ assignments.length }}</div>
            </div>
            <v-icon icon="mdi-clipboard-text-outline" size="36" color="color-bg-2"></v-icon>
          </div>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">In Progress</div>
              <div class="text-h4 font-weight-bold">{{ inProgressCount }}</div>
            </div>
            <v-icon icon="mdi-progress-clock" size="36" color="primary"></v-icon>
          </div>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">Completed</div>
              <div class="text-h4 font-weight-bold">{{ completedCount }}</div>
            </div>
            <v-icon icon="mdi-check-circle-outline" size="36" color="success"></v-icon>
          </div>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">Overdue</div>
              <div class="text-h4 font-weight-bold">{{ overdueCount }}</div>
            </div>
            <v-icon icon="mdi-alert-circle-outline" size="36" color="error"></v-icon>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Filters -->
    <v-card class="pa-4 rounded-lg mb-4" elevation="1">
      <div class="d-flex flex-wrap align-center justify-space-between ga-3">
        <v-btn-toggle
          v-model="statusFilter"
          mandatory
          density="compact"
          rounded="lg"
          variant="outlined"
          class="filter-toggle"
        >
          <v-btn value="all" size="small" class="text-none">All ({{ assignments.length }})</v-btn>
          <v-btn value="pending" size="small" class="text-none">Pending</v-btn>
          <v-btn value="in_progress" size="small" class="text-none">In Progress</v-btn>
          <v-btn value="completed" size="small" class="text-none">Completed</v-btn>
        </v-btn-toggle>

        <div class="d-flex align-center ga-3 flex-wrap">
          <v-select
            v-model="moduleFilter"
            :items="moduleFilterOptions"
            item-title="title"
            item-value="id"
            density="compact"
            variant="outlined"
            rounded="lg"
            hide-details
            style="width: 210px;"
          ></v-select>
          <v-text-field
            v-model="searchQuery"
            prepend-inner-icon="mdi-magnify"
            placeholder="Search assignments..."
            density="compact"
            variant="outlined"
            rounded="lg"
            hide-details
            clearable
            style="width: 220px;"
          ></v-text-field>
        </div>
      </div>
    </v-card>

    <!-- Loading -->
    <div v-if="isLoading" class="py-12 text-center">
      <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
      <p class="mt-3 text-caption text-grey">Loading assignments...</p>
    </div>

    <!-- Empty -->
    <v-card v-else-if="filteredAssignments.length === 0" class="pa-12 text-center rounded-lg" variant="tonal">
      <v-icon icon="mdi-clipboard-check-outline" size="56" color="grey" class="mb-3"></v-icon>
      <h3 class="text-h6 font-weight-medium mb-1">
        {{ assignments.length === 0 ? 'No assignments yet' : 'No assignments match these filters' }}
      </h3>
      <p class="text-body-2 text-grey mb-4">
        {{ assignments.length === 0 ? 'Create an assignment to start tracking your coursework progress.' : 'Try a different status, module, or search term.' }}
      </p>
      <v-btn
        v-if="assignments.length === 0"
        color="color-bg-2"
        variant="flat"
        prepend-icon="mdi-plus"
        class="text-none rounded-lg"
        @click="openCreateDialog"
      >
        Create Assignment
      </v-btn>
    </v-card>

    <!-- Task list -->
    <div v-else class="d-flex flex-column ga-3">
      <v-card
        v-for="assign in filteredAssignments"
        :key="assign.id"
        class="pa-4 rounded-lg task-card"
        elevation="1"
      >
        <div class="d-flex justify-space-between align-start flex-wrap ga-3">
          <div class="d-flex align-start ga-3 flex-grow-1 task-main">
            <v-btn
              :icon="assign.status === 'completed' ? 'mdi-checkbox-marked-circle' : 'mdi-checkbox-blank-circle-outline'"
              :color="assign.status === 'completed' ? 'success' : 'grey'"
              variant="text"
              density="compact"
              size="small"
              class="mt-1"
              :title="assign.status === 'completed' ? 'Mark Incomplete' : 'Mark Complete'"
              @click.stop="setStatus(assign.id, assign.status === 'completed' ? 'pending' : 'completed')"
            ></v-btn>
            <div class="flex-grow-1">
              <div class="d-flex align-center ga-2 flex-wrap mb-1">
                <span
                  class="module-dot"
                  :style="{ backgroundColor: assign.module?.color || '#7A8B69' }"
                ></span>
                <span class="text-caption font-weight-bold" :style="{ color: assign.module?.color || '#7A8B69' }">
                  {{ assign.module?.code || 'Module' }}
                </span>
                <span
                  class="font-weight-bold text-subtitle-1"
                  :class="{ 'text-decoration-line-through text-grey': assign.status === 'completed' }"
                >{{ assign.title }}</span>
                <v-chip v-if="assign.assignmentNumber" size="x-small" variant="tonal" color="grey">
                  #{{ assign.assignmentNumber }}
                </v-chip>
                <v-chip size="x-small" :color="getPriorityColor(assign.priority)" variant="flat" class="text-uppercase font-weight-bold">
                  {{ assign.priority }}
                </v-chip>
                <v-chip size="x-small" :color="getStatusColor(assign.status)" variant="tonal" class="text-uppercase font-weight-medium">
                  {{ formatStatus(assign.status) }}
                </v-chip>
                <v-chip v-if="isOverdue(assign)" size="x-small" color="error" variant="flat" class="font-weight-bold">
                  Overdue
                </v-chip>
              </div>

              <p v-if="assign.description" class="text-body-2 text-grey-darken-1 mt-1 mb-0">
                {{ assign.description }}
              </p>

              <div class="text-caption text-grey mt-2 d-flex align-center flex-wrap ga-3">
                <span v-if="assign.dueDate" class="d-flex align-center ga-1 font-weight-medium" :class="isOverdue(assign) ? 'text-error' : 'text-primary'">
                  <v-icon icon="mdi-calendar-clock" size="14"></v-icon>
                  Due: {{ formatDueDate(assign.dueDate) }}
                </span>
                <span v-else class="d-flex align-center ga-1">
                  <v-icon icon="mdi-calendar-remove-outline" size="14"></v-icon>
                  No due date
                </span>
                <span>Updated {{ formatDate(assign.updatedAt) }}</span>
              </div>
            </div>
          </div>

          <div class="d-flex align-center ga-2">
            <v-btn
              size="small"
              variant="flat"
              color="color-bg-2"
              prepend-icon="mdi-open-in-new"
              class="text-none"
              @click="openAssignmentDetails(assign)"
            >
              Open
            </v-btn>
            <v-btn
              size="small"
              variant="tonal"
              color="color-bg-2"
              prepend-icon="mdi-progress-pencil"
              class="text-none"
              @click="openProgressDialog(assign)"
            >
              Progress
            </v-btn>
            <v-btn
              v-if="assign.status !== 'completed'"
              size="small"
              variant="tonal"
              color="success"
              prepend-icon="mdi-check-circle-outline"
              class="text-none"
              @click="setStatus(assign.id, 'completed')"
            >
              Mark Complete
            </v-btn>
            <v-btn
              v-else
              size="small"
              variant="text"
              color="grey"
              prepend-icon="mdi-undo"
              class="text-none"
              @click="setStatus(assign.id, 'pending')"
            >
              Mark Incomplete
            </v-btn>
            <v-btn
              icon="mdi-pencil-outline"
              size="small"
              variant="text"
              title="Edit assignment"
              @click="openEditDialog(assign)"
            ></v-btn>
            <v-btn
              icon="mdi-delete-outline"
              size="small"
              variant="text"
              color="error"
              title="Delete assignment"
              @click="confirmDelete(assign)"
            ></v-btn>
          </div>
        </div>

        <div class="mt-3">
          <div class="d-flex justify-space-between align-center text-caption mb-1">
            <span class="text-grey">Progress</span>
            <span class="font-weight-bold">{{ assign.progress || 0 }}%</span>
          </div>
          <v-progress-linear
            :model-value="assign.progress || 0"
            :color="getProgressColor(assign.progress)"
            height="8"
            rounded
          ></v-progress-linear>
        </div>
      </v-card>
    </div>

    <!-- Update Progress Dialog -->
    <v-dialog v-model="progressDialog" width="440">
      <v-card v-if="progressTask" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-1">Update Progress</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4 text-truncate">
          {{ progressTask.title }}
        </p>

        <div class="d-flex justify-space-between align-center mb-2">
          <span class="text-caption text-grey font-weight-medium">Completion</span>
          <span class="text-h6 font-weight-bold">{{ progressValue }}%</span>
        </div>
        <v-slider
          v-model="progressValue"
          :min="0"
          :max="100"
          :step="5"
          color="color-bg-2"
          track-color="grey-lighten-3"
          thumb-label
          class="mb-2"
        ></v-slider>

        <div class="d-flex ga-2 mb-5 flex-wrap">
          <v-chip
            v-for="preset in progressPresets"
            :key="preset"
            size="small"
            variant="tonal"
            :color="progressValue === preset ? 'color-bg-2' : 'grey'"
            class="cursor-pointer font-weight-medium"
            @click="progressValue = preset"
          >
            {{ preset }}%
          </v-chip>
        </div>

        <v-alert density="compact" variant="tonal" color="info" class="text-caption mb-4">
          Progress above 0% marks the assignment as in progress, and 100% marks it completed.
        </v-alert>

        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="progressDialog = false">Cancel</v-btn>
          <v-btn color="color-bg-2" variant="flat" :loading="isSavingProgress" @click="saveProgress">
            Save Progress
          </v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Create / Edit Task Dialog -->
    <v-dialog v-model="taskDialog" width="500">
      <v-card class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-3">{{ editingTask ? 'Edit Assignment' : 'New Assignment' }}</h3>
        <v-form @submit.prevent="saveTask">
          <v-select
            v-model="taskForm.moduleId"
            :items="moduleOptions"
            item-title="title"
            item-value="id"
            label="Module"
            variant="outlined"
            density="comfortable"
            :disabled="!!editingTask"
            :rules="[(v: string) => !!v || 'Module is required']"
            class="mb-2"
            required
          ></v-select>
          <v-text-field
            v-model="taskForm.title"
            label="Title"
            variant="outlined"
            density="comfortable"
            :rules="[(v: string) => !!(v && v.trim()) || 'Title is required']"
            class="mb-2"
            required
          ></v-text-field>
          <v-textarea
            v-model="taskForm.description"
            label="Description"
            variant="outlined"
            density="comfortable"
            rows="2"
            class="mb-2"
          ></v-textarea>
          <v-row>
            <v-col cols="6">
              <v-select
                v-model="taskForm.priority"
                label="Priority"
                :items="['low', 'medium', 'high']"
                variant="outlined"
                density="comfortable"
              ></v-select>
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="taskForm.status"
                label="Status"
                :items="statusOptions"
                variant="outlined"
                density="comfortable"
              ></v-select>
            </v-col>
          </v-row>
          <v-text-field
            v-model="taskForm.dueDate"
            type="date"
            label="Due Date"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-text-field>

          <div class="d-flex justify-end ga-2">
            <v-btn variant="text" @click="taskDialog = false">Cancel</v-btn>
            <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isSavingTask">
              {{ editingTask ? 'Save Changes' : 'Create Assignment' }}
            </v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- Assignment details / new file -->
    <AssignmentDetailsDialog
      v-model="assignmentDetailDialog"
      :assignment="detailAssignment"
      @edit="openEditDialog"
      @files-changed="handleAssignmentFilesChanged"
      @status-changed="fetchData"
    />

    <!-- Delete confirmation -->
    <ConfirmDialog
      v-model="confirmDialog"
      title="Delete assignment?"
      :message="confirmMessage"
      confirm-text="Delete assignment"
      tone="error"
      icon="mdi-clipboard-remove-outline"
      :loading="confirmLoading"
      @confirm="runConfirmedDelete"
    />

    <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">

      {{ alertMsg }}
    </v-snackbar>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { models } from '../../wailsjs/go/models'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import AssignmentDetailsDialog from '../components/AssignmentDetailsDialog.vue'
import { GetAllModules } from '../../wailsjs/go/controllers/ModuleController'
import {
  CreateAssignment,
  GetAllAssignments,
  UpdateAssignment,
  UpdateAssignmentProgress,
  UpdateAssignmentStatus,
  DeleteAssignment
} from '../../wailsjs/go/controllers/AssignmentController'
import { eventBus, SYNC_EVENTS } from '../utils/eventBus'

const assignments = ref<models.Assignment[]>([])
const modules = ref<models.Module[]>([])
const isLoading = ref(false)

const statusFilter = ref('all')
const moduleFilter = ref('all')
const searchQuery = ref('')

const alertMsg = ref('')
const alertType = ref('success')
const showAlert = ref(false)

// Assignment details
const assignmentDetailDialog = ref(false)
const detailAssignment = ref<models.Assignment | null>(null)

// Delete confirmation
const confirmDialog = ref(false)
const confirmMessage = ref('')
const confirmLoading = ref(false)
const assignmentToDelete = ref<models.Assignment | null>(null)

const statusOptions = [
  { title: 'Pending', value: 'pending' },
  { title: 'In Progress', value: 'in_progress' },
  { title: 'Completed', value: 'completed' }
]

const progressPresets = [0, 25, 50, 75, 100]

// Update progress dialog
const progressDialog = ref(false)
const progressTask = ref<models.Assignment | null>(null)
const progressValue = ref(0)
const isSavingProgress = ref(false)

// Create / edit dialog
const taskDialog = ref(false)
const isSavingTask = ref(false)
const editingTask = ref<models.Assignment | null>(null)
const taskForm = ref({
  moduleId: '',
  title: '',
  description: '',
  priority: 'medium',
  status: 'pending',
  dueDate: ''
})

const moduleOptions = computed(() =>
  modules.value.map(m => ({ title: `${m.code} — ${m.name}`, id: m.id }))
)

const moduleFilterOptions = computed(() => [
  { title: 'All modules', id: 'all' },
  ...modules.value.map(m => ({ title: m.code, id: m.id }))
])

const inProgressCount = computed(() => assignments.value.filter(a => a.status === 'in_progress').length)
const completedCount = computed(() => assignments.value.filter(a => a.status === 'completed').length)
const overdueCount = computed(() => assignments.value.filter(a => isOverdue(a)).length)

const filteredAssignments = computed(() => {
  let list = assignments.value
  if (statusFilter.value !== 'all') {
    list = list.filter(a => a.status === statusFilter.value)
  }
  if (moduleFilter.value !== 'all') {
    list = list.filter(a => a.moduleId === moduleFilter.value)
  }
  const q = searchQuery.value?.trim().toLowerCase()
  if (q) {
    list = list.filter(a =>
      a.title.toLowerCase().includes(q) ||
      (a.description || '').toLowerCase().includes(q) ||
      (a.module?.code || '').toLowerCase().includes(q)
    )
  }
  return list
})

async function fetchData() {
  try {
    isLoading.value = true
    const [assigns, mods] = await Promise.all([GetAllAssignments(), GetAllModules()])
    assignments.value = assigns || []
    modules.value = mods || []
  } catch (err: any) {
    showMessage('Failed to load assignments: ' + err, 'error')
  } finally {
    isLoading.value = false
  }
}

function showMessage(msg: string, type: string) {
  alertMsg.value = msg
  alertType.value = type
  showAlert.value = true
}

function openCreateDialog() {
  editingTask.value = null
  taskForm.value = {
    moduleId: modules.value[0]?.id || '',
    title: '',
    description: '',
    priority: 'medium',
    status: 'pending',
    dueDate: ''
  }
  taskDialog.value = true
}

function openAssignmentDetails(assign: models.Assignment) {
  detailAssignment.value = assign
  assignmentDetailDialog.value = true
}

async function handleAssignmentFilesChanged() {
  const id = detailAssignment.value?.id
  await fetchData()
  if (id) {
    detailAssignment.value = assignments.value.find(a => a.id === id) || detailAssignment.value
  }
}

function openEditDialog(assign: models.Assignment) {
  assignmentDetailDialog.value = false
  editingTask.value = assign
  let formattedDate = ''
  if (assign.dueDate) {
    const d = new Date(assign.dueDate)
    if (!isNaN(d.getTime())) {
      formattedDate = d.toISOString().split('T')[0]
    }
  }
  taskForm.value = {
    moduleId: assign.moduleId,
    title: assign.title,
    description: assign.description || '',
    priority: assign.priority || 'medium',
    status: assign.status || 'pending',
    dueDate: formattedDate
  }
  taskDialog.value = true
}

async function saveTask() {
  const title = taskForm.value.title.trim()
  if (!title) return
  if (!editingTask.value && !taskForm.value.moduleId) return

  try {
    isSavingTask.value = true
    if (editingTask.value) {
      await UpdateAssignment(
        editingTask.value.id,
        title,
        taskForm.value.description,
        taskForm.value.status,
        taskForm.value.priority,
        taskForm.value.dueDate
      )
      showMessage('Assignment updated!', 'success')
    } else {
      const moduleTasks = assignments.value.filter(a => a.moduleId === taskForm.value.moduleId)
      await CreateAssignment(
        taskForm.value.moduleId,
        title,
        taskForm.value.description,
        moduleTasks.length + 1,
        taskForm.value.status,
        taskForm.value.priority,
        taskForm.value.dueDate
      )
      showMessage('Assignment created!', 'success')
    }
    taskDialog.value = false
    await fetchData()
  } catch (err: any) {
    showMessage('Failed to save assignment: ' + err, 'error')
  } finally {
    isSavingTask.value = false
  }
}

function openProgressDialog(assign: models.Assignment) {
  progressTask.value = assign
  progressValue.value = assign.progress || 0
  progressDialog.value = true
}

async function saveProgress() {
  if (!progressTask.value) return
  try {
    isSavingProgress.value = true
    await UpdateAssignmentProgress(progressTask.value.id, progressValue.value)
    progressDialog.value = false
    showMessage(`Progress updated to ${progressValue.value}%`, 'success')
    await fetchData()
  } catch (err: any) {
    showMessage('Failed to update progress: ' + err, 'error')
  } finally {
    isSavingProgress.value = false
  }
}

async function setStatus(id: string, status: string) {
  try {
    await UpdateAssignmentStatus(id, status)
    await fetchData()
    eventBus.emit(SYNC_EVENTS.ASSIGNMENTS_CHANGED, { id, status })
  } catch (err: any) {
    showMessage('Failed to update status: ' + err, 'error')
  }
}

function confirmDelete(assign: models.Assignment) {
  assignmentToDelete.value = assign
  confirmMessage.value = `"${assign.title}" will be permanently removed, including its uploaded instructions.`
  confirmDialog.value = true
}

async function runConfirmedDelete() {
  if (!assignmentToDelete.value) return
  confirmLoading.value = true
  try {
    await DeleteAssignment(assignmentToDelete.value.id)
    confirmDialog.value = false
    assignmentToDelete.value = null
    showMessage('Assignment deleted', 'success')
    await fetchData()
  } catch (err: any) {
    confirmDialog.value = false
    showMessage('Failed to delete assignment: ' + err, 'error')
  } finally {
    confirmLoading.value = false
  }
}

function isOverdue(assign: models.Assignment): boolean {
  if (!assign.dueDate || assign.status === 'completed') return false
  const due = new Date(assign.dueDate)
  return !isNaN(due.getTime()) && due.getTime() < Date.now()
}

function formatDueDate(dateStr: any): string {
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatDate(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

function formatStatus(status?: string): string {
  if (status === 'in_progress') return 'In Progress'
  if (status === 'completed') return 'Completed'
  return 'Pending'
}

function getStatusColor(status?: string): string {
  switch (status) {
    case 'completed': return 'success'
    case 'in_progress': return 'primary'
    case 'pending': return 'warning'
    default: return 'grey'
  }
}

function getPriorityColor(priority?: string): string {
  switch ((priority || '').toLowerCase()) {
    case 'high': return 'error'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'grey'
  }
}

function getProgressColor(progress?: number): string {
  const p = progress || 0
  if (p >= 100) return 'success'
  if (p >= 50) return 'color-bg-2'
  if (p > 0) return 'warning'
  return 'grey'
}

let unsubAssignments: (() => void) | null = null

onMounted(() => {
  fetchData()
  unsubAssignments = eventBus.on(SYNC_EVENTS.ASSIGNMENTS_CHANGED, () => {
    fetchData()
  })
})

onUnmounted(() => {
  if (unsubAssignments) unsubAssignments()
})
</script>

<style scoped>
.tasks-page {
  width: 100%;
}

.tasks-hero {
  border-left: 6px solid var(--color-bg-2, #7A8B69);
}

.stat-card {
  transition: transform 0.15s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.task-card {
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.task-card:hover {
  border-color: var(--color-bg-2, #7A8B69);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.task-main {
  flex: 1 1 420px;
  min-width: 0;
}

.module-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
