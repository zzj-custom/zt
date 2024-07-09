import {createApp} from 'vue'
import App from './App.vue'
import '@/scss/style.scss';

import router from './router';
import vuetify from './plugins/vuetify';

import PerfectScrollbar from 'vue3-perfect-scrollbar';
import VueTableIcons from 'vue-tabler-icons';

import {createPinia} from 'pinia'

// 全局消息组件
import messagePlugin  from './plugins/messagePlugin';

const app = createApp(App);
const pinia = createPinia(); // Create a new Pinia
app.use(vuetify);
app.use(PerfectScrollbar);
app.use(VueTableIcons);
app.use(pinia)
app.use(messagePlugin);
app.use(router).mount('#app');
