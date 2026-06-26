import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authService } from '../services/api'
import { useAuth } from '../context/AuthContext'
import Particles from '../components/common/Particles'
import styles from './Auth.module.css'

export function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(false)
  const [sesionExpirada, setSesionExpirada] = useState(false)

  useEffect(() => {
    if (localStorage.getItem('st_sesion_expirada') === 'true') {
      setSesionExpirada(true)
    }
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const res = await authService.login(form)
      localStorage.removeItem('st_sesion_expirada')
      login(res.data.token, res.data.usuario)
      navigate(res.data.usuario.rol === 'administrador' ? '/admin' : '/')
    } catch (err) {
      setError(err.response?.data?.error || 'Credenciales inválidas.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <Particles />
      <div className={styles.card}>
        <Link to="/" className={styles.logo}>Smart<span>Ticket</span></Link>
        <h1 className={styles.titulo}>Iniciar sesión</h1>
        {sesionExpirada && (
          <p className={styles.error}>Tu sesión expiró. Por favor iniciá sesión de nuevo.</p>
        )}
        <p className={styles.sub}>Ingresá con tu cuenta para comprar entradas</p>
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label>Email</label>
            <input type="email" placeholder="tu@email.com" value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })} required />
          </div>
          <div className={styles.field}>
            <label>Contraseña</label>
            <input type="password" placeholder="••••••••" value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })} required />
          </div>
          {error && <p className={styles.error}>{error}</p>}
          <button type="submit" disabled={loading} className={styles.btnSubmit}>
            {loading ? 'Ingresando...' : 'Ingresar'}
          </button>
        </form>
        <p className={styles.switch}>
          ¿No tenés cuenta? <Link to="/register">Registrate</Link>
        </p>
      </div>
    </div>
  )
}

export function Register() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ nombre: '', email: '', password: '' })
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const res = await authService.register(form)
      login(res.data.token, res.data.usuario)
      navigate('/')
    } catch (err) {
      setError(err.response?.data?.error || 'Error al registrarse.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <Particles />
      <div className={styles.card}>
        <Link to="/" className={styles.logo}>Smart<span>Ticket</span></Link>
        <h1 className={styles.titulo}>Crear cuenta</h1>
        <p className={styles.sub}>Registrate gratis y empezá a comprar entradas</p>
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label>Nombre completo</label>
            <input type="text" placeholder="Juan Pérez" value={form.nombre}
              onChange={(e) => setForm({ ...form, nombre: e.target.value })} required />
          </div>
          <div className={styles.field}>
            <label>Email</label>
            <input type="email" placeholder="tu@email.com" value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })} required />
          </div>
          <div className={styles.field}>
            <label>Contraseña</label>
            <input type="password" placeholder="Mínimo 6 caracteres" value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })} required minLength={6} />
          </div>
          {error && <p className={styles.error}>{error}</p>}
          <button type="submit" disabled={loading} className={styles.btnSubmit}>
            {loading ? 'Creando cuenta...' : 'Crear cuenta'}
          </button>
        </form>
        <p className={styles.switch}>
          ¿Ya tenés cuenta? <Link to="/login">Iniciá sesión</Link>
        </p>
      </div>
    </div>
  )
}

export default Login
