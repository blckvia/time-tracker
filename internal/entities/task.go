package entities

import "time"

type Task struct {
	ID          int       `json:"-"`
	UserID      int       `json:"user_id"`
	Task        string    `json:"task"`
	Description string    `json:"description"`
	Timer       bool      `json:"timer"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}
