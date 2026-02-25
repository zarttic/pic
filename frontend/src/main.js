import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import { setupLazyload } from './directives/lazyload'

// 导入全局样式（Vuetify 样式已在 vuetify.js 中导入）
import './styles/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(vuetify)

// 注册全局指令
setupLazyload(app)

app.mount('#app')
