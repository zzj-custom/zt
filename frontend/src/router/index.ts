import { createRouter, createWebHistory,RouteRecordRaw } from 'vue-router';
import MainRoutes from './MainRoutes';
import AuthRoutes from './AuthRoutes';
import {useAuthStore} from '@/pinia/useAuthStore'

const routers:Array<RouteRecordRaw> = [
    {
        path: '/:pathMatch(.*)*',
        component: () => import('@/views/pages/Error404.vue')
    },
    MainRoutes,
    AuthRoutes
]



const router = createRouter({
    history: createWebHistory(),
    routes:routers
});

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();
    if (!authStore.authorization) {
        if (to.path !== '/auth/login') {
            // 如果用户未登录，且不是访问登录页面，重定向到登录页面
            next('/auth/login');
        } else {
            // 如果用户未登录，且是访问登录页面，允许访问
            next();
        }
    } else {
        if (to.path === '/auth/login') {
            // 如果用户已登录，且是访问登录页面，重定向到首页或其他页面
            next('/');
        } else {
            // 其他情况，允许访问
            next();
        }
    }
});

export default router;