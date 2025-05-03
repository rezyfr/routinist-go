package seed

import (
	"routinist/internal/domain/model"
	"routinist/pkg/logger"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB, l *logger.Logger) {
	seedUnits(db, l)
	seedHabits(db, l)
	seedHabitUnits(db, l)
}

func seedUnits(db *gorm.DB, l *logger.Logger) {
	var count int64
	db.Model(&model.Unit{}).Count(&count)
	if count > 0 {
		l.Info("Units already seeded")
		return
	}

	units := []model.Unit{
		{Name: "Minutes", Symbol: "min", Measurement: model.MeasurementTime},
		{Name: "Hours", Symbol: "h", Measurement: model.MeasurementTime},
		{Name: "Meters", Symbol: "m", Measurement: model.MeasurementDistance},
		{Name: "Kilometers", Symbol: "km", Measurement: model.MeasurementDistance},
		{Name: "Steps", Symbol: "steps", Measurement: model.MeasurementCount},
		{Name: "Kilograms", Symbol: "kg", Measurement: model.MeasurementWeight},
		{Name: "Grams", Symbol: "g", Measurement: model.MeasurementWeight},
		{Name: "Pounds", Symbol: "lb", Measurement: model.MeasurementWeight},
		{Name: "Ounces", Symbol: "oz", Measurement: model.MeasurementWeight},
		{Name: "Chapter", Symbol: "ch", Measurement: model.MeasurementCount},
		{Name: "Item", Symbol: "item", Measurement: model.MeasurementCount},
		{Name: "Episode", Symbol: "eps", Measurement: model.MeasurementCount},
		{Name: "Pages", Symbol: "page", Measurement: model.MeasurementCount},
		{Name: "Time", Symbol: "time", Measurement: model.MeasurementCount},
		{Name: "Volume", Symbol: "l", Measurement: model.MeasurementVolume},
	}

	if err := db.Create(&units).Error; err != nil {
		l.Fatal("failed to seed units: %v", err)
	}
	l.Info("Seeded units")
}

func seedHabits(db *gorm.DB, l *logger.Logger) {
	var count int64
	db.Model(&model.Habit{}).Count(&count)
	if count > 0 {
		l.Info("Habits already seeded")
		return
	}

	// Fetch all units
	var units []model.Unit
	db.Find(&units)

	// Map symbols to unit IDs
	unitMap := map[string]uint{}
	for _, u := range units {
		unitMap[u.Symbol] = u.ID
	}

	// Seed data
	habitSeed := []struct {
		Name        string
		Icon        string
		Measurement model.Measurement
		Units       []string
		DefaultGoal float64
	}{
		{Name: "Run", Icon: "🏃", Measurement: model.MeasurementDistance, Units: []string{"km", "m"}, DefaultGoal: 5},
		{Name: "Read Book", Icon: "📚", Measurement: model.MeasurementTime, Units: []string{"min", "h"}, DefaultGoal: 60},
		{Name: "Meditate", Icon: "🧘", Measurement: model.MeasurementTime, Units: []string{"min", "h"}, DefaultGoal: 60},
		{Name: "Study", Icon: "👨‍💻", Measurement: model.MeasurementTime, Units: []string{"min", "h"}, DefaultGoal: 60},
		{Name: "Journal", Icon: "📓", Measurement: model.MeasurementCount, Units: []string{"page"}, DefaultGoal: 3},
		{Name: "Water Plant", Icon: "🌿", Measurement: model.MeasurementCount, Units: []string{"time"}, DefaultGoal: 2},
		{Name: "Walk", Icon: "🚶", Measurement: model.MeasurementCount, Units: []string{"steps"}, DefaultGoal: 10000},
		{Name: "Drink Water", Icon: "💧", Measurement: model.MeasurementVolume, Units: []string{"l"}, DefaultGoal: 2},
	}

	for _, h := range habitSeed {
		habit := model.Habit{
			Name:        h.Name,
			Icon:        h.Icon,
			Measurement: h.Measurement,
			DefaultGoal: h.DefaultGoal,
		}

		if err := db.Create(&habit).Error; err != nil {
			l.Fatal("failed to seed habit: %v", err)
		}

		for _, symbol := range h.Units {
			if unitID, ok := unitMap[symbol]; ok {
				habitUnit := model.HabitUnit{
					HabitID:     habit.ID,
					UnitID:      unitID,
					DefaultGoal: h.DefaultGoal,
				}
				if err := db.Create(&habitUnit).Error; err != nil {
					l.Fatal("failed to seed habit_unit: %v", err)
				}
			}
		}
	}

	l.Info("Seeded habits")
}

func seedHabitUnits(db *gorm.DB, l *logger.Logger) {
	var count int64
	db.Model(&model.HabitUnit{}).Count(&count)
	if count > 0 {
		l.Info("HabitUnits already seeded")
		return
	}

	var habits []model.Habit
	if err := db.Preload("Units").Find(&habits).Error; err != nil {
		l.Fatal("failed to fetch habits for HabitUnit seeding: %v", err)
	}

	for _, habit := range habits {
		for _, unit := range habit.Units {
			habitUnit := model.HabitUnit{
				HabitID:     habit.ID,
				UnitID:      unit.ID,
				DefaultGoal: habit.DefaultGoal,
			}

			if err := db.Create(&habitUnit).Error; err != nil {
				l.Fatal("failed to seed HabitUnit for habit %d and unit %d: %v", habit.ID, unit.ID, err)
			}
		}
	}

	l.Info("Seeded HabitUnits")
}
