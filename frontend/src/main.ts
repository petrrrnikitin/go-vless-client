import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'
import { useAppStore } from './stores/app'

const pinia = createPinia()
createApp(App).use(pinia).mount('#app')

// init вызывается ровно один раз при загрузке страницы.
// main.ts не участвует в Vite HMR — компоненты перезагружаются без него.
useAppStore(pinia).init()
