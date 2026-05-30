import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Navbar from '../components/common/Navbar'
import { eventoService } from '../services/api'
import styles from './AdminPanel.module.css'

export default function AdminPanel() {
  const navigate = useNavigate()
  const [eventos, setEventos] = useState([])
  const [loading, setLoading] = useState(true)
  const [msg, setMsg] = useState(null)

  const cargar = async () => {
    setLoading(true)
    try {
      const res = await eventoService.getAll()
      setEventos(res.data || [])
    } catch {
      setMsg({ tipo: 'err', texto: 'Error al cargar eventos.' })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { cargar() }, [])

  const handleCancelar = async (id, titulo) => {
    if (!confirm(`¿Cancelar el evento "${titulo}"? Esta acción no se puede deshacer.`)) return
    try {
      await eventoService.cancelar(id)
      setMsg({ tipo: 'ok', texto: 'Evento cancelado correctamente.' })
      cargar()
    } catch (err) {
      setMsg({ tipo: 'err', texto: err.response?.data?.error || 'Error al cancelar.' })
    }
  }

  const fechaFormateada = (fecha) =>
    new Date(fecha).toLocaleDateString('es-AR', { day: 'numeric', month: 'short', year: 'numeric' })

  return (
    <div className={styles.page}>
      <Navbar />
      <div className={styles.container}>
        <div className={styles.header}>
          <div>
            <h1 className={styles.titulo}>Panel de administración</h1>
            <p className={styles.sub}>Gestioná todos los eventos del sistema</p>
          </div>
          <button onClick={() => navigate('/admin/eventos/nuevo')} className={styles.btnNuevo}>
            + Nuevo evento
          </button>
        </div>

        {msg && (
          <div className={msg.tipo === 'ok' ? styles.msgOk : styles.msgErr}>
            {msg.texto}
            <button onClick={() => setMsg(null)} className={styles.msgClose}>✕</button>
          </div>
        )}

        {loading && <div className={styles.loading}><div className={styles.spinner} /></div>}

        {!loading && (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Evento</th>
                  <th>Fecha</th>
                  <th>Categoría</th>
                  <th>Capacidad</th>
                  <th>Estado</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {eventos.map((ev) => (
                  <tr key={ev.id} className={ev.estado !== 'activo' ? styles.rowInactiva : ''}>
                    <td className={styles.tdTitulo}>{ev.titulo}</td>
                    <td>{fechaFormateada(ev.fecha)}</td>
                    <td><span className={styles.catBadge}>{ev.categoria}</span></td>
                    <td>{ev.capacidad_total?.toLocaleString('es-AR')}</td>
                    <td>
                      <span className={`${styles.estadoBadge} ${styles[`estado_${ev.estado}`]}`}>
                        {ev.estado}
                      </span>
                    </td>
                    <td>
                      <div className={styles.acciones}>
                        <button onClick={() => navigate(`/admin/eventos/${ev.id}/reporte`)} className={styles.btnAccion}>
                          Reporte
                        </button>
                        <button onClick={() => navigate(`/admin/eventos/${ev.id}/editar`)} className={styles.btnAccion}>
                          Editar
                        </button>
                        {ev.estado === 'activo' && (
                          <button onClick={() => handleCancelar(ev.id, ev.titulo)} className={`${styles.btnAccion} ${styles.btnDanger}`}>
                            Cancelar
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {eventos.length === 0 && (
              <p className={styles.empty}>No hay eventos cargados.</p>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
