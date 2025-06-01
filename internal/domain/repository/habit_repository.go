package repository

import (
	"routinist/internal/domain/model"
	"time"
)

type HabitRepository interface {
	CreateUserHabit(userId uint, habitId uint, unitId *uint, goal *float64) (*model.UserHabit, error)
	GetRandomHabits() (*[]model.Habit, error)
	GetTodayHabits(userId uint) ([]model.UserHabit, error)
	GetUserHabit(userId uint, userHabitId uint) (*model.UserHabit, error)
	GetUserHabits(userId uint) ([]*model.UserHabit, error)
	CreateProgress(userHabitId uint, value float64) (*model.HabitProgress, error)
	UpdateProgress(progressId uint, value float64) (*model.HabitProgress, error)
	GetProgress(userHabitId uint) (float64, error)
	GetProgressSummary(userHabitID uint, from, to time.Time) (completed int64, total int64, err error)
	EnsureTodayProgressForUser(userId uint) error
	GetTodayHabitProgress(userHabitId uint) (*model.HabitProgress, error)
	GetTodayHabitProgresses(userHabitId []uint) ([]model.HabitProgress, error)
	GetUserHabitProgresses(userId uint, userHabitId uint, from, to time.Time) ([]model.HabitProgress, error)
}
