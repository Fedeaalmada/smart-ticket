import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' }
})

// Interceptor: agrega el token JWT a cada request si existe
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('st_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Interceptor: si el token expiró, limpia la sesión
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('st_token')
      localStorage.removeItem('st_usuario')
      localStorage.setItem('st_sesion_expirada', 'true')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

// ── Auth ──────────────────────────────────────────────────
export const authService = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
}

// ── Eventos ───────────────────────────────────────────────
export const eventoService = {
  getAll: (filtros = {}) => api.get('/eventos', { params: filtros }),
  getById: (id) => api.get(`/eventos/${id}`),
  crear: (data) => api.post('/admin/eventos', data),
  actualizar: (id, data) => api.put(`/admin/eventos/${id}`, data),
  cancelar: (id) => api.delete(`/admin/eventos/${id}`),
  getReporte: (id) => api.get(`/admin/eventos/${id}/reporte`),
}

// ── Entradas ──────────────────────────────────────────────
export const entradaService = {
  comprar: (data) => api.post('/entradas/comprar', data),
  misEntradas: () => api.get('/entradas/mis-entradas'),
  cancelar: (id) => api.delete(`/entradas/${id}`),
  transferir: (id, data) => api.post(`/entradas/${id}/transferir`, data),
}

export default api
