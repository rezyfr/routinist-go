package response

type ActivitySummaryDto struct {
	SuccessRate   float64 `json:"success_rate"`
	Completed     uint    `json:"completed"`
	Failed        uint    `json:"failed"`
	UserHabitId   uint    `json:"user_habit_id"`
	UserHabitName string  `json:"user_habit_name"`
}
