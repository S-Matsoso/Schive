<template>
  <v-navigation-drawer permanent width="230" class="custom-sidebar">
    <div>
      <!-- App Brand -->
      <div class="sidebar-header px-4 pt-4 pb-2 d-flex align-center justify-space-between">
        <div class="d-flex align-center ga-2 cursor-pointer" @click="router.push('/')">
          <v-icon icon="mdi-notebook-outline" color="color-bg-2" size="24"></v-icon>
          <span class="sidebar-brand-text">SCHIVE</span>
        </div>

        <!-- Global Search Icon Button -->
        <v-btn
          icon="mdi-magnify"
          size="small"
          variant="text"
          title="Search (Ctrl+K)"
          @click="openSearchDialog"
        ></v-btn>
      </div>

      <!-- Main Navigation Items -->
      <div class="px-2 py-2 d-flex flex-column ga-1">
        <router-link to="/" class="nav-item" exact-active-class="active">
          <v-icon icon="mdi-view-dashboard-outline" size="20"></v-icon>
          <span class="label">Dashboard</span>
        </router-link>

        <router-link to="/timetable" class="nav-item" active-class="active">
          <v-icon icon="mdi-table-large" size="20"></v-icon>
          <span class="label">Timetable</span>
        </router-link>

        <router-link to="/tasks" class="nav-item" active-class="active">
          <v-icon icon="mdi-clipboard-check-outline" size="20"></v-icon>
          <span class="label">Assignments</span>
        </router-link>

        <div class="nav-item cursor-pointer" @click="openSearchDialog">
          <v-icon icon="mdi-magnify" size="20"></v-icon>
          <span class="label">Search</span>
          <span class="search-kbd-pill ms-auto">Ctrl+K</span>
        </div>
      </div>

      <hr class="sidebar-divider my-2" />

      <!-- Tasks Section (Collapsible - sorted by deadline then importance) -->
      <div class="sidebar-section px-2 mb-2">
        <div
          class="d-flex justify-space-between align-center px-2 py-1 mb-1 cursor-pointer sidebar-section-header rounded-lg"
          @click="toggleTasksCollapse"
        >
          <div class="d-flex align-center ga-1">
            <v-icon
              :icon="isTasksCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-down'"
              size="15"
              color="grey"
            ></v-icon>
            <span class="sidebar-subtitle">TASKS</span>
          </div>
          <v-chip size="x-small" variant="flat" color="color-bg-3" class="font-weight-bold">
            {{ sortedAllTasks.length }}
          </v-chip>
        </div>

        <div v-show="!isTasksCollapsed" class="tasks-sidebar-list d-flex flex-column ga-1">
          <div v-if="sortedAllTasks.length === 0" class="text-caption text-grey px-2 py-1 font-italic">
            No tasks found
          </div>
          <div
            v-for="task in sortedAllTasks"
            :key="task.id"
            class="sidebar-task-item d-flex align-center justify-space-between px-2 py-1 rounded-lg cursor-pointer"
            :class="{ 'task-completed': task.completed }"
            :title="`${task.title} • Priority: ${task.priority || 'medium'}${task.dueDate ? ' • Due: ' + formatDeadlineFull(task.dueDate) : ''}`"
            @click="navigateToTask(task)"
          >
            <div class="d-flex align-center ga-2 overflow-hidden flex-grow-1 me-1">
              <v-btn
                :icon="task.completed ? 'mdi-checkbox-marked-circle' : 'mdi-checkbox-blank-circle-outline'"
                :color="task.completed ? 'success' : 'grey'"
                variant="text"
                density="compact"
                size="small"
                class="pa-0 flex-shrink-0"
                :title="task.completed ? 'Mark Incomplete' : 'Mark Complete'"
                @click.stop="toggleTaskComplete(task)"
              ></v-btn>
              <span
                class="task-priority-dot"
                :style="{ backgroundColor: getPriorityColor(task.priority) }"
              ></span>
              <span
                class="task-title text-truncate"
                :class="{ 'text-decoration-line-through text-grey': task.completed }"
              >
                {{ task.title }}
              </span>
            </div>

            <div class="d-flex align-center ga-1 flex-shrink-0">
              <span
                v-if="task.dueDate"
                class="task-deadline-badge"
                :class="getDeadlineClass(task.dueDate)"
              >
                {{ formatShortDeadline(task.dueDate) }}
              </span>
              <span
                v-if="task.moduleCode"
                class="task-module-pill"
                :style="{ color: task.moduleColor || 'inherit' }"
              >
                {{ task.moduleCode }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <hr class="sidebar-divider my-2" />

      <!-- Modules Section -->
      <nav class="nav-list px-2">
        <div class="d-flex justify-space-between align-center px-2 py-1 mb-1">
          <span class="sidebar-subtitle">MODULES</span>
          <div class="d-flex align-center">
            <v-btn
              icon="mdi-folder-plus"
              size="x-small"
              variant="text"
              color="color-bg-2"
              title="New folder"
              @click="openCreateFolderDialog"
            ></v-btn>
            <v-btn
              icon="mdi-plus"
              size="x-small"
              variant="text"
              color="color-bg-2"
              title="Add module"
              @click="dialog = true"
            ></v-btn>
          </div>
        </div>

        <div v-if="isLoading" class="px-3 py-1 text-caption text-grey">Loading modules...</div>

        <template v-else>
          <div v-for="group in moduleGroups" :key="group.id || 'unfiled'" class="module-group">
            <!-- Folder header -->
            <div
              v-if="group.id || folders.length"
              class="folder-row d-flex align-center ga-1 rounded-lg px-2 py-1 cursor-pointer"
              @click="group.id ? toggleFolder(group.id) : undefined"
            >
              <v-icon
                :icon="group.id && collapsedFolders[group.id] ? 'mdi-chevron-right' : 'mdi-chevron-down'"
                size="14"
              ></v-icon>
              <v-icon :icon="group.id ? 'mdi-folder-outline' : 'mdi-folder-remove-outline'" size="15"></v-icon>
              <span class="folder-label text-truncate" :title="group.name">{{ group.name }}</span>
              <span class="folder-count ms-auto">{{ group.modules.length }}</span>

              <v-menu v-if="group.id" location="end">
                <template #activator="{ props }">
                  <v-btn
                    v-bind="props"
                    icon="mdi-dots-vertical"
                    size="x-small"
                    variant="text"
                    class="module-more-btn"
                    @click.stop
                  ></v-btn>
                </template>
                <v-list density="compact" class="rounded-lg" elevation="3">
                  <v-list-item
                    prepend-icon="mdi-pencil-outline"
                    title="Rename Folder"
                    @click="openRenameFolderDialog(group)"
                  ></v-list-item>
                  <v-list-item
                    prepend-icon="mdi-delete-outline"
                    title="Delete Folder"
                    class="text-error"
                    @click="confirmDeleteFolder(group)"
                  ></v-list-item>
                </v-list>
              </v-menu>
            </div>

            <div v-show="!group.id || !collapsedFolders[group.id]" :class="{ 'folder-module-list': !!group.id }">
              <div
                v-for="mod in group.modules"
                :key="mod.id"
                class="module-row d-flex align-center justify-space-between rounded-lg"
              >
                <router-link
                  :to="`/module/${mod.code}`"
                  class="nav-item flex-grow-1 overflow-hidden"
                  active-class="active"
                >
                  <span class="module-dot" :style="{ backgroundColor: mod.color }"></span>
                  <span class="label text-truncate" :title="mod.name || mod.code">
                    {{ mod.name?.trim() ? mod.name : mod.code }}
                  </span>
                </router-link>

                <!-- Module Actions Menu -->
                <v-menu location="end">
                  <template #activator="{ props }">
                    <v-btn
                      v-bind="props"
                      icon="mdi-dots-vertical"
                      size="x-small"
                      variant="text"
                      class="module-more-btn"
                    ></v-btn>
                  </template>
                  <v-list density="compact" class="rounded-lg" elevation="3">
                    <v-list-item prepend-icon="mdi-pencil-outline" title="Edit Module" @click="openEditDialog(mod)"></v-list-item>
                    <v-list-item prepend-icon="mdi-folder-move-outline" title="Move to Folder" @click="openMoveDialog(mod)"></v-list-item>
                    <v-list-item prepend-icon="mdi-delete-outline" title="Delete Module" class="text-error" @click="confirmDeleteModule(mod)"></v-list-item>
                  </v-list>
                </v-menu>
              </div>
            </div>
          </div>
        </template>
      </nav>
    </div>

    <!-- Bottom Actions: Dark Mode & Backup (Pinned to the very bottom) -->
    <template #append>
      <div class="sidebar-bottom pa-3 border-top">
        <div class="d-flex align-center justify-space-between mb-2">
          <span class="text-caption text-grey font-weight-medium">THEME</span>
          <v-btn
            :icon="isDark ? 'mdi-weather-night' : 'mdi-weather-sunny'"
            size="small"
            variant="text"
            :color="isDark ? 'warning' : 'color-bg-2'"
            :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
            @click="toggleTheme"
          ></v-btn>
        </div>

        <v-btn
          variant="text"
          block
          size="small"
          prepend-icon="mdi-cloud-download-outline"
          class="text-none justify-start px-2 backup-btn"
          :loading="isBackingUp"
          @click="handleBackup"
        >
          Backup Data
        </v-btn>

        <div v-if="lastBackupDateFormatted" class="backup-date-display mt-1 px-2 d-flex align-center ga-1 text-grey text-caption">
          <v-icon icon="mdi-history" size="12"></v-icon>
          <span>Last: {{ lastBackupDateFormatted }}</span>
        </div>
      </div>
    </template>
  </v-navigation-drawer>

  <!-- Create Module Dialog -->
  <v-dialog persistent v-model="dialog" width="420">
    <v-card class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-4">New module</h3>

      <v-form ref="formRef" @submit.prevent="handleCreateModule">
        <v-text-field
          v-model="moduleData.code"
          label="Module code"
          placeholder="e.g. CS101"
          variant="outlined"
          :rules="codeRules"
          density="comfortable"
          class="mb-2"
          required
        ></v-text-field>

        <v-text-field
          v-model="moduleData.name"
          label="Module name"
          placeholder="e.g. Intro to Computer Science"
          variant="outlined"
          :rules="nameRules"
          density="comfortable"
          class="mb-2"
          required
        ></v-text-field>

        <v-select
          v-model="moduleData.folderId"
          :items="folders"
          item-title="name"
          item-value="id"
          label="Folder (optional)"
          variant="outlined"
          density="comfortable"
          clearable
          class="mb-2"
        ></v-select>

        <v-text-field
          v-model="moduleData.color"
          type="color"
          label="Module color"
          variant="outlined"
          density="comfortable"
          class="mb-4"
        ></v-text-field>

        <div class="d-flex gap-2 justify-end">
          <v-btn variant="text" @click="handleCancelCreateModule">Cancel</v-btn>
          <v-btn type="submit" color="color-bg-2" variant="flat">Create module</v-btn>
        </div>
      </v-form>
    </v-card>
  </v-dialog>

  <!-- Edit Module Dialog -->
  <v-dialog persistent v-model="editDialog" width="420">
    <v-card v-if="moduleToEdit" class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-4">Edit Module</h3>

      <v-form ref="editFormRef" @submit.prevent="handleUpdateModule">
        <v-text-field
          v-model="editData.code"
          label="Module code"
          variant="outlined"
          :rules="codeRules"
          density="comfortable"
          class="mb-2"
          required
        ></v-text-field>

        <v-text-field
          v-model="editData.name"
          label="Module name"
          variant="outlined"
          :rules="nameRules"
          density="comfortable"
          class="mb-2"
          required
        ></v-text-field>

        <v-select
          v-model="editData.folderId"
          :items="folders"
          item-title="name"
          item-value="id"
          label="Folder (optional)"
          variant="outlined"
          density="comfortable"
          clearable
          class="mb-2"
        ></v-select>

        <v-text-field
          v-model="editData.color"
          type="color"
          label="Module color"
          variant="outlined"
          density="comfortable"
          class="mb-4"
        ></v-text-field>

        <div class="d-flex gap-2 justify-end">
          <v-btn variant="text" @click="editDialog = false">Cancel</v-btn>
          <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isUpdatingModule">Save Changes</v-btn>
        </div>
      </v-form>
    </v-card>
  </v-dialog>

  <!-- Move Module Dialog -->
  <v-dialog persistent v-model="moveDialog" width="380">
    <v-card v-if="moduleToMove" class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-2">Move module</h3>
      <p class="text-body-2 text-grey-darken-1 mb-4">
        Choose a folder for <strong>{{ moduleToMove.code }}</strong>.
      </p>

      <v-select
        v-model="moveFolderId"
        :items="folders"
        item-title="name"
        item-value="id"
        label="Folder"
        variant="outlined"
        density="comfortable"
        clearable
        placeholder="No folder"
        class="mb-4"
      ></v-select>

      <div class="d-flex justify-end ga-2">
        <v-btn variant="text" @click="moveDialog = false">Cancel</v-btn>
        <v-btn color="color-bg-2" variant="flat" :loading="isMovingModule" @click="handleMoveModule">Move</v-btn>
      </div>
    </v-card>
  </v-dialog>

  <!-- Delete Module Confirmation Dialog -->
  <v-dialog v-model="deleteDialog" width="420">
    <v-card v-if="moduleToDelete" class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-2 text-error">Delete Module?</h3>
      <p class="text-body-2 text-grey-darken-1 mb-4">
        Are you sure you want to delete <strong>{{ moduleToDelete.code }} ({{ moduleToDelete.name }})</strong>?
        <span class="text-error d-block mt-2">
          <strong>Warning:</strong> All files, folders, assignments, and timetable slots for this module will be permanently removed.
        </span>
      </p>
      <div class="d-flex justify-end ga-2">
        <v-btn variant="text" @click="deleteDialog = false">Cancel</v-btn>
        <v-btn color="error" variant="flat" :loading="isDeletingModule" @click="executeDeleteModule">Delete</v-btn>
      </div>
    </v-card>
  </v-dialog>

  <!-- Create Folder Dialog -->
  <v-dialog persistent v-model="folderDialog" width="380">
    <v-card class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-4">New folder</h3>
      <v-form ref="folderFormRef" @submit.prevent="handleCreateFolder">
        <v-text-field
          v-model="folderName"
          label="Folder name"
          placeholder="e.g. 2026 Modules"
          variant="outlined"
          density="comfortable"
          :rules="folderNameRules"
          class="mb-4"
          autofocus
          required
        ></v-text-field>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="folderDialog = false">Cancel</v-btn>
          <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isSavingFolder">Create folder</v-btn>
        </div>
      </v-form>
    </v-card>
  </v-dialog>

  <!-- Rename Folder Dialog -->
  <v-dialog persistent v-model="renameFolderDialog" width="380">
    <v-card class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-4">Rename folder</h3>
      <v-form @submit.prevent="handleRenameFolder">
        <v-text-field
          v-model="renameFolderName"
          label="Folder name"
          variant="outlined"
          density="comfortable"
          :rules="folderNameRules"
          class="mb-4"
          required
        ></v-text-field>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="renameFolderDialog = false">Cancel</v-btn>
          <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isSavingFolder">Save</v-btn>
        </div>
      </v-form>
    </v-card>
  </v-dialog>

  <!-- Delete Folder Confirmation Dialog -->
  <v-dialog v-model="deleteFolderDialog" width="400">
    <v-card v-if="folderToDelete" class="pa-6 rounded-lg" elevation="4">
      <h3 class="text-h6 font-weight-medium mb-2 text-error">Delete folder?</h3>
      <p class="text-body-2 text-grey-darken-1 mb-4">
        Delete <strong>{{ folderToDelete.name }}</strong>? Modules inside it will not be deleted, they will move back to the top level.
      </p>
      <div class="d-flex justify-end ga-2">
        <v-btn variant="text" @click="deleteFolderDialog = false">Cancel</v-btn>
        <v-btn color="error" variant="flat" :loading="isSavingFolder" @click="executeDeleteFolder">Delete</v-btn>
      </div>
    </v-card>
  </v-dialog>

  <!-- Global Spotlight Search Dialog (Ctrl+K) -->
  <v-dialog v-model="searchDialog" width="560">
    <v-card class="pa-4 rounded-lg global-search-card" elevation="6">
      <v-text-field
        v-model="searchQuery"
        prepend-inner-icon="mdi-magnify"
        placeholder="Search modules, files, assignments, timetable..."
        variant="outlined"
        density="comfortable"
        autofocus
        clearable
        hide-details
        class="mb-3"
      ></v-text-field>

      <div class="search-results-list" style="max-height: 400px; overflow-y: auto;">
        <div v-if="!searchQuery.trim()" class="py-6 text-center text-grey">
          <v-icon icon="mdi-keyboard-outline" size="32" class="mb-1"></v-icon>
          <div class="text-caption">Type to search modules, documents, and assignments</div>
        </div>

        <div v-else-if="filteredSearchResults.length === 0" class="py-6 text-center text-grey">
          <div class="text-caption">No results found for "{{ searchQuery }}"</div>
        </div>

        <template v-else>
          <div
            v-for="(item, idx) in filteredSearchResults"
            :key="idx"
            class="search-result-item pa-2 rounded-lg d-flex align-center justify-space-between cursor-pointer mb-1"
            @click="handleSelectSearchResult(item)"
          >
            <div class="d-flex align-center ga-3 overflow-hidden">
              <v-icon :icon="item.icon" size="20" :color="item.color || 'color-bg-2'"></v-icon>
              <div class="overflow-hidden">
                <div class="text-body-2 font-weight-bold text-truncate">{{ item.title }}</div>
                <div class="text-caption text-grey text-truncate">{{ item.subtitle }}</div>
              </div>
            </div>
            <v-chip size="x-small" variant="tonal" :color="item.color || 'grey'" class="text-uppercase font-weight-medium">
              {{ item.category }}
            </v-chip>
          </div>
        </template>
      </div>
    </v-card>
  </v-dialog>

  <!-- Alerts -->
  <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
    {{ alertMsg }}
  </v-snackbar>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from 'vuetify'
import {
  GetAllModules,
  CreateModule,
  UpdateModule,
  DeleteModule,
  GetAllModuleFolders,
  CreateModuleFolder,
  UpdateModuleFolder,
  DeleteModuleFolder,
} from '../../wailsjs/go/controllers/ModuleController'
import { GetAllAssignments, UpdateAssignmentStatus } from '../../wailsjs/go/controllers/AssignmentController'
import { SetTaskDone } from '../../wailsjs/go/controllers/TaskController'
import { CreateBackup } from '../../wailsjs/go/controllers/BackupController'
import { OpenFile } from '../../wailsjs/go/controllers/FileController'
import { models } from '../../wailsjs/go/models'
import { eventBus, SYNC_EVENTS } from '../utils/eventBus'

const router = useRouter()
const theme = useTheme()

// Collapsible tasks in sidebar
const allAssignments = ref<models.Assignment[]>([])
const isTasksCollapsed = ref(localStorage.getItem('schive_sidebar_tasks_collapsed') === 'true')

function toggleTasksCollapse() {
  isTasksCollapsed.value = !isTasksCollapsed.value
  localStorage.setItem('schive_sidebar_tasks_collapsed', String(isTasksCollapsed.value))
}

function getPriorityWeight(p?: string): number {
  switch ((p || '').toLowerCase()) {
    case 'high': return 3
    case 'medium': return 2
    case 'low': return 1
    default: return 0
  }
}

function compareTasks(a: any, b: any): number {
  // 1. Incomplete before completed
  const aDone = !!a.completed
  const bDone = !!b.completed
  if (aDone !== bDone) {
    return aDone ? 1 : -1
  }

  // 2. Deadline (earliest deadline first)
  const aTime = a.dueDate ? new Date(a.dueDate).getTime() : NaN
  const bTime = b.dueDate ? new Date(b.dueDate).getTime() : NaN
  const aHasDeadline = !isNaN(aTime)
  const bHasDeadline = !isNaN(bTime)

  if (aHasDeadline && bHasDeadline) {
    if (aTime !== bTime) {
      return aTime - bTime
    }
    // Same deadline -> sort by importance (priority)
    const pDiff = getPriorityWeight(b.priority) - getPriorityWeight(a.priority)
    if (pDiff !== 0) return pDiff
  } else if (aHasDeadline && !bHasDeadline) {
    return -1
  } else if (!aHasDeadline && bHasDeadline) {
    return 1
  } else {
    // Neither has a deadline -> sort by importance
    const pDiff = getPriorityWeight(b.priority) - getPriorityWeight(a.priority)
    if (pDiff !== 0) return pDiff
  }

  // 3. Fallback: alphabetical by title
  return (a.title || '').localeCompare(b.title || '')
}

const sortedAllTasks = computed(() => {
  const items: any[] = []

  // From assignments
  for (const a of allAssignments.value) {
    items.push({
      id: `assign-${a.id}`,
      originalId: a.id,
      title: a.title,
      dueDate: a.dueDate,
      priority: a.priority || 'medium',
      status: a.status,
      completed: a.status === 'completed',
      moduleCode: a.module?.code || modules.value.find(m => m.id === a.moduleId)?.code || '',
      moduleColor: a.module?.color || modules.value.find(m => m.id === a.moduleId)?.color || '#7A8B69',
      isModuleTodo: false
    })
  }

  // From module tasks (if any)
  for (const mod of modules.value) {
    if (mod.tasks && mod.tasks.length > 0) {
      for (const t of mod.tasks) {
        if (!items.some(i => i.originalId === t.id)) {
          items.push({
            id: `task-${t.id}`,
            originalId: t.id,
            title: t.title,
            dueDate: t.dueDate,
            priority: 'medium',
            status: t.done ? 'completed' : 'pending',
            completed: !!t.done,
            moduleCode: mod.code,
            moduleColor: mod.color || '#7A8B69',
            isModuleTodo: true
          })
        }
      }
    }
  }

  return items.sort(compareTasks)
})

function formatShortDeadline(dateStr?: string | Date): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  const now = new Date()
  const todayDate = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const targetDate = new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const diffDays = Math.round((targetDate - todayDate) / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return 'Overdue'
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Tmrw'
  if (diffDays < 7) return `${diffDays}d`
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

function getDeadlineClass(dateStr?: string | Date): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  const now = new Date()
  const todayDate = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const targetDate = new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const diffDays = Math.round((targetDate - todayDate) / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return 'deadline-overdue'
  if (diffDays === 0) return 'deadline-today'
  if (diffDays === 1) return 'deadline-soon'
  return 'deadline-normal'
}

function getPriorityColor(priority?: string): string {
  switch ((priority || '').toLowerCase()) {
    case 'high': return '#E53935'
    case 'medium': return '#FB8C00'
    case 'low': return '#1E88E5'
    default: return '#9E9E9E'
  }
}

function formatDeadlineFull(dateStr?: string | Date): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
}

function navigateToTask(task: any) {
  if (task.moduleCode) {
    router.push(`/module/${task.moduleCode}?tab=${task.isModuleTodo ? 'tasks' : 'assignments'}`)
  } else {
    router.push('/tasks')
  }
}

async function toggleTaskComplete(task: any) {
  const willBeCompleted = !task.completed
  try {
    if (task.isModuleTodo) {
      await SetTaskDone(task.originalId, willBeCompleted)
      task.completed = willBeCompleted
      task.status = willBeCompleted ? 'completed' : 'pending'
      eventBus.emit(SYNC_EVENTS.TASKS_CHANGED, { id: task.originalId, done: willBeCompleted })
    } else {
      const newStatus = willBeCompleted ? 'completed' : 'pending'
      await UpdateAssignmentStatus(task.originalId, newStatus)
      task.completed = willBeCompleted
      task.status = newStatus
      eventBus.emit(SYNC_EVENTS.ASSIGNMENTS_CHANGED, { id: task.originalId, status: newStatus })
    }
    await fetchModules()
  } catch (err: any) {
    showAlertMsg('Failed to update task: ' + err, 'error')
  }
}

// Last backup timestamp
const lastBackupDate = ref<string | null>(localStorage.getItem('schive_last_backup_date'))

const lastBackupDateFormatted = computed(() => {
  if (!lastBackupDate.value) return null
  const d = new Date(lastBackupDate.value)
  if (isNaN(d.getTime())) return null
  return d.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

const isDark = ref(theme.global.current.value.dark)

function toggleTheme() {
  const nextTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = nextTheme
  isDark.value = !isDark.value
  localStorage.setItem('schive-theme', nextTheme)
  if (nextTheme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

const formRef = ref<any>(null)
const editFormRef = ref<any>(null)

const moduleData = ref(models.Module.createFrom({
  code: '',
  name: '',
  color: '#7A8B69',
  folderId: undefined,
}))

const editData = ref({
  code: '',
  name: '',
  color: '#7A8B69',
  folderId: undefined as string | undefined,
})

const codeRules = [
  (value: string) => {
    if (value && value.trim().length > 0) return true
    return 'Module Code is required.'
  }
]

const nameRules = [
  (value: string) => {
    if (value && value.trim().length > 0) return true
    return 'Module Name is required.'
  }
]

const dialog = ref(false)
const editDialog = ref(false)
const deleteDialog = ref(false)
const isUpdatingModule = ref(false)
const isDeletingModule = ref(false)
const isBackingUp = ref(false)

const moduleToEdit = ref<models.Module | null>(null)
const moduleToDelete = ref<models.Module | null>(null)

const alertMsg = ref('')
const alertType = ref('success')
const showAlert = ref(false)

const modules = ref<models.Module[]>([])
const folders = ref<models.ModuleFolder[]>([])
const isLoading = ref(true)
const collapsedFolders = ref<Record<string, boolean>>({})

const moduleGroups = computed(() => {
  const groups: { id: string | null; name: string; modules: models.Module[] }[] = []

  for (const folder of folders.value) {
    groups.push({
      id: folder.id,
      name: folder.name,
      modules: modules.value.filter(m => m.folderId === folder.id),
    })
  }

  const unfiled = modules.value.filter(m => !m.folderId || !folders.value.some(f => f.id === m.folderId))
  if (unfiled.length > 0) {
    groups.push({ id: null, name: folders.value.length ? 'Unfiled' : 'Modules', modules: unfiled })
  }

  return groups
})

function toggleFolder(id: string) {
  collapsedFolders.value[id] = !collapsedFolders.value[id]
}

// Global Spotlight Search
const searchDialog = ref(false)
const searchQuery = ref('')

function openSearchDialog() {
  searchQuery.value = ''
  searchDialog.value = true
}

const filteredSearchResults = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return []
  const results: any[] = []

  for (const mod of modules.value) {
    // Search module
    if (mod.code.toLowerCase().includes(q) || mod.name.toLowerCase().includes(q)) {
      results.push({
        category: 'Module',
        title: `${mod.code} — ${mod.name}`,
        subtitle: 'Navigate to module details',
        icon: 'mdi-book-outline',
        color: mod.color,
        action: () => router.push(`/module/${mod.code}`)
      })
    }

    // Search files
    if (mod.files) {
      for (const f of mod.files) {
        if (f.name.toLowerCase().includes(q)) {
          results.push({
            category: 'File',
            title: f.name,
            subtitle: `In ${mod.code}`,
            icon: 'mdi-file-document-outline',
            color: 'primary',
            action: async () => {
              await OpenFile(f.id)
            }
          })
        }
      }
    }

    // Search assignments
    if (mod.assignments) {
      for (const a of mod.assignments) {
        if (a.title.toLowerCase().includes(q) || (a.description && a.description.toLowerCase().includes(q))) {
          results.push({
            category: 'Assignment',
            title: a.title,
            subtitle: `In ${mod.code} (${a.status})`,
            icon: 'mdi-clipboard-text-outline',
            color: 'warning',
            action: () => router.push(`/module/${mod.code}`)
          })
        }
      }
    }

    // Search timetable slots
    if (mod.classSlots) {
      for (const s of mod.classSlots) {
        if ((s.day && s.day.toLowerCase().includes(q)) || (s.room && s.room.toLowerCase().includes(q))) {
          results.push({
            category: 'Timetable',
            title: `${mod.code} — ${s.day} @ ${s.startHour}:00`,
            subtitle: `Room: ${s.room || 'TBA'}`,
            icon: 'mdi-calendar-clock',
            color: mod.color,
            action: () => router.push('/timetable')
          })
        }
      }
    }
  }

  return results.slice(0, 15)
})

function handleSelectSearchResult(item: any) {
  searchDialog.value = false
  if (item.action) {
    item.action()
  }
}

function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    searchDialog.value = !searchDialog.value
  }
}

async function fetchModules() {
  try {
    isLoading.value = true
    const [mods, folderList, assigns] = await Promise.all([
      GetAllModules(),
      GetAllModuleFolders(),
      GetAllAssignments().catch(() => []),
    ])
    modules.value = mods || []
    folders.value = folderList || []
    allAssignments.value = assigns || []
    isLoading.value = false
  } catch (err) {
    alertMsg.value = 'Failed to get available modules ' + (err as string)
    alertType.value = 'error'
    showAlert.value = true
  }
}

async function handleCreateModule() {
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) return
  }

  const trimmedCode = moduleData.value.code?.trim() || ''
  const trimmedName = moduleData.value.name?.trim() || ''

  if (!trimmedCode || !trimmedName) {
    alertMsg.value = 'Module Code and Module Name are required.'
    alertType.value = 'error'
    showAlert.value = true
    return
  }

  try {
    moduleData.value.code = trimmedCode
    moduleData.value.name = trimmedName
    await CreateModule(moduleData.value)
    await fetchModules()
    handleCancelCreateModule()
    alertMsg.value = 'Module created successfully'
    alertType.value = 'success'
    showAlert.value = true
  } catch (err) {
    alertMsg.value = 'Failed to create module: ' + (err as string)
    alertType.value = 'error'
    showAlert.value = true
  }
}

function handleCancelCreateModule() {
  dialog.value = false
  if (formRef.value) {
    formRef.value.resetValidation()
  }
  moduleData.value.code = ''
  moduleData.value.name = ''
  moduleData.value.color = '#7A8B69'
  moduleData.value.folderId = undefined
}

function openEditDialog(mod: models.Module) {
  moduleToEdit.value = mod
  editData.value.code = mod.code
  editData.value.name = mod.name
  editData.value.color = mod.color || '#7A8B69'
  editData.value.folderId = mod.folderId
  editDialog.value = true
}

async function handleUpdateModule() {
  if (!moduleToEdit.value) return
  if (editFormRef.value) {
    const { valid } = await editFormRef.value.validate()
    if (!valid) return
  }

  try {
    isUpdatingModule.value = true
    const updated = models.Module.createFrom({
      id: moduleToEdit.value.id,
      code: editData.value.code.trim(),
      name: editData.value.name.trim(),
      color: editData.value.color.trim(),
      folderId: editData.value.folderId || undefined,
    })
    await UpdateModule(moduleToEdit.value.code, updated)
    editDialog.value = false
    alertMsg.value = 'Module updated successfully!'
    alertType.value = 'success'
    showAlert.value = true
    await fetchModules()
  } catch (err: any) {
    alertMsg.value = 'Failed to update module: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isUpdatingModule.value = false
  }
}

// Move a module into / out of a folder
const moveDialog = ref(false)
const moduleToMove = ref<models.Module | null>(null)
const moveFolderId = ref<string | null>(null)
const isMovingModule = ref(false)

function openMoveDialog(mod: models.Module) {
  moduleToMove.value = mod
  moveFolderId.value = mod.folderId || null
  moveDialog.value = true
}

async function handleMoveModule() {
  if (!moduleToMove.value) return
  try {
    isMovingModule.value = true
    const updated = models.Module.createFrom({
      id: moduleToMove.value.id,
      code: moduleToMove.value.code,
      name: moduleToMove.value.name,
      color: moduleToMove.value.color,
      folderId: moveFolderId.value || undefined,
    })
    await UpdateModule(moduleToMove.value.code, updated)
    moveDialog.value = false

    const folderName = folders.value.find(f => f.id === moveFolderId.value)?.name
    alertMsg.value = folderName
      ? `${moduleToMove.value.code} moved to ${folderName}`
      : `${moduleToMove.value.code} removed from its folder`
    alertType.value = 'success'
    showAlert.value = true
    await fetchModules()
  } catch (err: any) {
    alertMsg.value = 'Failed to move module: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isMovingModule.value = false
  }
}

function confirmDeleteModule(mod: models.Module) {
  moduleToDelete.value = mod
  deleteDialog.value = true
}

async function executeDeleteModule() {
  if (!moduleToDelete.value) return
  try {
    isDeletingModule.value = true
    await DeleteModule(moduleToDelete.value.code)
    deleteDialog.value = false
    alertMsg.value = `Deleted module ${moduleToDelete.value.code}`
    alertType.value = 'success'
    showAlert.value = true
    await fetchModules()
    router.push('/')
  } catch (err: any) {
    alertMsg.value = 'Failed to delete module: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isDeletingModule.value = false
  }
}

// --- Module folders ---

const folderDialog = ref(false)
const renameFolderDialog = ref(false)
const deleteFolderDialog = ref(false)
const folderName = ref('')
const renameFolderId = ref('')
const renameFolderName = ref('')
const folderToDelete = ref<models.ModuleFolder | null>(null)
const isSavingFolder = ref(false)
const folderFormRef = ref<any>(null)

const folderNameRules = [
  (value: string) => {
    if (value && value.trim().length > 0) return true
    return 'Folder name is required.'
  }
]

function openCreateFolderDialog() {
  folderName.value = ''
  folderDialog.value = true
}

async function handleCreateFolder() {
  if (folderFormRef.value) {
    const { valid } = await folderFormRef.value.validate()
    if (!valid) return
  }

  try {
    isSavingFolder.value = true
    await CreateModuleFolder(folderName.value.trim())
    folderDialog.value = false
    alertMsg.value = 'Folder created successfully'
    alertType.value = 'success'
    showAlert.value = true
    await fetchModules()
  } catch (err: any) {
    alertMsg.value = 'Failed to create folder: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingFolder.value = false
  }
}

function openRenameFolderDialog(group: { id: string | null; name: string }) {
  if (!group.id) return
  renameFolderId.value = group.id
  renameFolderName.value = group.name
  renameFolderDialog.value = true
}

async function handleRenameFolder() {
  if (!renameFolderId.value) return
  try {
    isSavingFolder.value = true
    await UpdateModuleFolder(renameFolderId.value, renameFolderName.value.trim())
    renameFolderDialog.value = false
    alertMsg.value = 'Folder renamed successfully'
    alertType.value = 'success'
    showAlert.value = true
    await fetchModules()
  } catch (err: any) {
    alertMsg.value = 'Failed to rename folder: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingFolder.value = false
  }
}

function confirmDeleteFolder(group: { id: string | null; name: string }) {
  if (!group.id) return
  folderToDelete.value = folders.value.find(f => f.id === group.id) || null
  deleteFolderDialog.value = true
}

async function executeDeleteFolder() {
  if (!folderToDelete.value) return
  try {
    isSavingFolder.value = true
    await DeleteModuleFolder(folderToDelete.value.id)
    deleteFolderDialog.value = false
    alertMsg.value = `Deleted folder ${folderToDelete.value.name}`
    alertType.value = 'success'
    showAlert.value = true
    await fetchModules()
  } catch (err: any) {
    alertMsg.value = 'Failed to delete folder: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingFolder.value = false
  }
}

async function handleBackup() {
  try {
    isBackingUp.value = true
    const backupPath = await CreateBackup()
    const nowIso = new Date().toISOString()
    lastBackupDate.value = nowIso
    localStorage.setItem('schive_last_backup_date', nowIso)
    alertMsg.value = `Backup created successfully at: ${backupPath}`
    alertType.value = 'success'
    showAlert.value = true
  } catch (err: any) {
    alertMsg.value = 'Failed to create backup: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isBackingUp.value = false
  }
}

watch(
  () => router.currentRoute.value.fullPath,
  () => {
    fetchModules()
  }
)

let unsubTasks: (() => void) | null = null
let unsubAssignments: (() => void) | null = null
let unsubModules: (() => void) | null = null

onMounted(() => {
  fetchModules()
  window.addEventListener('keydown', handleGlobalKeydown)

  unsubTasks = eventBus.on(SYNC_EVENTS.TASKS_CHANGED, () => {
    fetchModules()
  })
  unsubAssignments = eventBus.on(SYNC_EVENTS.ASSIGNMENTS_CHANGED, () => {
    fetchModules()
  })
  unsubModules = eventBus.on(SYNC_EVENTS.MODULES_CHANGED, () => {
    fetchModules()
  })
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
  if (unsubTasks) unsubTasks()
  if (unsubAssignments) unsubAssignments()
  if (unsubModules) unsubModules()
})
</script>

<style scoped>
.custom-sidebar {
  background-color: var(--color-secondary) !important;
  border-right: 2px solid var(--color-bg-4) !important;
  color: var(--color-text-light);
}

.sidebar-brand-text {
  font-size: 16px;
  font-weight: 800;
  letter-spacing: 1px;
  color: var(--color-text-dark);
}

.sidebar-subtitle {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: var(--color-text-light);
}

.module-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.nav-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 8px;
  text-decoration: none;
  color: var(--color-text-dark);
  white-space: nowrap;
  font-size: 13.5px;
  font-weight: 600;
  transition: all 0.15s ease;
}

.nav-item:hover {
  background: var(--color-bg-3);
}

.nav-item.active {
  background: var(--color-bg-3);
  color: var(--color-bg-2);
}

.module-group {
  margin-bottom: 4px;
}

.folder-row {
  color: var(--color-text-light);
  transition: background-color 0.15s ease;
}

.folder-row:hover {
  background: var(--color-bg-3);
}

.folder-label {
  font-size: 11.5px;
  font-weight: 700;
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.folder-count {
  font-size: 10.5px;
  font-weight: 700;
  color: var(--color-text-light);
  opacity: 0.8;
}

.folder-module-list {
  padding-left: 10px;
}

.folder-row:hover .module-more-btn {
  opacity: 1;
}

.module-row {
  transition: background-color 0.15s ease;
}

.module-row:hover .module-more-btn {
  opacity: 1;
}

.module-more-btn {
  opacity: 0.4;
  transition: opacity 0.15s ease;
}

.search-kbd-pill {
  font-size: 10px;
  padding: 2px 5px;
  border-radius: 4px;
  background: var(--color-bg-4);
  color: var(--color-text-light);
  font-family: monospace;
}

.sidebar-divider {
  width: 90%;
  margin: 10px auto;
  border: none;
  border-top: 1px solid var(--color-bg-4);
}

.sidebar-bottom {
  border-top: 1px solid var(--color-bg-4);
}

.backup-btn:hover {
  background-color: var(--color-bg-3);
}

.search-result-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.sidebar-section-header:hover {
  background: var(--color-bg-3);
}

.tasks-sidebar-list {
  max-height: 240px;
  overflow-y: auto;
}

.sidebar-task-item {
  transition: background-color 0.15s ease;
  min-height: 28px;
}

.sidebar-task-item:hover {
  background: var(--color-bg-3);
}

.sidebar-task-item.task-completed {
  opacity: 0.55;
}

.sidebar-task-item.task-completed .task-title {
  text-decoration: line-through;
}

.task-priority-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.task-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-dark);
}

.task-deadline-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 4px;
  white-space: nowrap;
}

.deadline-overdue {
  background-color: rgba(229, 57, 53, 0.15);
  color: #E53935;
}

.deadline-today {
  background-color: rgba(251, 140, 0, 0.15);
  color: #FB8C00;
}

.deadline-soon {
  background-color: rgba(30, 136, 229, 0.15);
  color: #1E88E5;
}

.deadline-normal {
  background-color: var(--color-bg-4);
  color: var(--color-text-light);
}

.task-module-pill {
  font-size: 9.5px;
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 3px;
  background-color: var(--color-bg-4);
  letter-spacing: 0.2px;
}

.backup-date-display {
  font-size: 11px;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
