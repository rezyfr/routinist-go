package seed

import (
	"routinist/internal/domain/model"
	"routinist/pkg/logger"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB, l *logger.Logger) {
	seedUnits(db, l)
	seedHabits(db, l)
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
		{Name: "Time", Symbol: "page", Measurement: model.MeasurementCount},
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

	// Find unit IDs
	var units []model.Unit
	db.Find(&units)

	// Map unit symbols to their IDs
	unitMap := map[string]uint{}
	for _, u := range units {
		unitMap[u.Symbol] = u.ID
	}

	habitSeed := []struct {
		Name        string
		Icon        string
		Measurement model.Measurement
		UnitSymbols []string
	}{
		{Name: "Run", Icon: "🏃", Measurement: model.MeasurementDistance, UnitSymbols: []string{"m", "km"}},
		{Name: "Read Book", Icon: "📚", Measurement: model.MeasurementTime, UnitSymbols: []string{"min", "h"}},
		{Name: "Meditate", Icon: "🧘", Measurement: model.MeasurementTime, UnitSymbols: []string{"min", "h"}},
		{Name: "Study", Icon: "👨‍💻", Measurement: model.MeasurementTime, UnitSymbols: []string{"min", "h"}},
		{Name: "Journal", Icon: "📓", Measurement: model.MeasurementCount, UnitSymbols: []string{"page"}},
		{Name: "Water Plant", Icon: "🌿", Measurement: model.MeasurementCount, UnitSymbols: []string{"time"}},
		{Name: "Walk", Icon: "🚶", Measurement: model.MeasurementDistance, UnitSymbols: []string{"m", "km"}},
		{Name: "Drink Water", Icon: "💧", Measurement: model.MeasurementVolume, UnitSymbols: []string{"time"}},
	}

	for _, h := range habitSeed {
		habit := model.Habit{
			Name:        h.Name,
			Icon:        h.Icon,
			Measurement: h.Measurement,
		}

		// Create habit first so we get the ID
		if err := db.Create(&habit).Error; err != nil {
			l.Fatal("failed to seed habit: %v", err)
		}

		// Attach units
		var unitsToAttach []model.Unit
		for _, symbol := range h.UnitSymbols {
			if id, ok := unitMap[symbol]; ok {
				unitsToAttach = append(unitsToAttach, model.Unit{ID: id})
			}
		}
		if err := db.Model(&habit).Association("Units").Append(unitsToAttach); err != nil {
			l.Fatal("failed to attach units to habit: %v", err)
		}
	}

	l.Info("Seeded habits")
}
