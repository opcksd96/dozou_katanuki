import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import './style.css';
import { initWailsPolyfill } from './utils/wailsPolyfill';

initWailsPolyfill();

createApp(App).use(router).mount('#app');
