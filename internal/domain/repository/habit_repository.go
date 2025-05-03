package repository

import "routinist/internal/domain/model"

type HabitRepository interface {
	GetRandomHabits() (*[]model.Habit, error)
}
