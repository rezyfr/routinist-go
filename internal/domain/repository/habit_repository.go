package repository

import "routinist/internal/domain/model"

type HabitRepository interface {
	GetRandomHabits() (*[]model.Habit, error)
	GetTodayHabits(userId uint) ([]model.UserHabit, error)
	GetUserHabit(userId uint, userHabitId uint) (*model.UserHabit, error)
	CreateProgress(userHabitId uint, value float64) (string, error)
	UpdateProgress(progressId uint, value float64) (string, error)
	GetProgress(userHabitId uint) (float64, error)
}
