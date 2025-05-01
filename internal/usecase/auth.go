package usecase

import (
	"fmt"
	"routinist/internal/domain/model"
	"routinist/internal/domain/repository"
	"routinist/pkg/logger"
)

type AuthUseCase interface {
	Login(request *model.LoginRequestDTO) (*model.AuthResponseDTO, error)
	Register(request *model.RegisterRequestDTO) (*model.AuthResponseDTO, error)
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

func (uc *authUseCase) Register(request *model.RegisterRequestDTO) (*model.AuthResponseDTO, error) {

	token, err := uc.repo.Register(request)
	if err != nil {
		uc.logger.Error(err)
		return token, fmt.Errorf("failed to register: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Login(request *model.LoginRequestDTO) (*model.AuthResponseDTO, error) {
	token, err := uc.repo.Login(request)
	if err != nil {
		uc.logger.Error(err)
		return token, fmt.Errorf("failed to login: %w", err)
	}

	return token, nil
}
