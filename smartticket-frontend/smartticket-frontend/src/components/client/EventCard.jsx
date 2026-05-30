import { useNavigate } from 'react-router-dom'
import styles from './EventCard.module.css'

const ICONS = {
  Concierto: '🎵',
  Deporte: '⚽',
  Teatro: '🎭',
  Feria: '🎪',
  default: '🎫',
}

export default function EventCard({ evento }) {
  const navigate = useNavigate()
  const icon = ICONS[evento.categoria] || ICONS.default
  const precioMin = evento.sectores?.length
    ? Math.min(...evento.sectores.map((s) => s.precio))
    : null

  const fechaFormateada = new Date(evento.fecha).toLocaleDateString('es-AR', {
    day: 'numeric', month: 'short', year: 'numeric',
  })

  return (
    <div className={styles.card} onClick={() => navigate(`/eventos/${evento.id}`)}>
      <div className={styles.thumb}>
        <span className={styles.icon}>{icon}</span>
        <span className={styles.cat}>{evento.categoria}</span>
      </div>
      <div className={styles.body}>
        <div className={styles.name}>{evento.titulo}</div>
        <div className={styles.meta}>
          <span>📅 {fechaFormateada}</span>
          <span>🕐 {evento.horario}</span>
        </div>
        <div className={styles.footer}>
          <span className={styles.precio}>
            {precioMin ? `Desde $${precioMin.toLocaleString('es-AR')}` : 'Ver precios'}
          </span>
          <div className={styles.sectores}>
            {evento.sectores?.slice(0, 2).map((s) => (
              <span key={s.id} className={styles.sectorBadge}>{s.nombre}</span>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
