package repository

import (
	"gorm.io/gorm"
	"routinist/internal/dto/request"
)

type AuthRepository interface {
	Register(db *gorm.DB, e *request.RegisterRequestDTO) (*request.AuthResponseDTO, uint, error)
	Login(e *request.LoginRequestDTO) (*request.AuthResponseDTO, uint, error)
	GetDB() *gorm.DB
}
