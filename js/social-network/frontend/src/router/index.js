import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import FeedView from '../views/FeedView.vue'
import ExploreView from '../views/ExploreView.vue'
import ProfileView from '../views/ProfileView.vue'
import UserProfileView from '../views/UserProfileView.vue'
import NotificationsView from '../views/NotificationsView.vue'
import ChatView from '../views/ChatView.vue'
import GroupsView from '../views/GroupsView.vue'
import GroupDetailView from '../views/GroupDetailView.vue'
import { isLoggedIn, setLoggedIn } from '../stores/auth'

const routes = [
    { path: '/', redirect: '/feed' },
    { path: '/login', component: LoginView, meta: { public: true } },
    { path: '/register', component: RegisterView, meta: { public: true } },
    { path: '/feed', component: FeedView, meta: { public: true } },
    { path: '/explore', component: ExploreView },
    { path: '/profile', component: ProfileView },
    { path: '/profile/:id', component: UserProfileView },
    { path: '/notifications', component: NotificationsView },
    { path: '/chat', component: ChatView },
    { path: '/groups', component: GroupsView },
    { path: '/groups/:id', component: GroupDetailView },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior() {
        return { top: 0 }
    },
})

let isAuthChecked = false

router.beforeEach(async (to) => {
    if (!isLoggedIn.value && !isAuthChecked) {
        isAuthChecked = true
        try {
            const response = await fetch('/api/me', { credentials: 'include' })
            if (response.ok) {
                setLoggedIn(true)
            }
        } catch (e) {
            console.error("Auth check failed:", e)
        }
    }

    if (to.meta.public) return true
    if (!isLoggedIn.value) return '/login'
})

export default router