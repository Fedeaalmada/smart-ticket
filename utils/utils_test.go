package utils

import (
	"os"
	"testing"

	"ticketek/backend/domain"
)

func TestHashPassword_Exitoso(t *testing.T) {
	hash, err := HashPassword("miPassword123")
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if hash == "" {
		t.Fatal("el hash no debe estar vacío")
	}
	if hash == "miPassword123" {
		t.Fatal("el hash no debe ser igual a la contraseña original")
	}
}

func TestCheckPassword_Correcto(t *testing.T) {
	hash, _ := HashPassword("miPassword123")
	if !CheckPassword("miPassword123", hash) {
		t.Fatal("CheckPassword debería devolver true con la contraseña correcta")
	}
}

func TestCheckPassword_Incorrecto(t *testing.T) {
	hash, _ := HashPassword("miPassword123")
	if CheckPassword("otraPassword", hash) {
		t.Fatal("CheckPassword debería devolver false con una contraseña incorrecta")
	}
}

func TestCheckPassword_HashVacio(t *testing.T) {
	if CheckPassword("password", "") {
		t.Fatal("CheckPassword debería devolver false con hash vacío")
	}
}

func TestGenerarToken_Exitoso(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_test_123")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")

	usuario := &domain.Usuario{
		ID:    1,
		Email: "test@test.com",
		Rol:   domain.RolCliente,
	}

	token, err := GenerarToken(usuario)
	if err != nil {
		t.Fatalf("no se esperaba error, got: %v", err)
	}
	if token == "" {
		t.Fatal("el token no debe estar vacío")
	}
}

func TestParsearToken_Exitoso(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_test_123")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")

	usuario := &domain.Usuario{
		ID:    42,
		Email: "malena@test.com",
		Rol:   domain.RolAdministrador,
	}

	token, _ := GenerarToken(usuario)
	claims, err := ParsearToken(token)
	if err != nil {
		t.Fatalf("no se esperaba error al parsear, got: %v", err)
	}
	if claims.UsuarioID != 42 {
		t.Errorf("esperaba UsuarioID=42, got: %d", claims.UsuarioID)
	}
	if claims.Email != "malena@test.com" {
		t.Errorf("esperaba email malena@test.com, got: %s", claims.Email)
	}
	if claims.Rol != domain.RolAdministrador {
		t.Errorf("esperaba rol administrador, got: %s", claims.Rol)
	}
}

func TestParsearToken_TokenInvalido(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_test_123")

	_, err := ParsearToken("esto.no.es.un.token.valido")
	if err == nil {
		t.Fatal("se esperaba error con token inválido")
	}
}

func TestParsearToken_SecretoDistinto(t *testing.T) {
	os.Setenv("JWT_SECRET", "secreto_original")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")

	usuario := &domain.Usuario{ID: 1, Email: "a@a.com", Rol: domain.RolCliente}
	token, _ := GenerarToken(usuario)

	os.Setenv("JWT_SECRET", "secreto_diferente")
	_, err := ParsearToken(token)
	if err == nil {
		t.Fatal("se esperaba error al parsear con secreto distinto")
	}
}