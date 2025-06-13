package response

import "routinist/internal/domain/model"

type UserHabitProgressDto struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	Icon          string              `json:"icon"`
	Goal          float64             `json:"goal"`
	GoalFrequency model.GoalFrequency `json:"goal_frequency"`
	Unit          UnitDto             `json:"unit"`
	CreatedAt     string              `json:"created_at"`
	Progress      float64             `json:"progress"`
}

func ToUserHabitProgressDto(uh *model.UserHabit, p *model.HabitProgress) UserHabitProgressDto {
	return UserHabitProgressDto{
		ID:            uh.ID,
		Name:          uh.Habit.Name,
		Icon:          uh.Habit.Icon,
		Goal:          uh.Goal,
		GoalFrequency: uh.GoalFrequency,
		Unit:          toUnitDto(uh.Unit),
		CreatedAt:     p.Date.String(),
		Progress:      p.Value,
	}
}
