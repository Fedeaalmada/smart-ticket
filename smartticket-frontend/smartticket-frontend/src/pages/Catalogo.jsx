import { useState, useEffect } from 'react'
import Navbar from '../components/common/Navbar'
import Particles from '../components/common/Particles'
import EventCard from '../components/client/EventCard'
import { eventoService } from '../services/api'
import styles from './Catalogo.module.css'

const CATEGORIAS = ['Todos', 'Concierto', 'Deporte', 'Teatro', 'Feria']

export default function Catalogo() {
  const [eventos, setEventos] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [busqueda, setBusqueda] = useState('')
  const [categoria, setCategoria] = useState('Todos')

  const cargarEventos = async (filtros = {}) => {
    setLoading(true)
    setError(null)
    try {
      const res = await eventoService.getAll(filtros)
      setEventos(Array.isArray(res.data) ? res.data : [])
    } catch {
      setError('No se pudieron cargar los eventos.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { cargarEventos() }, [])

  const handleBuscar = (e) => {
    e.preventDefault()
    const filtros = {}
    if (busqueda) filtros.titulo = busqueda
    if (categoria !== 'Todos') filtros.categoria = categoria
    cargarEventos(filtros)
  }

  const handleCategoria = (cat) => {
    setCategoria(cat)
    const filtros = {}
    if (busqueda) filtros.titulo = busqueda
    if (cat !== 'Todos') filtros.categoria = cat
    cargarEventos(filtros)
  }

  return (
    <div className={styles.page}>
      <Navbar />

      {/* Hero */}
      <section className={styles.hero}>
        <Particles />
        <div className={styles.heroContent}>
          <p className={styles.eyebrow}>SmartTicket — Argentina</p>
          <h1 className={styles.heroTitle}>
            Tu próxima <span>experiencia</span><br />te espera.
          </h1>
          <p className={styles.heroSub}>Conciertos · Deportes · Teatro · Ferias</p>
          <form onSubmit={handleBuscar} className={styles.searchBar}>
            <input
              type="text"
              placeholder="Buscar eventos, artistas..."
              value={busqueda}
              onChange={(e) => setBusqueda(e.target.value)}
            />
            <button type="submit">Buscar</button>
          </form>
          <div className={styles.chips}>
            {CATEGORIAS.map((cat) => (
              <button
                key={cat}
                className={`${styles.chip} ${categoria === cat ? styles.chipActive : ''}`}
                onClick={() => handleCategoria(cat)}
              >
                {cat}
              </button>
            ))}
          </div>
        </div>
      </section>

      {/* Catálogo */}
      <section className={styles.catalogo}>
        <div className={styles.catalogoHeader}>
          <h2 className={styles.sectionTitle}>Próximos eventos</h2>
          {!loading && <span className={styles.count}>{eventos.length} eventos</span>}
        </div>

        {loading && (
          <div className={styles.loading}>
            <div className={styles.spinner} />
            <span>Cargando eventos...</span>
          </div>
        )}

        {error && <div className={styles.error}>{error}</div>}

        {!loading && !error && eventos.length === 0 && (
          <div className={styles.empty}>No se encontraron eventos para tu búsqueda.</div>
        )}

        {!loading && !error && eventos.length > 0 && (
          <div className={styles.grid}>
            {eventos.map((evento) => (
              <EventCard key={evento.id} evento={evento} />
            ))}
          </div>
        )}
      </section>
      <footer style={{
  background: '#0D0D0D',
  borderTop: '1px solid #1a1a1a',
  padding: '2.5rem 1.5rem 1.5rem',
  textAlign: 'center',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: '1rem'
}}>
  <div style={{ display: 'flex', gap: '1rem', fontSize: '1.3rem' }}>
    <span style={{ cursor: 'pointer', color: '#555' }}>📘</span>
    <span style={{ cursor: 'pointer', color: '#555' }}>📸</span>
    <span style={{ cursor: 'pointer', color: '#555' }}>🐦</span>
  </div>
  <div style={{ fontFamily: 'var(--font-display)', fontSize: '1rem', fontWeight: 700, color: '#333' }}>
    Smart<span style={{ color: '#E03030' }}>Ticket</span>
  </div>
  <p style={{ fontSize: '0.72rem', color: '#333', maxWidth: '500px', lineHeight: 1.6 }}>
    Importante: SmartTicket no se hace responsable de la calidad o satisfacción de los eventos publicados. SmartTicket es un sistema que presta el servicio de venta de entradas online. Al usar este sitio usted acepta los términos y condiciones de la aplicación. Copyright © 2026 SmartTicket.
  </p>
  <p style={{ fontSize: '0.7rem', color: '#2a2a2a' }}>
    Términos y condiciones
  </p>
</footer>
    </div>
  )
}
