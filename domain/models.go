package domain

import "time"

type RolUsuario string

const (
	RolCliente       RolUsuario = "cliente"
	RolAdministrador RolUsuario = "administrador"
)

type Usuario struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre       string     `gorm:"size:100;not null" json:"nombre"`
	Email        string     `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Rol          RolUsuario `gorm:"type:enum('cliente','administrador');default:'cliente'" json:"rol"`
	Activo       bool       `gorm:"default:true" json:"activo"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Usuario) TableName() string { return "usuarios" }

type EstadoEvento string

const (
	EstadoActivo    EstadoEvento = "activo"
	EstadoCancelado EstadoEvento = "cancelado"
	EstadoAgotado   EstadoEvento = "agotado"
)

type Evento struct {
	ID              uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo          string       `gorm:"size:200;not null" json:"titulo"`
	Descripcion     string       `gorm:"type:text" json:"descripcion"`
	FotoURL         string       `gorm:"size:500" json:"foto_url"`
	Fecha           time.Time    `gorm:"not null" json:"fecha"`
	Horario         string       `gorm:"size:8;not null" json:"horario"`
	DuracionMinutos uint         `gorm:"default:120" json:"duracion_minutos"`
	Categoria       string       `gorm:"size:100;not null" json:"categoria"`
	CapacidadTotal  uint         `gorm:"not null" json:"capacidad_total"`
	Estado          EstadoEvento `gorm:"type:enum('activo','cancelado','agotado');default:'activo'" json:"estado"`
	CreadoPor       uint         `gorm:"not null" json:"creado_por"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Sectores        []Sector     `gorm:"foreignKey:EventoID" json:"sectores,omitempty"`
}

func (Evento) TableName() string { return "eventos" }

type Sector struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventoID            uint      `gorm:"not null" json:"evento_id"`
	Nombre              string    `gorm:"size:100;not null" json:"nombre"`
	CapacidadMaxima     uint      `gorm:"not null" json:"capacidad_maxima"`
	CapacidadDisponible uint      `gorm:"not null" json:"capacidad_disponible"`
	Precio              float64   `gorm:"type:decimal(10,2);not null" json:"precio"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Evento              *Evento   `gorm:"foreignKey:EventoID" json:"evento,omitempty"`
}

func (Sector) TableName() string { return "sectores" }

type EstadoEntrada string

const (
	EntradaActiva      EstadoEntrada = "activa"
	EntradaCancelada   EstadoEntrada = "cancelada"
	EntradaTransferida EstadoEntrada = "transferida"
)

type Entrada struct {
	ID           uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UsuarioID    uint          `gorm:"not null" json:"usuario_id"`
	SectorID     uint          `gorm:"not null" json:"sector_id"`
	EventoID     uint          `gorm:"not null" json:"evento_id"`
	PrecioPagado float64       `gorm:"type:decimal(10,2);not null" json:"precio_pagado"`
	Estado       EstadoEntrada `gorm:"type:enum('activa','cancelada','transferida');default:'activa'" json:"estado"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Usuario      *Usuario      `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Sector       *Sector       `gorm:"foreignKey:SectorID" json:"sector,omitempty"`
	Evento       *Evento       `gorm:"foreignKey:EventoID" json:"evento,omitempty"`
}

func (Entrada) TableName() string { return "entradas" }

type Transferencia struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EntradaID      uint      `gorm:"not null" json:"entrada_id"`
	UsuarioOrigen  uint      `gorm:"not null" json:"usuario_origen"`
	UsuarioDestino uint      `gorm:"not null" json:"usuario_destino"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Transferencia) TableName() string { return "transferencias" }
