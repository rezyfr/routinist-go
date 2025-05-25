package request

import "time"

type GetActivitySummaryRequest struct {
	UserHabitId uint      `json:"user_habit_id"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
}
