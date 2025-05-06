package dto

import "routinist/internal/domain/model"

type UserHabitDto struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	Icon          string              `json:"icon"`
	Goal          float64             `json:"goal"`
	GoalFrequency model.GoalFrequency `json:"goal_frequency"`
	Unit          UnitDto             `json:"unit"`
	UpdatedAt     string              `json:"updated_at"`
}

func ToUserHabitDto(uh model.UserHabit) UserHabitDto {
	return UserHabitDto{
		ID:            uh.ID,
		Name:          uh.Habit.Name,
		Icon:          uh.Habit.Icon,
		Goal:          uh.Goal,
		GoalFrequency: uh.GoalFrequency,
		Unit:          toUnitDto(uh.Unit),
		UpdatedAt:     uh.UpdatedAt.String(),
	}
}
