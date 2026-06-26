import { createRouter, createWebHashHistory } from "vue-router";
import Home from '../views/Files.vue'
import Files from '../views/Home.vue'


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