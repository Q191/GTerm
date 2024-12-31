import type { RouteRecordRaw } from 'vue-router';
import { createRouter, createWebHashHistory } from 'vue-router';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    redirect: 'content/host',
    component: () => import('@/layouts/Layout.vue'),
    children: [
      {
        path: '/terminal',
        name: 'Terminal',
        component: () => import('@/views/contents/Terminal.vue'),
      },
      {
        path: 'content',
        name: 'Content',
        children: [
          {
            path: 'host',
            name: 'Host',
            component: () => import('@/views/contents/Host.vue'),
          },
          {
            path: 'credentials',
            name: 'Credentials',
            component: () => import('@/views/contents/Credentials.vue'),
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
