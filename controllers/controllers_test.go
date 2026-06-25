package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"ticketek/backend/clients"
	"ticketek/backend/domain"
	"ticketek/backend/middleware"
	"ticketek/backend/services"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
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

func setupRouter() *gin.Engine {
	r := gin.New()
	r.POST("/auth/register", Register)
	r.POST("/auth/login", Login)
	r.GET("/eventos", GetEventos)
	r.GET("/eventos/:id", GetEventoByID)
	auth := r.Group("/")
	auth.Use(middleware.AuthRequerido())
	{
		auth.POST("/entradas/comprar", ComprarEntrada)
		auth.GET("/entradas/mis-entradas", GetMisEntradas)
		auth.DELETE("/entradas/:id", CancelarEntrada)
		auth.POST("/entradas/:id/transferir", TransferirEntrada)
	}
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequerido(), middleware.SoloAdmin())
	{
		admin.PUT("/eventos/:id", ActualizarEvento)
		admin.DELETE("/eventos/:id", CancelarEvento)
		admin.GET("/eventos/:id/reporte", GetReporteEvento)
	}
	return r
}

func obtenerTokenTest(t *testing.T) string {
	email := fmt.Sprintf("ctrl_%d@test.com", time.Now().UnixNano())
	services.Registrar(domain.RegisterRequest{
		Nombre:   "Ctrl Test",
		Email:    email,
		Password: "password123",
	})
	resp, err := services.Login(domain.LoginRequest{
		Email:    email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("no se pudo obtener token: %v", err)
	}
	return resp.Token
}

func TestRegisterController_Exitoso(t *testing.T) {
	r := setupRouter()
	email := fmt.Sprintf("reg_%d@test.com", time.Now().UnixNano())
	body := map[string]string{"nombre": "Test User", "email": email, "password": "password123"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("esperaba 201, got: %d", w.Code)
	}
}

func TestRegisterController_SinDatos(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestRegisterController_EmailDuplicado(t *testing.T) {
	r := setupRouter()
	email := fmt.Sprintf("dup_%d@test.com", time.Now().UnixNano())
	body := map[string]string{"nombre": "Test", "email": email, "password": "password123"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	req2, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400 en duplicado, got: %d", w2.Code)
	}
}

func TestLoginController_Exitoso(t *testing.T) {
	r := setupRouter()
	email := fmt.Sprintf("login_%d@test.com", time.Now().UnixNano())
	services.Registrar(domain.RegisterRequest{Nombre: "Login Test", Email: email, Password: "password123"})
	body := map[string]string{"email": email, "password": "password123"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestLoginController_Invalido(t *testing.T) {
	r := setupRouter()
	body := map[string]string{"email": "noexiste@test.com", "password": "wrong"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestLoginController_SinDatos(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestGetEventos_Exitoso(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/eventos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestGetEventoByID_Existente(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/eventos/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestGetEventoByID_NoExiste(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/eventos/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperaba 404, got: %d", w.Code)
	}
}

func TestGetEventoByID_IDInvalido(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/eventos/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestComprarEntrada_SinToken(t *testing.T) {
	r := setupRouter()
	body := map[string]uint{"evento_id": 1, "sector_id": 1}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/entradas/comprar", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestGetMisEntradas_ConToken(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("GET", "/entradas/mis-entradas", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestGetMisEntradas_SinToken(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/entradas/mis-entradas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestCancelarEntrada_SinToken(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("DELETE", "/entradas/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestCancelarEntrada_IDInvalido(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("DELETE", "/entradas/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestActualizarEvento_SinToken(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("PUT", "/admin/eventos/1", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestCancelarEvento_SinToken(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("DELETE", "/admin/eventos/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestGetReporteEvento_SinToken(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/admin/eventos/1/reporte", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestTransferirEntrada_SinToken(t *testing.T) {
	r := setupRouter()
	body := map[string]string{"email_destino": "test@test.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/entradas/1/transferir", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}

func TestTransferirEntrada_IDInvalido(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	body := map[string]string{"email_destino": "test@test.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/entradas/abc/transferir", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestActualizarEvento_IDInvalido(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("PUT", "/admin/eventos/abc", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestCancelarEvento_IDInvalido(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("DELETE", "/admin/eventos/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestGetReporteEvento_IDInvalido(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("GET", "/admin/eventos/abc/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestComprarEntrada_SinDatos(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("POST", "/entradas/comprar", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestComprarEntrada_ConToken(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	body := map[string]uint{"evento_id": 1, "sector_id": 1}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/entradas/comprar", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 201 o 400, got: %d", w.Code)
	}
}

func TestCancelarEntrada_ConToken(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("DELETE", "/entradas/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestTransferirEntrada_ConToken(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	body := map[string]string{"email_destino": "noexiste@test.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/entradas/99999/transferir", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestTransferirEntrada_SinDatos(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("POST", "/entradas/1/transferir", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func obtenerTokenAdmin(t *testing.T) string {
	email := fmt.Sprintf("admin_%d@test.com", time.Now().UnixNano())
	resp, err := services.Registrar(domain.RegisterRequest{
		Nombre:   "Admin Test",
		Email:    email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("no se pudo registrar: %v", err)
	}
	clients.DB.Exec("UPDATE usuarios SET rol = 'administrador' WHERE id = ?", resp.Usuario.ID)
	loginResp, err := services.Login(domain.LoginRequest{Email: email, Password: "password123"})
	if err != nil {
		t.Fatalf("no se pudo loguear: %v", err)
	}
	return loginResp.Token
}

func TestActualizarEvento_ConTokenAdmin(t *testing.T) {
	evento, err := services.CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Actualizar Test",
		Fecha:     "2029-06-01",
		Horario:   "20:00",
		Categoria: "Teatro",
		Sectores:  []domain.CrearSectorRequest{{Nombre: "Platea", CapacidadMaxima: 5, Precio: 1000}},
	}, 1)
	if err != nil {
		t.Fatalf("no se pudo crear evento de prueba: %v", err)
	}
	r := setupRouter()
	token := obtenerTokenAdmin(t)
	body := map[string]string{"titulo": "Titulo Nuevo"}
	jsonBody, _ := json.Marshal(body)
	url := fmt.Sprintf("/admin/eventos/%d", evento.ID)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestActualizarEvento_SinCambios(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenAdmin(t)
	req, _ := http.NewRequest("PUT", "/admin/eventos/1", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestGetReporteEvento_ConAdmin(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenAdmin(t)
	req, _ := http.NewRequest("GET", "/admin/eventos/1/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestGetReporteEvento_NoExiste(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenAdmin(t)
	req, _ := http.NewRequest("GET", "/admin/eventos/99999/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperaba 404, got: %d", w.Code)
	}
}

func TestCancelarEvento_ConAdmin(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenAdmin(t)
	eventoResp, _ := services.CrearEvento(domain.CrearEventoRequest{
		Titulo:    "Evento Controller Test",
		Fecha:     "2029-01-01",
		Horario:   "20:00",
		Categoria: "Teatro",
		Sectores: []domain.CrearSectorRequest{
			{Nombre: "Platea", CapacidadMaxima: 5, Precio: 2000},
		},
	}, 1)
	url := fmt.Sprintf("/admin/eventos/%d", eventoResp.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestCancelarEvento_NoExiste(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenAdmin(t)
	req, _ := http.NewRequest("DELETE", "/admin/eventos/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperaba 404, got: %d", w.Code)
	}
}

func TestGetEventos_ConFiltros(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/eventos?categoria=Concierto", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got: %d", w.Code)
	}
}

func TestActualizarEvento_TokenCliente(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	body := map[string]string{"titulo": "Titulo Nuevo"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/admin/eventos/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestCancelarEvento_TokenCliente(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("DELETE", "/admin/eventos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestGetReporteEvento_TokenCliente(t *testing.T) {
	r := setupRouter()
	token := obtenerTokenTest(t)
	req, _ := http.NewRequest("GET", "/admin/eventos/1/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperaba 403, got: %d", w.Code)
	}
}

func TestCrearEvento_ConAdmin(t *testing.T) {
	r := gin.New()
	r.POST("/admin/eventos", middleware.AuthRequerido(), middleware.SoloAdmin(), CrearEvento)
	token := obtenerTokenAdmin(t)
	body := map[string]interface{}{
		"titulo":    "Evento Nuevo Controller",
		"fecha":     "2029-06-01",
		"horario":   "21:00",
		"categoria": "Concierto",
		"sectores": []map[string]interface{}{
			{"nombre": "Campo", "capacidad_maxima": 100, "precio": 1500.0},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("esperaba 201, got: %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCrearEvento_SinDatos(t *testing.T) {
	r := gin.New()
	r.POST("/admin/eventos", middleware.AuthRequerido(), middleware.SoloAdmin(), CrearEvento)
	token := obtenerTokenAdmin(t)
	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperaba 400, got: %d", w.Code)
	}
}

func TestCrearEvento_SinToken(t *testing.T) {
	r := gin.New()
	r.POST("/admin/eventos", middleware.AuthRequerido(), middleware.SoloAdmin(), CrearEvento)
	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got: %d", w.Code)
	}
}