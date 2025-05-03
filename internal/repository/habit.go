package repository

import (
	"routinist/internal/domain/model"
	"routinist/pkg/logger"

	"gorm.io/gorm"
)

type HabitRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewHabitRepo(db *gorm.DB, logger *logger.Logger) *HabitRepo {
	return &HabitRepo{db, logger}
}

func (r *HabitRepo) GetRandomHabits() (*[]model.Habit, error) {
	var habits []model.Habit
	if err := r.db.Preload("Units").
		Order("RANDOM()").
		Limit(8).
		Find(&habits).Error; err != nil {
		return nil, err
	}

	return &habits, nil
}
