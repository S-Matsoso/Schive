import { createRouter, createWebHashHistory } from "vue-router";
import Files from '../views/Files.vue'
import Home from '../views/Home.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home,
    },
    {
        path: '/files',
        name: 'files',
        component: Files
    },
    // Add this new dynamic route right here:
    {
        path: '/module/:id',
        name: 'module-details',
        // Using () => import() is called "lazy loading". 
        // It keeps your app fast by only loading this page when the user actually clicks it!
        component: () => import('../views/ModuleDetails.vue')
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router