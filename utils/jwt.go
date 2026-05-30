package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"ticketek/backend/domain"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UsuarioID uint               `json:"usuario_id"`
	Email     string             `json:"email"`
	Rol       domain.RolUsuario  `json:"rol"`
	jwt.RegisteredClaims
}

// GenerarToken crea un JWT firmado para el usuario dado.
func GenerarToken(usuario *domain.Usuario) (string, error) {
	horas, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		horas = 24
	}

	claims := Claims{
		UsuarioID: usuario.ID,
		Email:     usuario.Email,
		Rol:       usuario.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(horas) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ParsearToken valida y parsea un JWT, devuelve los claims.
func ParsearToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}
	return claims, nil
}
