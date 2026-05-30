import { useState, useEffect } from 'react'
import Navbar from '../components/common/Navbar'
import { entradaService } from '../services/api'
import styles from './MisEntradas.module.css'

export default function MisEntradas() {
  const [entradas, setEntradas] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [modalTransferir, setModalTransferir] = useState(null)
  const [emailDestino, setEmailDestino] = useState('')
  const [procesando, setProcesando] = useState(false)
  const [msg, setMsg] = useState(null)

  const cargar = async () => {
    setLoading(true)
    try {
      const res = await entradaService.misEntradas()
      setEntradas(res.data || [])
    } catch {
      setError('No se pudieron cargar tus entradas.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { cargar() }, [])

  const handleCancelar = async (id) => {
    if (!confirm('¿Estás seguro que querés cancelar esta entrada?')) return
    try {
      await entradaService.cancelar(id)
      setMsg({ tipo: 'ok', texto: 'Entrada cancelada correctamente.' })
      cargar()
    } catch (err) {
      setMsg({ tipo: 'err', texto: err.response?.data?.error || 'Error al cancelar.' })
    }
  }

  const handleTransferir = async (e) => {
    e.preventDefault()
    setProcesando(true)
    try {
      await entradaService.transferir(modalTransferir, { email_destino: emailDestino })
      setMsg({ tipo: 'ok', texto: 'Entrada transferida correctamente.' })
      setModalTransferir(null)
      setEmailDestino('')
      cargar()
    } catch (err) {
      setMsg({ tipo: 'err', texto: err.response?.data?.error || 'Error al transferir.' })
    } finally {
      setProcesando(false)
    }
  }

  const fechaFormateada = (fecha) =>
    new Date(fecha).toLocaleDateString('es-AR', { day: 'numeric', month: 'short', year: 'numeric' })

  return (
    <div className={styles.page}>
      <Navbar />
      <div className={styles.container}>
        <h1 className={styles.titulo}>Mis entradas</h1>

        {msg && (
          <div className={msg.tipo === 'ok' ? styles.msgOk : styles.msgErr}>
            {msg.texto}
            <button onClick={() => setMsg(null)} className={styles.msgClose}>✕</button>
          </div>
        )}

        {loading && <div className={styles.loading}><div className={styles.spinner} /></div>}
        {error && <div className={styles.error}>{error}</div>}

        {!loading && !error && entradas.length === 0 && (
          <div className={styles.empty}>
            <span>🎫</span>
            <p>No tenés entradas todavía.</p>
          </div>
        )}

        {!loading && entradas.map((entrada) => (
          <div key={entrada.id} className={`${styles.card} ${entrada.estado !== 'activa' ? styles.cardInactiva : ''}`}>
            <div className={styles.cardAccent} />
            <div className={styles.cardBody}>
              <div className={styles.cardTop}>
                <div>
                  <div className={styles.cardEvento}>{entrada.evento?.titulo || `Evento #${entrada.evento_id}`}</div>
                  <div className={styles.cardMeta}>
                    {entrada.sector?.nombre} · {entrada.evento ? fechaFormateada(entrada.evento.fecha) : ''} · {entrada.evento?.horario} hs
                  </div>
                </div>
                <div className={styles.cardRight}>
                  <span className={styles.cardPrecio}>${entrada.precio_pagado?.toLocaleString('es-AR')}</span>
                  <span className={`${styles.badge} ${styles[`badge_${entrada.estado}`]}`}>
                    {entrada.estado}
                  </span>
                </div>
              </div>

              {entrada.estado === 'activa' && (
                <div className={styles.cardActions}>
                  <button
                    className={styles.btnTransferir}
                    onClick={() => { setModalTransferir(entrada.id); setMsg(null) }}
                  >
                    Transferir
                  </button>
                  <button
                    className={styles.btnCancelar}
                    onClick={() => handleCancelar(entrada.id)}
                  >
                    Cancelar
                  </button>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* Modal transferir */}
      {modalTransferir && (
        <div className={styles.modalOverlay} onClick={() => setModalTransferir(null)}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <h3 className={styles.modalTitle}>Transferir entrada</h3>
            <p className={styles.modalSub}>Ingresá el email del usuario al que querés transferirle la entrada.</p>
            <form onSubmit={handleTransferir} className={styles.modalForm}>
              <input
                type="email"
                placeholder="email@ejemplo.com"
                value={emailDestino}
                onChange={(e) => setEmailDestino(e.target.value)}
                required
                className={styles.modalInput}
              />
              {msg?.tipo === 'err' && <p className={styles.modalError}>{msg.texto}</p>}
              <div className={styles.modalBtns}>
                <button type="submit" disabled={procesando} className={styles.btnPrimary}>
                  {procesando ? 'Transfiriendo...' : 'Confirmar transferencia'}
                </button>
                <button type="button" onClick={() => setModalTransferir(null)} className={styles.btnOutline}>
                  Cancelar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
