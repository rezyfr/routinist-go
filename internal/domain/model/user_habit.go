package model

import (
	"gorm.io/gorm"
)

type UserHabit struct {
	gorm.Model

	UserID        uint          `gorm:"not null" json:"user_id"`
	HabitID       uint          `gorm:"not null" json:"habit_id"`
	UnitID        uint          `gorm:"not null" json:"unit_id"`
	Goal          float64       `gorm:"not null" json:"goal"`
	GoalFrequency GoalFrequency `gorm:"type:varchar(10);default:'daily'" json:"goal_frequency"`

	User  User  `gorm:"foreignKey:UserID"`
	Habit Habit `gorm:"foreignKey:HabitID"`
	Unit  Unit  `gorm:"foreignKey:UnitID"`
}

type GoalFrequency string

const (
	FrequencyDaily   GoalFrequency = "daily"
	FrequencyWeekly  GoalFrequency = "weekly"
	FrequencyMonthly GoalFrequency = "monthly"
)
