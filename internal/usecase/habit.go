package usecase

import (
	"fmt"
	"routinist/internal/domain/model"
	"routinist/internal/domain/repository"
	"routinist/pkg/logger"
)

type HabitUsecase interface {
	GetRandomHabits() (*[]model.Habit, error)
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
