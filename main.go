package main

import (
	"api2/config"
	"api2/database"
	_ "api2/docs"
	"api2/handler"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
	"log"
)

// @title           Go Fiber Input
// @version         1.0
// @description     API for controlling the receipt and withdrawal of funds

// @host      127.0.0.1:3000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Bearer
// @description					registers by /register, login by /login and get access jwt token, later use token for add transaction

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Post("/login", handler.Login)

	app.Post("/register", handler.Register)

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Config("JWT_SECRET_KEY")),
	}))

	app.Post("/transaction", handler.Transaction)

	app.Get("/transaction", handler.GetTransaction)

	app.Get("/user", handler.GetUser)

	app.Listen(":3000")
}
