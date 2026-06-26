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
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router