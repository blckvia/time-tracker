package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type Users struct {
	ID             int    `json:"-"`
	PassportSeries string `json:"passport_series"`
	PassportNumber string `json:"passport_number"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type GetAllUsers struct {
	Meta  Meta    `json:"meta"`
	Users []Users `json:"users"`
}

type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UpdateUsers struct {
	PassportSeries string `json:"passport_series"`
	PassportNumber string `json:"passport_number"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type UserStats struct {
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Patronymic     string `json:"patronymic"`
	PassportSeries string `json:"passport_series"`
	PassportNumber string `json:"passport_number"`
	Address        string `json:"address"`
	Tasks          []Task `json:"tasks"`
	OverallTime    string `json:"overall_time"`
}

func (t *Task) MarshalJSON() ([]byte, error) {
	type TaskAlias Task
	return json.Marshal(&struct {
		*TaskAlias
		OverallTimeString string `json:"overall_time"`
	}{
		TaskAlias:         (*TaskAlias)(t),
		OverallTimeString: FormatDuration(t.OverallTime),
	})
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
