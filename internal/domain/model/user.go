package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string      `gorm:"unique;not null" json:"email"`
	Password   string      `gorm:"not null" json:"-"`
	Name       string      `gorm:"not null" json:"name"`
	Gender     string      `gorm:"not null" json:"gender"`
	UserHabits []UserHabit `gorm:"foreignKey:UserID"`
}

type Gender string

const (
	GenderMale        Gender = "male"
	GenderFemale      Gender = "female"
	GenderUnspecified Gender = "unspecified"
)
