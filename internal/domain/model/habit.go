package model

import "time"

type Habit struct {
	ID          uint        `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Name        string      `gorm:"not null" json:"name"`
	Icon        string      `gorm:"not null" json:"icon"`
	Measurement Measurement `gorm:"type:varchar(20);not null" json:"measurement"`
	Units       []Unit      `gorm:"many2many:habit_units;" json:"units"`
	DefaultGoal float64     `gorm:"not null" json:"default_goal"`
}
