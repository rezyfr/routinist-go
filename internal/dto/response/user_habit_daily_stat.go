package response

import "time"

type DailyHabitStat struct {
	Date    time.Time `json:"date"`
	Total   int       `json:"total"`
	Success int       `json:"success"`
}
