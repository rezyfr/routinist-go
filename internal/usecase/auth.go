package usecase

import (
	"fmt"
	"gorm.io/gorm"
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

func (uc *authUseCase) Register(req *request.RegisterRequestDTO) (*request.AuthResponseDTO, error) {
	var result *request.AuthResponseDTO
	var userId uint

	db := uc.repo.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		var err error

		result, userId, err = uc.repo.Register(tx, req)
		if err != nil {
			uc.logger.Error(err)
			return fmt.Errorf("failed to register: %w", err)
		}

		_, err = uc.habitRepo.CreateUserHabit(tx, userId, req.HabitID, nil, nil)
		if err != nil {
			uc.logger.Error(err)
			return fmt.Errorf("failed to create habit: %w", err)
		}

		return nil
	})

	if err != nil {
		uc.logger.Error(err)
		return nil, err // Both will be rolled back on any error!
	}
	return result, nil
}

func (uc *authUseCase) Login(request *request.LoginRequestDTO) (*request.AuthResponseDTO, error) {
	token, userId, err := uc.repo.Login(request)
	if err != nil {
		uc.logger.Error(err)
		return token, fmt.Errorf("failed to login: %w", err)
	}

	err = uc.habitRepo.EnsureTodayProgressForUser(userId)
	if err != nil {
		uc.logger.Error(err)
		return nil, err
	}

	return token, nil
}
