import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import Navbar from '../components/common/Navbar'
import { eventoService } from '../services/api'
import styles from './AdminReporte.module.css'

export default function AdminReporte() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [evento, setEvento] = useState(null)
  const [reporte, setReporte] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const cargar = async () => {
      try {
        const [evRes, repRes] = await Promise.all([
          eventoService.getById(id),
          eventoService.getReporte(id),
        ])
        setEvento(evRes.data)
        setReporte(repRes.data)
      } catch {
        setError('No se pudo cargar el reporte.')
      } finally {
        setLoading(false)
      }
    }
    cargar()
  }, [id])

  const pct = reporte && reporte.capacidad_total > 0
    ? Math.round((reporte.total_vendidas / reporte.capacidad_total) * 100)
    : 0

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

  return (
    <div className={styles.page}>
      <Navbar />
      <div className={styles.container}>
        <button onClick={() => navigate('/admin')} className={styles.btnVolver}>← Volver al panel</button>

        <div className={styles.header}>
          <span className={styles.cat}>{evento?.categoria}</span>
          <h1 className={styles.titulo}>{evento?.titulo}</h1>
          <p className={styles.sub}>
            {evento ? new Date(evento.fecha).toLocaleDateString('es-AR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }) : ''} · {evento?.horario} hs
          </p>
        </div>

        {/* Métricas */}
        <div className={styles.metricsGrid}>
          <div className={styles.metricCard}>
            <span className={styles.metricLabel}>Entradas vendidas</span>
            <span className={styles.metricValue}>{reporte?.total_vendidas?.toLocaleString('es-AR') || 0}</span>
          </div>
          <div className={styles.metricCard}>
            <span className={styles.metricLabel}>Capacidad total</span>
            <span className={styles.metricValue}>{reporte?.capacidad_total?.toLocaleString('es-AR') || 0}</span>
          </div>
          <div className={styles.metricCard}>
            <span className={styles.metricLabel}>Ingresos totales</span>
            <span className={`${styles.metricValue} ${styles.metricRed}`}>
              ${reporte?.ingresos_totales?.toLocaleString('es-AR') || 0}
            </span>
          </div>
          <div className={styles.metricCard}>
            <span className={styles.metricLabel}>Ocupación</span>
            <span className={`${styles.metricValue} ${pct > 75 ? styles.metricGreen : ''}`}>{pct}%</span>
          </div>
        </div>

        {/* Barra de progreso */}
        <div className={styles.progressSection}>
          <div className={styles.progressHeader}>
            <span>Ocupación del evento</span>
            <span>{reporte?.total_vendidas} / {reporte?.capacidad_total}</span>
          </div>
          <div className={styles.progressBar}>
            <div className={styles.progressFill} style={{ width: `${pct}%` }} />
          </div>
        </div>

        {/* Sectores */}
        {evento?.sectores?.length > 0 && (
          <div className={styles.sectoresSection}>
            <h2 className={styles.sectionTitle}>Desglose por sector</h2>
            <div className={styles.sectoresGrid}>
              {evento.sectores.map((s) => {
                const vendidas = s.capacidad_maxima - s.capacidad_disponible
                const pctSector = s.capacidad_maxima > 0 ? Math.round((vendidas / s.capacidad_maxima) * 100) : 0
                return (
                  <div key={s.id} className={styles.sectorCard}>
                    <div className={styles.sectorHeader}>
                      <span className={styles.sectorNombre}>{s.nombre}</span>
                      <span className={styles.sectorPrecio}>${s.precio?.toLocaleString('es-AR')}</span>
                    </div>
                    <div className={styles.sectorStats}>
                      <span>{vendidas} vendidas / {s.capacidad_maxima} total</span>
                      <span className={styles.sectorPct}>{pctSector}%</span>
                    </div>
                    <div className={styles.sectorBar}>
                      <div className={styles.sectorBarFill} style={{ width: `${pctSector}%` }} />
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
