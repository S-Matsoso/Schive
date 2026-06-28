import { createRouter, createWebHashHistory } from "vue-router";
import Home from '../views/Home.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home,
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