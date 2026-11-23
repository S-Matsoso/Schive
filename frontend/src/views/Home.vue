<template>
  <div class="dashboard-page pa-6">
    <!-- Welcome Header -->
    <v-card class="pa-6 rounded-lg mb-6 dashboard-hero" elevation="2">
      <div class="d-flex justify-space-between align-center flex-wrap ga-4">
        <div>
          <div class="text-caption text-grey-darken-1 font-weight-medium text-uppercase mb-1">
            {{ formattedToday }}
          </div>
          <h1 class="text-h4 font-weight-bold mb-1">
            {{ greeting }}, Student
          </h1>
          <p class="text-body-2 text-grey-darken-1">
            Here is your daily study overview for today.
          </p>
        </div>

        <div class="d-flex ga-2 flex-wrap">
          <v-btn
            to="/timetable"
            color="color-bg-2"
            variant="flat"
            prepend-icon="mdi-calendar-clock"
            class="text-none font-weight-medium rounded-lg"
          >
            View Timetable
          </v-btn>
          <v-btn
            variant="outlined"
            prepend-icon="mdi-refresh"
            size="small"
            class="text-none rounded-lg"
            :loading="isLoading"
            @click="fetchAllData"
          >
            Refresh
          </v-btn>
        </div>
      </div>
    </v-card>

    <!-- Stats Quick Cards -->
    <v-row class="mb-4">
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">Total Modules</div>
              <div class="text-h4 font-weight-bold">{{ modules.length }}</div>
            </div>
            <v-icon icon="mdi-book-education-outline" size="36" color="color-bg-2"></v-icon>
          </div>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-col cols="12" class="pa-0">
          <v-card class="pa-4 rounded-lg stat-card" elevation="1">
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-caption text-grey">Today's Classes</div>
                <div class="text-h4 font-weight-bold">{{ todayClasses.length }}</div>
              </div>
              <v-icon icon="mdi-calendar-today" size="36" color="primary"></v-icon>
            </div>
          </v-card>
        </v-col>
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">Pending Tasks</div>
              <div class="text-h4 font-weight-bold">{{ pendingAssignments.length }}</div>
            </div>
            <v-icon icon="mdi-clipboard-alert-outline" size="36" color="warning"></v-icon>
          </div>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4 rounded-lg stat-card" elevation="1">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-caption text-grey">Course Files</div>
              <div class="text-h4 font-weight-bold">{{ totalFilesCount }}</div>
            </div>
            <v-icon icon="mdi-file-document-multiple-outline" size="36" color="info"></v-icon>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Main Content Grid -->
    <v-row>
      <!-- LEFT COLUMN: Today's Classes & Pending Deadlines -->
      <v-col cols="12" lg="7">
        <!-- Today's Schedule -->
        <v-card class="pa-5 rounded-lg mb-6" elevation="2">
          <div class="d-flex justify-space-between align-center mb-4">
            <div class="d-flex align-center ga-2">
              <v-icon icon="mdi-clock-outline" color="color-bg-2"></v-icon>
              <h2 class="text-h6 font-weight-bold">Today's Schedule ({{ currentDayName }})</h2>
            </div>
            <div class="d-flex ga-2">
              <v-chip v-if="todayClasses.length" size="small" variant="tonal" color="color-bg-2">
                {{ todayClasses.length }} {{ todayClasses.length === 1 ? 'class' : 'classes' }}
              </v-chip>
              <v-chip v-if="todayAssignments.length" size="small" variant="tonal" color="warning">
                {{ todayAssignments.length }} due today
              </v-chip>
            </div>
          </div>

          <div v-if="todayClasses.length === 0 && todayAssignments.length === 0" class="py-8 text-center text-grey">
            <v-icon icon="mdi-calendar-check-outline" size="44" class="mb-2 text-grey-lighten-1"></v-icon>
            <div class="font-weight-medium">No classes or deadlines for today!</div>
            <div class="text-caption">Take time to catch up on revision or assignments.</div>
          </div>

          <div v-else class="d-flex flex-column ga-3">
            <!-- Class slots -->
            <v-card
              v-for="slot in todayClasses"
              :key="slot.id"
              class="pa-3 rounded-lg border d-flex align-center justify-space-between"
              variant="flat"
            >
              <div class="d-flex align-center ga-3">
                <span class="slot-color-bar" :style="{ backgroundColor: slot.module?.color || '#7A8B69' }"></span>
                <div>
                  <div class="font-weight-bold">{{ slot.module?.code || 'Module' }} — {{ slot.module?.name }}</div>
                  <div class="text-caption text-grey d-flex align-center ga-2 mt-1">
                    <span v-if="slot.room"><v-icon icon="mdi-map-marker-outline" size="14"></v-icon> {{ slot.room }}</span>
                    <span><v-icon icon="mdi-clock-time-four-outline" size="14"></v-icon> {{ formatTimeRange(slot.startHour, slot.lengthHrs) }}</span>
                  </div>
                </div>
              </div>

              <v-chip size="small" variant="tonal" :color="slot.module?.color || 'primary'">
                {{ slot.lengthHrs }}h
              </v-chip>
            </v-card>

            <!-- Assignments due today -->
            <v-card
              v-for="assign in todayAssignments"
              :key="'due-' + assign.id"
              class="pa-3 rounded-lg border d-flex align-center justify-space-between flex-wrap ga-2"
              variant="flat"
            >
              <div class="d-flex align-center ga-3">
                <span class="slot-color-bar" style="background-color: #D32F2F;"></span>
                <div>
                  <div class="d-flex align-center ga-2">
                    <v-icon icon="mdi-clipboard-alert-outline" size="18" color="error"></v-icon>
                    <span class="font-weight-bold">{{ assign.title }}</span>
                    <v-chip size="x-small" :color="getPriorityColor(assign.priority)" variant="flat" class="text-uppercase font-weight-bold">
                      {{ assign.priority }}
                    </v-chip>
                  </div>
                  <div class="text-caption text-grey mt-1">
                    <span class="font-weight-medium text-error me-2">Due Today</span>
                    <span class="text-grey-darken-2">{{ assign.module?.code }}</span>
                  </div>
                </div>
              </div>
              <v-btn
                size="small"
                variant="text"
                color="success"
                prepend-icon="mdi-check"
                class="text-none"
                @click="markAssignmentComplete(assign.id)"
              >
                Mark Done
              </v-btn>
            </v-card>
          </div>
        </v-card>

        <!-- Upcoming Deadlines & Assignments -->
        <v-card class="pa-5 rounded-lg mb-6" elevation="2">
          <div class="d-flex justify-space-between align-center mb-4">
            <div class="d-flex align-center ga-2">
              <v-icon icon="mdi-clipboard-text-clock-outline" color="warning"></v-icon>
              <h2 class="text-h6 font-weight-bold">Upcoming Deadlines</h2>
            </div>
            <div class="d-flex align-center ga-2">
              <v-chip size="small" variant="tonal" color="warning">
                {{ pendingAssignments.length }} pending
              </v-chip>
              <v-btn
                to="/tasks"
                size="small"
                variant="text"
                color="color-bg-2"
                append-icon="mdi-arrow-right"
                class="text-none font-weight-medium"
              >
                All Tasks
              </v-btn>
            </div>
          </div>

          <div v-if="pendingAssignments.length === 0" class="py-8 text-center text-grey">
            <v-icon icon="mdi-check-circle-outline" size="44" class="mb-2 text-success"></v-icon>
            <div class="font-weight-medium">All caught up!</div>
            <div class="text-caption">No pending assignments or deadlines on your plate.</div>
          </div>

          <div v-else class="d-flex flex-column ga-3">
            <v-card
              v-for="assign in pendingAssignments.slice(0, 5)"
              :key="assign.id"
              class="pa-3 rounded-lg border d-flex align-center justify-space-between flex-wrap ga-2"
              variant="flat"
            >
              <div class="flex-grow-1" style="min-width: 220px;">
                <div class="d-flex align-center ga-2">
                  <span class="font-weight-bold">{{ assign.title }}</span>
                  <v-chip size="x-small" :color="getPriorityColor(assign.priority)" variant="flat" class="text-uppercase font-weight-bold">
                    {{ assign.priority }}
                  </v-chip>
                </div>
                <div class="text-caption text-grey mt-1">
                  <span class="font-weight-medium text-grey-darken-2 me-2">{{ assign.module?.code }}</span>
                  <span v-if="assign.dueDate">Due: {{ formatDueDate(assign.dueDate) }}</span>
                </div>
                <div class="d-flex align-center ga-2 mt-2">
                  <v-progress-linear
                    :model-value="assign.progress || 0"
                    :color="getProgressColor(assign.progress)"
                    height="6"
                    rounded
                    class="flex-grow-1"
                  ></v-progress-linear>
                  <span class="text-caption text-grey font-weight-medium">{{ assign.progress || 0 }}%</span>
                </div>
              </div>

              <div class="d-flex align-center ga-2">
                <v-btn
                  size="small"
                  variant="text"
                  color="success"
                  prepend-icon="mdi-check"
                  class="text-none"
                  @click="markAssignmentComplete(assign.id)"
                >
                  Mark Done
                </v-btn>
              </div>
            </v-card>
          </div>
        </v-card>
      </v-col>

      <!-- RIGHT COLUMN: Module Quick Links & Recent Files -->
      <v-col cols="12" lg="5">
        <!-- Module Quick Access -->
        <v-card class="pa-5 rounded-lg mb-6" elevation="2">
          <div class="d-flex justify-space-between align-center mb-4">
            <h2 class="text-h6 font-weight-bold">Your Modules</h2>
            <span class="text-caption text-grey">{{ modules.length }} active</span>
          </div>

          <div class="modules-quick-grid">
            <v-card
              v-for="mod in modules"
              :key="mod.id"
              :to="`/module/${mod.code}`"
              class="pa-3 rounded-lg border module-quick-card cursor-pointer"
              variant="flat"
            >
              <div class="d-flex align-center ga-2 mb-2">
                <span class="module-dot" :style="{ backgroundColor: mod.color }"></span>
                <span class="font-weight-bold text-subtitle-2" :style="{ color: mod.color }">{{ mod.code }}</span>
              </div>
              <div class="text-caption font-weight-medium text-truncate mb-2" :title="mod.name">
                {{ mod.name }}
              </div>
              <div class="d-flex justify-space-between text-caption text-grey">
                <span>{{ mod.files?.length || 0 }} files</span>
                <span>{{ mod.assignments?.length || 0 }} tasks</span>
              </div>
            </v-card>
          </div>
        </v-card>

        <!-- Recent Documents -->
        <v-card class="pa-5 rounded-lg mb-6" elevation="2">
          <div class="d-flex justify-space-between align-center mb-4">
            <h2 class="text-h6 font-weight-bold">Recent Files</h2>
            <v-icon icon="mdi-history" color="grey"></v-icon>
          </div>

          <div v-if="recentFiles.length === 0" class="py-6 text-center text-grey">
            <v-icon icon="mdi-file-outline" size="36" class="mb-2"></v-icon>
            <div class="text-caption">No files uploaded yet.</div>
          </div>

          <div v-else class="d-flex flex-column ga-2">
            <div
              v-for="file in recentFiles.slice(0, 6)"
              :key="file.id"
              class="recent-file-item pa-2 rounded-lg d-flex align-center justify-space-between cursor-pointer"
              @click="handleOpenFile(file.id)"
            >
              <div class="d-flex align-center ga-3 overflow-hidden me-2">
                <v-icon :icon="getFileIcon(file.fileType)" size="22" color="color-bg-2"></v-icon>
                <div class="overflow-hidden">
                  <div class="text-body-2 font-weight-medium text-truncate" :title="file.name">
                    {{ file.name }}
                  </div>
                  <div class="text-caption text-grey">
                    {{ formatFileDate(file.updatedAt || file.createdAt) }}
                  </div>
                </div>
              </div>

              <v-btn
                icon="mdi-open-in-new"
                size="x-small"
                variant="text"
                color="color-bg-2"
                title="Open file"
                @click.stop="handleOpenFile(file.id)"
              ></v-btn>
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Snackbar Notifications -->
    <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
      {{ alertMsg }}
    </v-snackbar>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { models } from '../../wailsjs/go/models'
import { GetAllModules } from '../../wailsjs/go/controllers/ModuleController'
import { GetAllAssignments, UpdateAssignmentStatus } from '../../wailsjs/go/controllers/AssignmentController'
import { OpenFile } from '../../wailsjs/go/controllers/FileController'

const modules = ref<models.Module[]>([])
const allAssignments = ref<models.Assignment[]>([])
const isLoading = ref(false)

const alertMsg = ref('')
const alertType = ref('success')
const showAlert = ref(false)

const today = new Date()
const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
const currentDayName = computed(() => dayNames[today.getDay()])

// Map full day names to the abbreviations stored in ClassSlot.day
const dayAbbreviations: Record<string, string> = {
  'sunday': 'sun',
  'monday': 'mon',
  'tuesday': 'tue',
  'wednesday': 'wed',
  'thursday': 'thu',
  'friday': 'fri',
  'saturday': 'sat',
}
const currentDayAbbrev = computed(() => dayAbbreviations[currentDayName.value.toLowerCase()] || '')

const formattedToday = computed(() => {
  return today.toLocaleDateString(undefined, {
    weekday: 'long',
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
})

const greeting = computed(() => {
  const hour = today.getHours()
  if (hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
})

// Today's classes from all modules
const todayClasses = computed(() => {
  const abbrev = currentDayAbbrev.value
  const list: any[] = []

  for (const mod of modules.value) {
    if (mod.classSlots && mod.classSlots.length > 0) {
      for (const slot of mod.classSlots) {
        if (slot.day && slot.day.toLowerCase() === abbrev) {
          list.push({
            ...slot,
            module: mod
          })
        }
      }
    }
  }

  // Sort by start hour
  return list.sort((a, b) => a.startHour - b.startHour)
})

// Assignments due today
const todayAssignments = computed(() => {
  const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
  return allAssignments.value.filter(a => {
    if (!a.dueDate || a.status === 'completed') return false
    const d = new Date(a.dueDate)
    if (isNaN(d.getTime())) return false
    const dateStr = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    return dateStr === todayStr
  })
})

// Pending assignments
const pendingAssignments = computed(() => {
  return allAssignments.value.filter(a => a.status !== 'completed')
})

// Total files count across modules
const totalFilesCount = computed(() => {
  return modules.value.reduce((acc, m) => acc + (m.files?.length || 0), 0)
})

// Recent files across all modules
const recentFiles = computed(() => {
  const files: any[] = []
  for (const mod of modules.value) {
    if (mod.files && mod.files.length > 0) {
      for (const f of mod.files) {
        files.push({
          ...f,
          moduleCode: mod.code,
          moduleColor: mod.color
        })
      }
    }
  }

  return files.sort((a, b) => {
    const timeA = new Date(a.updatedAt || a.createdAt || 0).getTime()
    const timeB = new Date(b.updatedAt || b.createdAt || 0).getTime()
    return timeB - timeA
  })
})

async function fetchAllData() {
  try {
    isLoading.value = true
    const [mods, assigns] = await Promise.all([
      GetAllModules(),
      GetAllAssignments()
    ])
    modules.value = mods || []
    allAssignments.value = assigns || []
  } catch (err: any) {
    alertMsg.value = 'Failed to load dashboard data: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isLoading.value = false
  }
}

async function markAssignmentComplete(id: string) {
  try {
    await UpdateAssignmentStatus(id, 'completed')
    alertMsg.value = 'Assignment marked as completed!'
    alertType.value = 'success'
    showAlert.value = true
    await fetchAllData()
  } catch (err: any) {
    alertMsg.value = 'Failed to update assignment: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

async function handleOpenFile(fileId: string) {
  try {
    await OpenFile(fileId)
  } catch (err: any) {
    alertMsg.value = 'Failed to open file: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

function formatTimeRange(startHour: number, lengthHrs: number): string {
  const pad = (n: number) => n.toString().padStart(2, '0')
  const endHour = startHour + lengthHrs
  return `${pad(startHour)}:00 - ${pad(endHour)}:00`
}

function formatDueDate(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  const diffTime = d.getTime() - today.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return `Overdue by ${Math.abs(diffDays)}d`
  if (diffDays === 0) return 'Due today'
  if (diffDays === 1) return 'Due tomorrow'
  return `In ${diffDays} days (${d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })})`
}

function formatFileDate(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function getProgressColor(progress?: number): string {
  const p = progress || 0
  if (p >= 100) return 'success'
  if (p >= 50) return 'color-bg-2'
  if (p > 0) return 'warning'
  return 'grey'
}

function getPriorityColor(priority?: string): string {
  switch ((priority || '').toLowerCase()) {
    case 'high': return 'error'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'grey'
  }
}

function getFileIcon(type?: string): string {
  const t = (type || '').toLowerCase()
  if (t === 'pdf') return 'mdi-file-pdf-box'
  if (['doc', 'docx', 'odt'].includes(t)) return 'mdi-file-word-box'
  if (['xls', 'xlsx', 'csv'].includes(t)) return 'mdi-file-excel-box'
  if (['png', 'jpg', 'jpeg'].includes(t)) return 'mdi-file-image'
  return 'mdi-file-document-outline'
}

onMounted(() => {
  fetchAllData()
})
</script>

<style scoped>
.dashboard-page {
  width: 100%;
}

.dashboard-hero {
  border-left: 6px solid var(--color-bg-2, #7A8B69);
}

.stat-card {
  transition: transform 0.15s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.slot-color-bar {
  width: 6px;
  height: 38px;
  border-radius: 3px;
  flex-shrink: 0;
}

.module-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.modules-quick-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.module-quick-card {
  text-decoration: none;
  transition: all 0.15s ease;
}

.module-quick-card:hover {
  transform: translateY(-2px);
  border-color: var(--color-bg-2, #7A8B69);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.recent-file-item {
  transition: background-color 0.15s ease;
}

.recent-file-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.cursor-pointer {
  cursor: pointer;
}
</style>
