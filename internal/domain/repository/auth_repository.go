package repository

import (
	"routinist/internal/domain/model"
)

type AuthRepository interface {
	Register(e *model.RegisterRequestDTO) (*model.AuthResponseDTO, error)
	Login(e *model.LoginRequestDTO) (*model.AuthResponseDTO, error)
}
