package request

type CreateUserHabitRequestDTO struct {
	UnitId  uint    `json:"unit_id"`
	HabitId uint    `json:"habit_id"`
	Goal    float64 `json:"goal"`
}
