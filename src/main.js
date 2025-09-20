import './assets/styles.min.css'
import './assets/theme.scss'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Notifications from '@kyvg/vue3-notification'
import { createClient } from '@supabase/supabase-js'

import App from './App.vue'
import router from './router'

const app = createApp(App)
app.config.globalProperties.$supabase = createClient(
    'https://lqcqxkbjgpoygmczrlap.supabase.co',
    import.meta.env.VITE_ANON_KEY
)

app.use(createPinia())
app.use(router)
app.use(Notifications)

app.mount('#app')
