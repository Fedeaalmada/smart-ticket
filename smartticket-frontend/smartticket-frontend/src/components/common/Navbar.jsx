import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import styles from './Navbar.module.css'

export default function Navbar() {
  const { usuario, logout, esAdmin } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <nav className={styles.nav}>
      <Link to="/" className={styles.logo}>
        Smart<span>Ticket</span>
      </Link>

      <div className={styles.links}>
        <Link to="/" className={styles.link}>Eventos</Link>

        {usuario ? (
          <>
            {esAdmin() ? (
              <Link to="/admin" className={styles.link}>Panel Admin</Link>
            ) : (
              <Link to="/mis-entradas" className={styles.link}>Mis entradas</Link>
            )}
            <span className={styles.usuario}>{usuario.nombre}</span>
            <button onClick={handleLogout} className={styles.btnOutline}>
              Salir
            </button>
          </>
        ) : (
          <>
            <Link to="/login" className={styles.link}>Ingresar</Link>
            <Link to="/register" className={styles.btnSolid}>Registrarse</Link>
          </>
        )}
      </div>
    </nav>
  )
}
