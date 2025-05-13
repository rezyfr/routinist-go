package repository

import (
	"routinist/internal/dto/request"
)

type AuthRepository interface {
	Register(e *request.RegisterRequestDTO) (*request.AuthResponseDTO, uint, error)
	Login(e *request.LoginRequestDTO) (*request.AuthResponseDTO, error)
}
