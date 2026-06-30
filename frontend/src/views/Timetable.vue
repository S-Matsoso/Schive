<template>
    <div>
      <h1>THIS IS TIMETABLE</h1> 
        <vue-cal 
            v-bind="config"
            :views="['day', 'days']"
            :time-from="7 * 60"
            :time-to="20 * 60"
            :special-hours="timetableSpecialHours"
            :events="classes"
            />
    </div>
</template>

<script setup>
    import {VueCal} from 'vue-cal'
    import 'vue-cal/style'

    const timetableSpecialHours = {
        mon: {from: 7 * 60, to: 20 * 60, class:'classSlots'},
        tue: {from: 7 * 60, to: 20 * 60, class:'classSlots'},
        wed: {from: 7 * 60, to: 20 * 60, class:'classSlots'},
        thu: {from: 7 * 60, to: 20 * 60, class:'classSlots'},
        fri: {from: 7 * 60, to: 20 * 60, class:'classSlots'},

    }

    const classes = [
        {start: '2026-07-01 07:00', end: '2026-07-01 07:50', title: 'CSIS3891'},
        {start: '2026-07-02 09:00', end: '2026-07-02 09:50', title: 'CSIS3891'},
        {start: '2026-07-02 10:00', end: '2026-07-02 10:50', title: 'CSIS3891'},
        {start: '2026-07-04 07:00', end: '2026-07-04 07:50', title: 'CSIS3891'},
        {start: '2026-07-05 17:00', end: '2026-07-05 17:50', title: 'CSIS3891'},

    ]
  
</script>

<style>
    /* set colors to module colors */
  .classSlots {
        background: #00daff21;
    }
</style>
