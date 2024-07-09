const AuthRoutes = {
    path: '/auth',
    component: () => import('@/layout/BlankLayout.vue'),
    meta: {
        requiresAuth: true
    },
    children: [
        {
            name: 'Login',
            path: '/auth/login',
            component: () => import('@/views/auth/Login.vue')
        }
    ]
};

export default AuthRoutes;
