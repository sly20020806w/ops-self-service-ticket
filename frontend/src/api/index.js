import axios from 'axios'
import { ElMessage } from 'element-plus'

const http = axios.create({ baseURL: '/', timeout: 15000 })

http.interceptors.request.use((cfg) => {
  const token = localStorage.getItem('token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

http.interceptors.response.use(
  (r) => r,
  (err) => {
    const msg = err?.response?.data?.error || err.message || '请求失败'
    if (err?.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      if (!location.hash.includes('/login')) location.hash = '#/login'
    } else {
      ElMessage.error(msg)
    }
    return Promise.reject(err)
  }
)

export const api = {
  login: (username, password) => http.post('/api/login', { username, password }),
  me: () => http.get('/api/me'),
  dashboard: () => http.get('/api/dashboard'),
  tickets: (status) => http.get('/api/tickets', { params: { status } }),
  ticket: (id) => http.get(`/api/tickets/${id}`),
  createTicket: (data) => http.post('/api/tickets', data),
  action: (id, action, comment) => http.post(`/api/tickets/${id}/actions`, { action, comment }),
  forms: () => http.get('/api/forms'),
  flows: () => http.get('/api/flows'),
  tree: () => http.get('/api/tree'),
  audits: () => http.get('/api/audits'),
  notifies: () => http.get('/api/notifies'),
}

export default http
