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
    }
})


createApp(App).use(router).use(vuetify).mount('#app')