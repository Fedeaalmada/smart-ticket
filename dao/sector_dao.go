package dao

import (
	"ticketek/backend/clients"
	"ticketek/backend/domain"
)

func CrearSectores(sectores []domain.Sector) error {
	return clients.DB.Create(&sectores).Error
}

func ObtenerSectorPorID(id uint) (*domain.Sector, error) {
	var sector domain.Sector
	err := clients.DB.First(&sector, id).Error
	return &sector, err
}

func DecrementarCapacidad(sectorID uint) error {
	return clients.DB.Model(&domain.Sector{}).
		Where("id = ? AND capacidad_disponible > 0", sectorID).
		UpdateColumn("capacidad_disponible", clients.DB.Raw("capacidad_disponible - 1")).Error
}

func IncrementarCapacidad(sectorID uint) error {
	return clients.DB.Model(&domain.Sector{}).
		Where("id = ?", sectorID).
		UpdateColumn("capacidad_disponible", clients.DB.Raw("capacidad_disponible + 1")).Error
}

func ActualizarSector(id uint, precio float64, capacidadMaxima uint) error {
	cambios := map[string]interface{}{}
	if precio > 0 {
		cambios["precio"] = precio
	}
	if capacidadMaxima > 0 {
		cambios["capacidad_maxima"] = capacidadMaxima
	}
	if len(cambios) == 0 {
		return nil
	}
	return clients.DB.Model(&domain.Sector{}).Where("id = ?", id).Updates(cambios).Error
}

func EliminarSectoresPorEvento(eventoID uint) error {
	return clients.DB.Where("evento_id = ?", eventoID).Delete(&domain.Sector{}).Error
}
