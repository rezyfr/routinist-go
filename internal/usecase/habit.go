package usecase

import (
	"fmt"
	"routinist/internal/domain/model"
	"routinist/internal/domain/repository"
	"routinist/internal/dto/response"
	"routinist/pkg/logger"
)

type HabitUsecase interface {
	GetRandomHabits() (*[]model.Habit, error)
	GetTodayHabits(userId uint) ([]response.UserHabitDto, error)
	PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (string, error)
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

func (uc *habitUseCase) GetTodayHabits(userId uint) ([]response.UserHabitDto, error) {
	uh, err := uc.repo.GetTodayHabits(userId)

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get today habits: %w", err)
	}

	var result []response.UserHabitDto
	for _, u := range uh {
		p, err := uc.repo.GetProgress(u.ID)

		if err != nil {
			uc.logger.Error(err)
			return nil, fmt.Errorf("failed to get habit progress: %w", err)
		}
		result = append(result, response.ToUserHabitDto(u, p))
	}
	return result, nil
}

func (uc *habitUseCase) PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (string, error) {
	uh, err := uc.repo.GetUserHabit(userHabitId, userId)

	if err != nil {
		uc.logger.Error(err)
		return "", fmt.Errorf("failed to get habit: %w", err)
	}

	ph, e := uc.repo.CreateProgress(uh.ID, value)

	return ph, e
}
