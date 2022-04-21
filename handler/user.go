package handler

import (
	"api2/database"
	"api2/model"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

// GetUser handler
// @Summary get user by jwt auth
// @Description get user by jwt auth
// @ID GetUser
// @Accept  json
// @Produce  json
// @Success 200 {string} model.SystemUserBalance	"user"
// @Failure 400 {object} model.MessageModel "Server error"
// @Failure 401 {object} model.MessageModel "Bearer auth required"
// @Router /user [get]
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
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	case nil:
		log.Println("Ok")
	default:
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.MessageModel{Message: "error"})
	}

	return c.Status(fiber.StatusOK).JSON(SystemUserBalance)
}
