import { createRouter, createWebHashHistory } from 'vue-router'
import Login from './views/Login.vue'
import Layout from './views/Layout.vue'
import Dashboard from './views/Dashboard.vue'
import Tickets from './views/Tickets.vue'
import TicketDetail from './views/TicketDetail.vue'
import CreateTicket from './views/CreateTicket.vue'
import Workflows from './views/Workflows.vue'
import Forms from './views/Forms.vue'
import Audit from './views/Audit.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/login', component: Login },
    {
      path: '/',
      component: Layout,
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: Dashboard },
        { path: 'tickets', component: Tickets },
        { path: 'tickets/create', component: CreateTicket },
        { path: 'tickets/:id', component: TicketDetail },
        { path: 'workflows', component: Workflows },
        { path: 'forms', component: Forms },
        { path: 'audit', component: Audit },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  if (to.path !== '/login' && !localStorage.getItem('token')) next('/login')
  else next()
})

export default router
