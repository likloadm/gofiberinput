package handler

import (
	"api2/database"
	"api2/model"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

type TypeTransactionNotFound struct{}

func (m *TypeTransactionNotFound) Error() string {
	return "Type Transaction Not Found"
}

type ZeroAmount struct{}

func (m *ZeroAmount) Error() string {
	return "Zero Amount"
}

// type_transaction
// 1 - Ввод
// 2 - Вывод

// status_transaction
// 0 - внесена
// 1 - исполнена
// 2 - заблокирована

// Transaction handler
// @Summary Add a new transaction to the database and change balance
// @Description get id by type transaction and amount > 0
// @ID transaction
// @Accept  json
// @Produce  json
// @Param   amount      formData   number     true  "Transaction amount(>0)" 100
// @Param   type        formData   int         true  "1 output, 2 input" 1
// @Success 200 {object} model.SystemUserBalance	"ok"
// @Failure 400 {object} model.MessageModel "Server error"
// @Failure 401 {object} model.MessageModel "Bearer auth required"
// @Router /transaction [post]
func Transaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	p := new(model.TransactionUser)

	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	modelSystemUser := model.SystemUserBalance{}

	tx, err := database.DB.Begin()

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	row := tx.QueryRow("SELECT balance FROM system_user WHERE id=$1", id)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	// Throws Unauthorized error
	err = row.Scan(&modelSystemUser.Balance)

	switch err {
	case sql.ErrNoRows:
		log.Println("No rows were returned! Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(model.MessageModel{Message: "error"})
	case nil:
		log.Println("Ok")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	err = TransactionProcess(id, p, &modelSystemUser, tx)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	return c.Status(fiber.StatusOK).JSON(p)
}

// GetTransaction handler
// @Summary Get transaction by id
// @Description get transaction by id
// @ID GetTransaction
// @Accept  json
// @Produce  json
// @Param   id      formData   int     true  "Transaction id"
// @Success 200 {object} model.TransactionUserSystem	"ok"
// @Failure 400 {object} model.MessageModel "Server error"
// @Failure 401 {object} model.MessageModel "Bearer auth required"
// @Router /transaction [get]
func GetTransaction(c *fiber.Ctx) error {
	p := new(model.TransactionUser)
	TransactionSystemUser := model.TransactionUserSystem{}

	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	if p.Id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	row := database.DB.QueryRow("SELECT * FROM transaction_user WHERE id=$1", p.Id)

	// Throws Unauthorized error
	switch err := row.Scan(&TransactionSystemUser.Id, &TransactionSystemUser.Amount, &TransactionSystemUser.Type, &TransactionSystemUser.Status, &TransactionSystemUser.Sender); err {
	case sql.ErrNoRows:
		log.Println("transaction not found")
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	case nil:
		log.Println("Ok")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	return c.Status(fiber.StatusOK).JSON(TransactionSystemUser)
}

func TransactionProcess(id float64, p *model.TransactionUser, modelSystemUser *model.SystemUserBalance, tx *sql.Tx) error {
	transactionStatus := 0
	if p.Amount <= 0 {
		return &ZeroAmount{}
	}

	switch p.Type {
	case 1:
		if p.Amount > modelSystemUser.Balance {
			transactionStatus = 2
		} else {
			rows, err := tx.Query("UPDATE system_user set balance = balance - $1 where id = $2", p.Amount, id)
			rows.Close()
			transactionStatus = 1
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		row := tx.QueryRow("INSERT INTO transaction_user (amount, type_transaction, status, sender) VALUES ($1, $2, $3, $4) RETURNING id", p.Amount, p.Type, transactionStatus, id)

		switch err := row.Scan(&p.Id); err {
		case sql.ErrNoRows:
			tx.Rollback()
			return err
		case nil:
			log.Println("Ok")
		default:
			tx.Rollback()
			return err
		}
		err := tx.Commit()
		return err
	case 2:
		rows, err := tx.Query("UPDATE system_user set balance = balance + $1 where id = $2", p.Amount, id)
		rows.Close()
		if err != nil {
			tx.Rollback()
			return err
		}
		transactionStatus = 1

		row := tx.QueryRow("INSERT INTO transaction_user (amount, type_transaction, status, sender) VALUES ($1, $2, $3, $4) RETURNING id", p.Amount, p.Type, transactionStatus, id)
		switch err := row.Scan(&p.Id); err {
		case sql.ErrNoRows:
			tx.Rollback()
			return err
		case nil:
			log.Println("Ok")
		default:
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		return err
	default:
		return &TypeTransactionNotFound{}
	}
}
