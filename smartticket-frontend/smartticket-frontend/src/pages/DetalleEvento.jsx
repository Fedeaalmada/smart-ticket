import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import Navbar from '../components/common/Navbar'
import { eventoService, entradaService } from '../services/api'
import { useAuth } from '../context/AuthContext'
import styles from './DetalleEvento.module.css'

export default function DetalleEvento() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { usuario } = useAuth()

  const [evento, setEvento] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [sectorSeleccionado, setSectorSeleccionado] = useState(null)
  const [comprando, setComprando] = useState(false)
  const [exito, setExito] = useState(false)
  const [errorCompra, setErrorCompra] = useState(null)

  useEffect(() => {
    const cargar = async () => {
      try {
        const res = await eventoService.getById(id)
        setEvento(res.data)
      } catch {
        setError('No se pudo cargar el evento.')
      } finally {
        setLoading(false)
      }
    }
    cargar()
  }, [id])

  const handleComprar = async () => {
    if (!usuario) { navigate('/login'); return }
    if (!sectorSeleccionado) return
    setComprando(true)
    setErrorCompra(null)
    try {
      await entradaService.comprar({ evento_id: evento.id, sector_id: sectorSeleccionado.id })
      setExito(true)
    } catch (err) {
      setErrorCompra(err.response?.data?.error || 'Error al procesar la compra.')
    } finally {
      setComprando(false)
    }
  }

  const fechaFormateada = evento
    ? new Date(evento.fecha).toLocaleDateString('es-AR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
    : ''

  if (loading) return (
    <div className={styles.page}><Navbar />
      <div className={styles.center}><div className={styles.spinner} /></div>
    </div>
  )

  if (error) return (
    <div className={styles.page}><Navbar />
      <div className={styles.center}><p className={styles.errorMsg}>{error}</p></div>
    </div>
  )

  // Pantalla de éxito
  if (exito) return (
    <div className={styles.page}><Navbar />
      <div className={styles.exitoWrap}>
        <div className={styles.exitoCard}>
          <div className={styles.exitoIcon}>🎫</div>
          <h2 className={styles.exitoTitle}>¡Entrada confirmada!</h2>
          <p className={styles.exitoSub}>
            Compraste una entrada para <strong>{evento.titulo}</strong><br />
            Sector: <strong>{sectorSeleccionado.nombre}</strong> — ${sectorSeleccionado.precio.toLocaleString('es-AR')}
          </p>
          <div className={styles.exitoBtns}>
            <button onClick={() => navigate('/mis-entradas')} className={styles.btnPrimary}>
              Ver mis entradas
            </button>
            <button onClick={() => navigate('/')} className={styles.btnOutline}>
              Volver al catálogo
            </button>
          </div>
        </div>
      </div>
    </div>
  )

  return (
    <div className={styles.page}>
      <Navbar />
      <div className={styles.container}>

        {/* Header del evento */}
        <div className={styles.header}>
          <div className={styles.headerContent}>
            <span className={styles.cat}>{evento.categoria}</span>
            <h1 className={styles.titulo}>{evento.titulo}</h1>
            <p className={styles.descripcion}>{evento.descripcion}</p>
            <div className={styles.metaRow}>
              <span className={styles.metaItem}>📅 {fechaFormateada}</span>
              <span className={styles.metaItem}>🕐 {evento.horario} hs</span>
              <span className={styles.metaItem}>⏱ {evento.duracion_minutos} min</span>
              <span className={styles.metaItem}>👥 Capacidad: {evento.capacidad_total?.toLocaleString('es-AR')}</span>
            </div>
          </div>
        </div>

        {/* Selección de sector */}
        <div className={styles.sectoresSection}>
          <h2 className={styles.sectionLabel}>Elegí tu sector</h2>
          <div className={styles.sectoresList}>
            {evento.sectores?.map((sector) => (
              <div
                key={sector.id}
                className={`${styles.sectorItem} ${sectorSeleccionado?.id === sector.id ? styles.sectorSelected : ''} ${sector.capacidad_disponible === 0 ? styles.sectorAgotado : ''}`}
                onClick={() => sector.capacidad_disponible > 0 && setSectorSeleccionado(sector)}
              >
                <div className={styles.sectorInfo}>
                  <span className={styles.sectorNombre}>{sector.nombre}</span>
                  <span className={styles.sectorDisp}>
                    {sector.capacidad_disponible === 0
                      ? 'Agotado'
                      : `${sector.capacidad_disponible.toLocaleString('es-AR')} lugares disponibles`}
                  </span>
                </div>
                <span className={styles.sectorPrecio}>
                  ${sector.precio.toLocaleString('es-AR')}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Botón de compra */}
        <div className={styles.compraSection}>
          {errorCompra && <p className={styles.errorMsg}>{errorCompra}</p>}
          <button
            className={styles.btnComprar}
            onClick={handleComprar}
            disabled={!sectorSeleccionado || comprando}
          >
            {comprando
              ? 'Procesando...'
              : sectorSeleccionado
                ? `Comprar — ${sectorSeleccionado.nombre} $${sectorSeleccionado.precio.toLocaleString('es-AR')}`
                : 'Seleccioná un sector'}
          </button>
          {!usuario && (
            <p className={styles.loginHint}>
              Necesitás <span onClick={() => navigate('/login')}>iniciar sesión</span> para comprar.
            </p>
          )}
        </div>

      </div>
    </div>
  )
}
