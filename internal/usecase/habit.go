package usecase

import (
	"fmt"
	"routinist/internal/domain/model"
	"routinist/internal/domain/repository"
	"routinist/internal/dto/response"
	"routinist/pkg/logger"
	"time"
)

type HabitUsecase interface {
	GetRandomHabits() (*[]model.Habit, error)
	GetTodayHabits(userId uint) ([]response.UserHabitDto, error)
	PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (string, error)
	GetProgressSummary(userID uint, from, to time.Time) (*response.ProgressSummaryDto, error)
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
		p, err := uc.repo.GetTodayHabitProgress(u.ID)

		if err != nil {
			uc.logger.Error(err)
			return nil, fmt.Errorf("failed to get habit progress: %w", err)
		}
		result = append(result, response.ToUserHabitDto(u, p))
	}
	return result, nil
}

func (uc *habitUseCase) PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (string, error) {
	uh, err := uc.repo.GetUserHabit(userId, userHabitId)

	if err != nil {
		uc.logger.Error(err)
		return "", fmt.Errorf("failed to get habit: %w", err)
	}

	ph, e := uc.repo.CreateProgress(uh.ID, value)

	return ph, e
}

func (uc *habitUseCase) GetProgressSummary(userID uint, from, to time.Time) (*response.ProgressSummaryDto, error) {
	completed, total, err := uc.repo.GetProgressSummary(userID, from, to)
	if err != nil {
		uc.logger.Error(err)
		return nil, err
	}
	percentage := 0.0
	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	return &response.ProgressSummaryDto{
		CompletedHabit: float64(completed),
		TotalHabit:     float64(total),
		Percentage:     percentage,
	}, nil
}
