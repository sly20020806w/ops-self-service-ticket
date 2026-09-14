import axios from 'axios'
import { ElMessage } from 'element-plus'

const http = axios.create({ baseURL: '/api', timeout: 15000 })

http.interceptors.request.use((cfg) => {
  const token = localStorage.getItem('token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

http.interceptors.response.use(
  (r) => r.data,
  (err) => {
    const msg = err?.response?.data?.message || err.message || '请求失败'
    if (err?.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      if (!location.hash.includes('login')) location.hash = '#/login'
    }
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export const api = {
  login: (data) => http.post('/auth/login', data),
  demoUsers: () => http.get('/auth/demo-users'),
  me: () => http.get('/me'),
  dashboard: () => http.get('/dashboard'),
  tree: () => http.get('/tree/nodes'),
  forms: () => http.get('/forms'),
  form: (id) => http.get(`/forms/${id}`),
  workflows: () => http.get('/workflows'),
  workflow: (id) => http.get(`/workflows/${id}`),
  templates: () => http.get('/templates'),
  tickets: (params) => http.get('/tickets', { params }),
  ticket: (id) => http.get(`/tickets/${id}`),
  ticketActions: (id) => http.get(`/tickets/${id}/actions`),
  createTicket: (data) => http.post('/tickets', data),
  submit: (id) => http.post(`/tickets/${id}/submit`),
  approve: (id, data) => http.post(`/tickets/${id}/approve`, data),
  reject: (id, data) => http.post(`/tickets/${id}/reject`, data),
  execute: (id, data) => http.post(`/tickets/${id}/execute`, data),
  close: (id, data) => http.post(`/tickets/${id}/close`, data),
  cancel: (id, data) => http.post(`/tickets/${id}/cancel`, data),
  comment: (id, data) => http.post(`/tickets/${id}/comment`, data),
  audit: () => http.get('/audit'),
  notifyLogs: () => http.get('/notify-logs'),
  refreshSLA: () => http.post('/sla/refresh'),
}
