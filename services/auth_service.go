package services

import (
	"errors"

	"ticketek/backend/dao"
	"ticketek/backend/domain"
	"ticketek/backend/utils"
)

func Registrar(req domain.RegisterRequest) (*domain.LoginResponse, error) {
	// Verificar si el email ya existe
	_, err := dao.ObtenerUsuarioPorEmail(req.Email)
	if err == nil {
		return nil, errors.New("el email ya está registrado")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	usuario := &domain.Usuario{
		Nombre:       req.Nombre,
		Email:        req.Email,
		PasswordHash: hash,
		Rol:          domain.RolCliente,
	}

	if err := dao.CrearUsuario(usuario); err != nil {
		return nil, errors.New("error al crear el usuario")
	}

	token, err := utils.GenerarToken(usuario)
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	return &domain.LoginResponse{
		Token: token,
		Usuario: domain.UsuarioDTO{
			ID:     usuario.ID,
			Nombre: usuario.Nombre,
			Email:  usuario.Email,
			Rol:    usuario.Rol,
		},
	}, nil
}

func Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	usuario, err := dao.ObtenerUsuarioPorEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if !utils.CheckPassword(req.Password, usuario.PasswordHash) {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := utils.GenerarToken(usuario)
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	return &domain.LoginResponse{
		Token: token,
		Usuario: domain.UsuarioDTO{
			ID:     usuario.ID,
			Nombre: usuario.Nombre,
			Email:  usuario.Email,
			Rol:    usuario.Rol,
		},
	}, nil
}
