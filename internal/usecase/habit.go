package usecase

import (
	"fmt"
	"math/rand"
	"routinist/internal/domain/repository"
	"routinist/internal/dto/response"
	"routinist/pkg/logger"
	"time"
)

type HabitUsecase interface {
	CreateUserHabit(userId uint, habitId uint, unitId *uint, goal *float64) (string, error)
	GetRandomHabits() (*[]response.HabitDto, error)
	GetTodayHabits(userId uint) ([]response.UserHabitDto, error)
	PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (*response.CreateProgressDto, error)
	GetProgressSummary(userID uint, from, to time.Time) (*response.ProgressSummaryDto, error)
	GetActivitySummary(userID uint, userHabitId uint, from, to time.Time) (*response.ActivitySummaryDto, error)
}

type habitUseCase struct {
	repo     repository.HabitRepository
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewHabitUseCase(r repository.HabitRepository, u repository.UserRepository, l *logger.Logger) HabitUsecase {
	return &habitUseCase{r, u, l}
}

func (uc *habitUseCase) CreateUserHabit(userId uint, habitId uint, unitId *uint, goal *float64) (string, error) {
	uh, err := uc.repo.CreateUserHabit(userId, habitId, unitId, goal)
	if err != nil {
		uc.logger.Error(err)
		return "", fmt.Errorf("failed to create habit: %w", err)
	}

	err = uc.repo.EnsureTodayProgressForUser(uh.UserID)

	if err != nil {
		uc.logger.Error(err)
		return "", fmt.Errorf("failed to create habit progress: %w", err)
	}

	return "success to create user habit", err
}
func (uc *habitUseCase) GetRandomHabits() (*[]response.HabitDto, error) {
	habits, err := uc.repo.GetRandomHabits()
	var result []response.HabitDto

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get random habits: %w", err)
	}

	for _, h := range *habits {
		result = append(result, response.ToHabitDto(h, generateRandomColor()))
	}

	return &result, nil
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
		result = append(result, response.ToUserHabitDto(&u, p))
	}
	return result, nil
}

func (uc *habitUseCase) PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (*response.CreateProgressDto, error) {
	uh, err := uc.repo.GetUserHabit(userId, userHabitId)

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	c, e := uc.repo.CreateProgress(uh.ID, value)

	if e != nil {
		uc.logger.Error(e)
		return nil, fmt.Errorf("failed to create habit progress: %w", e)
	}

	if c.IsCompleted {
		m, err := uc.userRepo.UpdateMilestone(userId, 1)
		if err != nil {
			return nil, err
		}

		r := response.CreateProgressDto{
			Milestone: m,
		}
		return &r, nil
	}

	return &response.CreateProgressDto{}, nil
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

func (uc *habitUseCase) GetActivitySummary(userID uint, userHabitId uint, from, to time.Time) (*response.ActivitySummaryDto, error) {
	var habitName string
	var habitId uint

	completed, total, failed, err := uc.repo.GetActivitySummary(userID, userHabitId, from, to)
	if err != nil {
		uc.logger.Error(err)
		return nil, err
	}

	if userHabitId != 0 {
		uh, err := uc.repo.GetUserHabit(userID, userHabitId)

		if err != nil {
			uc.logger.Error(err)
			return nil, err
		}

		habitName = uh.Habit.Name
		habitId = uh.HabitID
	}

	percentage := 0.0
	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	return &response.ActivitySummaryDto{
		SuccessRate:   percentage,
		Completed:     uint(completed),
		Failed:        uint(failed),
		UserHabitName: habitName,
		UserHabitId:   habitId,
	}, nil
}

func generateRandomColor() float64 {
	colors := []float64{
		0xFFFFFFFF, 0xFFFCDCD3, 0xFFD7D9FF, 0xFFBBE5FA, 0xFFF7CECD,
		0xFFFFE6B6, 0xFFC3EBC0, 0xFFE8D3FF, 0xFFD5ECE0,
	}

	rand.Seed(time.Now().UnixNano())
	color := colors[rand.Intn(len(colors))]

	return color
}
