import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'

// Páginas
import Catalogo from './pages/Catalogo'
import DetalleEvento from './pages/DetalleEvento'
import Login from './pages/Login'
import Register from './pages/Register'
import MisEntradas from './pages/MisEntradas'
import AdminPanel from './pages/AdminPanel'
import AdminFormEvento from './pages/AdminFormEvento'
import AdminReporte from './pages/AdminReporte'

// Rutas protegidas
function RutaProtegida({ children }) {
  const { token, loading } = useAuth()
  if (loading) return <div style={{ color: '#888', padding: '2rem', textAlign: 'center' }}>Cargando...</div>
  return token ? children : <Navigate to="/login" replace />
}

function RutaAdmin({ children }) {
  const { token, esAdmin, loading } = useAuth()
  if (loading) return <div style={{ color: '#888', padding: '2rem', textAlign: 'center' }}>Cargando...</div>
  if (!token) return <Navigate to="/login" replace />
  if (!esAdmin()) return <Navigate to="/" replace />
  return children
}

function AppRoutes() {
  return (
    <Routes>
      {/* Públicas */}
      <Route path="/" element={<Catalogo />} />
      <Route path="/eventos/:id" element={<DetalleEvento />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      {/* Cliente autenticado */}
      <Route path="/mis-entradas" element={
        <RutaProtegida><MisEntradas /></RutaProtegida>
      } />

      {/* Admin */}
      <Route path="/admin" element={
        <RutaAdmin><AdminPanel /></RutaAdmin>
      } />
      <Route path="/admin/eventos/nuevo" element={
        <RutaAdmin><AdminFormEvento /></RutaAdmin>
      } />
      <Route path="/admin/eventos/:id/editar" element={
        <RutaAdmin><AdminFormEvento /></RutaAdmin>
      } />
      <Route path="/admin/eventos/:id/reporte" element={
        <RutaAdmin><AdminReporte /></RutaAdmin>
      } />

      {/* Fallback */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </AuthProvider>
  )
}
