import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import Navbar from '../components/common/Navbar'
import { eventoService } from '../services/api'
import styles from './AdminFormEvento.module.css'

const CATEGORIAS = ['Concierto', 'Deporte', 'Teatro', 'Feria', 'Conferencia', 'Otro']

const sectorVacio = () => ({ nombre: '', capacidad_maxima: '', precio: '' })

export default function AdminFormEvento() {
  const { id } = useParams()
  const navigate = useNavigate()
  const esEdicion = Boolean(id)

  const [form, setForm] = useState({
    titulo: '', descripcion: '', foto_url: '',
    fecha: '', horario: '', duracion_minutos: 120, categoria: 'Concierto',
  })
  const [sectores, setSectores] = useState([sectorVacio()])
  const [loading, setLoading] = useState(false)
  const [loadingData, setLoadingData] = useState(esEdicion)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (!esEdicion) return
    const cargar = async () => {
      try {
        const res = await eventoService.getById(id)
        const ev = res.data
        setForm({
          titulo: ev.titulo || '',
          descripcion: ev.descripcion || '',
          foto_url: ev.foto_url || '',
          fecha: ev.fecha ? ev.fecha.split('T')[0] : '',
          horario: ev.horario || '',
          duracion_minutos: ev.duracion_minutos || 120,
          categoria: ev.categoria || 'Concierto',
        })
      } catch {
        setError('No se pudo cargar el evento.')
      } finally {
        setLoadingData(false)
      }
    }
    cargar()
  }, [id, esEdicion])

  const handleField = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSector = (i, field, value) => {
    const nuevos = [...sectores]
    nuevos[i][field] = value
    setSectores(nuevos)
  }

  const agregarSector = () => setSectores([...sectores, sectorVacio()])
  const eliminarSector = (i) => setSectores(sectores.filter((_, idx) => idx !== i))

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      if (esEdicion) {
        await eventoService.actualizar(id, form)
      } else {
        const payload = {
          ...form,
          duracion_minutos: Number(form.duracion_minutos),
          sectores: sectores.map((s) => ({
            nombre: s.nombre,
            capacidad_maxima: Number(s.capacidad_maxima),
            precio: Number(s.precio),
          })),
        }
        await eventoService.crear(payload)
      }
      navigate('/admin')
    } catch (err) {
      setError(err.response?.data?.error || 'Error al guardar el evento.')
    } finally {
      setLoading(false)
    }
  }

  if (loadingData) return (
    <div className={styles.page}><Navbar />
      <div className={styles.center}><div className={styles.spinner} /></div>
    </div>
  )

  return (
    <div className={styles.page}>
      <Navbar />
      <div className={styles.container}>
        <div className={styles.header}>
          <button onClick={() => navigate('/admin')} className={styles.btnVolver}>← Volver</button>
          <h1 className={styles.titulo}>{esEdicion ? 'Editar evento' : 'Nuevo evento'}</h1>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          {error && <p className={styles.error}>{error}</p>}

          <section className={styles.section}>
            <h2 className={styles.sectionTitle}>Información del evento</h2>
            <div className={styles.grid2}>
              <div className={styles.field}>
                <label>Título *</label>
                <input name="titulo" value={form.titulo} onChange={handleField} required placeholder="Nombre del evento" />
              </div>
              <div className={styles.field}>
                <label>Categoría *</label>
                <select name="categoria" value={form.categoria} onChange={handleField}>
                  {CATEGORIAS.map((c) => <option key={c}>{c}</option>)}
                </select>
              </div>
              <div className={styles.field}>
                <label>Fecha *</label>
                <input type="date" name="fecha" value={form.fecha} onChange={handleField} required />
              </div>
              <div className={styles.field}>
                <label>Horario *</label>
                <input type="time" name="horario" value={form.horario} onChange={handleField} required />
              </div>
              <div className={styles.field}>
                <label>Duración (minutos)</label>
                <input type="number" name="duracion_minutos" value={form.duracion_minutos} onChange={handleField} min="1" />
              </div>
              <div className={styles.field}>
                <label>URL de foto</label>
                <input name="foto_url" value={form.foto_url} onChange={handleField} placeholder="https://..." />
              </div>
            </div>
            <div className={styles.field}>
              <label>Descripción</label>
              <textarea name="descripcion" value={form.descripcion} onChange={handleField} rows={3} placeholder="Descripción del evento..." />
            </div>
          </section>

          {!esEdicion && (
            <section className={styles.section}>
              <div className={styles.sectionHeader}>
                <h2 className={styles.sectionTitle}>Sectores y precios</h2>
                <button type="button" onClick={agregarSector} className={styles.btnAgregar}>
                  + Agregar sector
                </button>
              </div>
              {sectores.map((s, i) => (
                <div key={i} className={styles.sectorRow}>
                  <div className={styles.field}>
                    <label>Nombre del sector</label>
                    <input value={s.nombre} onChange={(e) => handleSector(i, 'nombre', e.target.value)} required placeholder="Ej: Campo, Platea, VIP" />
                  </div>
                  <div className={styles.field}>
                    <label>Capacidad máxima</label>
                    <input type="number" value={s.capacidad_maxima} onChange={(e) => handleSector(i, 'capacidad_maxima', e.target.value)} required min="1" placeholder="Ej: 500" />
                  </div>
                  <div className={styles.field}>
                    <label>Precio ($)</label>
                    <input type="number" value={s.precio} onChange={(e) => handleSector(i, 'precio', e.target.value)} required min="0" placeholder="Ej: 15000" />
                  </div>
                  {sectores.length > 1 && (
                    <button type="button" onClick={() => eliminarSector(i)} className={styles.btnEliminarSector}>✕</button>
                  )}
                </div>
              ))}
            </section>
          )}

          <div className={styles.formActions}>
            <button type="button" onClick={() => navigate('/admin')} className={styles.btnCancelar}>
              Cancelar
            </button>
            <button type="submit" disabled={loading} className={styles.btnGuardar}>
              {loading ? 'Guardando...' : esEdicion ? 'Guardar cambios' : 'Crear evento'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
