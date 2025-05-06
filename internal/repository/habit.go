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

func (r *HabitRepo) GetTodayHabits(userId uint) ([]model.UserHabit, error) {
	var userHabits []model.UserHabit
	err := r.db.Preload("Habit").
		Preload("Unit").
		Where("user_id = ?", userId).
		Where("goal_frequency = ?", model.FrequencyDaily).
		Find(&userHabits).Error

	if err != nil {
		r.logger.Error("failed to get today habits", err)
		return nil, err
	}

	habits := make([]model.UserHabit, 0, len(userHabits))
	for _, userHabit := range userHabits {
		habits = append(habits, userHabit)
	}

	return habits, err
}
