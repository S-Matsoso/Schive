import { createRouter, createWebHashHistory } from "vue-router";
import Timetable from '../views/Timetable.vue'

const routes = [
    {
        path: '/',
        alias: '/timetable',
        name: 'timetable',
        component: Timetable,
    },
    {
        path: '/module/:code',
        name: 'module-details',
        component: () => import('../views/ModuleDetails.vue')
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/'
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router