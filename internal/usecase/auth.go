package usecase

import (
	"fmt"
	"routinist/internal/domain/repository"
	"routinist/internal/dto/request"
	"routinist/pkg/logger"
)

type AuthUseCase interface {
	Login(request *request.LoginRequestDTO) (*request.AuthResponseDTO, error)
	Register(request *request.RegisterRequestDTO) (*request.AuthResponseDTO, error)
}

type authUseCase struct {
	repo   repository.AuthRepository
	logger *logger.Logger
}

func NewAuthUseCase(r repository.AuthRepository, l *logger.Logger) AuthUseCase {
	return &authUseCase{
		repo:   r,
		logger: l,
	}
}

func (uc *authUseCase) Register(request *request.RegisterRequestDTO) (*request.AuthResponseDTO, error) {

	token, err := uc.repo.Register(request)
	if err != nil {
		uc.logger.Error(err)
		return token, fmt.Errorf("failed to register: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Login(request *request.LoginRequestDTO) (*request.AuthResponseDTO, error) {
	token, err := uc.repo.Login(request)
	if err != nil {
		uc.logger.Error(err)
		return token, fmt.Errorf("failed to login: %w", err)
	}

	return token, nil
}
