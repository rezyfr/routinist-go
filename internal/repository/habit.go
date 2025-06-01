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

func (r *HabitRepo) CreateUserHabit(userId uint, habitId uint, unitId *uint, goal *float64) (*model.UserHabit, error) {
	var habit model.Habit
	var unit model.Unit
	var userHabit model.UserHabit

	if err := r.db.Preload("Units").Where("id = ?", habitId).First(&habit).Error; err != nil {
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

	result := r.db.Create(&userHabit)
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
	exists := r.db.
		Where("user_habit_id = ?", userHabitId).
		Where("date = ?", time.Now().Truncate(24*time.Hour)).
		Find(&p).RowsAffected > 0

	if exists {
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
	ph.Value = ph.Value + value

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
	// Total habit progress in range for user
	var debugTotal int64
	r.db.
		Model(&model.HabitProgress{}).
		Joins("JOIN user_habits ON habit_progresses.user_habit_id = user_habits.id").
		Where("user_habits.user_id = ? AND habit_progresses.date >= ? AND habit_progresses.date < ?", userId, from, to).Count(&debugTotal)

	r.logger.Info("Total habit progress in range for user: ", debugTotal)
	err = r.db.
		Model(&model.HabitProgress{}).
		Joins("JOIN user_habits ON habit_progresses.user_habit_id = user_habits.id").
		Where("user_habits.user_id = ? AND habit_progresses.date >= ? AND habit_progresses.date < ?", userId, from, to).
		Count(&total).Error

	if err != nil {
		r.logger.Error("failed to get habit progress summary", err)
		return
	}
	// Completed ones
	err = r.db.
		Model(&model.HabitProgress{}).
		Joins("JOIN user_habits ON habit_progresses.user_habit_id = user_habits.id").
		Where("user_habits.user_id = ? AND habit_progresses.is_completed = ? AND habit_progresses.date >= ? AND habit_progresses.date < ?", userId, true, from, to).
		Count(&completed).Error
	if err != nil {
		r.logger.Error("failed to get habit progress summary", err)
	}
	return
}

func (r *HabitRepo) EnsureTodayProgressForUser(userId uint) error {
	today := time.Now().Truncate(24 * time.Hour)
	var user model.User
	var userHabits []model.UserHabit

	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		r.logger.Error("failed to get user", err)
		return err
	}

	if err := r.db.Where("user_id = ?", user.ID).Find(&userHabits).Error; err != nil {
		r.logger.Error("failed to get user habits", err)
		return err
	}
	for _, uh := range userHabits {
		for i := 0; i < 7; i++ { // For multiple days if needed
			d := today.AddDate(0, 0, i)
			// Ensure d is always truncated to midnight!
			d = d.Truncate(24 * time.Hour)
			hp := model.HabitProgress{
				UserHabitID: uh.ID,
				Date:        d,
				Value:       0,
				IsCompleted: false,
			}
			err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&hp).Error
			if err != nil {
				r.logger.Error("failed to ensure progress for user habit", err)
				return err
			}
		}
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

func (r *HabitRepo) GetActivitySummary(userId uint, userHabitId uint, from, to time.Time) (completed int64, total int64, failed int64, err error) {
	var progresses []model.HabitProgress

	q := r.db.Where("date BETWEEN ? AND ?", from, to)

	if userHabitId == 0 {
		// Get all user_habit IDs for the user
		var pIds []uint
		r.db.Model(&model.UserHabit{}).Where("user_id = ?", userId).Pluck("id", &pIds)
		if len(pIds) == 0 {
			return 0, 0, 0, nil // no habits
		}
		q = q.Where("user_habit_id IN ?", pIds)
	} else {
		q = q.Where("user_habit_id = ?", userHabitId)
	}

	if err = q.Find(&progresses).Error; err != nil {
		return 0, 0, 0, err
	}

	var completedCount, failedCount int64
	for _, log := range progresses {
		if log.IsCompleted {
			completedCount++
		} else {
			failedCount++
		}
	}

	return completedCount, int64(len(progresses)), failedCount, nil
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
