package repository

import (
	"routinist/internal/domain/model"
	"routinist/pkg/logger"

	"gorm.io/gorm"
)

type habitRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewHabitRepo(db *gorm.DB, logger *logger.Logger) *habitRepo {
	return &habitRepo{db, logger}
}

func (r *habitRepo) GetRandomHabits() (*[]model.Habit, error) {
	var habits []model.Habit
	r.db.Order("RANDOM()").Limit(8).Find(&habits)
	return &habits, nil
}
