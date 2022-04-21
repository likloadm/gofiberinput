package model

type SystemUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Pass string `json:"password"`
}

type SystemUsers struct {
	SystemUsers []SystemUser `json:"system_users"`
}

type SystemUserBalance struct {
	Id      int
	Name    string  `json:"name"`
	Pass    string  `json:"password"`
	Balance float64 `json:"balance"`
}
