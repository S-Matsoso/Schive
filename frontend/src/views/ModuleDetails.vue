<template>
  <div class="module-details-page pa-6">
    <!-- Module Header -->
    <v-card v-if="currentModule" class="pa-6 rounded-lg mb-6" elevation="2">
      <div class="d-flex justify-space-between align-center">
        <div>
          <div class="d-flex align-center ga-3 mb-1">
            <span
              class="module-color-pill"
              :style="{ backgroundColor: currentModule.color }"
            ></span>
            <h1 class="text-h4 font-weight-bold" :style="{ color: currentModule.color }">
              {{ currentModule.code }}
            </h1>
          </div>
          <h2 class="text-h6 text-grey-darken-1 font-weight-regular">
            {{ currentModule.name || 'No module name specified' }}
          </h2>
        </div>

        <div class="d-flex ga-2">
          <v-chip color="color-bg-2" variant="tonal" prepend-icon="mdi-file-document-multiple">
            {{ files.length }} {{ files.length === 1 ? 'file' : 'files' }}
          </v-chip>
        </div>
      </div>
    </v-card>

    <v-card v-else class="pa-12 rounded-lg text-center" elevation="1">
      <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
      <p class="mt-4 text-grey">Loading module data...</p>
    </v-card>

    <!-- Files Section -->
    <div v-if="currentModule" class="files-container">
      <!-- Toolbar -->
      <div class="d-flex flex-wrap justify-space-between align-center ga-4 mb-4">
        <div class="d-flex align-center ga-3">
          <h2 class="text-h6 font-weight-bold">Module Files</h2>
          <v-chip size="small" variant="flat" color="color-bg-3" class="text-grey-darken-3 font-weight-medium">
            {{ filteredFiles.length }}
          </v-chip>
        </div>

        <div class="d-flex align-center ga-3 flex-wrap">
          <!-- Search input -->
          <v-text-field
            v-if="files.length > 0"
            v-model="searchQuery"
            prepend-inner-icon="mdi-magnify"
            placeholder="Search files..."
            density="compact"
            variant="outlined"
            rounded="lg"
            hide-details
            clearable
            style="width: 220px;"
          ></v-text-field>

          <!-- Grid / List Toggle -->
          <v-btn-toggle
            v-model="viewMode"
            mandatory
            divided
            rounded="lg"
            variant="outlined"
            class="view-toggle"
          >
            <v-btn value="grid" icon="mdi-view-grid" title="Grid View"></v-btn>
            <v-btn value="list" icon="mdi-format-list-bulleted" title="List View"></v-btn>
          </v-btn-toggle>

          <!-- Hidden File Input -->
          <input
            ref="fileInputRef"
            type="file"
            multiple
            style="display: none;"
            @change="handleFileSelected"
          />

          <!-- Upload Button -->
          <v-btn
            prepend-icon="mdi-upload"
            variant="flat"
            class="upload-btn text-none font-weight-medium rounded-lg"
            height="40"
            :loading="isUploading"
            @click="triggerFileInput"
          >
            Upload File
          </v-btn>
        </div>
      </div>

      <!-- Drag & Drop Overlay (active when dragging files) -->
      <div
        v-if="isDragging"
        class="drag-overlay pa-8 rounded-lg text-center mb-6"
        @dragover.prevent
        @dragleave.prevent="isDragging = false"
        @drop.prevent="handleDrop"
      >
        <v-icon icon="mdi-cloud-upload" size="48" color="color-bg-2" class="mb-2"></v-icon>
        <div class="text-subtitle-1 font-weight-bold">Drop files here to upload to {{ currentModule.code }}</div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoadingFiles" class="py-12 text-center">
        <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
        <p class="mt-3 text-caption text-grey">Loading files...</p>
      </div>

      <!-- Empty State -->
      <v-card
        v-else-if="files.length === 0"
        class="pa-12 text-center rounded-lg empty-card"
        variant="tonal"
        @dragover.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        @drop.prevent="handleDrop"
      >
        <v-icon icon="mdi-folder-open-outline" size="56" color="grey" class="mb-3"></v-icon>
        <h3 class="text-h6 font-weight-medium mb-1">No files yet for {{ currentModule.code }}</h3>
        <p class="text-body-2 text-grey-darken-1">
          Use the <strong>Upload File</strong> button above or drag and drop documents here.
        </p>
      </v-card>

      <!-- No Search Results State -->
      <div v-else-if="filteredFiles.length === 0" class="py-8 text-center text-grey">
        <v-icon icon="mdi-file-search-outline" size="36" class="mb-2"></v-icon>
        <p>No files match "{{ searchQuery }}"</p>
      </div>

      <!-- GRID VIEW -->
      <div v-else-if="viewMode === 'grid'" class="files-grid">
        <v-card
          v-for="file in filteredFiles"
          :key="file.id"
          class="file-card pa-4 rounded-lg d-flex flex-column justify-space-between"
          elevation="2"
          @click="handleOpenFile(file.id)"
        >
          <!-- Top file icon & badge -->
          <div>
            <div class="d-flex justify-space-between align-start mb-3">
              <v-icon
                :icon="getFileInfo(file.fileType).icon"
                :color="getFileInfo(file.fileType).color"
                size="44"
              ></v-icon>
              <v-chip
                size="x-small"
                variant="tonal"
                :color="getFileInfo(file.fileType).color"
                class="font-weight-bold"
              >
                {{ getFileInfo(file.fileType).label }}
              </v-chip>
            </div>

            <!-- File Name -->
            <div class="file-name font-weight-bold mb-1" :title="file.name">
              {{ file.name }}
            </div>

            <!-- File Meta -->
            <div class="text-caption text-grey">
              {{ formatFileSize(file.fileSize) }}
              <span v-if="file.createdAt"> • {{ formatDate(file.createdAt) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="file-card-actions d-flex justify-space-between align-center mt-4 pt-2 border-top">
            <v-btn
              variant="text"
              size="small"
              prepend-icon="mdi-open-in-new"
              color="color-bg-2"
              class="text-none"
              @click.stop="handleOpenFile(file.id)"
            >
              Open
            </v-btn>

            <v-btn
              variant="text"
              size="small"
              icon="mdi-delete-outline"
              color="error"
              title="Delete File"
              @click.stop="confirmDeleteFile(file)"
            ></v-btn>
          </div>
        </v-card>
      </div>

      <!-- LIST VIEW -->
      <v-card v-else class="rounded-lg overflow-hidden" elevation="2">
        <v-table hover>
          <thead>
            <tr>
              <th style="width: 50px;"></th>
              <th>Name</th>
              <th style="width: 120px;">Type</th>
              <th style="width: 120px;">Size</th>
              <th style="width: 160px;">Date Uploaded</th>
              <th style="width: 130px;" class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="file in filteredFiles"
              :key="file.id"
              class="file-table-row cursor-pointer"
              @click="handleOpenFile(file.id)"
            >
              <td>
                <v-icon
                  :icon="getFileInfo(file.fileType).icon"
                  :color="getFileInfo(file.fileType).color"
                  size="26"
                ></v-icon>
              </td>
              <td>
                <span class="font-weight-medium">{{ file.name }}</span>
              </td>
              <td>
                <v-chip
                  size="x-small"
                  variant="tonal"
                  :color="getFileInfo(file.fileType).color"
                  class="font-weight-bold"
                >
                  {{ getFileInfo(file.fileType).label }}
                </v-chip>
              </td>
              <td class="text-caption text-grey-darken-1">
                {{ formatFileSize(file.fileSize) }}
              </td>
              <td class="text-caption text-grey-darken-1">
                {{ formatDate(file.createdAt) }}
              </td>
              <td class="text-right">
                <v-btn
                  icon="mdi-open-in-new"
                  variant="text"
                  size="small"
                  color="color-bg-2"
                  title="Open file"
                  @click.stop="handleOpenFile(file.id)"
                ></v-btn>
                <v-btn
                  icon="mdi-delete-outline"
                  variant="text"
                  size="small"
                  color="error"
                  title="Delete file"
                  @click.stop="confirmDeleteFile(file)"
                ></v-btn>
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </div>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" width="380">
      <v-card v-if="fileToDelete" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-2">Delete File?</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">
          Are you sure you want to delete <strong>{{ fileToDelete.name }}</strong>? This action cannot be undone.
        </p>
        <div class="d-flex justify-end ga-2">
          <v-btn type="button" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="flat" :loading="isDeleting" @click="executeDeleteFile">
            Delete
          </v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Alert Snackbar -->
    <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
      {{ alertMsg }}
    </v-snackbar>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { models } from '../../wailsjs/go/models'
import { GetModule } from '../../wailsjs/go/controllers/ModuleController'
import { UploadFile, GetFilesByModule, DeleteFile, OpenFile } from '../../wailsjs/go/controllers/FileController'

const route = useRoute()
const moduleCode = computed(() => route.params.code as string)

const currentModule = ref<models.Module | null>(null)
const files = ref<models.File[]>([])
const viewMode = ref<'grid' | 'list'>('grid')
const searchQuery = ref('')

const isUploading = ref(false)
const isLoadingFiles = ref(false)
const isDragging = ref(false)
const isDeleting = ref(false)

const fileInputRef = ref<HTMLInputElement | null>(null)
const deleteDialog = ref(false)
const fileToDelete = ref<models.File | null>(null)

const alertMsg = ref('')
const alertType = ref('success')
const showAlert = ref(false)

const filteredFiles = computed(() => {
  if (!searchQuery.value?.trim()) return files.value
  const q = searchQuery.value.toLowerCase().trim()
  return files.value.filter(f => f.name.toLowerCase().includes(q))
})

async function fetchModuleDetails() {
  try {
    currentModule.value = await GetModule(moduleCode.value)
    if (currentModule.value?.id) {
      await fetchFiles()
    }
  } catch (err: any) {
    alertMsg.value = 'Failed to get module details: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

async function fetchFiles() {
  if (!currentModule.value?.id) return
  try {
    isLoadingFiles.value = true
    const res = await GetFilesByModule(currentModule.value.id)
    files.value = res || []
  } catch (err: any) {
    alertMsg.value = 'Failed to load module files: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isLoadingFiles.value = false
  }
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

async function handleFileSelected(event: Event) {
  const target = event.target as HTMLInputElement
  const fileList = target.files
  if (!fileList || fileList.length === 0) return

  await uploadFilesList(fileList)
  target.value = ''
}

function handleDrop(event: DragEvent) {
  isDragging.value = false
  const fileList = event.dataTransfer?.files
  if (!fileList || fileList.length === 0) return
  uploadFilesList(fileList)
}

async function uploadFilesList(fileList: FileList) {
  if (!currentModule.value?.id) return

  isUploading.value = true
  let successCount = 0

  for (let i = 0; i < fileList.length; i++) {
    const file = fileList[i]
    try {
      const base64Data = await readFileAsBase64(file)
      await UploadFile(currentModule.value.id, file.name, base64Data)
      successCount++
    } catch (err: any) {
      alertMsg.value = `Failed to upload ${file.name}: ` + err
      alertType.value = 'error'
      showAlert.value = true
    }
  }

  isUploading.value = false

  if (successCount > 0) {
    alertMsg.value = successCount === 1 ? 'File uploaded successfully!' : `${successCount} files uploaded successfully!`
    alertType.value = 'success'
    showAlert.value = true
    await fetchFiles()
  }
}

function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      const commaIndex = result.indexOf(',')
      if (commaIndex !== -1) {
        resolve(result.substring(commaIndex + 1))
      } else {
        resolve(result)
      }
    }
    reader.onerror = error => reject(error)
    reader.readAsDataURL(file)
  })
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

function confirmDeleteFile(file: models.File) {
  fileToDelete.value = file
  deleteDialog.value = true
}

async function executeDeleteFile() {
  if (!fileToDelete.value) return
  try {
    isDeleting.value = true
    await DeleteFile(fileToDelete.value.id)
    deleteDialog.value = false
    alertMsg.value = `Deleted "${fileToDelete.value.name}"`
    alertType.value = 'success'
    showAlert.value = true
    await fetchFiles()
  } catch (err: any) {
    alertMsg.value = 'Failed to delete file: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isDeleting.value = false
  }
}

function formatFileSize(bytes: number): string {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatDate(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function getFileInfo(ext: string) {
  const e = (ext || '').toLowerCase()
  switch (e) {
    case 'pdf':
      return { icon: 'mdi-file-pdf-box', color: '#D32F2F', label: 'PDF' }
    case 'doc':
    case 'docx':
      return { icon: 'mdi-file-word-box', color: '#1976D2', label: 'DOC' }
    case 'xls':
    case 'xlsx':
    case 'csv':
      return { icon: 'mdi-file-excel-box', color: '#388E3C', label: 'SHEET' }
    case 'ppt':
    case 'pptx':
      return { icon: 'mdi-file-powerpoint-box', color: '#F57C00', label: 'SLIDES' }
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'gif':
    case 'svg':
    case 'webp':
      return { icon: 'mdi-file-image', color: '#00897B', label: 'IMAGE' }
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return { icon: 'mdi-folder-zip', color: '#FFA000', label: 'ARCHIVE' }
    case 'txt':
    case 'md':
      return { icon: 'mdi-file-document-outline', color: '#546E7A', label: 'TEXT' }
    default:
      return { icon: 'mdi-file-outline', color: '#757575', label: e.toUpperCase() || 'FILE' }
  }
}

onMounted(() => {
  fetchModuleDetails()
})

watch(moduleCode, (newCode, oldCode) => {
  if (newCode && newCode !== oldCode) {
    currentModule.value = null
    files.value = []
    searchQuery.value = ''
    fetchModuleDetails()
  }
})
</script>

<style scoped>
.module-details-page {
  width: 100%;
}

.module-color-pill {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  display: inline-block;
}

/* Grid / list toggle */
.view-toggle {
  height: 40px !important;
  border: 1px solid rgba(0, 0, 0, 0.38);
  overflow: hidden;
}

.view-toggle .v-btn {
  height: 100% !important;
  width: 44px !important;
  border: none !important;
  color: rgba(0, 0, 0, 0.6);
}

.view-toggle .v-btn.v-btn--active {
  background-color: var(--color-bg-2, #7A8B69) !important;
  color: #ffffff !important;
}

/* Upload button */
.upload-btn {
  background-color: var(--color-bg-2, #7A8B69) !important;
  color: #ffffff !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1) !important;
}

.upload-btn:hover {
  filter: brightness(0.93);
}

.drag-overlay {
  border: 2px dashed var(--color-bg-2);
  background-color: rgba(122, 139, 105, 0.15);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(210px, 1fr));
  gap: 16px;
}

.file-card {
  min-height: 165px;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.file-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.12) !important;
}

.file-name {
  font-size: 14px;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-word;
}

.border-top {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}

.file-table-row:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.cursor-pointer {
  cursor: pointer;
}
</style>