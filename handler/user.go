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
// @Success 200 {string} string	"user"
// @Failure 400 {string} string "Server error"
// @Failure 401 {string} string "Bearer auth required"
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
