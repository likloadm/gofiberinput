package handler

import (
	"api2/config"
	"api2/database"
	"api2/model"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type TypeTransactionNotFound struct{}

func (m *TypeTransactionNotFound) Error() string {
	return "Type Transaction Not Found"
}

type ZeroAmount struct{}

func (m *ZeroAmount) Error() string {
	return "Zero Amount"
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	p := new(model.SystemUser)
	user := model.SystemUser{}

	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	row := database.DB.QueryRow("SELECT id, password_hash FROM system_user WHERE name=$1", p.Name)

	// Throws Unauthorized error
	switch err := row.Scan(&user.Id, &user.Pass); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	case nil:
		log.Println("Rows detected")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	if !CheckPasswordHash(p.Pass, user.Pass) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Config("JWT_SECRET_KEY")))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	return c.JSON(fiber.Map{"token": t})
}

func Accessible(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": nil,
	})
}

func Register(c *fiber.Ctx) error {
	p := new(model.SystemUser)

	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	if p.Pass == "" || p.Name == "" {
		log.Println("Name or pass not found")
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "Name or pass not found",
		})
	}

	PasswordHash, err := HashPassword(p.Pass)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	_, err = database.DB.Query("INSERT INTO system_user (name, password_hash, balance) VALUES ($1, $2, $3)", p.Name, PasswordHash, 0)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": nil,
	})
}

// type_transaction
// 1 - Ввод
// 2 - Вывод

// status_transaction
// 0 - внесена
// 1 - исполнена
// 2 - заблокирована

func Transaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	p := new(model.TransactionUser)

	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	modelSystemUser := model.SystemUserBalance{}

	tx, err := database.DB.Begin()

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	row := tx.QueryRow("SELECT balance FROM system_user WHERE id=$1", id)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	// Throws Unauthorized error
	err = row.Scan(&modelSystemUser.Balance)

	switch err {
	case sql.ErrNoRows:
		log.Println("No rows were returned! Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	case nil:
		log.Println("Ok")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	err = TransactionProcess(id, p, &modelSystemUser, tx)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": modelSystemUser.Id,
	})
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

		switch err := row.Scan(&modelSystemUser.Id); err {
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
		switch err := row.Scan(&modelSystemUser.Id); err {
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

func GetTransaction(c *fiber.Ctx) error {
	p := new(model.TransactionUser)
	TransactionSystemUser := model.TransactionUserSystem{}

	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	if p.Id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "transaction id not found",
		})
	}

	row := database.DB.QueryRow("SELECT * FROM transaction_user WHERE id=$1", p.Id)

	// Throws Unauthorized error
	switch err := row.Scan(&TransactionSystemUser.Id, &TransactionSystemUser.Amount, &TransactionSystemUser.Type, &TransactionSystemUser.Status, &TransactionSystemUser.Sender); err {
	case sql.ErrNoRows:
		log.Println("transaction not found")
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "transaction not found",
		})
	case nil:
		log.Println("Ok")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"amount": TransactionSystemUser.Amount,
		"id":     TransactionSystemUser.Id,
		"type":   TransactionSystemUser.Type,
		"status": TransactionSystemUser.Status,
		"sender": TransactionSystemUser.Sender,
	})
}

func GetUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	SystemUserBalance := model.SystemUserBalance{}

	row := database.DB.QueryRow("SELECT id, balance, name FROM system_user WHERE id=$1", id)

	// Throws Unauthorized error
	switch err := row.Scan(&SystemUserBalance.Id, &SystemUserBalance.Balance, &SystemUserBalance.Name); err {
	case sql.ErrNoRows:
		log.Println("user not found")
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "transaction not found",
		})
	case nil:
		log.Println("Ok")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"balance": SystemUserBalance.Balance,
		"id":      SystemUserBalance.Id,
		"name":    SystemUserBalance.Name,
	})
}
