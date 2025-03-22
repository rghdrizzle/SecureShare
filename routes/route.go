package routes

import (
	"github.com/gofiber/fiber/v2"
	"rghdrizzle/SecureShare/controllers"
)

func UserRoute(app *fiber.App){
	app.Post("/upload",controller.FileUpload)
}