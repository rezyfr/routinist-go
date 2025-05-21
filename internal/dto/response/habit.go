package response

import "routinist/internal/domain/model"

type HabitDto struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	Measurement model.Measurement `json:"measurement"`
	Units       []UnitDto         `json:"units"`
	DefaultGoal float64           `json:"default_goal"`
	Color       float64           `json:"color"`
}

func ToHabitDto(h model.Habit, color float64) HabitDto {
	units := make([]UnitDto, len(h.Units))
	for i, u := range h.Units {
		units[i] = toUnitDto(u)
	}
	return HabitDto{
		ID:          h.ID,
		Name:        h.Name,
		Icon:        h.Icon,
		Measurement: h.Measurement,
		Units:       units,
		DefaultGoal: h.DefaultGoal,
		Color:       color,
	}
}
