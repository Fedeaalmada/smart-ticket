package services

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"ticketek/backend/clients"
	"ticketek/backend/domain"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	clients.ConnectDB()

	code := m.Run()

	// Limpiar datos creados por los tests (respetando el orden de FK)
	clients.DB.Exec("DELETE FROM transferencias WHERE entrada_id IN (SELECT id FROM entradas WHERE evento_id > 3)")
	clients.DB.Exec("DELETE FROM transferencias WHERE entrada_id IN (SELECT id FROM entradas WHERE usuario_id IN (SELECT id FROM usuarios WHERE email LIKE '%@test.com'))")
	clients.DB.Exec("DELETE FROM entradas WHERE evento_id > 3")
	clients.DB.Exec("DELETE FROM entradas WHERE usuario_id IN (SELECT id FROM usuarios WHERE email LIKE '%@test.com')")
	clients.DB.Exec("DELETE FROM sectores WHERE evento_id > 3")
	clients.DB.Exec("DELETE FROM eventos WHERE id > 3")
	clients.DB.Exec("DELETE FROM usuarios WHERE email LIKE '%@test.com'")

	os.Exit(code)
}

func TestRegistrar_Exitoso(t *testing.T) {
	email := fmt.Sprintf("test_%d@test.com", time.Now().UnixNano())
	req := domain.RegisterRequest{
		Nombre:   "Test User",
		Email:    email,
		Password: "password123",
	}
	resp, err := Registrar(req)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("el token no debe estar vacío")
	}
}

func TestRegistrar_EmailDuplicado(t *testing.T) {
	email := fmt.Sprintf("dup_%d@test.com", time.Now().UnixNano())
	req := domain.RegisterRequest{
		Nombre:   "Test User",
		Email:    email,
		Password: "password123",
	}
	Registrar(req)
	_, err := Registrar(req)
	if err == nil {
		t.Fatal("se esperaba error con email duplicado")
	}
}

func TestLogin_Exitoso(t *testing.T) {
	email := fmt.Sprintf("login_%d@test.com", time.Now().UnixNano())
	Registrar(domain.RegisterRequest{
		Nombre:   "Login Test",
		Email:    email,
		Password: "password123",
	})
	resp, err := Login(domain.LoginRequest{
		Email:    email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("el token no debe estar vacío")
	}
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	_, err := Login(domain.LoginRequest{
		Email:    "noexiste@test.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("se esperaba error con credenciales inválidas")
	}
}

func TestLogin_PasswordIncorrecta(t *testing.T) {
	email := fmt.Sprintf("wrong_%d@test.com", time.Now().UnixNano())
	Registrar(domain.RegisterRequest{
		Nombre:   "Wrong Pass",
		Email:    email,
		Password: "correcta123",
	})
	_, err := Login(domain.LoginRequest{
		Email:    email,
		Password: "incorrecta123",
	})
	if err == nil {
		t.Fatal("se esperaba error con password incorrecta")
	}
}

func TestObtenerCatalogoEventos_SinFiltros(t *testing.T) {
	eventos, err := ObtenerCatalogoEventos(map[string]string{})
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if len(eventos) == 0 {
		t.Fatal("se esperaba al menos un evento")
	}
}

func TestObtenerDetalleEvento_Existente(t *testing.T) {
	evento, err := ObtenerDetalleEvento(1)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if evento.ID != 1 {
		t.Errorf("esperaba evento con ID 1, got: %d", evento.ID)
	}
}

func TestObtenerDetalleEvento_NoExiste(t *testing.T) {
	_, err := ObtenerDetalleEvento(99999)
	if err == nil {
		t.Fatal("se esperaba error con evento inexistente")
	}
}

func TestObtenerMisEntradas_UsuarioExistente(t *testing.T) {
	_, err := ObtenerMisEntradas(2)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
}

func TestCancelarEntrada_NoExiste(t *testing.T) {
	err := CancelarEntrada(99999, 1)
	if err == nil {
		t.Fatal("se esperaba error con entrada inexistente")
	}
}

func TestCancelarEntrada_NoEsTuEntrada(t *testing.T) {
	err := CancelarEntrada(1, 99999)
	if err == nil {
		t.Fatal("se esperaba error cuando el usuario no es dueño")
	}
}

func TestTransferirEntrada_NoExiste(t *testing.T) {
	err := TransferirEntrada(99999, 1, domain.TransferirEntradaRequest{
		EmailDestino: "juan@email.com",
	})
	if err == nil {
		t.Fatal("se esperaba error con entrada inexistente")
	}
}

func TestComprarEntrada_EventoNoExiste(t *testing.T) {
	_, err := ComprarEntrada(1, domain.ComprarEntradaRequest{
		EventoID: 99999,
		SectorID: 1,
	})
	if err == nil {
		t.Fatal("se esperaba error con evento inexistente")
	}
}

func TestCancelarEvento_NoExiste(t *testing.T) {
	err := CancelarEvento(99999)
	if err == nil {
		t.Fatal("se esperaba error con evento inexistente")
	}
}

func TestCrearEvento_Exitoso(t *testing.T) {
	req := domain.CrearEventoRequest{
		Titulo:    "Evento Test Auto",
		Fecha:     "2027-01-01",
		Horario:   "20:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 10, Precio: 1000},
			{Nombre: "VIP", CapacidadMaxima: 5, Precio: 5000},
		},
	}
	evento, err := CrearEvento(req, 1)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if evento.ID == 0 {
		t.Fatal("el evento debe tener un ID")
	}
	if len(evento.Sectores) != 2 {
		t.Errorf("esperaba 2 sectores, got: %d", len(evento.Sectores))
	}
}

func TestActualizarEvento_ConCambios(t *testing.T) {
	req := domain.ActualizarEventoRequest{
		Titulo: "Titulo Actualizado Test",
	}
	err := ActualizarEvento(1, req)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
}

func TestActualizarEvento_ConSectores(t *testing.T) {
	req := domain.ActualizarEventoRequest{
		Sectores: []domain.ActualizarSectorRequest{
			{ID: 1, Precio: 16000},
		},
	}
	err := ActualizarEvento(1, req)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
}

func TestActualizarEvento_FechaInvalida(t *testing.T) {
	req := domain.ActualizarEventoRequest{
		Fecha: "fecha-invalida",
	}
	err := ActualizarEvento(1, req)
	if err == nil {
		t.Fatal("se esperaba error con fecha inválida")
	}
}

func TestObtenerReporteEvento_Existente(t *testing.T) {
	reporte, err := ObtenerReporteEvento(1)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if reporte == nil {
		t.Fatal("el reporte no debe ser nil")
	}
}

func TestComprarEntrada_Exitoso(t *testing.T) {
	evento, err := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Compra Test",
		Fecha:     "2027-06-01",
		Horario:   "21:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 5, Precio: 2000},
		},
	}, 1)
	if err != nil {
		t.Fatalf("no se pudo crear evento: %v", err)
	}
	req := domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	}
	entrada, err := ComprarEntrada(2, req)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if entrada.ID == 0 {
		t.Fatal("la entrada debe tener un ID")
	}
}

func TestComprarEntrada_EventoCancelado(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Cancelado Test",
		Fecha:     "2027-07-01",
		Horario:   "21:00",
		Categoria: "Teatro",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Platea", CapacidadMaxima: 5, Precio: 3000},
		},
	}, 1)
	CancelarEvento(evento.ID)
	_, err := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	if err == nil {
		t.Fatal("se esperaba error al comprar en evento cancelado")
	}
}

func TestComprarEntrada_SectorNoExiste(t *testing.T) {
	_, err := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: 1,
		SectorID: 99999,
	})
	if err == nil {
		t.Fatal("se esperaba error con sector inexistente")
	}
}

func TestComprarEntrada_SinCupo(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Sin Cupo",
		Fecha:     "2027-08-01",
		Horario:   "21:00",
		Categoria: "Feria",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "General", CapacidadMaxima: 1, Precio: 500},
		},
	}, 1)
	sectorID := evento.Sectores[0].ID
	ComprarEntrada(2, domain.ComprarEntradaRequest{EventoID: evento.ID, SectorID: sectorID})
	_, err := ComprarEntrada(3, domain.ComprarEntradaRequest{EventoID: evento.ID, SectorID: sectorID})
	if err == nil {
		t.Fatal("se esperaba error por falta de cupo")
	}
}

func TestCancelarEntrada_Exitosa(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Cancelar Entrada",
		Fecha:     "2027-09-01",
		Horario:   "20:00",
		Categoria: "Deporte",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Popular", CapacidadMaxima: 5, Precio: 1500},
		},
	}, 1)
	entrada, _ := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	err := CancelarEntrada(entrada.ID, 2)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
}

func TestCancelarEntrada_EntradaNoActiva(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento No Activa",
		Fecha:     "2027-10-01",
		Horario:   "20:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 5, Precio: 1000},
		},
	}, 1)
	entrada, _ := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	CancelarEntrada(entrada.ID, 2)
	err := CancelarEntrada(entrada.ID, 2)
	if err == nil {
		t.Fatal("se esperaba error al cancelar entrada ya cancelada")
	}
}

func TestTransferirEntrada_UsuarioDestinoNoExiste(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Transferir",
		Fecha:     "2027-11-01",
		Horario:   "20:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 5, Precio: 1000},
		},
	}, 1)
	entrada, _ := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	err := TransferirEntrada(entrada.ID, 2, domain.TransferirEntradaRequest{
		EmailDestino: "noexiste@noexiste.com",
	})
	if err == nil {
		t.Fatal("se esperaba error con usuario destino inexistente")
	}
}

func TestTransferirEntrada_AMismoUsuario(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Mismo Usuario",
		Fecha:     "2027-12-01",
		Horario:   "20:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 5, Precio: 1000},
		},
	}, 1)
	entrada, _ := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	err := TransferirEntrada(entrada.ID, 2, domain.TransferirEntradaRequest{
		EmailDestino: "juan@email.com",
	})
	if err == nil {
		t.Fatal("se esperaba error al transferir al mismo usuario")
	}
}

func TestTransferirEntrada_Exitosa(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Transferir Exitoso",
		Fecha:     "2028-01-01",
		Horario:   "20:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 5, Precio: 1000},
		},
	}, 1)
	entrada, _ := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	err := TransferirEntrada(entrada.ID, 2, domain.TransferirEntradaRequest{
		EmailDestino: "maria@email.com",
	})
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
}

func TestTransferirEntrada_NoEsTuya(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento No Tuya",
		Fecha:     "2028-02-01",
		Horario:   "20:00",
		Categoria: "Concierto",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Campo", CapacidadMaxima: 5, Precio: 1000},
		},
	}, 1)
	entrada, _ := ComprarEntrada(2, domain.ComprarEntradaRequest{
		EventoID: evento.ID,
		SectorID: evento.Sectores[0].ID,
	})
	err := TransferirEntrada(entrada.ID, 3, domain.TransferirEntradaRequest{
		EmailDestino: "maria@email.com",
	})
	if err == nil {
		t.Fatal("se esperaba error cuando el usuario no es dueño")
	}
}

func TestCancelarEvento_Exitoso(t *testing.T) {
	evento, _ := CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento A Cancelar",
		Fecha:     "2028-03-01",
		Horario:   "20:00",
		Categoria: "Teatro",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Platea", CapacidadMaxima: 5, Precio: 2000},
		},
	}, 1)
	err := CancelarEvento(evento.ID)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
}