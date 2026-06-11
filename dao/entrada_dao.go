package dao

import (
	"ticketek/backend/clients"
	"ticketek/backend/domain"
)

func CrearEntrada(entrada *domain.Entrada) error {
	return clients.DB.Create(entrada).Error
}

func ObtenerEntradasPorUsuario(usuarioID uint) ([]domain.Entrada, error) {
	var entradas []domain.Entrada
	err := clients.DB.
		Preload("Evento").
		Preload("Sector").
		Where("usuario_id = ?", usuarioID).
		Find(&entradas).Error
	return entradas, err
}

func ObtenerEntradaPorID(id uint) (*domain.Entrada, error) {
	var entrada domain.Entrada
	err := clients.DB.Preload("Sector").First(&entrada, id).Error
	return &entrada, err
}

func CancelarEntrada(id uint) error {
	return clients.DB.Model(&domain.Entrada{}).
		Where("id = ?", id).
		Update("estado", domain.EntradaCancelada).Error
}

func TransferirEntrada(entradaID, nuevoUsuarioID uint) error {
	// El ticket cambia de titular pero sigue activo: el nuevo dueño
	// debe poder cancelarlo o volver a transferirlo. El historial de
	// quién lo tuvo antes queda registrado en la tabla "transferencias".
	return clients.DB.Model(&domain.Entrada{}).
		Where("id = ?", entradaID).
		Updates(map[string]interface{}{
			"usuario_id": nuevoUsuarioID,
			"estado":     domain.EntradaActiva,
		}).Error
}

func RegistrarTransferencia(t *domain.Transferencia) error {
	return clients.DB.Create(t).Error
}

// Para el reporte del admin
type ReporteEvento struct {
	TotalVendidas  int64   `json:"total_vendidas"`
	CapacidadTotal uint    `json:"capacidad_total"`
	IngresosTotal  float64 `json:"ingresos_totales"`
}

func ObtenerReporteEvento(eventoID uint) (*ReporteEvento, error) {
	var reporte ReporteEvento
	err := clients.DB.Model(&domain.Entrada{}).
		Select("COUNT(*) as total_vendidas, SUM(precio_pagado) as ingresos_total").
		Where("evento_id = ? AND estado = 'activa'", eventoID).
		Scan(&reporte).Error
	if err != nil {
		return nil, err
	}

	evento, err := ObtenerEventoPorID(eventoID)
	if err != nil {
		return nil, err
	}
	reporte.CapacidadTotal = evento.CapacidadTotal
	return &reporte, nil
}
