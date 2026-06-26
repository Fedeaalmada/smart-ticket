package utils

import (
	"os"
	"testing"

	"ticketek/backend/domain"
)

func TestGenerarToken_Admin(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_test_123")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")
	u := &domain.Usuario{ID: 2, Email: "admin@test.com", Rol: domain.RolAdministrador}
	token, err := GenerarToken(u)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if token == "" {
		t.Fatal("el token no debe estar vacío")
	}
}

func TestGenerarToken_SinExpiracion(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_test_123")
	os.Setenv("JWT_EXPIRATION_HOURS", "invalido")
	u := &domain.Usuario{ID: 1, Email: "test@test.com", Rol: domain.RolCliente}
	token, err := GenerarToken(u)
	if err != nil {
		t.Fatalf("no se esperaba error con horas invalidas, got: %v", err)
	}
	if token == "" {
		t.Fatal("el token no debe estar vacío")
	}
	os.Setenv("JWT_EXPIRATION_HOURS", "24")
}

func TestParsearToken_TokenVacio(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_test_123")
	_, err := ParsearToken("")
	if err == nil {
		t.Fatal("se esperaba error con token vacío")
	}
}
