<template>
  <div class="timetable-page">
    <div class="d-flex justify-space-between align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">Timetable</h1>
        <p class="text-subtitle-2 text-grey-darken-1">View and manage your weekly module schedule</p>
      </div>
      <v-btn
        color="color-bg-2"
        prepend-icon="mdi-plus"
        variant="flat"
        class="text-none font-weight-medium rounded-lg"
        @click="openAddDialog"
      >
        Add Class Slot
      </v-btn>
    </div>

    <!-- Calendar View -->
    <v-card class="pa-2 rounded-lg" elevation="2">
      <div v-if="isLoading" class="d-flex justify-center align-center py-12">
        <v-progress-circular indeterminate color="color-bg-2"></v-progress-circular>
      </div>
      <vue-cal
        v-else
        class="custom-vuecal"
        style="height: 650px;"
        active-view="week"
        :views="['week', 'day']"
        :time-from="7 * 60"
        :time-to="20 * 60"
        :time-step="60"
        :time-cell-height="55"
        :hide-weekends="true"
        :events="calendarEvents"
      >
        <template #event="{ event }">
          <div
            class="timetable-slot"
            :style="{ backgroundColor: event.color }"
            @click.stop="openDetails(event)"
          >
            <div class="slot-code font-weight-bold">{{ event.title }}</div>
            <div v-if="event.room" class="slot-room text-caption">📍 {{ event.room }}</div>
            <div class="slot-time text-caption">{{ event.timeFormatted }}</div>
          </div>
        </template>
      </vue-cal>
    </v-card>

    <!-- Dialog: Add Class Slot -->
    <v-dialog v-model="addDialog" width="460">
      <v-card class="pa-6 rounded-lg" elevation="4">
        <h3 class="text-h6 font-weight-medium mb-4">Add Class Slot</h3>
        <v-form @submit.prevent="handleCreateSlot">
          <v-select
            v-model="slotForm.moduleId"
            :items="moduleOptions"
            item-title="label"
            item-value="id"
            label="Module"
            variant="outlined"
            density="comfortable"
            :rules="[v => !!v || 'Module is required']"
            class="mb-2"
          ></v-select>

          <v-select
            v-model="slotForm.day"
            :items="dayOptions"
            item-title="name"
            item-value="value"
            label="Day of Week"
            variant="outlined"
            density="comfortable"
            :rules="[v => !!v || 'Day is required']"
            class="mb-2"
          ></v-select>

          <div class="d-flex gap-3">
            <v-select
              v-model="slotForm.startHour"
              :items="timeOptions"
              item-title="label"
              item-value="value"
              label="Start Time"
              variant="outlined"
              density="comfortable"
              class="mb-2 flex-1-1"
            ></v-select>

            <v-select
              v-model="slotForm.lengthHrs"
              :items="durationOptions"
              item-title="label"
              item-value="value"
              label="Duration"
              variant="outlined"
              density="comfortable"
              class="mb-2 flex-1-1"
            ></v-select>
          </div>

          <v-text-field
            v-model="slotForm.room"
            label="Room / Venue (optional)"
            placeholder="e.g. Science Lab 2"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-text-field>

          <div class="d-flex justify-end gap-2">
            <v-btn type="button" variant="text" @click="closeAddDialog">Cancel</v-btn>
            <v-btn type="submit" color="color-bg-2" variant="flat" :loading="isSaving">
              Save Slot
            </v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- Dialog: Slot Details & Delete -->
    <v-dialog v-model="detailsDialog" width="380">
      <v-card v-if="selectedEvent" class="pa-6 rounded-lg" elevation="4">
        <div class="d-flex align-center justify-space-between mb-3">
          <div class="d-flex align-center gap-2">
            <span
              class="module-dot"
              :style="{ backgroundColor: selectedEvent.color }"
            ></span>
            <h3 class="text-h6 font-weight-bold">{{ selectedEvent.title }}</h3>
          </div>
          <v-chip size="small" variant="tonal" :color="selectedEvent.color">
            {{ selectedEvent.dayName }}
          </v-chip>
        </div>

        <p v-if="selectedEvent.moduleName" class="text-body-2 text-grey-darken-1 mb-2">
          {{ selectedEvent.moduleName }}
        </p>

        <v-divider class="my-3"></v-divider>

        <div class="text-body-2 mb-2">
          <strong>Time:</strong> {{ selectedEvent.timeFormatted }}
        </div>
        <div v-if="selectedEvent.room" class="text-body-2 mb-4">
          <strong>Room:</strong> {{ selectedEvent.room }}
        </div>

        <div class="d-flex justify-space-between align-center mt-4">
          <v-btn
            color="error"
            variant="text"
            prepend-icon="mdi-delete"
            @click="handleDeleteSlot(selectedEvent.slotId)"
          >
            Delete Slot
          </v-btn>
          <v-btn type="button" variant="text" @click="detailsDialog = false">Close</v-btn>
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
import { ref, computed, onMounted } from 'vue'
import { VueCal } from 'vue-cal'
import 'vue-cal/style'
import { GetAllModules } from '../../wailsjs/go/controllers/ModuleController'
import { GetAllClassSlots, CreateClassSlot, DeleteClassSlot } from '../../wailsjs/go/controllers/TimetableController'
import { models } from '../../wailsjs/go/models'

const modules = ref<models.Module[]>([])
const classSlots = ref<models.ClassSlot[]>([])
const calendarEvents = ref<any[]>([])
const isLoading = ref(true)
const isSaving = ref(false)

const addDialog = ref(false)
const detailsDialog = ref(false)
const selectedEvent = ref<any>(null)

const showAlert = ref(false)
const alertMsg = ref('')
const alertType = ref('success')

const dayOptions = [
  { name: 'Monday', value: 'mon' },
  { name: 'Tuesday', value: 'tue' },
  { name: 'Wednesday', value: 'wed' },
  { name: 'Thursday', value: 'thu' },
  { name: 'Friday', value: 'fri' },
]

const timeOptions = [
  { label: '07:00 AM', value: 7 },
  { label: '08:00 AM', value: 8 },
  { label: '09:00 AM', value: 9 },
  { label: '10:00 AM', value: 10 },
  { label: '11:00 AM', value: 11 },
  { label: '12:00 PM', value: 12 },
  { label: '01:00 PM (13:00)', value: 13 },
  { label: '02:00 PM (14:00)', value: 14 },
  { label: '03:00 PM (15:00)', value: 15 },
  { label: '04:00 PM (16:00)', value: 16 },
  { label: '05:00 PM (17:00)', value: 17 },
  { label: '06:00 PM (18:00)', value: 18 },
  { label: '07:00 PM (19:00)', value: 19 },
]

const durationOptions = [
  { label: '1 Hour', value: 1 },
  { label: '2 Hours', value: 2 },
  { label: '3 Hours', value: 3 },
  { label: '4 Hours', value: 4 },
]

const slotForm = ref({
  moduleId: '',
  day: 'mon',
  startHour: 8,
  lengthHrs: 1,
  room: '',
})

const moduleOptions = computed(() => {
  return modules.value.map(m => ({
    id: m.id,
    label: m.name?.trim() ? `${m.code} - ${m.name}` : m.code,
    code: m.code,
    name: m.name,
    color: m.color || '#C9622C',
  }))
})

function getMonday(date = new Date()) {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  d.setDate(diff)
  d.setHours(0, 0, 0, 0)
  return d
}

const dayOffsets: Record<string, number> = {
  mon: 0,
  tue: 1,
  wed: 2,
  thu: 3,
  fri: 4,
}

const dayNames: Record<string, string> = {
  mon: 'Monday',
  tue: 'Tuesday',
  wed: 'Wednesday',
  thu: 'Thursday',
  fri: 'Friday',
}

function formatHour(h: number) {
  const period = h >= 12 ? 'PM' : 'AM'
  const displayHour = h % 12 || 12
  return `${displayHour}:00 ${period}`
}

function buildCalendarEvents() {
  const events: any[] = []
  const baseMonday = getMonday()

  // Generate recurring events for weeks around the current date
  for (let weekOffset = -4; weekOffset <= 12; weekOffset++) {
    for (const slot of classSlots.value) {
      const dayKey = (slot.day || 'mon').toLowerCase()
      const offset = dayOffsets[dayKey] ?? 0

      const slotDate = new Date(baseMonday)
      slotDate.setDate(slotDate.getDate() + weekOffset * 7 + offset)

      const year = slotDate.getFullYear()
      const month = String(slotDate.getMonth() + 1).padStart(2, '0')
      const day = String(slotDate.getDate()).padStart(2, '0')
      const dateStr = `${year}-${month}-${day}`

      const startHour = Number(slot.startHour) || 7
      const lengthHrs = Number(slot.lengthHrs) || 1
      const endHour = startHour + lengthHrs

      const startHourStr = String(startHour).padStart(2, '0')
      const endHourStr = String(endHour).padStart(2, '0')

      const mod = slot.module || modules.value.find(m => m.id === slot.moduleId)

      events.push({
        id: `${slot.id}-${weekOffset}`,
        slotId: slot.id,
        start: `${dateStr} ${startHourStr}:00`,
        end: `${dateStr} ${endHourStr}:00`,
        title: mod?.code || 'Module',
        moduleName: mod?.name || '',
        room: slot.room,
        color: mod?.color || '#C9622C',
        day: slot.day,
        dayName: dayNames[dayKey] || slot.day,
        timeFormatted: `${formatHour(startHour)} - ${formatHour(endHour)}`,
      })
    }
  }

  calendarEvents.value = events
}

async function fetchData() {
  try {
    isLoading.value = true
    const [mods, slots] = await Promise.all([
      GetAllModules(),
      GetAllClassSlots(),
    ])
    modules.value = mods || []
    classSlots.value = slots || []
    buildCalendarEvents()
  } catch (err: any) {
    alertMsg.value = 'Failed to load timetable data: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isLoading.value = false
  }
}

function openAddDialog() {
  if (modules.value.length === 0) {
    alertMsg.value = 'Please create at least one module first before adding timetable events.'
    alertType.value = 'warning'
    showAlert.value = true
    return
  }
  slotForm.value.moduleId = modules.value[0]?.id || ''
  slotForm.value.day = 'mon'
  slotForm.value.startHour = 8
  slotForm.value.lengthHrs = 1
  slotForm.value.room = ''
  addDialog.value = true
}

function closeAddDialog() {
  addDialog.value = false
}

async function handleCreateSlot() {
  if (!slotForm.value.moduleId) return

  try {
    isSaving.value = true
    await CreateClassSlot(new models.ClassSlot({
      id: '',
      moduleId: slotForm.value.moduleId,
      day: slotForm.value.day,
      startHour: Number(slotForm.value.startHour),
      lengthHrs: Number(slotForm.value.lengthHrs),
      room: slotForm.value.room?.trim() || '',
    }))

    addDialog.value = false
    alertMsg.value = 'Class slot added successfully!'
    alertType.value = 'success'
    showAlert.value = true
    await fetchData()
  } catch (err: any) {
    alertMsg.value = 'Failed to add class slot: ' + err
    alertType.value = 'error'
    showAlert.value = true
  } finally {
    isSaving.value = false
  }
}

function openDetails(event: any) {
  selectedEvent.value = event
  detailsDialog.value = true
}

async function handleDeleteSlot(slotId: string) {
  try {
    await DeleteClassSlot(slotId)
    detailsDialog.value = false
    alertMsg.value = 'Class slot deleted successfully!'
    alertType.value = 'success'
    showAlert.value = true
    await fetchData()
  } catch (err: any) {
    alertMsg.value = 'Failed to delete class slot: ' + err
    alertType.value = 'error'
    showAlert.value = true
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.timetable-page {
  width: 100%;
}

.module-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.timetable-slot {
  width: 100%;
  height: 100%;
  padding: 6px 8px;
  border-radius: 6px;
  color: #fff;
  cursor: pointer;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.timetable-slot:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.slot-code {
  font-size: 13px;
  line-height: 1.2;
}

.slot-room {
  font-size: 11px;
  opacity: 0.9;
  margin-top: 2px;
}

.slot-time {
  font-size: 10px;
  opacity: 0.85;
  margin-top: auto;
}

:deep(.vuecal__event) {
  background-color: transparent !important;
  border: none !important;
  padding: 2px !important;
}

:deep(.vuecal__menu) {
  background-color: var(--color-bg-3) !important;
}

:deep(.vuecal__view-btn--active) {
  border-bottom-color: var(--color-bg-2) !important;
  color: var(--color-bg-2) !important;
}
</style>
