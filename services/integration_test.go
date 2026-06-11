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
	os.Exit(m.Run())
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