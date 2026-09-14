import { createRouter, createWebHashHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import Tickets from '../views/Tickets.vue'
import TicketDetail from '../views/TicketDetail.vue'
import Flows from '../views/Flows.vue'
import Audits from '../views/Audits.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/', component: Dashboard },
    { path: '/tickets', component: Tickets },
    { path: '/tickets/:id', component: TicketDetail },
    { path: '/flows', component: Flows },
    { path: '/audits', component: Audits },
  ],
})

router.beforeEach((to) => {
  if (to.path !== '/login' && !localStorage.getItem('token')) return '/login'
})

export default router
