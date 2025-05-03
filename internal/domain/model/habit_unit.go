package model

type HabitUnit struct {
	HabitID     uint    `gorm:"primaryKey"`
	UnitID      uint    `gorm:"primaryKey"`
	DefaultGoal float64 `gorm:"not null"`

	Habit Habit `gorm:"foreignKey:HabitID"`
	Unit  Unit  `gorm:"foreignKey:UnitID"`
}

func (HabitUnit) TableName() string {
	return "habit_units" // important for GORM to match your existing join table
}
