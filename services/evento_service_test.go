package services

import (
	"testing"

	"ticketek/backend/domain"
)

func TestCrearEvento_FechaInvalida(t *testing.T) {
	req := domain.CrearEventoRequest{
		Titulo:    "Evento Test",
		Fecha:     "fecha-mal-escrita",
		Horario:   "21:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 100, Precio: 1000},
		},
	}
	_, err := CrearEvento(req, 1)
	if err == nil {
		t.Fatal("se esperaba error con fecha inválida")
	}
}

func TestActualizarEvento_SinCambios(t *testing.T) {
	req := domain.ActualizarEventoRequest{}
	err := ActualizarEvento(1, req)
	if err == nil {
		t.Fatal("se esperaba error cuando no se envían campos")
	}
	if err.Error() != "no se enviaron campos para actualizar" {
		t.Errorf("mensaje de error incorrecto: %s", err.Error())
	}
}