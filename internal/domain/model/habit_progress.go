package model

import "time"

type HabitProgress struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UserHabitID uint      `gorm:"index:idx_userhabit_date,unique"`
	Date        time.Time `gorm:"index:idx_userhabit_date,unique"`
	Value       float64
	IsCompleted bool
}
