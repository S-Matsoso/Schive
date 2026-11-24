<template>
  <div class="module-details-page pa-6">
    <!-- Module Header -->
    <v-card v-if="currentModule" class="pa-6 rounded-lg mb-6" elevation="2">
      <div class="d-flex justify-space-between align-center flex-wrap ga-3">
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

        <div class="d-flex ga-2 flex-wrap">
          <v-chip color="color-bg-2" variant="tonal" prepend-icon="mdi-folder-multiple">
            {{ folders.length }} {{ folders.length === 1 ? 'folder' : 'folders' }}
          </v-chip>
          <v-chip color="color-bg-2" variant="tonal" prepend-icon="mdi-file-document-multiple">
            {{ files.length }} {{ files.length === 1 ? 'file' : 'files' }}
          </v-chip>
          <v-chip color="warning" variant="tonal" prepend-icon="mdi-clipboard-text-clock-outline">
            {{ assignments.length }} {{ assignments.length === 1 ? 'assignment' : 'assignments' }}
          </v-chip>
          <v-chip color="success" variant="tonal" prepend-icon="mdi-checkbox-marked-circle-outline">
            {{ openTasksCount }} {{ openTasksCount === 1 ? 'task' : 'tasks' }}
          </v-chip>
          <v-chip color="info" variant="tonal" prepend-icon="mdi-cards-outline">
            {{ flashcards.length }} cards
          </v-chip>
        </div>
      </div>
    </v-card>

    <v-card v-else class="pa-12 rounded-lg text-center" elevation="1">
      <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
      <p class="mt-4 text-grey">Loading module data...</p>
    </v-card>

    <!-- Main Navigation Tabs -->
    <div v-if="currentModule" class="content-container">
      <v-tabs v-model="activeTab" color="color-bg-2" class="mb-4 module-tabs" rounded="lg">
        <v-tab value="files" prepend-icon="mdi-folder-multiple-outline" class="text-none font-weight-bold">
          Files & Folders
          <v-chip size="x-small" class="ms-2 font-weight-bold" variant="flat">{{ files.length }}</v-chip>
        </v-tab>
        <v-tab value="assignments" prepend-icon="mdi-clipboard-text-clock-outline" class="text-none font-weight-bold">
          Assignments
          <v-chip size="x-small" color="warning" class="ms-2 font-weight-bold" variant="flat">{{ pendingAssignmentsCount }}</v-chip>
        </v-tab>
        <v-tab value="tasks" prepend-icon="mdi-checkbox-marked-circle-outline" class="text-none font-weight-bold">
          Tasks
          <v-chip size="x-small" color="success" class="ms-2 font-weight-bold" variant="flat">{{ openTasksCount }}</v-chip>
        </v-tab>
        <v-tab value="flashcards" prepend-icon="mdi-cards-outline" class="text-none font-weight-bold">
          Flashcards
          <v-chip size="x-small" color="info" class="ms-2 font-weight-bold" variant="flat">{{ flashcards.length }}</v-chip>
        </v-tab>
      </v-tabs>

      <v-window v-model="activeTab">
        <!-- ==================== TAB 1: FILES & FOLDERS ==================== -->
        <v-window-item value="files">
          <!-- Breadcrumb Navigation Bar -->
          <v-card class="pa-3 mb-4 rounded-lg d-flex align-center justify-space-between flex-wrap ga-2" variant="flat" color="color-bg-3">
            <div class="d-flex align-center ga-2 flex-wrap">
              <v-btn
                variant="text"
                size="small"
                class="text-none font-weight-bold px-2"
                :class="{ 'text-primary': !selectedFolderId }"
                prepend-icon="mdi-home-outline"
                @click="selectedFolderId = null"
              >
                All Files & Folders
              </v-btn>

              <template v-if="currentFolder">
                <v-icon icon="mdi-chevron-right" size="18" color="grey"></v-icon>
                <v-chip
                  color="color-bg-2"
                  variant="flat"
                  prepend-icon="mdi-folder-open"
                  class="font-weight-bold"
                >
                  {{ currentFolder.name }}
                </v-chip>

                <v-btn
                  variant="text"
                  size="small"
                  class="text-none text-caption text-grey-darken-2 ms-2"
                  prepend-icon="mdi-arrow-left"
                  @click="selectedFolderId = null"
                >
                  Back to root
                </v-btn>
              </template>
            </div>

            <!-- Upload Button -->
            <div class="d-flex align-center ga-2">
              <v-btn
                prepend-icon="mdi-upload"
                variant="flat"
                size="small"
                class="upload-btn text-none font-weight-medium rounded-lg"
                :loading="isUploading"
                @click="triggerFileInput"
              >
                Upload to {{ currentFolder ? currentFolder.name : 'Module' }}
              </v-btn>
            </div>
          </v-card>

          <!-- FOLDERS SECTION (visible at root view) -->
          <div v-if="!selectedFolderId" class="folders-section mb-6">
            <div class="d-flex justify-space-between align-center mb-3">
              <div class="d-flex align-center ga-2">
                <h2 class="text-subtitle-1 font-weight-bold">Folders</h2>
                <v-chip size="x-small" variant="flat" color="color-bg-3" class="font-weight-medium">
                  {{ folders.length }}
                </v-chip>
              </div>
            </div>

            <div v-if="isLoadingFolders" class="py-4 text-center">
              <v-progress-circular indeterminate size="24" color="color-bg-2"></v-progress-circular>
            </div>

            <div v-else class="folders-grid">
              <!-- Folder Card -->
              <v-card
                v-for="folder in folders"
                :key="folder.id"
                class="folder-card pa-3 rounded-lg d-flex align-center justify-space-between cursor-pointer"
                elevation="1"
                @click="openFolder(folder.id)"
              >
                <div class="d-flex align-center ga-3 overflow-hidden me-2">
                  <v-icon icon="mdi-folder" size="36" color="color-bg-2"></v-icon>
                  <div class="overflow-hidden">
                    <div class="font-weight-bold text-truncate text-body-2" :title="folder.name">
                      {{ folder.name }}
                    </div>
                    <div class="text-caption text-grey">
                      {{ folder.fileCount }} {{ folder.fileCount === 1 ? 'file' : 'files' }}
                    </div>
                  </div>
                </div>

                <v-btn
                  icon="mdi-delete-outline"
                  size="x-small"
                  variant="text"
                  color="grey-darken-1"
                  title="Delete folder"
                  @click.stop="confirmDeleteFolder(folder)"
                ></v-btn>
              </v-card>

              <!-- Add Folder Quick Card -->
              <div
                class="add-folder-card pa-3 rounded-lg d-flex align-center justify-center ga-2 cursor-pointer"
                @click="openCreateFolderDialog"
              >
                <v-icon icon="mdi-plus" color="color-bg-2"></v-icon>
                <span class="text-body-2 font-weight-medium text-grey-darken-2">Add Folder</span>
              </div>
            </div>
          </div>

          <!-- FILES SECTION -->
          <div class="files-container">
            <!-- Toolbar -->
            <div class="d-flex flex-wrap justify-space-between align-center ga-4 mb-4">
              <div class="d-flex align-center ga-3 flex-wrap">
                <h2 class="text-subtitle-1 font-weight-bold">
                  {{ currentFolder ? currentFolder.name + ' Files' : 'All Module Files' }}
                </h2>
                <v-chip size="x-small" variant="flat" color="color-bg-3" class="text-grey-darken-3 font-weight-medium">
                  {{ displayedFiles.length }}
                </v-chip>

                <!-- Unfiled filter toggle at root -->
                <v-btn-toggle
                  v-if="!selectedFolderId && folders.length > 0"
                  v-model="rootFilterMode"
                  mandatory
                  density="compact"
                  rounded="lg"
                  variant="outlined"
                  class="filter-toggle"
                >
                  <v-btn value="all" size="x-small" class="text-none">All</v-btn>
                  <v-btn value="unfiled" size="x-small" class="text-none">Unsorted only</v-btn>
                </v-btn-toggle>
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
              </div>
            </div>

            <!-- Drag & Drop Overlay -->
            <div
              v-if="isDragging"
              class="drag-overlay pa-8 rounded-lg text-center mb-6"
              @dragover.prevent
              @dragleave.prevent="isDragging = false"
              @drop.prevent="handleDrop"
            >
              <v-icon icon="mdi-cloud-upload" size="48" color="color-bg-2" class="mb-2"></v-icon>
              <div class="text-subtitle-1 font-weight-bold">
                Drop files here to upload to {{ currentFolder ? currentFolder.name : currentModule.code }}
              </div>
            </div>

            <!-- Loading State -->
            <div v-if="isLoadingFiles" class="py-12 text-center">
              <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
              <p class="mt-3 text-caption text-grey">Loading files...</p>
            </div>

            <!-- Empty State -->
            <v-card
              v-else-if="displayedFiles.length === 0"
              class="pa-12 text-center rounded-lg empty-card"
              variant="tonal"
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop.prevent="handleDrop"
            >
              <v-icon icon="mdi-folder-open-outline" size="56" color="grey" class="mb-3"></v-icon>
              <h3 class="text-h6 font-weight-medium mb-1">
                {{ currentFolder ? `No files in folder "${currentFolder.name}" yet` : `No files yet for ${currentModule.code}` }}
              </h3>
              <p class="text-body-2 text-grey-darken-1 mb-4">
                Upload files like documents, notes, or PDFs, or drag and drop files here.
              </p>
              <v-btn
                prepend-icon="mdi-upload"
                variant="flat"
                class="upload-btn text-none rounded-lg"
                @click="triggerFileInput"
              >
                Upload File
              </v-btn>
            </v-card>

            <!-- No Search Results State -->
            <div v-else-if="searchQuery && displayedFiles.length === 0" class="py-8 text-center text-grey">
              <v-icon icon="mdi-file-search-outline" size="36" class="mb-2"></v-icon>
              <p>No files match "{{ searchQuery }}"</p>
            </div>

            <!-- GRID VIEW -->
            <div v-else-if="viewMode === 'grid'" class="files-grid">
              <v-card
                v-for="file in displayedFiles"
                :key="file.id"
                class="file-card pa-3 rounded-lg d-flex flex-column justify-space-between"
                elevation="2"
                @click="handleOpenFile(file.id)"
              >
                <div>
                  <div class="d-flex justify-space-between align-start mb-2">
                    <v-icon
                      :icon="getFileInfo(file.fileType).icon"
                      :color="getFileInfo(file.fileType).color"
                      size="40"
                    ></v-icon>
                    <div class="d-flex align-center ga-1 flex-wrap justify-end">
                      <v-chip
                        v-if="!selectedFolderId && getFolderName(file.folderId)"
                        size="x-small"
                        variant="tonal"
                        color="color-bg-2"
                        prepend-icon="mdi-folder-outline"
                        class="font-weight-medium cursor-pointer"
                        title="Click to view folder"
                        @click.stop="openFolder(file.folderId)"
                      >
                        {{ getFolderName(file.folderId) }}
                      </v-chip>
                      <v-chip
                        size="x-small"
                        variant="tonal"
                        :color="getFileInfo(file.fileType).color"
                        class="font-weight-bold"
                      >
                        {{ getFileInfo(file.fileType).label }}
                      </v-chip>
                    </div>
                  </div>

                  <div class="file-name font-weight-bold mb-1" :title="file.name">
                    {{ file.name }}
                  </div>

                  <div class="text-caption text-grey">
                    <div>{{ formatFileSize(file.fileSize) }} • Uploaded {{ formatDate(file.createdAt) }}</div>
                    <div
                      v-if="isEditable(file.fileType)"
                      class="mt-1 d-flex align-center ga-1 text-primary font-weight-medium modification-date"
                      :title="'Last modified: ' + formatDateTime(file.updatedAt || file.createdAt)"
                    >
                      <v-icon icon="mdi-file-document-edit-outline" size="14"></v-icon>
                      <span>Modified: {{ formatDateTime(file.updatedAt || file.createdAt) }}</span>
                    </div>
                  </div>
                </div>

                <!-- Actions -->
                <div class="file-card-actions d-flex justify-space-between align-center mt-3 pt-2 border-top flex-wrap ga-1">
                  <div class="d-flex align-center ga-1 flex-nowrap">
                    <v-btn
                      variant="text"
                      size="small"
                      density="compact"
                      prepend-icon="mdi-open-in-new"
                      color="color-bg-2"
                      class="text-none px-2 font-weight-medium"
                      title="Open in default application"
                      @click.stop="handleOpenFile(file.id)"
                    >
                      Open
                    </v-btn>
                    <v-btn
                      variant="text"
                      size="small"
                      density="compact"
                      icon="mdi-eye-outline"
                      color="info"
                      title="Quick Preview"
                      @click.stop="openPreview(file)"
                    ></v-btn>
                    <v-btn
                      variant="text"
                      size="small"
                      density="compact"
                      icon="mdi-download-outline"
                      color="grey-darken-2"
                      title="Export a copy"
                      @click.stop="handleExportFile(file)"
                    ></v-btn>
                  </div>

                  <div class="d-flex align-center ga-1 flex-nowrap">
                    <v-btn
                      variant="text"
                      size="small"
                      density="compact"
                      icon="mdi-pencil-outline"
                      color="grey-darken-2"
                      title="Rename file"
                      @click.stop="openRenameDialog(file)"
                    ></v-btn>
                    <v-btn
                      variant="text"
                      size="small"
                      density="compact"
                      icon="mdi-folder-move-outline"
                      color="grey-darken-2"
                      title="Move to folder"
                      @click.stop="openMoveDialog(file)"
                    ></v-btn>
                    <v-btn
                      variant="text"
                      size="small"
                      density="compact"
                      icon="mdi-delete-outline"
                      color="error"
                      title="Delete File"
                      @click.stop="confirmDeleteFile(file)"
                    ></v-btn>
                  </div>
                </div>
              </v-card>
            </div>

            <!-- LIST VIEW -->
            <v-card v-else class="rounded-lg overflow-hidden" elevation="2">
              <v-table hover>
                <thead>
                  <tr>
                    <th style="width: 44px;"></th>
                    <th>Name</th>
                    <th v-if="!selectedFolderId" style="width: 140px;">Folder</th>
                    <th style="width: 140px;">Date Uploaded</th>
                    <th style="width: 190px;">Last Modified</th>
                    <th style="width: 220px; min-width: 220px;" class="text-right text-no-wrap">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="file in displayedFiles"
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
                    <td v-if="!selectedFolderId">
                      <v-chip
                        v-if="getFolderName(file.folderId)"
                        size="x-small"
                        variant="tonal"
                        color="color-bg-2"
                        prepend-icon="mdi-folder-outline"
                        class="cursor-pointer"
                        @click.stop="openFolder(file.folderId)"
                      >
                        {{ getFolderName(file.folderId) }}
                      </v-chip>
                      <span v-else class="text-caption text-grey">—</span>
                    </td>
                    <td class="text-caption text-grey-darken-1">
                      {{ formatDate(file.createdAt) }}
                    </td>
                    <td class="text-caption">
                      <div
                        v-if="isEditable(file.fileType)"
                        class="font-weight-medium text-primary d-flex align-center ga-1"
                      >
                        <v-icon icon="mdi-file-document-edit-outline" size="14"></v-icon>
                        <span>{{ formatDateTime(file.updatedAt || file.createdAt) }}</span>
                      </div>
                      <span v-else class="text-grey-darken-1">
                        {{ formatDate(file.updatedAt || file.createdAt) }}
                      </span>
                    </td>
                    <td class="text-right text-no-wrap" style="white-space: nowrap;">
                      <div class="d-inline-flex align-center justify-end ga-1 flex-nowrap">
                        <v-btn
                          icon="mdi-eye-outline"
                          variant="text"
                          size="small"
                          density="compact"
                          color="info"
                          title="Quick preview"
                          @click.stop="openPreview(file)"
                        ></v-btn>
                        <v-btn
                          icon="mdi-download-outline"
                          variant="text"
                          size="small"
                          density="compact"
                          color="grey-darken-2"
                          title="Export a copy"
                          @click.stop="handleExportFile(file)"
                        ></v-btn>
                        <v-btn
                          icon="mdi-pencil-outline"
                          variant="text"
                          size="small"
                          density="compact"
                          color="grey-darken-2"
                          title="Rename file"
                          @click.stop="openRenameDialog(file)"
                        ></v-btn>
                        <v-btn
                          icon="mdi-folder-move-outline"
                          variant="text"
                          size="small"
                          density="compact"
                          color="grey-darken-2"
                          title="Move to folder"
                          @click.stop="openMoveDialog(file)"
                        ></v-btn>
                        <v-btn
                          icon="mdi-delete-outline"
                          variant="text"
                          size="small"
                          density="compact"
                          color="error"
                          title="Delete file"
                          @click.stop="confirmDeleteFile(file)"
                        ></v-btn>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </v-table>
            </v-card>
          </div>
        </v-window-item>

        <!-- ==================== TAB 2: ASSIGNMENTS ==================== -->
        <v-window-item value="assignments">
          <div class="assignments-tab">
            <!-- Assignments Toolbar -->
            <div class="d-flex justify-space-between align-center flex-wrap ga-3 mb-4">
              <div class="d-flex align-center ga-2">
                <v-btn-toggle
                  v-model="assignmentFilter"
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
              </div>

              <v-btn
                prepend-icon="mdi-plus"
                color="color-bg-2"
                variant="flat"
                class="text-none font-weight-medium rounded-lg"
                @click="openCreateAssignmentDialog"
              >
                New Assignment
              </v-btn>
            </div>

            <!-- Empty State -->
            <v-card v-if="filteredAssignments.length === 0" class="pa-12 text-center rounded-lg" variant="tonal">
              <v-icon icon="mdi-clipboard-check-outline" size="56" color="grey" class="mb-3"></v-icon>
              <h3 class="text-h6 font-weight-medium mb-1">No assignments found</h3>
              <p class="text-body-2 text-grey mb-4">Track homework, lab reports, and coursework deadlines for this module.</p>
              <v-btn color="color-bg-2" variant="flat" prepend-icon="mdi-plus" class="text-none rounded-lg" @click="openCreateAssignmentDialog">
                Create Assignment
              </v-btn>
            </v-card>

            <!-- Assignments List -->
            <div v-else class="d-flex flex-column ga-3">
              <v-card
                v-for="assign in filteredAssignments"
                :key="assign.id"
                class="pa-4 rounded-lg d-flex justify-space-between align-center flex-wrap ga-3"
                elevation="1"
              >
                <div class="d-flex align-start ga-3 flex-grow-1">
                  <v-btn
                    :icon="assign.status === 'completed' ? 'mdi-checkbox-marked-circle' : 'mdi-checkbox-blank-circle-outline'"
                    :color="assign.status === 'completed' ? 'success' : 'grey'"
                    variant="text"
                    density="compact"
                    size="small"
                    class="mt-1"
                    :title="assign.status === 'completed' ? 'Mark Incomplete' : 'Mark Complete'"
                    @click.stop="updateStatus(assign.id, assign.status === 'completed' ? 'pending' : 'completed')"
                  ></v-btn>
                  <div>
                    <div class="d-flex align-center ga-2 flex-wrap">
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
                        {{ assign.status }}
                      </v-chip>
                    </div>

                    <p v-if="assign.description" class="text-body-2 text-grey-darken-1 mt-1 mb-0">
                      {{ assign.description }}
                    </p>

                    <div v-if="assign.files?.length" class="mt-2 d-flex flex-wrap ga-2">
                      <v-chip
                        v-for="f in assign.files"
                        :key="f.id"
                        size="small"
                        variant="tonal"
                        color="color-bg-2"
                        prepend-icon="mdi-paperclip"
                        class="cursor-pointer"
                        @click="handleOpenFile(f.id)"
                      >
                        {{ f.name }}
                      </v-chip>
                    </div>

                    <div class="text-caption text-grey mt-2 d-flex align-center ga-3">
                      <span v-if="assign.dueDate" class="d-flex align-center ga-1 font-weight-medium" :class="isOverdue(assign.dueDate) ? 'text-error' : 'text-primary'">
                        <v-icon icon="mdi-calendar-clock" size="14"></v-icon>
                        Due: {{ formatDueDateDisplay(assign.dueDate) }}
                      </span>
                      <span>Created {{ formatDate(assign.createdAt) }}</span>
                    </div>

                    <div class="mt-2" style="min-width: 240px; max-width: 360px;">
                      <div class="d-flex justify-space-between text-caption mb-1">
                        <span class="text-grey">Progress</span>
                        <span class="font-weight-bold">{{ assign.progress || 0 }}%</span>
                      </div>
                      <v-progress-linear
                        :model-value="assign.progress || 0"
                        :color="getProgressColor(assign.progress)"
                        height="6"
                        rounded
                      ></v-progress-linear>
                    </div>
                  </div>
                </div>

                <!-- Assignment Actions -->
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
                    @click="updateStatus(assign.id, 'completed')"
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
                    @click="updateStatus(assign.id, 'pending')"
                  >
                    Mark Incomplete
                  </v-btn>

                  <v-btn
                    icon="mdi-pencil-outline"
                    size="small"
                    variant="text"
                    title="Edit assignment"
                    @click="openEditAssignmentDialog(assign)"
                  ></v-btn>
                  <v-btn
                    icon="mdi-delete-outline"
                    size="small"
                    variant="text"
                    color="error"
                    title="Delete assignment"
                    @click="confirmDeleteAssignment(assign)"
                  ></v-btn>
                </div>
              </v-card>
            </div>
          </div>
        </v-window-item>

        <!-- ==================== TAB 3: TASKS ==================== -->
        <v-window-item value="tasks">
          <div class="tasks-tab">
            <div class="d-flex justify-space-between align-center flex-wrap ga-3 mb-4">
              <v-btn-toggle
                v-model="taskFilter"
                mandatory
                density="compact"
                rounded="lg"
                variant="outlined"
                class="filter-toggle"
              >
                <v-btn value="open" size="small" class="text-none">Open ({{ openTasksCount }})</v-btn>
                <v-btn value="done" size="small" class="text-none">Done ({{ tasks.length - openTasksCount }})</v-btn>
                <v-btn value="all" size="small" class="text-none">All ({{ tasks.length }})</v-btn>
              </v-btn-toggle>

              <v-btn
                prepend-icon="mdi-plus"
                color="color-bg-2"
                variant="flat"
                class="text-none font-weight-medium rounded-lg"
                @click="openCreateTaskDialog"
              >
                New Task
              </v-btn>
            </div>

            <v-card v-if="filteredTasks.length === 0" class="pa-12 text-center rounded-lg" variant="tonal">
              <v-icon icon="mdi-checkbox-marked-circle-outline" size="56" color="grey" class="mb-3"></v-icon>
              <h3 class="text-h6 font-weight-medium mb-1">
                {{ tasks.length === 0 ? 'No tasks yet' : 'Nothing to show' }}
              </h3>
              <p class="text-body-2 text-grey mb-4">
                Quick to-dos for this module, kept separate from graded assignments and their instructions.
              </p>
              <v-btn color="color-bg-2" variant="flat" prepend-icon="mdi-plus" class="text-none rounded-lg" @click="openCreateTaskDialog">
                Create Task
              </v-btn>
            </v-card>

            <div v-else class="d-flex flex-column ga-2">
              <v-card
                v-for="task in filteredTasks"
                :key="task.id"
                class="pa-3 rounded-lg task-row"
                elevation="1"
              >
                <div class="d-flex align-center ga-2">
                  <v-checkbox-btn
                    :model-value="task.done"
                    color="success"
                    hide-details
                    @update:model-value="toggleTask(task)"
                  ></v-checkbox-btn>

                  <div class="flex-grow-1 overflow-hidden">
                    <div class="d-flex align-center ga-2 flex-wrap">
                      <span class="font-weight-medium" :class="{ 'task-title--done': task.done }">{{ task.title }}</span>
                      <v-chip
                        v-if="task.dueDate"
                        size="x-small"
                        variant="tonal"
                        :color="isTaskOverdue(task) ? 'error' : 'grey'"
                        prepend-icon="mdi-calendar-clock"
                      >
                        {{ formatDueDateDisplay(task.dueDate) }}{{ isTaskOverdue(task) ? ' · overdue' : '' }}
                      </v-chip>
                    </div>
                    <p v-if="task.notes" class="text-caption text-grey mt-1 mb-0 task-notes">{{ task.notes }}</p>
                  </div>

                  <v-btn icon="mdi-pencil-outline" size="small" variant="text" title="Edit task" @click="openEditTaskDialog(task)"></v-btn>
                  <v-btn icon="mdi-delete-outline" size="small" variant="text" color="error" title="Delete task" @click="askDeleteTask(task)"></v-btn>
                </div>
              </v-card>
            </div>
          </div>
        </v-window-item>

        <!-- ==================== TAB 4: FLASHCARDS ==================== -->
        <v-window-item value="flashcards">
          <div class="flashcards-tab">
            <!-- Flashcards Toolbar -->
            <div class="d-flex justify-space-between align-center flex-wrap ga-3 mb-4">
              <div class="d-flex align-center ga-3">
                <v-chip color="info" variant="tonal" prepend-icon="mdi-cards-outline" class="font-weight-medium">
                  {{ flashcards.length }} {{ flashcards.length === 1 ? 'Card' : 'Cards' }}
                </v-chip>
                <v-chip color="success" variant="tonal" prepend-icon="mdi-check-decagram-outline" class="font-weight-medium">
                  {{ masteredFlashcardsCount }} Mastered
                </v-chip>
              </div>

              <div class="d-flex ga-2">
                <v-btn
                  v-if="flashcards.length > 0"
                  color="color-bg-2"
                  variant="flat"
                  prepend-icon="mdi-play-circle-outline"
                  class="text-none font-weight-medium rounded-lg"
                  @click="startStudySession"
                >
                  Study Deck
                </v-btn>
                <v-btn
                  prepend-icon="mdi-plus"
                  variant="outlined"
                  class="text-none font-weight-medium rounded-lg"
                  @click="openCreateFlashcardDialog"
                >
                  Add Card
                </v-btn>
              </div>
            </div>

            <!-- Empty State -->
            <v-card v-if="flashcards.length === 0" class="pa-12 text-center rounded-lg" variant="tonal">
              <v-icon icon="mdi-cards-outline" size="56" color="grey" class="mb-3"></v-icon>
              <h3 class="text-h6 font-weight-medium mb-1">No flashcards yet</h3>
              <p class="text-body-2 text-grey mb-4">Create questions and answers to review key concepts and study for exams.</p>
              <v-btn color="color-bg-2" variant="flat" prepend-icon="mdi-plus" class="text-none rounded-lg" @click="openCreateFlashcardDialog">
                Create First Flashcard
              </v-btn>
            </v-card>

            <!-- Flashcards Grid -->
            <div v-else class="flashcards-grid">
              <v-card
                v-for="card in flashcards"
                :key="card.id"
                class="flashcard-preview-card pa-4 rounded-lg d-flex flex-column justify-space-between"
                :class="{ 'card-mastered': card.mastered }"
                elevation="1"
              >
                <div>
                  <div class="d-flex justify-space-between align-center mb-2">
                    <v-chip
                      size="x-small"
                      :color="card.mastered ? 'success' : 'grey'"
                      variant="flat"
                      class="font-weight-bold"
                    >
                      {{ card.mastered ? 'Mastered' : 'To Review' }}
                    </v-chip>

                    <div class="d-flex ga-1">
                      <v-btn
                        icon="mdi-pencil-outline"
                        size="x-small"
                        variant="text"
                        @click="openEditFlashcardDialog(card)"
                      ></v-btn>
                      <v-btn
                        icon="mdi-delete-outline"
                        size="x-small"
                        variant="text"
                        color="error"
                        @click="confirmDeleteFlashcard(card)"
                      ></v-btn>
                    </div>
                  </div>

                  <div class="font-weight-bold text-body-1 mb-2">
                    {{ card.question }}
                  </div>

                  <div class="text-body-2 text-grey-darken-1 card-answer-text">
                    {{ card.answer }}
                  </div>
                </div>

                <div class="d-flex justify-end mt-3 pt-2 border-top">
                  <v-btn
                    size="small"
                    variant="text"
                    :color="card.mastered ? 'success' : 'grey-darken-1'"
                    :prepend-icon="card.mastered ? 'mdi-check-circle' : 'mdi-circle-outline'"
                    class="text-none"
                    @click="toggleCardMastered(card.id)"
                  >
                    {{ card.mastered ? 'Mastered' : 'Mark as Mastered' }}
                  </v-btn>
                </div>
              </v-card>
            </div>
          </div>
        </v-window-item>
      </v-window>
    </div>

    <!-- ==================== DIALOGS ==================== -->

    <!-- Interactive Study Deck Modal -->
    <v-dialog v-model="studyModal" width="560">
      <v-card v-if="studyCards.length > 0" class="pa-6 rounded-xl study-session-card" elevation="8">
        <div class="d-flex justify-space-between align-center mb-3">
          <span class="text-caption font-weight-bold text-grey">CARD {{ currentStudyIndex + 1 }} OF {{ studyCards.length }}</span>
          <v-btn icon="mdi-close" size="small" variant="text" @click="studyModal = false"></v-btn>
        </div>

        <v-progress-linear
          :model-value="((currentStudyIndex + 1) / studyCards.length) * 100"
          color="color-bg-2"
          rounded
          height="6"
          class="mb-4"
        ></v-progress-linear>

        <!-- Interactive Flip Card -->
        <div
          class="study-flip-card pa-8 rounded-xl d-flex flex-column align-center justify-center text-center cursor-pointer mb-6"
          @click="showStudyAnswer = !showStudyAnswer"
        >
          <div class="text-caption font-weight-bold text-uppercase mb-2 text-grey">
            {{ showStudyAnswer ? 'Answer' : 'Question (Click card to reveal answer)' }}
          </div>
          <div class="text-h5 font-weight-bold">
            {{ showStudyAnswer ? studyCards[currentStudyIndex].answer : studyCards[currentStudyIndex].question }}
          </div>
        </div>

        <!-- Controls -->
        <div class="d-flex justify-space-between align-center ga-3">
          <v-btn
            variant="outlined"
            :disabled="currentStudyIndex === 0"
            prepend-icon="mdi-arrow-left"
            class="text-none"
            @click="prevStudyCard"
          >
            Previous
          </v-btn>

          <v-btn
            color="success"
            variant="tonal"
            prepend-icon="mdi-check"
            class="text-none"
            @click="markCurrentStudyMastered"
          >
            I Know This
          </v-btn>

          <v-btn
            color="color-bg-2"
            variant="flat"
            append-icon="mdi-arrow-right"
            class="text-none"
            @click="nextStudyCard"
          >
            {{ currentStudyIndex === studyCards.length - 1 ? 'Finish' : 'Next' }}
          </v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- In-App Quick File Preview Dialog -->
    <v-dialog v-model="previewDialog" width="700" scrollable>
      <v-card v-if="previewFile" class="rounded-xl overflow-hidden" elevation="6">
        <v-card-title class="d-flex justify-space-between align-center pa-4 border-bottom">
          <div class="d-flex align-center ga-2 text-truncate">
            <v-icon :icon="getFileInfo(previewFile.fileType).icon" :color="getFileInfo(previewFile.fileType).color"></v-icon>
            <span class="text-subtitle-1 font-weight-bold text-truncate">{{ previewFile.name }}</span>
          </div>
          <div class="d-flex align-center ga-2">
            <v-btn size="small" variant="tonal" prepend-icon="mdi-open-in-new" color="color-bg-2" class="text-none" @click="handleOpenFile(previewFile.id)">
              Open in App
            </v-btn>
            <v-btn icon="mdi-close" size="small" variant="text" @click="previewDialog = false"></v-btn>
          </div>
        </v-card-title>

        <v-card-text class="pa-4 preview-content-container">
          <div v-if="isLoadingPreview" class="py-12 text-center">
            <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
            <p class="mt-2 text-caption text-grey">Loading preview...</p>
          </div>

          <!-- Image Preview -->
          <div v-else-if="['png', 'jpg', 'jpeg', 'gif', 'svg', 'webp'].includes(previewFile.fileType)" class="text-center">
            <img :src="'data:image/' + previewFile.fileType + ';base64,' + previewContent" class="preview-image rounded-lg" />
          </div>

          <!-- PDF Preview -->
          <div v-else-if="previewFile.fileType === 'pdf'" style="height: 520px;">
            <iframe :src="'data:application/pdf;base64,' + previewContent" width="100%" height="100%" style="border: none;"></iframe>
          </div>

          <!-- Text / Code Preview -->
          <div v-else-if="['txt', 'md', 'csv', 'json', 'py', 'js', 'html', 'css'].includes(previewFile.fileType)">
            <pre class="preview-text pa-4 rounded-lg">{{ decodeBase64ToText(previewContent) }}</pre>
          </div>

          <!-- Fallback -->
          <div v-else class="py-12 text-center text-grey">
            <v-icon icon="mdi-file-eye-outline" size="48" class="mb-2"></v-icon>
            <div class="font-weight-medium">Quick preview is not available for .{{ previewFile.fileType }} files.</div>
            <div class="text-caption mb-4">Click below to open in your system's default viewer or editor.</div>
            <v-btn color="color-bg-2" variant="flat" prepend-icon="mdi-open-in-new" class="text-none rounded-lg" @click="handleOpenFile(previewFile.id)">
              Open in System Editor
            </v-btn>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Rename File Dialog -->
    <v-dialog v-model="renameDialog" width="400">
      <v-card v-if="fileToRename" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-3">Rename File</h3>
        <v-text-field
          v-model="newFileName"
          label="File name"
          variant="outlined"
          density="comfortable"
          class="mb-4"
          autofocus
        ></v-text-field>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="renameDialog = false">Cancel</v-btn>
          <v-btn color="color-bg-2" variant="flat" :loading="isRenaming" @click="executeRenameFile">Save</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Assignment Details Dialog -->
    <AssignmentDetailsDialog
      v-model="assignmentDetailDialog"
      :assignment="detailAssignment"
      @edit="openEditAssignmentDialog"
      @files-changed="handleAssignmentFilesChanged"
      @status-changed="fetchAssignments"
    />

    <!-- Create / Edit Assignment Dialog -->
    <v-dialog v-model="assignmentDialog" width="480">
      <v-card class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-3">{{ editingAssignment ? 'Edit Assignment' : 'New Assignment' }}</h3>
        <v-form @submit.prevent="saveAssignment">
          <v-text-field v-model="assignmentForm.title" label="Title" variant="outlined" density="comfortable" class="mb-2" required></v-text-field>
          <v-textarea v-model="assignmentForm.description" label="Description" variant="outlined" density="comfortable" rows="2" class="mb-2"></v-textarea>
          <v-row>
            <v-col cols="6">
              <v-select v-model="assignmentForm.priority" label="Priority" :items="['low', 'medium', 'high']" variant="outlined" density="comfortable"></v-select>
            </v-col>
            <v-col cols="6">
              <v-select v-model="assignmentForm.status" label="Status" :items="['pending', 'in_progress', 'completed']" variant="outlined" density="comfortable"></v-select>
            </v-col>
          </v-row>
          <v-text-field v-model="assignmentForm.dueDate" type="date" label="Due Date" variant="outlined" density="comfortable" class="mb-4"></v-text-field>

          <v-file-input
            v-if="!editingAssignment"
            v-model="assignmentFiles"
            label="Assignment instructions (optional)"
            variant="outlined"
            density="comfortable"
            prepend-icon=""
            prepend-inner-icon="mdi-paperclip"
            show-size
            hint="Saved to the Assignments folder in Files"
            persistent-hint
            class="mb-5"
          ></v-file-input>

          <div class="d-flex justify-end ga-2">
            <v-btn variant="text" @click="assignmentDialog = false">Cancel</v-btn>
            <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isSavingAssignment">Save</v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- Update Assignment Progress Dialog -->
    <v-dialog v-model="progressDialog" width="440">
      <v-card v-if="progressTask" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-1">Update Progress</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4 text-truncate">{{ progressTask.title }}</p>

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
            v-for="preset in [0, 25, 50, 75, 100]"
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

        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="progressDialog = false">Cancel</v-btn>
          <v-btn color="color-bg-2" variant="flat" :loading="isSavingProgress" @click="saveProgress">
            Save Progress
          </v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Create / Edit Flashcard Dialog -->
    <v-dialog v-model="flashcardDialog" width="480">
      <v-card class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-3">{{ editingFlashcard ? 'Edit Flashcard' : 'New Flashcard' }}</h3>
        <v-form @submit.prevent="saveFlashcard">
          <v-textarea v-model="flashcardForm.question" label="Question / Prompt" variant="outlined" density="comfortable" rows="3" class="mb-2" required></v-textarea>
          <v-textarea v-model="flashcardForm.answer" label="Answer / Solution" variant="outlined" density="comfortable" rows="3" class="mb-4" required></v-textarea>
          <div class="d-flex justify-end ga-2">
            <v-btn variant="text" @click="flashcardDialog = false">Cancel</v-btn>
            <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isSavingFlashcard">Save Card</v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- Create Folder Dialog -->
    <v-dialog v-model="createFolderDialog" width="400">
      <v-card class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-1">New Folder</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">Organise files in {{ currentModule?.code }} (e.g. Notes, Assignments, Flash Cards).</p>
        <v-form ref="folderFormRef" @submit.prevent="handleCreateFolder">
          <v-text-field v-model="newFolderName" label="Folder name" variant="outlined" density="comfortable" :rules="folderNameRules" class="mb-4" autofocus required></v-text-field>
          <div class="d-flex justify-end ga-2">
            <v-btn variant="text" @click="createFolderDialog = false">Cancel</v-btn>
            <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isCreatingFolder">Create Folder</v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- Move File Dialog -->
    <v-dialog v-model="moveDialog" width="400">
      <v-card v-if="fileToMove" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-1">Move File</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4 text-truncate">Move <strong>{{ fileToMove.name }}</strong> to:</p>
        <v-select v-model="targetFolderSelection" label="Destination Folder" variant="outlined" density="comfortable" :items="folderSelectOptions" item-title="title" item-value="id" class="mb-4"></v-select>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="moveDialog = false">Cancel</v-btn>
          <v-btn color="color-bg-2" variant="flat" :loading="isMoving" @click="executeMoveFile">Move File</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Delete Folder Confirmation Dialog -->
    <v-dialog v-model="deleteFolderDialog" width="420">
      <v-card v-if="folderToDelete" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-2 text-error">Delete Folder?</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">
          Are you sure you want to delete folder <strong>{{ folderToDelete.name }}</strong>?
          <span v-if="folderToDelete.fileCount > 0" class="text-error d-block mt-2">
            <strong>Warning:</strong> All {{ folderToDelete.fileCount }} files inside this folder will also be deleted.
          </span>
        </p>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="deleteFolderDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="flat" :loading="isDeletingFolder" @click="executeDeleteFolder">Delete Folder</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Delete File Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" width="380">
      <v-card v-if="fileToDelete" class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-2">Delete File?</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">
          Are you sure you want to delete <strong>{{ fileToDelete.name }}</strong>? This action cannot be undone.
        </p>
        <div class="d-flex justify-end ga-2">
          <v-btn type="button" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="flat" :loading="isDeleting" @click="executeDeleteFile">Delete</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Create / Edit Task Dialog -->
    <v-dialog v-model="taskDialog" width="460">
      <v-card class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-1">{{ editingTask ? 'Edit Task' : 'New Task' }}</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">
          A simple to-do for this module. Graded work belongs in the Assignments tab.
        </p>
        <v-form @submit.prevent="saveTask">
          <v-text-field
            v-model="taskForm.title"
            label="Title"
            variant="outlined"
            density="comfortable"
            class="mb-2"
            autofocus
          ></v-text-field>
          <v-textarea
            v-model="taskForm.notes"
            label="Notes"
            variant="outlined"
            density="comfortable"
            rows="2"
            class="mb-2"
          ></v-textarea>
          <v-text-field
            v-model="taskForm.dueDate"
            type="date"
            label="Due date (optional)"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-text-field>

          <div class="d-flex justify-end ga-2">
            <v-btn variant="text" class="text-none" @click="taskDialog = false">Cancel</v-btn>
            <v-btn type="submit" color="color-bg-2" variant="flat" class="text-none" :loading="isSavingTask">
              {{ editingTask ? 'Save Changes' : 'Create Task' }}
            </v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- Confirmation Dialog -->
    <ConfirmDialog
      v-model="confirmDialog"
      :title="confirmTitle"
      :message="confirmMessage"
      :confirm-text="confirmText"
      :tone="confirmTone"
      :icon="confirmIcon"
      :loading="confirmLoading"
      @confirm="runConfirmedAction"
    />

    <!-- Snackbar -->
    <v-snackbar v-model="showAlert" :color="alertType" timeout="3000">
      {{ alertMsg }}
    </v-snackbar>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { models } from '../../wailsjs/go/models'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import AssignmentDetailsDialog from '../components/AssignmentDetailsDialog.vue'
import { exportSchiveFile } from '../utils/exportFile'
import { GetModule } from '../../wailsjs/go/controllers/ModuleController'
import {
  UploadFile,
  UploadAssignmentInstructions,
  GetFilesByModule,
  DeleteFile,
  OpenFile,
  CreateFolder,
  GetFoldersByModule,
  DeleteFolder,
  MoveFile,
  RenameFile,
  GetFileContent
} from '../../wailsjs/go/controllers/FileController'
import {
  CreateAssignment,
  GetAssignmentsByModule,
  UpdateAssignment,
  UpdateAssignmentProgress,
  UpdateAssignmentStatus,
  DeleteAssignment
} from '../../wailsjs/go/controllers/AssignmentController'
import {
  CreateTask,
  GetTasksByModule,
  UpdateTask,
  SetTaskDone,
  DeleteTask
} from '../../wailsjs/go/controllers/TaskController'
import {
  CreateFlashcard,
  GetFlashcardsByModule,
  ToggleMastered,
  UpdateFlashcard,
  DeleteFlashcard
} from '../../wailsjs/go/controllers/FlashcardController'
import { eventBus, SYNC_EVENTS } from '../utils/eventBus'

const route = useRoute()
const moduleCode = computed(() => route.params.code as string)

const activeTab = ref('files')
const currentModule = ref<models.Module | null>(null)
const files = ref<models.File[]>([])
const folders = ref<models.Folder[]>([])
const assignments = ref<models.Assignment[]>([])
const flashcards = ref<models.Flashcard[]>([])

const selectedFolderId = ref<string | null>(null)
const rootFilterMode = ref<'all' | 'unfiled'>('all')
const viewMode = ref<'grid' | 'list'>('grid')
const searchQuery = ref('')

const isUploading = ref(false)
const isLoadingFiles = ref(false)
const isLoadingFolders = ref(false)
const isDragging = ref(false)
const isDeleting = ref(false)

const fileInputRef = ref<HTMLInputElement | null>(null)
const deleteDialog = ref(false)
const fileToDelete = ref<models.File | null>(null)

// Folder dialogs & state
const createFolderDialog = ref(false)
const isCreatingFolder = ref(false)
const newFolderName = ref('')
const folderFormRef = ref<any>(null)
const deleteFolderDialog = ref(false)
const isDeletingFolder = ref(false)
const folderToDelete = ref<models.Folder | null>(null)

// Move file dialog & state
const moveDialog = ref(false)
const isMoving = ref(false)
const fileToMove = ref<models.File | null>(null)
const targetFolderSelection = ref<string>('')

// Rename file dialog
const renameDialog = ref(false)
const isRenaming = ref(false)
const fileToRename = ref<models.File | null>(null)
const newFileName = ref('')

// In-app preview
const previewDialog = ref(false)
const previewFile = ref<models.File | null>(null)
const previewContent = ref('')
const isLoadingPreview = ref(false)

// Assignments state & dialog
const assignmentFilter = ref('all')
const assignmentDialog = ref(false)
const isSavingAssignment = ref(false)
const editingAssignment = ref<models.Assignment | null>(null)
const progressDialog = ref(false)
const progressTask = ref<models.Assignment | null>(null)
const progressValue = ref(0)
const isSavingProgress = ref(false)
const assignmentForm = ref({
  title: '',
  description: '',
  priority: 'medium',
  status: 'pending',
  dueDate: ''
})
const assignmentFiles = ref<File[] | File | null>(null)

// Assignment details state
const assignmentDetailDialog = ref(false)
const detailAssignment = ref<models.Assignment | null>(null)

// Tasks state & dialog (separate from assignments)
const tasks = ref<models.Task[]>([])
const taskFilter = ref<'open' | 'done' | 'all'>('open')
const taskDialog = ref(false)
const isSavingTask = ref(false)
const editingTask = ref<models.Task | null>(null)
const taskForm = ref({
  title: '',
  notes: '',
  dueDate: ''
})

// Generic confirmation dialog state
const confirmDialog = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmText = ref('Delete')
const confirmTone = ref<'error' | 'warning'>('error')
const confirmIcon = ref('')
const confirmLoading = ref(false)
let pendingConfirm: (() => Promise<void>) | null = null

// Flashcards state & dialog
const flashcardDialog = ref(false)
const isSavingFlashcard = ref(false)
const editingFlashcard = ref<models.Flashcard | null>(null)
const flashcardForm = ref({
  question: '',
  answer: ''
})

// Study session state
const studyModal = ref(false)
const studyCards = ref<models.Flashcard[]>([])
const currentStudyIndex = ref(0)
const showStudyAnswer = ref(false)

const alertMsg = ref('')
const alertType = ref('success')
const showAlert = ref(false)

const EDITABLE_EXTENSIONS = new Set([
  'doc', 'docx', 'odt', 'rtf', 'txt', 'md',
  'xls', 'xlsx', 'ods', 'csv',
  'ppt', 'pptx', 'odp'
])

function isEditable(ext: string): boolean {
  return EDITABLE_EXTENSIONS.has((ext || '').toLowerCase())
}

const currentFolder = computed(() => {
  if (!selectedFolderId.value) return null
  return folders.value.find(f => f.id === selectedFolderId.value) || null
})

function getFolderName(folderId?: string): string {
  if (!folderId) return ''
  const folder = folders.value.find(f => f.id === folderId)
  return folder ? folder.name : ''
}

const folderSelectOptions = computed(() => {
  return [
    { title: '📁 (None / Root)', id: '' },
    ...folders.value.map(f => ({ title: `📁 ${f.name}`, id: f.id }))
  ]
})

const folderNameRules = [
  (v: string) => (v && v.trim().length > 0) || 'Folder name is required',
  (v: string) => {
    if (!v) return true
    const trimmed = v.trim().toLowerCase()
    const duplicate = folders.value.some(f => f.name.toLowerCase() === trimmed)
    return !duplicate || 'A folder with this name already exists'
  }
]

const displayedFiles = computed(() => {
  let list = files.value
  if (selectedFolderId.value) {
    list = list.filter(f => f.folderId === selectedFolderId.value)
  } else if (rootFilterMode.value === 'unfiled') {
    list = list.filter(f => !f.folderId)
  }
  if (searchQuery.value?.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter(f => f.name.toLowerCase().includes(q))
  }
  return list
})

const pendingAssignmentsCount = computed(() => {
  return assignments.value.filter(a => a.status !== 'completed').length
})

const openTasksCount = computed(() => tasks.value.filter(t => !t.done).length)

const filteredTasks = computed(() => {
  if (taskFilter.value === 'open') return tasks.value.filter(t => !t.done)
  if (taskFilter.value === 'done') return tasks.value.filter(t => t.done)
  return tasks.value
})

const filteredAssignments = computed(() => {
  if (assignmentFilter.value === 'all') return assignments.value
  return assignments.value.filter(a => a.status === assignmentFilter.value)
})

const masteredFlashcardsCount = computed(() => {
  return flashcards.value.filter(c => c.mastered).length
})

async function fetchModuleDetails() {
  try {
    currentModule.value = await GetModule(moduleCode.value)
    if (currentModule.value?.id) {
      await Promise.all([
        fetchFolders(),
        fetchFiles(),
        fetchAssignments(),
        fetchTasks(),
        fetchFlashcards()
      ])
    }
  } catch (err: any) {
    alertMsg.value = 'Failed to get module details: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

async function fetchFolders() {
  if (!currentModule.value?.id) return
  try {
    isLoadingFolders.value = true
    folders.value = await GetFoldersByModule(currentModule.value.id) || []
  } catch (err: any) {
    console.error('Failed to load folders:', err)
  } finally {
    isLoadingFolders.value = false
  }
}

async function fetchFiles() {
  if (!currentModule.value?.id) return
  try {
    isLoadingFiles.value = true
    files.value = await GetFilesByModule(currentModule.value.id) || []
  } catch (err: any) {
    console.error('Failed to load files:', err)
  } finally {
    isLoadingFiles.value = false
  }
}

async function fetchAssignments() {
  if (!currentModule.value?.id) return
  try {
    assignments.value = await GetAssignmentsByModule(currentModule.value.id) || []
    if (detailAssignment.value) {
      const id = detailAssignment.value.id
      detailAssignment.value = assignments.value.find(a => a.id === id) || detailAssignment.value
    }
  } catch (err: any) {
    console.error('Failed to load assignments:', err)
  }
}

async function fetchTasks() {
  if (!currentModule.value?.id) return
  try {
    tasks.value = await GetTasksByModule(currentModule.value.id) || []
  } catch (err: any) {
    console.error('Failed to load tasks:', err)
  }
}

async function fetchFlashcards() {
  if (!currentModule.value?.id) return
  try {
    flashcards.value = await GetFlashcardsByModule(currentModule.value.id) || []
  } catch (err: any) {
    console.error('Failed to load flashcards:', err)
  }
}

function openFolder(folderId?: string) {
  if (folderId) {
    selectedFolderId.value = folderId
  }
}

function openCreateFolderDialog() {
  newFolderName.value = ''
  if (folderFormRef.value) folderFormRef.value.resetValidation()
  createFolderDialog.value = true
}

async function handleCreateFolder() {
  if (!currentModule.value?.id) return
  if (folderFormRef.value) {
    const { valid } = await folderFormRef.value.validate()
    if (!valid) return
  }
  const cleanName = newFolderName.value.trim()
  if (!cleanName) return

  try {
    isCreatingFolder.value = true
    const newFolder = await CreateFolder(currentModule.value.id, cleanName)
    createFolderDialog.value = false
    alertMsg.value = `Folder "${newFolder.name}" created`
    alertType.value = 'success'
    showAlert.value = true
    await fetchFolders()
  } catch (err: any) {
    alertMsg.value = 'Failed to create folder: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isCreatingFolder.value = false
  }
}

function confirmDeleteFolder(folder: models.Folder) {
  folderToDelete.value = folder
  deleteFolderDialog.value = true
}

async function executeDeleteFolder() {
  if (!folderToDelete.value) return
  try {
    isDeletingFolder.value = true
    await DeleteFolder(folderToDelete.value.id)
    deleteFolderDialog.value = false
    if (selectedFolderId.value === folderToDelete.value.id) {
      selectedFolderId.value = null
    }
    alertMsg.value = `Folder "${folderToDelete.value.name}" deleted`
    alertType.value = 'success'
    showAlert.value = true
    await Promise.all([fetchFolders(), fetchFiles()])
  } catch (err: any) {
    alertMsg.value = 'Failed to delete folder: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isDeletingFolder.value = false
  }
}

function openMoveDialog(file: models.File) {
  fileToMove.value = file
  targetFolderSelection.value = file.folderId || ''
  moveDialog.value = true
}

async function executeMoveFile() {
  if (!fileToMove.value) return
  try {
    isMoving.value = true
    await MoveFile(fileToMove.value.id, targetFolderSelection.value)
    moveDialog.value = false
    alertMsg.value = 'File moved successfully'
    alertType.value = 'success'
    showAlert.value = true
    await Promise.all([fetchFolders(), fetchFiles()])
  } catch (err: any) {
    alertMsg.value = 'Failed to move file: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isMoving.value = false
  }
}

function openRenameDialog(file: models.File) {
  fileToRename.value = file
  newFileName.value = file.name
  renameDialog.value = true
}

async function executeRenameFile() {
  if (!fileToRename.value) return
  const clean = newFileName.value.trim()
  if (!clean) return

  try {
    isRenaming.value = true
    await RenameFile(fileToRename.value.id, clean)
    renameDialog.value = false
    alertMsg.value = 'File renamed successfully'
    alertType.value = 'success'
    showAlert.value = true
    await fetchFiles()
  } catch (err: any) {
    alertMsg.value = 'Failed to rename file: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isRenaming.value = false
  }
}

async function openPreview(file: models.File) {
  previewFile.value = file
  previewContent.value = ''
  previewDialog.value = true
  try {
    isLoadingPreview.value = true
    previewContent.value = await GetFileContent(file.id)
  } catch (err: any) {
    console.error('Failed to get preview content:', err)
  } finally {
    isLoadingPreview.value = false
  }
}

function decodeBase64ToText(b64: string): string {
  try {
    return atob(b64)
  } catch {
    return 'Unable to decode text'
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
  const targetFolderId = selectedFolderId.value || ''

  for (let i = 0; i < fileList.length; i++) {
    const file = fileList[i]
    try {
      const base64Data = await readFileAsBase64(file)
      await UploadFile(currentModule.value.id, targetFolderId, file.name, base64Data)
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
    await Promise.all([fetchFolders(), fetchFiles()])
  }
}

function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      const commaIndex = result.indexOf(',')
      resolve(commaIndex !== -1 ? result.substring(commaIndex + 1) : result)
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
    await Promise.all([fetchFolders(), fetchFiles()])
  } catch (err: any) {
    alertMsg.value = 'Failed to delete file: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isDeleting.value = false
  }
}

// Assignment actions
function openCreateAssignmentDialog() {
  editingAssignment.value = null
  assignmentForm.value = {
    title: '',
    description: '',
    priority: 'medium',
    status: 'pending',
    dueDate: ''
  }
  assignmentFiles.value = null
  assignmentDialog.value = true
}

function getSelectedAssignmentFile(): File | null {
  const selected = assignmentFiles.value
  if (!selected) return null
  return Array.isArray(selected) ? selected[0] ?? null : selected
}

function openEditAssignmentDialog(assign: models.Assignment) {
  assignmentDetailDialog.value = false
  editingAssignment.value = assign
  let formattedDate = ''
  if (assign.dueDate) {
    const d = new Date(assign.dueDate)
    if (!isNaN(d.getTime())) {
      formattedDate = d.toISOString().split('T')[0]
    }
  }
  assignmentForm.value = {
    title: assign.title,
    description: assign.description || '',
    priority: assign.priority || 'medium',
    status: assign.status || 'pending',
    dueDate: formattedDate
  }
  assignmentDialog.value = true
}

async function saveAssignment() {
  if (!currentModule.value?.id) return
  const title = assignmentForm.value.title.trim()
  if (!title) return

  let uploadFailed = false
  try {
    isSavingAssignment.value = true
    if (editingAssignment.value) {
      await UpdateAssignment(
        editingAssignment.value.id,
        title,
        assignmentForm.value.description,
        assignmentForm.value.status,
        assignmentForm.value.priority,
        assignmentForm.value.dueDate
      )
      alertMsg.value = 'Assignment updated!'
    } else {
      const created = await CreateAssignment(
        currentModule.value.id,
        title,
        assignmentForm.value.description,
        assignments.value.length + 1,
        assignmentForm.value.status,
        assignmentForm.value.priority,
        assignmentForm.value.dueDate
      )
      const instructions = getSelectedAssignmentFile()
      if (instructions) {
        try {
          const base64Data = await readFileAsBase64(instructions)
          await UploadAssignmentInstructions(currentModule.value.id, created.id, instructions.name, base64Data)
          alertMsg.value = 'Assignment created and instructions uploaded!'
        } catch (uploadErr: any) {
          uploadFailed = true
          alertMsg.value = 'Assignment created, but instructions upload failed: ' + uploadErr
        }
      } else {
        alertMsg.value = 'Assignment created!'
      }
    }
    assignmentDialog.value = false
    alertType.value = uploadFailed ? 'error' : 'success'
    showAlert.value = true
    await fetchAssignments()
    await fetchFolders()
    await fetchFiles()
  } catch (err: any) {
    alertMsg.value = 'Failed to save assignment: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingAssignment.value = false
  }
}

// Assignment details / created files
function openAssignmentDetails(assign: models.Assignment) {
  detailAssignment.value = assign
  assignmentDetailDialog.value = true
}

async function handleExportFile(file: models.File) {
  try {
    const destination = await exportSchiveFile(file)
    if (!destination) return
    alertMsg.value = `Exported ${file.name} to ${destination}`
    alertType.value = 'success'
    showAlert.value = true
  } catch (err: any) {
    alertMsg.value = 'Failed to export file: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

async function handleAssignmentFilesChanged() {
  await fetchFiles()
  await fetchFolders()
  await fetchAssignments()
}

function getProgressColor(progress?: number): string {
  const p = progress || 0
  if (p >= 100) return 'success'
  if (p >= 50) return 'color-bg-2'
  if (p > 0) return 'warning'
  return 'grey'
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
    alertMsg.value = `Progress updated to ${progressValue.value}%`
    alertType.value = 'success'
    showAlert.value = true
    await fetchAssignments()
  } catch (err: any) {
    alertMsg.value = 'Failed to update progress: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingProgress.value = false
  }
}

async function updateStatus(id: string, status: string) {
  try {
    await UpdateAssignmentStatus(id, status)
    await fetchAssignments()
    eventBus.emit(SYNC_EVENTS.ASSIGNMENTS_CHANGED, { id, status })
  } catch (err: any) {
    alertMsg.value = 'Failed to update status: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

function confirmDeleteAssignment(assign: models.Assignment) {
  askConfirm(
    {
      title: 'Delete assignment?',
      message: `"${assign.title}" will be permanently removed, including its uploaded instructions.`,
      confirmText: 'Delete assignment',
      tone: 'error',
      icon: 'mdi-clipboard-remove-outline'
    },
    async () => {
      await DeleteAssignment(assign.id)
      notifySuccess('Assignment deleted')
      await fetchAssignments()
    }
  )
}

// Task actions
function openCreateTaskDialog() {
  editingTask.value = null
  taskForm.value = { title: '', notes: '', dueDate: '' }
  taskDialog.value = true
}

function openEditTaskDialog(task: models.Task) {
  editingTask.value = task
  let formattedDate = ''
  if (task.dueDate) {
    const d = new Date(task.dueDate)
    if (!isNaN(d.getTime())) {
      formattedDate = d.toISOString().split('T')[0]
    }
  }
  taskForm.value = {
    title: task.title,
    notes: task.notes || '',
    dueDate: formattedDate
  }
  taskDialog.value = true
}

async function saveTask() {
  if (!currentModule.value?.id) return
  const title = taskForm.value.title.trim()
  if (!title) return

  try {
    isSavingTask.value = true
    if (editingTask.value) {
      await UpdateTask(editingTask.value.id, title, taskForm.value.notes, taskForm.value.dueDate)
      notifySuccess('Task updated')
    } else {
      await CreateTask(currentModule.value.id, title, taskForm.value.notes, taskForm.value.dueDate)
      notifySuccess('Task created')
    }
    taskDialog.value = false
    await fetchTasks()
  } catch (err: any) {
    alertMsg.value = 'Failed to save task: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingTask.value = false
  }
}

async function toggleTask(task: models.Task) {
  try {
    const nextDone = !task.done
    await SetTaskDone(task.id, nextDone)
    await fetchTasks()
    eventBus.emit(SYNC_EVENTS.TASKS_CHANGED, { id: task.id, done: nextDone })
  } catch (err: any) {
    alertMsg.value = 'Failed to update task: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

function askDeleteTask(task: models.Task) {
  askConfirm(
    {
      title: 'Delete task?',
      message: `"${task.title}" will be permanently removed.`,
      confirmText: 'Delete task',
      tone: 'error',
      icon: 'mdi-trash-can-outline'
    },
    async () => {
      await DeleteTask(task.id)
      notifySuccess('Task deleted')
      await fetchTasks()
    }
  )
}

function isTaskOverdue(task: models.Task): boolean {
  if (!task.dueDate || task.done) return false
  const due = new Date(task.dueDate)
  return !isNaN(due.getTime()) && due.getTime() < Date.now()
}

// Confirmation dialog helpers
function askConfirm(
  options: {
    title: string
    message: string
    confirmText?: string
    tone?: 'error' | 'warning'
    icon?: string
  },
  action: () => Promise<void>
) {
  confirmTitle.value = options.title
  confirmMessage.value = options.message
  confirmText.value = options.confirmText || 'Delete'
  confirmTone.value = options.tone || 'error'
  confirmIcon.value = options.icon || ''
  pendingConfirm = action
  confirmDialog.value = true
}

async function runConfirmedAction() {
  if (!pendingConfirm) return
  const action = pendingConfirm
  pendingConfirm = null
  confirmLoading.value = true
  try {
    await action()
    confirmDialog.value = false
  } catch (err: any) {
    confirmDialog.value = false
    alertMsg.value = 'Action failed: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    confirmLoading.value = false
  }
}

function notifySuccess(message: string) {
  alertMsg.value = message
  alertType.value = 'success'
  showAlert.value = true
}

// Flashcard actions
function openCreateFlashcardDialog() {
  editingFlashcard.value = null
  flashcardForm.value = { question: '', answer: '' }
  flashcardDialog.value = true
}

function openEditFlashcardDialog(card: models.Flashcard) {
  editingFlashcard.value = card
  flashcardForm.value = { question: card.question, answer: card.answer }
  flashcardDialog.value = true
}

async function saveFlashcard() {
  if (!currentModule.value?.id) return
  const q = flashcardForm.value.question.trim()
  const a = flashcardForm.value.answer.trim()
  if (!q || !a) return

  try {
    isSavingFlashcard.value = true
    if (editingFlashcard.value) {
      await UpdateFlashcard(editingFlashcard.value.id, q, a)
      alertMsg.value = 'Flashcard updated!'
    } else {
      await CreateFlashcard(currentModule.value.id, q, a)
      alertMsg.value = 'Flashcard created!'
    }
    flashcardDialog.value = false
    alertType.value = 'success'
    showAlert.value = true
    await fetchFlashcards()
  } catch (err: any) {
    alertMsg.value = 'Failed to save flashcard: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSavingFlashcard.value = false
  }
}

async function toggleCardMastered(id: string) {
  try {
    await ToggleMastered(id)
    await fetchFlashcards()
  } catch (err: any) {
    console.error('Failed to toggle mastered:', err)
  }
}

function confirmDeleteFlashcard(card: models.Flashcard) {
  askConfirm(
    {
      title: 'Delete flashcard?',
      message: 'This flashcard will be permanently removed from your deck.',
      confirmText: 'Delete flashcard',
      tone: 'error',
      icon: 'mdi-cards-outline'
    },
    async () => {
      await DeleteFlashcard(card.id)
      notifySuccess('Flashcard deleted')
      await fetchFlashcards()
    }
  )
}

// Study deck actions
function startStudySession() {
  studyCards.value = [...flashcards.value]
  currentStudyIndex.value = 0
  showStudyAnswer.value = false
  studyModal.value = true
}

function nextStudyCard() {
  if (currentStudyIndex.value < studyCards.value.length - 1) {
    currentStudyIndex.value++
    showStudyAnswer.value = false
  } else {
    studyModal.value = false
    alertMsg.value = 'Study session complete! Great job.'
    alertType.value = 'success'
    showAlert.value = true
  }
}

function prevStudyCard() {
  if (currentStudyIndex.value > 0) {
    currentStudyIndex.value--
    showStudyAnswer.value = false
  }
}

async function markCurrentStudyMastered() {
  const card = studyCards.value[currentStudyIndex.value]
  if (card && !card.mastered) {
    await toggleCardMastered(card.id)
    card.mastered = true
  }
  nextStudyCard()
}

// Helpers
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
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatDateTime(dateStr: any): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatDueDateDisplay(dateStr: any): string {
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

function getPriorityColor(p?: string) {
  switch ((p || '').toLowerCase()) {
    case 'high': return 'error'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'grey'
  }
}

function getStatusColor(s?: string) {
  switch ((s || '').toLowerCase()) {
    case 'completed': return 'success'
    case 'in_progress': return 'warning'
    default: return 'grey'
  }
}

function getFileInfo(ext: string) {
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
    case 'txt':
    case 'md': return { icon: 'mdi-file-document-outline', color: '#546E7A', label: 'TEXT' }
    default: return { icon: 'mdi-file-outline', color: '#757575', label: e.toUpperCase() || 'FILE' }
  }
}

function onWindowFocus() {
  if (currentModule.value?.id) {
    fetchFiles()
    fetchFolders()
    fetchAssignments()
    fetchFlashcards()
  }
}

let unsubTasks: (() => void) | null = null
let unsubAssignments: (() => void) | null = null

onMounted(() => {
  fetchModuleDetails()
  window.addEventListener('focus', onWindowFocus)

  unsubTasks = eventBus.on(SYNC_EVENTS.TASKS_CHANGED, () => {
    fetchTasks()
  })
  unsubAssignments = eventBus.on(SYNC_EVENTS.ASSIGNMENTS_CHANGED, () => {
    fetchAssignments()
  })
})

onUnmounted(() => {
  window.removeEventListener('focus', onWindowFocus)
  if (unsubTasks) unsubTasks()
  if (unsubAssignments) unsubAssignments()
})

watch(moduleCode, (newCode, oldCode) => {
  if (newCode && newCode !== oldCode) {
    currentModule.value = null
    files.value = []
    folders.value = []
    assignments.value = []
    tasks.value = []
    flashcards.value = []
    selectedFolderId.value = null
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

.module-tabs {
  background-color: var(--color-bg-3);
}

/* Folders Grid */
.folders-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.folder-card {
  border: 1px solid rgba(0, 0, 0, 0.08);
  transition: transform 0.15s ease, box-shadow 0.15s ease, border-color 0.15s ease;
}

.folder-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1) !important;
  border-color: var(--color-bg-2, #7A8B69);
}

.add-folder-card {
  border: 2px dashed rgba(0, 0, 0, 0.18);
  min-height: 60px;
  transition: all 0.15s ease;
}

.add-folder-card:hover {
  border-color: var(--color-bg-2, #7A8B69);
  background-color: rgba(122, 139, 105, 0.08);
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

.filter-toggle {
  height: 32px !important;
}

.task-row {
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.task-row:hover {
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.07) !important;
}

.task-title--done {
  text-decoration: line-through;
  color: var(--color-text-light, #9e9e9e);
}

.task-notes {
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.upload-btn {
  background-color: var(--color-bg-2, #7A8B69) !important;
  color: #ffffff !important;
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
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.file-card {
  min-height: 175px;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.file-card-actions {
  min-height: 36px;
}

.file-card-actions .v-btn {
  min-width: 0 !important;
}

.file-card-actions .v-btn--icon.v-btn--density-compact {
  width: 28px !important;
  height: 28px !important;
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

.modification-date {
  font-size: 11.5px;
}

.border-top {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}

.border-bottom {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.file-table-row:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

/* Flashcards Grid */
.flashcards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.flashcard-preview-card {
  min-height: 180px;
  transition: all 0.15s ease;
}

.flashcard-preview-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.1) !important;
}

.card-mastered {
  border-left: 5px solid rgb(var(--v-theme-success)) !important;
}

.card-answer-text {
  border-left: 2px solid var(--color-bg-2);
  padding-left: 8px;
  margin-top: 6px;
}

/* Study Deck */
.study-flip-card {
  min-height: 220px;
  background-color: var(--color-bg-3);
  border: 2px dashed rgba(0, 0, 0, 0.15);
  transition: all 0.2s ease;
}

.study-flip-card:hover {
  border-color: var(--color-bg-2);
  transform: scale(1.01);
}

.preview-image {
  max-width: 100%;
  max-height: 500px;
  object-fit: contain;
}

.preview-text {
  background-color: var(--color-bg-3);
  font-family: monospace;
  font-size: 13px;
  max-height: 480px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
