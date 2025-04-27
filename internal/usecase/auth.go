package usecase

import (
	"context"
	"routinist/internal/entity"
	"routinist/pkg/logger"
)

type AuthUseCase struct {
	repo   AuthRepo
	logger *logger.Logger
}

func NewAuthUseCase(r AuthRepo, l *logger.Logger) *AuthUseCase {
	return &AuthUseCase{
		repo:   r,
		logger: l,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, request entity.RegisterRequestDTO) (entity.AuthResponseDTO, error) {

	token, err := uc.repo.Register(ctx, request)
	if err != nil {
		uc.logger.Error(err)
		return token, err
	}

	return token, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, request entity.LoginRequestDTO) (entity.AuthResponseDTO, error) {
	token, err := uc.repo.Login(ctx, request)
	if err != nil {
		uc.logger.Error(err)
		return token, err
	}

	return token, nil
}
