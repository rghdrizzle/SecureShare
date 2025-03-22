package main

import (
	// "rghdrizzle/SecureShare/controllers"
	"rghdrizzle/SecureShare/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main(){

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders:  "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, OPTIONS",
	}))
	routes.UserRoute(app)
	app.Listen(":3333")
}