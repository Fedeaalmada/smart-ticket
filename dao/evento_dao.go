package dao

import (
	"ticketek/backend/clients"
	"ticketek/backend/domain"
)

func ObtenerEventos(filtros map[string]string) ([]domain.Evento, error) {
	var eventos []domain.Evento
	query := clients.DB.Preload("Sectores").Where("estado = 'activo'")

	if categoria, ok := filtros["categoria"]; ok && categoria != "" {
		query = query.Where("categoria = ?", categoria)
	}
	if fecha, ok := filtros["fecha"]; ok && fecha != "" {
		query = query.Where("fecha = ?", fecha)
	}
	if titulo, ok := filtros["titulo"]; ok && titulo != "" {
		query = query.Where("titulo LIKE ?", "%"+titulo+"%")
	}

	err := query.Find(&eventos).Error
	return eventos, err
}

func ObtenerEventoPorID(id uint) (*domain.Evento, error) {
	var evento domain.Evento
	err := clients.DB.Preload("Sectores").First(&evento, id).Error
	return &evento, err
}

func CrearEvento(evento *domain.Evento) error {
	return clients.DB.Create(evento).Error
}

func ActualizarEvento(id uint, cambios map[string]interface{}) error {
	return clients.DB.Model(&domain.Evento{}).Where("id = ?", id).Updates(cambios).Error
}

func EliminarEvento(id uint) error {
	// Soft delete: cambia estado a cancelado en lugar de borrar la fila
	return clients.DB.Model(&domain.Evento{}).Where("id = ?", id).
		Update("estado", domain.EstadoCancelado).Error
}

func ObtenerTodosLosEventos() ([]domain.Evento, error) {
	var eventos []domain.Evento
	err := clients.DB.Preload("Sectores").Find(&eventos).Error
	return eventos, err
}
