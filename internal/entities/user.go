package entities

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
