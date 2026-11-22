import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index.js';
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import './style.css';

import { createVuetify } from 'vuetify';

const savedTheme = localStorage.getItem('schive-theme') || 'light'
if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
} else {
    document.documentElement.classList.remove('dark')
}

const vuetify = createVuetify({
    theme: {
        defaultTheme: savedTheme,
        themes: {
            light: {
                dark: false,
                colors: {
                    background: '#F7F5F0',
                    surface: '#FFFFFF',
                    'color-primary': '#F7F5F0',
                    'color-secondary': '#FBFAF7',
                    'color-bg-1': '#C9622C',
                    'color-bg-2': '#7A8B69',
                    'color-bg-3': '#EFEBE0',
                    'color-bg-4': '#EDEAE2',
                    'color-text-dark': '#3A372F',
                    'color-text-light': '#6B665A',
                    'color-text-blue': '#E8E8E7',
                    'color-text-purple': '#8C6FAE',
                    'color-text-green': '#7A8B69',
                    'color-text-orange': '#C9622C',
                }
            },
            dark: {
                dark: true,
                colors: {
                    background: '#191A1D',
                    surface: '#24272C',
                    'color-primary': '#191A1D',
                    'color-secondary': '#202327',
                    'color-bg-1': '#D36F3D',
                    'color-bg-2': '#7A8B69',
                    'color-bg-3': '#2E3339',
                    'color-bg-4': '#3C424A',
                    'color-text-dark': '#EBE9E4',
                    'color-text-light': '#A29D92',
                    'color-text-blue': '#4D90FE',
                    'color-text-purple': '#A689CA',
                    'color-text-green': '#889B77',
                    'color-text-orange': '#D36F3D',
                }
            }
        }
    }
})

createApp(App).use(router).use(vuetify).mount('#app')