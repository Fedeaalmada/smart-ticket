package services

import (
	"errors"
	"time"

	"ticketek/backend/dao"
	"ticketek/backend/domain"
)

func ObtenerCatalogoEventos(filtros map[string]string) ([]domain.Evento, error) {
	return dao.ObtenerEventos(filtros)
}

func ObtenerDetalleEvento(id uint) (*domain.Evento, error) {
	return dao.ObtenerEventoPorID(id)
}

func CrearEvento(req domain.CrearEventoRequest, adminID uint) (*domain.Evento, error) {
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		return nil, errors.New("formato de fecha inválido, usar YYYY-MM-DD")
	}

	// Calcular capacidad total como suma de sectores
	var capacidadTotal uint
	for _, s := range req.Sectores {
		capacidadTotal += s.CapacidadMaxima
	}

	duracion := req.DuracionMinutos
	if duracion == 0 {
		duracion = 120
	}

	evento := &domain.Evento{
		Titulo:          req.Titulo,
		Descripcion:     req.Descripcion,
		FotoURL:         req.FotoURL,
		Fecha:           fecha,
		Horario:         req.Horario,
		DuracionMinutos: duracion,
		Categoria:       req.Categoria,
		CapacidadTotal:  capacidadTotal,
		Estado:          domain.EstadoActivo,
		CreadoPor:       adminID,
	}

	if err := dao.CrearEvento(evento); err != nil {
		return nil, errors.New("error al crear el evento")
	}

	// Crear los sectores asociados
	sectores := make([]domain.Sector, len(req.Sectores))
	for i, s := range req.Sectores {
		sectores[i] = domain.Sector{
			EventoID:            evento.ID,
			Nombre:              s.Nombre,
			CapacidadMaxima:     s.CapacidadMaxima,
			CapacidadDisponible: s.CapacidadMaxima, // al inicio disponible = máxima
			Precio:              s.Precio,
		}
	}

	if err := dao.CrearSectores(sectores); err != nil {
		return nil, errors.New("error al crear los sectores")
	}

	evento.Sectores = sectores
	return evento, nil
}

func ActualizarEvento(id uint, req domain.ActualizarEventoRequest) error {
	cambios := map[string]interface{}{}

	if req.Titulo != "" {
		cambios["titulo"] = req.Titulo
	}
	if req.Descripcion != "" {
		cambios["descripcion"] = req.Descripcion
	}
	if req.FotoURL != "" {
		cambios["foto_url"] = req.FotoURL
	}
	if req.Categoria != "" {
		cambios["categoria"] = req.Categoria
	}
	if req.Fecha != "" {
		fecha, err := time.Parse("2006-01-02", req.Fecha)
		if err != nil {
			return errors.New("formato de fecha inválido")
		}
		cambios["fecha"] = fecha
	}
	if req.Horario != "" {
		cambios["horario"] = req.Horario
	}
	if req.DuracionMinutos != 0 {
		cambios["duracion_minutos"] = req.DuracionMinutos
	}
	if req.Estado != "" {
		cambios["estado"] = req.Estado
	}

	if len(cambios) == 0 && len(req.Sectores) == 0 {
		return errors.New("no se enviaron campos para actualizar")
	}

	if len(cambios) > 0 {
		if err := dao.ActualizarEvento(id, cambios); err != nil {
			return err
		}
	}

	for _, s := range req.Sectores {
		if s.ID == 0 {
			continue
		}
		if err := dao.ActualizarSector(s.ID, s.Precio, s.CapacidadMaxima); err != nil {
			return err
		}
	}

	return nil
}

func CancelarEvento(id uint) error {
	_, err := dao.ObtenerEventoPorID(id)
	if err != nil {
		return errors.New("evento no encontrado")
	}
	return dao.EliminarEvento(id)
}

func ObtenerReporteEvento(id uint) (*dao.ReporteEvento, error) {
	return dao.ObtenerReporteEvento(id)
}
