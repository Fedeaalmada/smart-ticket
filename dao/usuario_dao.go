package dao

import (
	"ticketek/backend/clients"
	"ticketek/backend/domain"
)

func CrearUsuario(u *domain.Usuario) error {
	return clients.DB.Create(u).Error
}

func ObtenerUsuarioPorEmail(email string) (*domain.Usuario, error) {
	var u domain.Usuario
	err := clients.DB.Where("email = ? AND activo = true", email).First(&u).Error
	return &u, err
}

func ObtenerUsuarioPorID(id uint) (*domain.Usuario, error) {
	var u domain.Usuario
	err := clients.DB.First(&u, id).Error
	return &u, err
}
