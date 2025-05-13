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
	repo      repository.AuthRepository
	habitRepo repository.HabitRepository
	logger    *logger.Logger
}

func NewAuthUseCase(r repository.AuthRepository, habitRepo repository.HabitRepository, l *logger.Logger) AuthUseCase {
	return &authUseCase{
		repo:      r,
		habitRepo: habitRepo,
		logger:    l,
	}
}

func (uc *authUseCase) Register(request *request.RegisterRequestDTO) (*request.AuthResponseDTO, error) {

	result, userId, err := uc.repo.Register(request)
	if err != nil {
		uc.logger.Error(err)
		return result, fmt.Errorf("failed to register: %w", err)
	}

	_, err = uc.habitRepo.CreateUserHabit(userId, request.HabitID, nil, nil)

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return result, nil
}

func (uc *authUseCase) Login(request *request.LoginRequestDTO) (*request.AuthResponseDTO, error) {
	token, err := uc.repo.Login(request)
	if err != nil {
		uc.logger.Error(err)
		return token, fmt.Errorf("failed to login: %w", err)
	}

	err = uc.habitRepo.EnsureTodayProgressForUser(request.Email)
	if err != nil {
		uc.logger.Error(err)
		return nil, err
	}

	return token, nil
}
