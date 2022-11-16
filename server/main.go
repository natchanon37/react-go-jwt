package main

import (
	"react-go-jwt/server/database"
	"react-go-jwt/server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main(){

	//Connect database
	database.Connect()

	//Init fiber
	app:=fiber.New()

	//Middleware
	app.Use((cors.New(cors.Config{
		AllowCredentials: true,
	})))

	routes.Setup(app)

	app.Listen(":8000")
}