import type { RouteRecordRaw } from 'vue-router';
import { createRouter, createWebHashHistory } from 'vue-router';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    redirect: 'content/host',
    component: () => import('@/layouts/BasicLayout.vue'),
    children: [
      {
        path: '/terminal',
        name: 'Terminal',
        component: () => import('@/views/body/Terminal.vue'),
      },
      {
        path: '/sftp',
        name: 'SFTP',
        component: () => import('@/views/body/SFTP.vue'),
      },
      {
        path: 'content',
        name: 'Content',
        component: () => import('@/layouts/Main.vue'),
        children: [
          {
            path: 'host',
            name: 'Host',
            component: () => import('@/views/body/HostInventory.vue'),
          },
          {
            path: 'credentials',
            name: 'Credentials',
            component: () => import('@/views/body/Credentials.vue'),
          },
          {
            path: 'port-forwarding',
            name: 'PortForwarding',
            component: () => import('@/views/body/PortForwarding.vue'),
          },
          {
            path: 'snippets',
            name: 'Snippets',
            component: () => import('@/views/body/Snippets.vue'),
          },
          {
            path: 'history',
            name: 'History',
            component: () => import('@/views/body/History.vue'),
          },
        ],
      },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
