package services

import (
	"errors"

	"ticketek/backend/dao"
	"ticketek/backend/domain"
)

func ComprarEntrada(usuarioID uint, req domain.ComprarEntradaRequest) (*domain.Entrada, error) {
	// Validar que el evento existe y está activo
	evento, err := dao.ObtenerEventoPorID(req.EventoID)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}
	if evento.Estado != domain.EstadoActivo {
		return nil, errors.New("el evento no está disponible para compra")
	}

	// Validar que el sector pertenece al evento y tiene cupo
	sector, err := dao.ObtenerSectorPorID(req.SectorID)
	if err != nil {
		return nil, errors.New("sector no encontrado")
	}
	if sector.EventoID != req.EventoID {
		return nil, errors.New("el sector no pertenece a este evento")
	}
	if sector.CapacidadDisponible == 0 {
		return nil, errors.New("no hay lugares disponibles en este sector")
	}

	// Decrementar cupo del sector
	if err := dao.DecrementarCapacidad(req.SectorID); err != nil {
		return nil, errors.New("error al reservar el lugar")
	}

	// Crear la entrada con el precio histórico del sector
	entrada := &domain.Entrada{
		UsuarioID:    usuarioID,
		SectorID:     req.SectorID,
		EventoID:     req.EventoID,
		PrecioPagado: sector.Precio,
		Estado:       domain.EntradaActiva,
	}

	if err := dao.CrearEntrada(entrada); err != nil {
		// Rollback del cupo si falla la creación
		_ = dao.IncrementarCapacidad(req.SectorID)
		return nil, errors.New("error al registrar la compra")
	}

	return entrada, nil
}

func ObtenerMisEntradas(usuarioID uint) ([]domain.Entrada, error) {
	return dao.ObtenerEntradasPorUsuario(usuarioID)
}

func CancelarEntrada(entradaID, usuarioID uint) error {
	entrada, err := dao.ObtenerEntradaPorID(entradaID)
	if err != nil {
		return errors.New("entrada no encontrada")
	}
	if entrada.UsuarioID != usuarioID {
		return errors.New("no tenés permiso para cancelar esta entrada")
	}
	if entrada.Estado != domain.EntradaActiva {
		return errors.New("solo se pueden cancelar entradas activas")
	}

	if err := dao.CancelarEntrada(entradaID); err != nil {
		return errors.New("error al cancelar la entrada")
	}

	// Devolver el cupo al sector
	_ = dao.IncrementarCapacidad(entrada.SectorID)
	return nil
}

func TransferirEntrada(entradaID, usuarioOrigenID uint, req domain.TransferirEntradaRequest) error {
	entrada, err := dao.ObtenerEntradaPorID(entradaID)
	if err != nil {
		return errors.New("entrada no encontrada")
	}
	if entrada.UsuarioID != usuarioOrigenID {
		return errors.New("no tenés permiso para transferir esta entrada")
	}
	if entrada.Estado != domain.EntradaActiva {
		return errors.New("solo se pueden transferir entradas activas")
	}

	// Buscar usuario destino por email
	usuarioDestino, err := dao.ObtenerUsuarioPorEmail(req.EmailDestino)
	if err != nil {
		return errors.New("usuario destino no encontrado")
	}
	if usuarioDestino.ID == usuarioOrigenID {
		return errors.New("no podés transferirte la entrada a vos mismo")
	}

	// Transferir
	if err := dao.TransferirEntrada(entradaID, usuarioDestino.ID); err != nil {
		return errors.New("error al transferir la entrada")
	}

	// Registrar historial de transferencia
	_ = dao.RegistrarTransferencia(&domain.Transferencia{
		EntradaID:      entradaID,
		UsuarioOrigen:  usuarioOrigenID,
		UsuarioDestino: usuarioDestino.ID,
	})

	return nil
}
