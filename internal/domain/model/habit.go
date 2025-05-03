package model

import (
	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model
	Name        string      `gorm:"not null" json:"name"`
	Icon        string      `gorm:"not null" json:"icon"`
	Measurement Measurement `gorm:"type:varchar(20);not null" json:"measurement"`
	Units       []Unit      `gorm:"many2many:habit_units;" json:"units"`
	DefaultGoal float64     `gorm:"not null" json:"default_goal"`
}
