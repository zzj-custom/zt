import path from 'path';

const MainRoutes = {
    path: '/main',
    meta: {
        requiresAuth: true
    },
    redirect: '/main',
    component: () => import('@/layout/Layout.vue'),
    children: [
        {
            name: 'Dashboard',
            path: '/',
            component: () => import('@/views/dashboard/index.vue')
        },
        {
            name: '网易云音乐格式转换',
            path: '/netease/index',
            component: () => import('@/views/netease/index.vue')
        },
        {
            name: 'BING图片',
            path: '/bing/images',
            component: () => import('@/views/bing/index.vue')
        },
        {
            name: '视频下载',
            path: '/video/download',
            component: () => import('@/views/video/index.vue')
        }
    ]
};

export default MainRoutes;
