package database

func CreateSystemUser() {
	DB.Query(`CREATE TABLE IF NOT EXISTS system_user (
    id SERIAL PRIMARY KEY,
    balance float8,
    name text UNIQUE,
    password_hash text
)
`)
}

// 1 - Вывод
// 2 - Ввод

func CreateTransactionsUser() {
	DB.Query(`CREATE TABLE IF NOT EXISTS transaction_user (
    id SERIAL PRIMARY KEY,
    amount float8,
    type_transaction integer,
    status integer,
    sender integer,
    FOREIGN KEY (sender) REFERENCES system_user (id) ON DELETE CASCADE
)
`)
}
