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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login handler
// @Summary login user by name and password
// @Description login user by name and password
// @ID Login
// @Accept  json
// @Produce  json
// @Param   name      body   string     true  "User name"
// @Param   pass      body   string     true  "User password"
// @Success 200 {string} string	"access JWT token"
// @Failure 400 {string} string "Server error"
// @Failure 401 {string} string "Bearer auth required"
// @Router /login [post]
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

// Register handler
// @Summary register user by name and password
// @Description register user by name and password
// @ID Register
// @Accept  json
// @Produce  json
// @Param   name      body   string     true  "User name"
// @Param   pass      body   string     true  "User password"
// @Success 200 {string} string	"access JWT token"
// @Failure 400 {string} string "Server error"
// @Failure 401 {string} string "Bearer auth required"
// @Router /register [post]
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
