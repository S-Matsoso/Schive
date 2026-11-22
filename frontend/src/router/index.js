import { createRouter, createWebHashHistory } from "vue-router";
import Home from '../views/Home.vue'
import Timetable from '../views/Timetable.vue'
import Tasks from '../views/Tasks.vue'

const routes = [
    {
        path: '/',
        name: 'dashboard',
        component: Home,
    },
    {
        path: '/timetable',
        name: 'timetable',
        component: Timetable,
    },
    {
        path: '/tasks',
        name: 'tasks',
        component: Tasks,
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