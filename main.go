package main

import (
	"api2/config"
	"api2/database"
	"api2/handler"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Post("/login", handler.Login)

	app.Post("/register", handler.Register)

	app.Get("/", handler.Accessible)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Config("JWT_SECRET_KEY")),
	}))

	app.Post("/transaction", handler.Transaction)

	app.Get("/transaction", handler.GetTransaction)

	app.Get("/user", handler.GetUser)

	app.Listen(":3000")
}
