package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ticketek/backend/domain"
	"ticketek/backend/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret_key_para_tests")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")
}

func generarTokenTest(rol domain.RolUsuario) string {
	u := &domain.Usuario{ID: 1, Email: "test@test.com", Rol: rol}
	token, _ := utils.GenerarToken(u)
	return token
}

func TestAuthRequerido_SinToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestAuthRequerido_TokenInvalido(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token.invalido")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestAuthRequerido_TokenValido(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	token := generarTokenTest(domain.RolCliente)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestAuthRequerido_SinPrefixBearer(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	token := generarTokenTest(domain.RolCliente)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestSoloAdmin_ConRolAdmin(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), SoloAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	token := generarTokenTest(domain.RolAdministrador)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestSoloAdmin_ConRolCliente(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), SoloAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	token := generarTokenTest(domain.RolCliente)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestSoloAdmin_SinToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), SoloAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestObtenerUsuarioID_Exitoso(t *testing.T) {
	r := gin.New()
	r.GET("/test", AuthRequerido(), func(c *gin.Context) {
		id := ObtenerUsuarioID(c)
		if id == 0 {
			t.Error("el ID no debe ser 0")
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	token := generarTokenTest(domain.RolCliente)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}
