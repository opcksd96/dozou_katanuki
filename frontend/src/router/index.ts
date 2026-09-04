import { createRouter, createWebHashHistory } from 'vue-router';
import WebUI from '../views/webui/WebUI.vue';
import AdminConsole from '../views/admin/AdminConsole.vue';
import MobileConsole from '../views/mobile/MobileConsole.vue';

const routes = [
  {
    path: '/',
    redirect: '/webui'
  },
  {
    path: '/webui',
    name: 'WebUI',
    component: WebUI
  },
  {
    path: '/admin',
    name: 'Admin',
    component: AdminConsole
  },
  {
    path: '/mobile',
    name: 'Mobile',
    component: MobileConsole
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
