package model

import "time"

type HabitProgress struct {
	ID          uint `gorm:"primaryKey"`
	UserHabitID uint
	Date        time.Time
	Value       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsCompleted bool
}
