package usecase

import (
	"fmt"
	"math/rand"
	"routinist/internal/domain/model"
	"routinist/internal/domain/repository"
	"routinist/internal/dto/response"
	"routinist/internal/util"
	"routinist/pkg/logger"
	"time"
)

type HabitUsecase interface {
	CreateUserHabit(userId uint, habitId uint, unitId *uint, goal *float64) (string, error)
	GetRandomHabits() (*[]response.HabitDto, error)
	GetTodayHabitProgresses(userId uint) ([]response.UserHabitProgressDto, error)
	PostCreateHabitProgress(userId uint, userHabitId uint, value float64) (*response.CreateProgressDto, error)
	GetProgressSummary(userID uint, from, to time.Time) (*response.ProgressSummaryDto, error)
	GetActivitySummary(userID uint, userHabitId uint, from, to time.Time) (*response.ActivitySummaryDto, error)
	GetUserHabits(userId uint) ([]response.UserHabitDto, error)
	GetUserHabitDailyStats(userID uint, from, to time.Time) ([]response.DailyHabitStat, error)
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

	if err := uc.repo.EnsureTodayProgressForUser(uh.UserID); err != nil {
		uc.logger.Error(err)
		return "", fmt.Errorf("failed to prepare today's habit progress: %w", err)
	}

	return "success to create user habit", err
}
func (uc *habitUseCase) GetRandomHabits() (*[]response.HabitDto, error) {
	habits, err := uc.repo.GetRandomHabits()

	if habits == nil {
		return &[]response.HabitDto{}, nil
	}

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

func (uc *habitUseCase) GetTodayHabitProgresses(userId uint) ([]response.UserHabitProgressDto, error) {
	userHabits, err := uc.repo.GetTodayHabits(userId)

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get today habits: %w", err)
	}

	var uids []uint
	for _, uh := range userHabits {
		uids = append(uids, uh.ID)
	}

	progresses, err := uc.repo.GetTodayHabitProgresses(uids)
	if err != nil {
		uc.logger.Error("Failed to fetch progress records: ", err)
		return nil, err
	}

	progressMap := make(map[uint]model.HabitProgress)
	for _, p := range progresses {
		progressMap[p.UserHabitID] = p
	}

	var result []response.UserHabitProgressDto

	for _, u := range userHabits {
		progress := progressMap[u.ID]

		result = append(result, response.ToUserHabitProgressDto(&u, &progress))
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
	var habitIcon string
	var completedCount, failedCount int64

	hp, err := uc.repo.GetUserHabitProgresses(userID, userHabitId, from, to)

	if err != nil {
		uc.logger.Error(err)
		return nil, err
	}

	for _, log := range hp {
		if log.IsCompleted {
			completedCount++
		} else {
			failedCount++
		}
	}

	if userHabitId != 0 {
		uh, err := uc.repo.GetUserHabit(userID, userHabitId)

		if err != nil {
			uc.logger.Error(err)
			return nil, err
		}

		habitName = uh.Habit.Name
		habitId = uh.HabitID
		habitIcon = uh.Habit.Icon
	}

	percentage := 0.0
	completed := len(hp)
	if completed > 0 {
		percentage = float64(completedCount) / float64(completed) * 100
	}

	return &response.ActivitySummaryDto{
		SuccessRate:   util.RoundFloat(percentage, 2),
		Completed:     uint(completed),
		Failed:        uint(failedCount),
		UserHabitName: habitName,
		UserHabitId:   habitId,
		UserHabitIcon: habitIcon,
	}, nil
}

func (uc *habitUseCase) GetUserHabits(userId uint) ([]response.UserHabitDto, error) {
	uh, err := uc.repo.GetUserHabits(userId)

	if err != nil {
		uc.logger.Error(err)
		return nil, fmt.Errorf("failed to get user habit: %w", err)
	}

	var result []response.UserHabitDto

	for _, u := range uh {
		result = append(result, response.ToUserHabitDto(u))
	}

	return result, nil
}

func (uc *habitUseCase) GetUserHabitDailyStats(userID uint, from, to time.Time) ([]response.DailyHabitStat, error) {
	end := to
	hp, err := uc.repo.GetUserHabitProgresses(userID, 0, from, to)

	if err != nil {
		uc.logger.Error("Failed to get progress: ", err)
		return nil, err
	}

	var result []response.DailyHabitStat

	// For each day in the 7-day range (oldest → newest)
	for i := 6; i >= 0; i-- {
		r := response.DailyHabitStat{}
		day := end.AddDate(0, 0, -i)
		dayKey := day.Format("2006-01-02")

		for _, p := range hp {
			if p.Date.Format("2006-01-02") == dayKey {
				r.Total++

				if p.IsCompleted {
					r.Success++
				}
			}
		}

		result = append(result, r)
	}

	return result, nil
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
