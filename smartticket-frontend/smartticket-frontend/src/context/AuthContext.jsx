import { createContext, useContext, useState, useEffect } from 'react'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [usuario, setUsuario] = useState(null)
  const [token, setToken] = useState(null)
  const [loading, setLoading] = useState(true)

  const tokenExpirado = (token) => {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.exp * 1000 < Date.now()
    } catch {
      return true
    }
  }

  useEffect(() => {
    const savedToken = localStorage.getItem('st_token')
    const savedUsuario = localStorage.getItem('st_usuario')
    if (savedToken && savedUsuario) {
      if (tokenExpirado(savedToken)) {
        localStorage.removeItem('st_token')
        localStorage.removeItem('st_usuario')
        localStorage.setItem('st_sesion_expirada', 'true')
      } else {
        setToken(savedToken)
        setUsuario(JSON.parse(savedUsuario))
      }
    }
    setLoading(false)
  }, [])

  useEffect(() => {
    const intervalo = setInterval(() => {
      const t = localStorage.getItem('st_token')
      if (t && tokenExpirado(t)) {
        localStorage.removeItem('st_token')
        localStorage.removeItem('st_usuario')
        localStorage.setItem('st_sesion_expirada', 'true')
        setToken(null)
        setUsuario(null)
        window.location.href = '/login'
      }
    }, 10000)
    return () => clearInterval(intervalo)
  }, [])

  const login = (tokenData, usuarioData) => {
    localStorage.setItem('st_token', tokenData)
    localStorage.setItem('st_usuario', JSON.stringify(usuarioData))
    setToken(tokenData)
    setUsuario(usuarioData)
  }

  const logout = () => {
    localStorage.removeItem('st_token')
    localStorage.removeItem('st_usuario')
    setToken(null)
    setUsuario(null)
  }

  const esAdmin = () => usuario?.rol === 'administrador'
  const esCliente = () => usuario?.rol === 'cliente'

  return (
    <AuthContext.Provider value={{ usuario, token, loading, login, logout, esAdmin, esCliente }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
