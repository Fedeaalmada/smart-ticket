package domain

// ── Auth ─────────────────────────────────────────────────────

type RegisterRequest struct {
	Nombre   string `json:"nombre"   binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string   `json:"token"`
	Usuario UsuarioDTO `json:"usuario"`
}

type UsuarioDTO struct {
	ID     uint       `json:"id"`
	Nombre string     `json:"nombre"`
	Email  string     `json:"email"`
	Rol    RolUsuario `json:"rol"`
}

// ── Eventos ──────────────────────────────────────────────────

type CrearEventoRequest struct {
	Titulo          string         `json:"titulo"           binding:"required"`
	Descripcion     string         `json:"descripcion"`
	FotoURL         string         `json:"foto_url"`
	Fecha           string         `json:"fecha"            binding:"required"` // "2026-08-15"
	Horario         string         `json:"horario"          binding:"required"` // "21:00"
	DuracionMinutos uint           `json:"duracion_minutos"`
	Categoria       string         `json:"categoria"        binding:"required"`
	Sectores        []CrearSectorRequest `json:"sectores"   binding:"required,min=1"`
}

type ActualizarEventoRequest struct {
	Titulo          string `json:"titulo"`
	Descripcion     string `json:"descripcion"`
	FotoURL         string `json:"foto_url"`
	Fecha           string `json:"fecha"`
	Horario         string `json:"horario"`
	DuracionMinutos uint   `json:"duracion_minutos"`
	Categoria       string `json:"categoria"`
	Estado          string `json:"estado"`
}

// ── Sectores ─────────────────────────────────────────────────

type CrearSectorRequest struct {
	Nombre          string  `json:"nombre"           binding:"required"`
	CapacidadMaxima uint    `json:"capacidad_maxima" binding:"required,min=1"`
	Precio          float64 `json:"precio"           binding:"required,min=0"`
}

// ── Compra ───────────────────────────────────────────────────

type ComprarEntradaRequest struct {
	EventoID uint `json:"evento_id" binding:"required"`
	SectorID uint `json:"sector_id" binding:"required"`
}

// ── Transferencia ─────────────────────────────────────────────

type TransferirEntradaRequest struct {
	EmailDestino string `json:"email_destino" binding:"required,email"`
}

// ── Responses genéricas ───────────────────────────────────────

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
