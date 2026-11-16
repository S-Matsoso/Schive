import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index.js';
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import './style.css';

import { createVuetify } from 'vuetify';

const vuetify = createVuetify({
    theme: {
        defaultTheme: 'light', // This locks the entire app into light mode
        themes: {
            light: {
                colors: {
                    'color-primary': '#F7F5F0',
                    'color-secondary': "#FBFAF7",
                    'color-bg-1': "#C9622C",
                    'color-bg-2': "#7A8B69",
                    'color-bg-3': "#EFEBE0",
                    'color-bg-4': "#EDEAE2",
                    'color-text-dark': "#3A372F",
                    'color-text-light': "#6B665A",
                    'color-text-blue': "#E8E8E7",
                    'color-text-purple': "#8C6FAE",
                    'color-text-green': "#7A8B69",
                    'color-text-orange': "#C9622C",
                }
            }
        }
    }
})


createApp(App).use(router).use(vuetify).mount('#app')