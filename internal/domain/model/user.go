package model

import (
	"time"
)

type User struct {
	ID         uint        `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Email      string      `gorm:"unique;not null" json:"email"`
	Password   string      `gorm:"not null" json:"-"`
	Name       string      `gorm:"not null" json:"name"`
	Gender     string      `gorm:"not null" json:"gender"`
	UserHabits []UserHabit `gorm:"foreignKey:UserID"`
	Milestone  uint        `json:"milestone" gorm:"default:0;not null"`
}

type Gender string

const (
	GenderMale        Gender = "male"
	GenderFemale      Gender = "female"
	GenderUnspecified Gender = "unspecified"
)
