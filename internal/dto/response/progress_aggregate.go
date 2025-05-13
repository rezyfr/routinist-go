package response

type ProgressSummaryDto struct {
	CompletedHabit float64 `json:"completed_habit"`
	TotalHabit     float64 `json:"total_habit"`
	Percentage     float64 `json:"percentage"`
}
