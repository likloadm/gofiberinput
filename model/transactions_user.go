package model

type TransactionUser struct {
	Amount float64 `json:"amount"`
	Type   int     `json:"type"`
	Id     int     `json:"id"`
}

type TransactionUserSystem struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
	Type   int     `json:"type"`
	Status int     `json:"status"`
	Sender int     `json:"sender"`
}

type TransactionsUser struct {
	TransactionsUser []TransactionUser `json:"transactions"`
}
