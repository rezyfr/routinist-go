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

func (r *HabitRepo) GetUserHabit(userId uint, userHabitId uint) (*model.UserHabit, error) {
	var habit model.UserHabit

	err := r.db.Preload("Habit").
		Where("id = ?", userHabitId).
		Where("user_id = ?", userId).
		First(&habit).Error

	if err != nil {
		r.logger.Error("failed to get habit", err)
		return nil, err
	}

	return &habit, nil
}

func (r *HabitRepo) CreateProgress(userHabitId uint, value float64) (string, error) {
	var uh model.UserHabit
	err := r.db.Preload("Habit").
		Where("id = ?", userHabitId).
		First(&uh).Error

	if err != nil {
		r.logger.Error("failed to get habit", err)
		return "", err
	}

	// If habit has progress, update it
	p := model.HabitProgress{}
	exists := r.db.Where("user_habit_id = ?", userHabitId).Find(&p).RowsAffected > 0

	if exists {
		return r.UpdateProgress(p.ID, value)
	}

	ph := model.HabitProgress{
		UserHabitID: uh.ID,
		Value:       value,
	}

	result := r.db.Create(&ph)

	if result.Error != nil {
		return "", result.Error
	}

	return "Progress created successfully", nil
}

func (r *HabitRepo) UpdateProgress(progressId uint, value float64) (string, error) {
	var uh model.UserHabit
	err := r.db.Preload("HabitProgress").Where("id = ?", progressId).First(&uh).Error

	if err != nil {
		r.logger.Error("failed to get habit progress", err)
		return "", err
	}

	ph := model.HabitProgress{}
	ph.Value = ph.Value + value
	ph.ID = progressId

	result := r.db.Save(&ph)

	if result.Error != nil {
		return "", result.Error
	}

	return "Progress updated successfully", nil
}

func (r *HabitRepo) GetProgress(userHabitId uint) (float64, error) {
	var ph model.HabitProgress
	err := r.db.Where("user_habit_id = ?", userHabitId).First(&ph).Error

	if err != nil {
		r.logger.Error("failed to get habit progress", err)
		return 0, err
	}

	return ph.Value, nil
}
