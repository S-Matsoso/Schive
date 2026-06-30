import { createRouter, createWebHashHistory } from "vue-router";
import Timetable from '../views/Timetable.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: Timetable,
    },
    {
        path: '/module/:code',
        name: 'module-details',
        component: () => import('../views/ModuleDetails.vue')
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router