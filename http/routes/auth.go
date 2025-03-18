package routes

import (
	"github.com/fevzisahinler/Appointfy-Backend/http/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	app.Post("auth/login", controllers.Login)
	app.Post("auth/register", controllers.Register)

}
