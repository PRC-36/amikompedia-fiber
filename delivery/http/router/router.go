package router

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	AuthMiddleware     fiber.Handler
	RegisterController controller.RegisterController
	SurveyController   controller.SurveyController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {

	// Register
	c.App.Post("/api/v1/users/_register", c.RegisterController.Register)

	// Survey
	c.App.Post("/api/v1/surveys", c.SurveyController.Create)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

}
