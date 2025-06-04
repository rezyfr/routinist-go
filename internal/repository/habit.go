package repository

import (
	"errors"
	"gorm.io/gorm/clause"
	"routinist/internal/domain/model"
	"routinist/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type HabitRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewHabitRepo(db *gorm.DB, logger *logger.Logger) *HabitRepo {
	return &HabitRepo{db, logger}
}

func (r *HabitRepo) CreateUserHabit(db *gorm.DB, userId uint, habitId uint, unitId *uint, goal *float64) (*model.UserHabit, error) {
	var habit model.Habit
	var unit model.Unit
	var userHabit model.UserHabit

	if err := db.Preload("Units").Where("id = ?", habitId).First(&habit).Error; err != nil {
		r.logger.Error("failed to get habit", err)
		return nil, err
	}

	if unitId != nil {
		for _, u := range habit.Units {
			if *unitId == u.ID {
				unit = u
				break
			}

		}
	} else {
		unit = habit.Units[0]
	}

	if goal == nil {
		goal = &habit.DefaultGoal
	}

	userHabit = model.UserHabit{
		UserID:  userId,
		HabitID: habitId,
		UnitID:  unit.ID,
		Goal:    *goal,
	}

	result := db.Create(&userHabit)
	if result.Error != nil {
		r.logger.Error("failed to create user habit", result.Error)
		return nil, result.Error
	}

	return &userHabit, nil
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

	return userHabits, err
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

func (r *HabitRepo) CreateProgress(userHabitId uint, value float64) (*model.HabitProgress, error) {
	var uh model.UserHabit
	err := r.db.Preload("Habit").
		Where("id = ?", userHabitId).
		First(&uh).Error

	if err != nil {
		r.logger.Error("failed to get habit", err)
		return nil, err
	}

	// If habit has progress, update it
	p := model.HabitProgress{}
	today := time.Now().Truncate(24 * time.Hour)
	err = r.db.Where("user_habit_id = ? AND date = ?", userHabitId, today).First(&p).Error

	if err == nil {
		r.logger.Info("habit progress already exists for today, updating it, id = ", p.ID, " value = ", value, "")
		ph, err := r.UpdateProgress(p.ID, value)
		if err != nil {
			r.logger.Error("failed to update habit progress", err)
			return nil, err
		}
		return ph, nil
	}

	ph := model.HabitProgress{
		UserHabitID: uh.ID,
		Value:       value,
		IsCompleted: value >= uh.Goal,
		Date:        time.Now(),
	}

	result := r.db.Create(&ph)

	if result.Error != nil {
		return nil, result.Error
	}

	return &ph, nil
}

func (r *HabitRepo) UpdateProgress(progressId uint, value float64) (*model.HabitProgress, error) {
	var ph model.HabitProgress
	err := r.db.Where("id = ?", progressId).First(&ph).Error

	if err != nil {
		r.logger.Error("failed to get habit progress", err)
		return nil, err
	}

	var uh model.UserHabit
	err = r.db.Where("id = ?", ph.UserHabitID).First(&uh).Error
	if err != nil {
		r.logger.Error("failed to get user habit", err)
		return nil, err
	}

	ph.IsCompleted = ph.Value+value >= uh.Goal
	ph.Value += value

	result := r.db.Save(&ph)

	if result.Error != nil {
		r.logger.Error("failed to update habit progress", result.Error)
		return nil, result.Error
	}

	return &ph, nil
}

func (r *HabitRepo) GetProgress(userHabitId uint) (float64, error) {
	var ph model.HabitProgress
	err := r.db.Where("user_habit_id = ?", userHabitId).First(&ph).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		r.logger.Error("failed to get habit progress", err)
		return 0, err
	}

	return ph.Value, nil
}

func (r *HabitRepo) GetProgressSummary(userId uint, from, to time.Time) (completed int64, total int64, err error) {
	type SummaryResult struct {
		Total     int64
		Completed int64
	}
	var result SummaryResult

	err = r.db.
		Model(&model.HabitProgress{}).
		Select(`
			COUNT(*) AS total,
			SUM(CASE WHEN is_completed THEN 1 ELSE 0 END) AS completed
		`).
		Joins("JOIN user_habits ON user_habits.id = habit_progresses.user_habit_id").
		Where("user_habits.user_id = ? AND habit_progresses.date >= ? AND habit_progresses.date < ?", userId, from, to).
		Scan(&result).Error

	if err != nil {
		r.logger.Error("failed to get habit progress summary", err)
		return
	}

	return result.Completed, result.Total, nil
}

func (r *HabitRepo) EnsureTodayProgressForUser(userId uint) error {
	today := time.Now().Truncate(24 * time.Hour)

	var userHabits []model.UserHabit
	if err := r.db.Where("user_id = ?", userId).Find(&userHabits).Error; err != nil {
		return err
	}

	if len(userHabits) == 0 {
		return nil
	}

	var records []model.HabitProgress
	for _, uh := range userHabits {
		records = append(records, model.HabitProgress{
			UserHabitID: uh.ID,
			Date:        today,
			Value:       0,
			IsCompleted: false,
		})
	}

	err := r.db.
		Model(&model.HabitProgress{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_habit_id"}, {Name: "date"}},
			DoNothing: true,
		}).
		Create(&records).Error

	if err != nil {
		r.logger.Error("failed to ensure progress for user habit", err)
		return err
	}

	return nil
}

func (r *HabitRepo) GetTodayHabitProgress(userHabitId uint) (*model.HabitProgress, error) {
	var habits model.HabitProgress
	err := r.db.Where("user_habit_id = ?", userHabitId).
		Where("date = ?", time.Now().Truncate(24*time.Hour)).
		Find(&habits).Error

	if err != nil {
		r.logger.Error("failed to get habit progress", err)
		return nil, err
	}

	return &habits, nil
}

func (r *HabitRepo) GetTodayHabitProgresses(userHabitId []uint) ([]model.HabitProgress, error) {
	var habits []model.HabitProgress
	err := r.db.Where("user_habit_id IN ?", userHabitId).
		Where("date = ?", time.Now().Truncate(24*time.Hour)).
		Find(&habits).Error

	if err != nil {
		r.logger.Error("failed to get habit progress", err)
		return nil, err
	}

	return habits, nil
}

func (r *HabitRepo) GetUserHabitProgresses(userId uint, userHabitId uint, from, to time.Time) ([]model.HabitProgress, error) {
	var progresses []model.HabitProgress
	var q *gorm.DB

	if userHabitId == 0 {
		// Get all user_habit IDs for the user
		var pIds []uint
		r.db.Model(&model.UserHabit{}).Where("user_id = ?", userId).Pluck("id", &pIds)
		q = r.db.Where("user_habit_id IN ?", pIds).Where("date BETWEEN ? AND ?", from, to)
	} else {
		q = r.db.Where("user_habit_id = ?", userHabitId).Where("date BETWEEN ? AND ?", from, to)
	}

	if err := q.Find(&progresses).Error; err != nil {
		return nil, err
	}

	return progresses, nil
}

func (r *HabitRepo) GetUserHabits(userId uint) ([]*model.UserHabit, error) {
	var userHabits []*model.UserHabit
	err := r.db.Preload("Habit").
		Preload("Unit").
		Where("user_id = ?", userId).
		Find(&userHabits).Error

	if err != nil {
		r.logger.Error("failed to get habit", err)
		return nil, err
	}

	return userHabits, nil
}

func (r *HabitRepo) GetDB() *gorm.DB {
	return r.db
}
