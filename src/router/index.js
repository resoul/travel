import { createRouter, createWebHistory } from 'vue-router'
import { loadLayoutMiddleware } from "@/router/loadLayoutMiddleware";

const routes = [
    {
        path: "/",
        children: [
            {
                path: "/app",
                name: "app",
                component: () => import("../views/pages/app/Main.vue"),
                meta: {
                    title: 'App'
                }
            }
        ],
        meta: {
            layout: 'Main'
        }
    },
    {
        path: '/',
        children: [
            {
                path: "/404",
                name: 'error404',
                component: () => import('../views/pages/errors/404.vue'),
                meta: {
                    title: 'Error'
                }
            }
        ],
        meta: {
            layout: 'Error'
        }
    },
    {
        path: "/:pathMatch(.*)*",
        redirect: '/404'
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: routes
})

router.beforeEach(loadLayoutMiddleware);

export default router