import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import { initWailsPolyfill } from './utils/wailsPolyfill';

initWailsPolyfill();

createApp(App).mount('#app');
