package usecase

import (
	"fmt"
	"routinist/internal/domain/model"
	"routinist/internal/domain/repository"
	"routinist/internal/dto"
	"routinist/pkg/logger"
)

type HabitUsecase interface {
	GetRandomHabits() (*[]model.Habit, error)
	GetTodayHabits(userId uint) ([]dto.UserHabitDto, error)
}

type habitUseCase struct {
	repo   repository.HabitRepository
	logger *logger.Logger
}

func NewHabitUseCase(r repository.HabitRepository, l *logger.Logger) HabitUsecase {
	return &habitUseCase{r, l}
}

func (uc *habitUseCase) GetRandomHabits() (*[]model.Habit, error) {
	habits, err := uc.repo.GetRandomHabits()
	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get random habits: %w", err)
	}

	return habits, nil
}

func (uc *habitUseCase) GetTodayHabits(userId uint) ([]dto.UserHabitDto, error) {
	uh, err := uc.repo.GetTodayHabits(userId)

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get today habits: %w", err)
	}

	var result []dto.UserHabitDto
	for _, u := range uh {
		result = append(result, dto.ToUserHabitDto(u))
	}
	return result, nil
}
